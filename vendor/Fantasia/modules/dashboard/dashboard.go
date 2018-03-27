package dashboard

import (
	"Fantasia/system"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/mem"

	"github.com/shirou/gopsutil/cpu"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

const trackerSleepDuration = time.Second * 1

//genmodules:config
//go:generate go-bindata-assetfs -pkg dashboard assets/index.html assets/dist/build.js

// Config ...
type Config struct {
	// Address to host the server on
	Address string

	// Password used to log into the dashboard
	Password string

	// Set to true if you want to use a custom asset directory
	CustomAssets bool

	// AssetDirectory contains the server files
	AssetDirectory string

	// LogRequests logs requests to STDOUT
	LogRequests bool

	// Time format charts will use
	TimeFormat string
}

// NewConfig returns the default config
func NewConfig() *Config {
	c := &Config{
		Address:    "9090",
		Password:   "remilia",
		TimeFormat: "15:04",
	}
	return c
}

// Module ...
type Module struct {
	Config *Config
	Server http.Server
	// tmpl   *template.Template

	Stats []*Stats
}

// Log logs data
func (m *Module) Log(data ...interface{}) {
	log.Println(append([]interface{}{"DASHBOARD: "}, data...)...)
}

// Build ...
func (m *Module) Build(sys *system.System) {
	m.Log("Dashboard is a WIP")

	r := mux.NewRouter()
	m.ConstructRoutes(r)
	m.TrackStats()

	m.Server = http.Server{
		Addr:    ":" + m.Config.Address,
		Handler: r,
	}

	go m.Server.ListenAndServe()
}

// ConstructRoutes constructs the dashboard's routes
func (m *Module) ConstructRoutes(r *mux.Router) {
	var assetdir http.FileSystem

	if m.Config.LogRequests {
		r.Use(func(h http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, h)
		})
	}

	// Static file server
	if m.Config.CustomAssets {
		assetdir = http.Dir(m.Config.AssetDirectory)
		m.Log("Custom asset directory set to: ", m.Config.AssetDirectory)
	} else {
		assetdir = assetFS()
	}

	r.HandleFunc("/api/stats/{name}/", m.statsHandler)
	r.PathPrefix("/").Handler(http.FileServer(assetdir))
}

// TrackStats ...
func (m *Module) TrackStats() {
	statsLimit := 1000
	m.Stats = append(m.Stats,
		NewStats("mem", statsLimit),
		NewStats("cpu", statsLimit),
		NewStats("messages", statsLimit),
	)
	go m.TrackCPU()
	go m.TrackMem()
}

// TrackCPU ...
func (m *Module) TrackCPU() {
	c := m.findStats("cpu")
	for {
		percent, err := cpu.Percent(0, false)
		if err != nil {
			m.Log("error getting CPU percentage")
			continue
		}

		c.Push(int(percent[0]), time.Now().Format("15:04"))

		time.Sleep(trackerSleepDuration)
	}
}

// TrackMem ...
func (m *Module) TrackMem() {
	c := m.findStats("mem")
	for {
		memory, err := mem.VirtualMemory()
		if err != nil {
			m.Log("error getting mem percentage")
			continue
		}

		c.Push(int(memory.UsedPercent), time.Now().Format("15:04"))

		time.Sleep(trackerSleepDuration)
	}
}

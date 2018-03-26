package dashboard

import (
	"Fantasia/system"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

//genmodules:config
//go:generate go-bindata-assetfs -pkg dashboard assets/...

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
}

// NewConfig returns the default config
func NewConfig() *Config {
	c := &Config{
		Address:  "9090",
		Password: "remilia",
	}
	return c
}

// Module ...
type Module struct {
	Config *Config
	Server http.Server
	// tmpl   *template.Template
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

	m.Server = http.Server{
		Addr:    ":" + m.Config.Address,
		Handler: r,
	}

	go m.Server.ListenAndServe()
}

// ConstructRoutes constructs the dashboard's routes
func (m *Module) ConstructRoutes(r *mux.Router) {
	var assetdir http.FileSystem

	r.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	})

	// Static file server
	if m.Config.CustomAssets {
		assetdir = http.Dir(m.Config.AssetDirectory)
		m.Log("Custom asset directory set to: ", m.Config.AssetDirectory)
	} else {
		assetdir = assetFS()
	}

	r.PathPrefix("/").Handler(http.FileServer(assetdir))
}

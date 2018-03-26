package dashboard

import (
	"Fantasia/system"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//genmodules:config
//go:generate go-bindata -debug -pkg dashboard assets/...

// Config ...
type Config struct {
	// Address to host the server on
	Address string

	// Password used to log into the dashboard
	Password string

	// TemplateDir is the directory to parse templates from.
	// If none is supplied, the embedded assets will be used.
	// EX: templates/* will parse all files in the dir templates.
	TemplatePath string
}

// NewConfig returns the default config
func NewConfig() *Config {
	c := &Config{
		Address:  ":9090",
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
}

// ConstructRoutes constructs the dashboard's routes
func (m *Module) ConstructRoutes(r *mux.Router) {

}

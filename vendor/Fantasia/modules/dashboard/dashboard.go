package dashboard

import (
	"Fantasia/system"
	"log"
)

//genmodules:config

// Config ...
type Config struct {
	// Password used to log into the dashboard
	Password string
}

// NewConfig returns the default config
func NewConfig() *Config {
	c := &Config{}
	return c
}

// Module ...
type Module struct {
}

// Build ...
func (m *Module) Build(sys *system.System) {
	log.Println("Dashboard is a WIP")
}

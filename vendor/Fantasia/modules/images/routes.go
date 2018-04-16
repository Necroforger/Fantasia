package images

// CreateCommands adds the image commands
func (m *Module) CreateCommands() {
	r := m.Sys.CommandRouter
	r.On("hue", m.CmdHue).Set("", "adjusts the hue of the supplied image; ex: hue [degree]")
}

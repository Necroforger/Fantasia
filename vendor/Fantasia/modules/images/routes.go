package images

// CreateCommands adds the image commands
func (m *Module) CreateCommands() {
	r := m.Sys.CommandRouter

	// Adjustments
	r.On("hue", m.CmdHue).Set("", "adjusts the hue of the supplied image;\nex: `hue [degree]`")
	r.On("saturation", m.CmdSaturation).Set("", "Adjusts the saturation of an image;\nex: `saturation [value]`")
	r.On("contrast", m.CmdContrast).Set("", "Adjusts the contrast of an image;\nex: `contrast [value]`")
	r.On("gamma", m.CmdGamma).Set("", "Adjusts the gamma of an image;\nex: `gamma [value]`")
	r.On("brightness", m.CmdBrightness).Set("", "Adjusts the brightness of an image;\nex: `brightness [value]`")

	// Effects
	r.On("invert", m.CmdInvert).Set("", "Inverts an image")
}

package information

import "github.com/Necroforger/Fantasia/system"

// Module ...
type Module struct{}

// Build adds this modules commands to the system's router
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.CurrentCategory = "information"

	r.On("help", m.Help).
		Set("", "Lists the available commands")
}

// Help ...
func (m *Module) Help(ctx *system.Context) {
	ctx.ReplyStatus(system.StatusWarning, "Help command not yet implemented")
}

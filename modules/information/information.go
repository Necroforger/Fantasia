package information

import (
	"fmt"

	"github.com/Necroforger/Fantasia/system"
)

// Module ...
type Module struct{}

// Build adds this modules commands to the system's router
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.SetCategory("information")

	r.On("help", m.Help).
		Set("", "Lists the available commands")

	sub, _ := system.NewSubCommandRouter("test")
	sub.Router.On("args", m.Argtest)

	r.AddSubrouter(sub)

}

// Help ...
func (m *Module) Help(ctx *system.Context) {
	ctx.ReplyStatus(system.StatusWarning, "Help command not yet implemented")
}

// Argtest ...
func (m *Module) Argtest(ctx *system.Context) {
	ctx.ReplyStatus(system.StatusNotify, fmt.Sprint(ctx.Args))
}

package eval

import "Fantasia/system"

// Module ...
type Module struct{}

// Build ...
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.On("evaljs", m.EvalJS).Set("", "evaluates javascript code\nusage: `evaljs [text]`\nDo not enter any arguments to enter REPL mode")
}

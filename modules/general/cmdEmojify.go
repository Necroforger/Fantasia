package general

import (
	"strings"
	"time"

	"github.com/Necroforger/Fantasia/system"
)

func (m *Module) emojifyCommand(ctx *system.Context) {
	text := ""
	for _, v := range ctx.Args.After() {
		if len(text)+len(":regional_indicator__:") >= 1998 {
			ctx.Reply(text)
			text = ""
			time.Sleep(time.Millisecond * 500)
		}
		if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') {
			text += ":regional_indicator_" + strings.ToLower(string(v)) + ":"
		}
		if v == ' ' {
			text += "   "
		}
	}
	if text == "" {
		return
	}
	ctx.Reply(text)
}

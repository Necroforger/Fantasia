package general

import (
	"math/rand"
	"strings"

	"github.com/Necroforger/Fantasia/system"
)

// CmdChoose chooses between the listed options
func CmdChoose(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter your choices")
		return
	}

	choices := strings.Split(ctx.Args.After(), ",")
	ctx.ReplyNotify("I choose **" + strings.TrimSpace(choices[int(rand.Float64()*float64(len(choices)))]) + "**")
}

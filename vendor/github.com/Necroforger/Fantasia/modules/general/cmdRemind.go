package general

import (
	"strconv"
	"time"

	"github.com/Necroforger/Fantasia/system"
)

// CmdRemind reminds the user their input.
func CmdRemind(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter how long to wait in seconds and then a query")
		return
	}

	duration, err := strconv.Atoi(ctx.Args.Get(0))
	if err != nil {
		ctx.ReplyError("Error parsing duration")
		return
	}

	ctx.ReplyNotify("<@"+ctx.Msg.Author.ID+">", " I will notify you in ", duration, " seconds")
	time.Sleep(time.Second * time.Duration(duration))

	ctx.ReplyNotify("<@"+ctx.Msg.Author.ID+">\n", ctx.Args.AfterN(1))
}

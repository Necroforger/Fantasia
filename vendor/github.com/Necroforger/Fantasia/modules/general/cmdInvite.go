package general

import (
	"fmt"

	"github.com/Necroforger/Fantasia/system"
)

// CmdInviteURL produces an invite for the bot
func CmdInviteURL(ctx *system.Context) {
	if ctx.Ses.DG.State.User == nil {
		ctx.ReplyError("The bot user is nil. You should never see this error")
		return
	}

	ctx.ReplyNotify(fmt.Sprintf("https://discordapp.com/api/oauth2/authorize?client_id=%s&permissions=0&scope=bot", ctx.Ses.DG.State.User.ID))
}

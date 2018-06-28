package general

import (
	"github.com/Necroforger/Fantasia/system"
	"github.com/bwmarrin/discordgo"
	"github.com/Necroforger/dream"
)

// CmdAvatar returns the user's avatar
func CmdAvatar(ctx *system.Context) {
	var (
		user *discordgo.User
		err  error
	)

	if len(ctx.Msg.Mentions) > 0 {
		user = ctx.Msg.Mentions[0]
	} else if ctx.Args.After() != "" {
		user, err = ctx.Ses.DG.User(ctx.Args.After())
		if err != nil {
			ctx.ReplyError("Could not find userid")
			return
		}
	} else {
		user = ctx.Msg.Author
	}

	ctx.ReplyEmbed(dream.NewEmbed().
		SetImage(user.AvatarURL("2048")).
		SetTitle(user.Username).
		SetColor(system.StatusNotify).
		MessageEmbed)
}

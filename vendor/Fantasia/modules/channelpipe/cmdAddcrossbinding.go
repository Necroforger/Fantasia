package channelpipe

import (
	"Fantasia/system"
)

// TODO

// CmdAddCrossBinding binds two channels together
func (m *Module) CmdAddCrossBinding(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("You need to supply two channel IDs to use thing command")
		return
	}

	guildID, channelID, dstID, err := GetBindingArguments(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	bindingTo, err := CreateBinding(ctx.Ses, guildID, channelID, dstID)
	if err != nil {
		ctx.ReplyError("Error creating binding from channel1 to channel2; ", err)
		return
	}

	guildToID, err := ctx.Ses.GuildID(dstID)
	if err != nil {
		ctx.ReplyError("Error finding destination guild; Is the bot in the guild you want to bind to?")
		return
	}

	bindingFrom, err := CreateBinding(ctx.Ses, guildToID, dstID, channelID)
	if err != nil {
		ctx.ReplyError("Error creating binding from channel2 to channel1; ", err)
		return
	}

	err = m.AddBinding(bindingFrom)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	err = m.AddBinding(bindingTo)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	ctx.ReplySuccess("Crossbound " + channelID + " to " + dstID)

	m.SaveBindings()
}

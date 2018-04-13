package channelpipe

import (
	"Fantasia/system"
)

// CmdAddBinding binds a source to a sink. (channel to webhook, or other channel)
func (m *Module) CmdAddBinding(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("You need to enter the destination webhook or channel to bind to")
	}

	guildID, channelID, dstID, err := GetBindingArguments(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	binding, err := CreateBinding(ctx.Ses, guildID, channelID, dstID)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	err = m.AddBinding(binding)
	if err != nil {
		ctx.ReplyError("A binding from source " + binding.Source.ChannelID + " to destination " + binding.Sink.GetDest() + " already exists")
		return
	}

	ctx.ReplySuccess("Created binding from channel [" + binding.Source.ChannelID + "] to channel [" + binding.Sink.ChannelID() + "]")

	m.SaveBindings()
}

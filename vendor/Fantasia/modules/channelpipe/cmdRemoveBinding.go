package channelpipe

import "Fantasia/system"

// CmdRemoveBinding removes a binding from a channel
func (m *Module) CmdRemoveBinding(ctx *system.Context) {
	var (
		channelID string
		dstID     string
	)

	// Determine channelID and dstID
	if len(ctx.Args) == 1 { //         [dst]
		channelID = ctx.Msg.ChannelID
		dstID = ctx.Args.After()
	} else if len(ctx.Args) == 2 { //  [channelid] [dst]
		channelID = ctx.Args.Get(0)
		dstID = ctx.Args.Get(1)
	}

	err := m.RemoveBinding(channelID, dstID)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	ctx.ReplySuccess("Removed binding")
	m.SaveBindings()
}

package channelpipe

import (
	"github.com/Necroforger/Fantasia/system"
)

// CmdRemoveCrossBinding removes a crossbinding
func (m *Module) CmdRemoveCrossBinding(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please supply a channelID(s)")
		return
	}

	_, channelID, dstID, err := GetBindingArguments(ctx)
	if err != nil {
		ctx.ReplyError(err)
	}

	// check that the bindings exist
	_, err = m.Binding(channelID, dstID)
	if err != nil {
		ctx.ReplyError("Could not find bindingTo")
		return
	}

	_, err = m.Binding(dstID, channelID)
	if err != nil {
		ctx.ReplyError("Could not find bindingFrom")
		return
	}

	// Remove bindings from slice
	m.RemoveBinding(channelID, dstID)
	m.RemoveBinding(dstID, channelID)
	m.SaveBindings()

	// Delete webhooks
	err = DeleteChannelWebhookByName(ctx.Ses, channelID, dstID)
	if err != nil {
		ctx.ReplyError("Error deleting webhook 1: ", err)
		return
	}

	err = DeleteChannelWebhookByName(ctx.Ses, dstID, channelID)
	if err != nil {
		ctx.ReplyError("Error deleting webhook 2: ", err)
		return
	}

	ctx.ReplySuccess("Removed crossbinding")
}

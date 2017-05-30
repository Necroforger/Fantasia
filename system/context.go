package system

import (
	"fmt"

	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
)

//////////////////////////////////
// 		Context
/////////////////////////////////

// Context contains information about the command.
type Context struct {
	Msg          *discordgo.Message
	Ses          *dream.Bot
	System       *System
	CommandRoute *CommandRoute
	Args         Args
}

// ReplyStatus sends an embed to the message channel the command was received on
// Coloured with the given status code.
//		status: 		Colour code of the message to send
// 		notification: 	The content of the status message
func (c *Context) ReplyStatus(status int, notification string) (*discordgo.Message, error) {
	return c.Ses.DG.ChannelMessageSendEmbed(c.Msg.ChannelID,
		dream.
			NewEmbed().
			SetDescription(notification).
			SetColor(status).MessageEmbed,
	)
}

// Reply replys to the channel the context originated from
//		text: Content of the message to send
func (c *Context) Reply(i ...interface{}) (*discordgo.Message, error) {
	return c.Ses.DG.ChannelMessageSend(c.Msg.ChannelID, fmt.Sprint(i...))
}

// ReplyEmbed replys to the channel the context originated from with the given embed
//		embed: the discordgo messageembed to reply with
func (c *Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Ses.DG.ChannelMessageSendEmbed(c.Msg.ChannelID, embed)
}

// ReplyError replys with the given error value
func (c *Context) ReplyError(i interface{}) (*discordgo.Message, error) {
	return c.ReplyStatus(StatusError, fmt.Sprint(i))
}

// SendStatus sends a status embed to the given channel
//		Status:			The colour code of the embed to send
//		channelID:		The ID of the channel to send to
//		notification:	The content of the status message
func (c *Context) SendStatus(status int, channelID, notification string) (*discordgo.Message, error) {
	return c.Ses.DG.ChannelMessageSendEmbed(channelID,
		dream.
			NewEmbed().
			SetDescription(notification).
			SetColor(status).MessageEmbed,
	)
}

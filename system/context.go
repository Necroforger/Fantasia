package system

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/Necroforger/dream"
)

//////////////////////////////////
// 		Context
/////////////////////////////////

// Context contains information about the command.
type Context struct {
	Msg          *discordgo.Message
	Ses          *dream.Session
	System       *System
	CommandRoute *CommandRoute
	Args         Args
}

// Guild obtains the guild of the message sent over the context
func (c *Context) Guild() (*discordgo.Guild, error) {
	s := c.Ses.DG

	channel, err := s.State.Channel(c.Msg.ChannelID)
	if err != nil {
		channel, err = s.Channel(c.Msg.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		guild, err = s.Guild(channel.GuildID)
		if err != nil {
			return nil, err
		}
	}

	return guild, nil
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
//	i... : The objects to send through the embed
func (c *Context) ReplyError(i ...interface{}) (*discordgo.Message, error) {
	return c.ReplyStatus(StatusError, fmt.Sprint(i...))
}

// ReplyNotify replys with the given error value
//	i... : The objects to send through the embed
func (c *Context) ReplyNotify(i ...interface{}) (*discordgo.Message, error) {
	return c.ReplyStatus(StatusNotify, fmt.Sprint(i...))
}

// ReplyWarning replys with the given error value
//	i... : The objects to send through the embed
func (c *Context) ReplyWarning(i ...interface{}) (*discordgo.Message, error) {
	return c.ReplyStatus(StatusWarning, fmt.Sprint(i...))
}

// ReplySuccess replys with the given error value
//	i... : The objects to send through the embed
func (c *Context) ReplySuccess(i ...interface{}) (*discordgo.Message, error) {
	return c.ReplyStatus(StatusSuccess, fmt.Sprint(i...))
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

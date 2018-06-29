package system

import (
	"fmt"

	"github.com/Necroforger/dream"
	"github.com/bwmarrin/discordgo"
)

//////////////////////////////////
// 		Context
/////////////////////////////////

// Context contains information about the command.
type Context struct {
	data         map[string]interface{}
	Msg          *discordgo.Message
	Ses          *dream.Session
	System       *System
	CommandRoute *CommandRoute
	Args         Args
}

// Set saves a value
func (c *Context) Set(key string, data interface{}) {
	if c.data == nil {
		c.data = map[string]interface{}{}
	}

	c.data[key] = data
}

// Get retrieves a value
func (c *Context) Get(key string) interface{} {
	if c.data == nil {
		return nil
	}
	if v, ok := c.data[key]; ok {
		return v
	}
	return nil
}

// IsAdmin checks if the message author has administrator privileges
func (c *Context) IsAdmin() (bool, error) {
	isAdminInGuild := func() (bool, error) {
		return MemberHasPermission(c.System.Dream.DG, c.Msg.GuildID, c.Msg.Author.ID, discordgo.PermissionAdministrator)
	}

	gs, err := c.System.DB.GetGuild(c.Msg.GuildID)
	if err != nil {
		if err != ErrNotFound {
			return false, err
		}
		b, err := isAdminInGuild()
		if err != nil {
			return false, err
		}
		return c.System.IsAdmin(c.Msg.Author.ID) || b, nil
	}

	b, err := isAdminInGuild()
	if err != nil {
		return false, err
	}

	return gs.IsAdmin(c.Msg.Author.ID) || c.System.IsAdmin(c.Msg.Author.ID) || b, nil
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

// Channel obtains the channel of the message sent
func (c *Context) Channel() (*discordgo.Channel, error) {
	s := c.Ses.DG
	channel, err := s.State.Channel(c.Msg.ChannelID)
	if err != nil {
		channel, err = s.Channel(c.Msg.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	return channel, nil
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

package channelpipe

import (
	"Fantasia/system"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Necroforger/dream"

	"github.com/bwmarrin/discordgo"
)

// Errors
var (
	ErrWebhookNotFound = errors.New("Error: Webhook not found")
)

// ContentFromMessage extracts creates a Content struct from a discordgo message
func ContentFromMessage(m *discordgo.Message) *discordgo.WebhookParams {
	c := &discordgo.WebhookParams{}
	c.Username = m.Author.Username
	c.AvatarURL = m.Author.AvatarURL("")
	c.Embeds = m.Embeds
	c.Content = m.Content
	for _, v := range m.Attachments {
		c.Content += v.URL + "\n"
	}
	return c
}

// FetchWebhook ...
func FetchWebhook(path string) (*discordgo.Webhook, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Webpage did not return status 200")
	}

	// Decode returned JSON data
	var w discordgo.Webhook
	err = json.NewDecoder(resp.Body).Decode(&w)
	return &w, err
}

// DeleteChannelWebhookByName deletes a channel webhook by name
func DeleteChannelWebhookByName(s *dream.Session, channelID, webhookname string) error {
	webhook, err := ChannelWebhookByName(s, channelID, webhookname)
	if err != nil {
		return err
	}
	return s.DG.WebhookDelete(webhook.ID)
}

// CreateBinding creates a binding with the given options
func CreateBinding(s *dream.Session, guildID, channelID, dstID string) (*Binding, error) {
	var sink Sink
	if strings.HasPrefix(dstID, "http") {
		hook, err := FetchWebhook(dstID)
		if err != nil {
			return nil, err
		}
		sink = NewWebhookSink(hook)
	} else {
		hook, err := CreateOrGetChannelWebhook(s, dstID, channelID)
		if err != nil {
			return nil, err
		}
		sink = NewWebhookSink(hook)
	}

	// Create and save binding
	binding := &Binding{
		Source: Source{channelID, guildID},
		Sink:   sink,
	}
	return binding, nil
}

// CreateChannelWebhook creates a webhook on a channel
func CreateChannelWebhook(s *dream.Session, channelID string, webhookname string) (*discordgo.Webhook, error) {
	return s.DG.WebhookCreate(channelID, webhookname, "")
}

// CreateOrGetChannelWebhook creates or gets an existing webhook with name webhookname on the given channelID
func CreateOrGetChannelWebhook(s *dream.Session, channelID string, webhookname string) (*discordgo.Webhook, error) {
	w, err := ChannelWebhookByName(s, channelID, webhookname)
	if err != nil {
		if err == ErrWebhookNotFound {
			return CreateChannelWebhook(s, channelID, webhookname)
		}
		return nil, err
	}
	return w, nil
}

// ChannelWebhookByName returns a channel webhook by name
func ChannelWebhookByName(s *dream.Session, channelID string, webhookname string) (*discordgo.Webhook, error) {
	webhooks, err := s.DG.ChannelWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, w := range webhooks {
		if w.Name == webhookname {
			return w, nil
		}
	}

	return nil, ErrWebhookNotFound
}

// CrossBind binds two bindings so that messages are passed back and forth between both channels.
func CrossBind(s *dream.Session, from, to *Binding) error {
	return nil
}

// GetBindingArguments returns the arguments from a binding command
func GetBindingArguments(ctx *system.Context) (guildID, channelID, dstID string, err error) {
	// Determine channelID and dstID
	if len(ctx.Args) == 1 { //         [dst]
		channelID = ctx.Msg.ChannelID
		dstID = ctx.Args.After()
	} else if len(ctx.Args) == 2 { //  [channelid] [dst]
		channelID = ctx.Args.Get(0)
		dstID = ctx.Args.Get(1)
	}

	var c *discordgo.Channel
	// Verify channelID
	c, err = ctx.Ses.Channel(channelID)
	if err != nil {
		err = errors.New("invalid channel ID")
		return
	}
	guildID = c.GuildID

	return
}

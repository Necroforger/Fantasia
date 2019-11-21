package channelpipe

import (
	"github.com/bwmarrin/discordgo"
)

// WebhookSink ...
type WebhookSink struct {
	Webhook *discordgo.Webhook `json:"webhook"`
}

// NewWebhookSink webhook address
//   dst : webhook URL
func NewWebhookSink(hook *discordgo.Webhook) *WebhookSink {
	return &WebhookSink{hook}
}

// Send sends content over the webhook
func (w *WebhookSink) Send(s *discordgo.Session, message *discordgo.WebhookParams) (*discordgo.Message, error) {
	return s.WebhookExecute(w.Webhook.ID, w.Webhook.Token, false, message)
}

// GetDest returns the sink destination
func (w *WebhookSink) GetDest() string {
	return discordgo.EndpointWebhookToken(w.Webhook.ID, w.Webhook.Token)
}

// ChannelID returns the ChannelID of the webhook sink
func (w *WebhookSink) ChannelID() string {
	return w.Webhook.ChannelID
}

// ID returns the webhook ID
func (w *WebhookSink) ID() string {
	return w.Webhook.ID
}

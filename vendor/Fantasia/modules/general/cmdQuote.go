package general

import (
	"errors"

	"Fantasia/system"
	"github.com/Necroforger/dream"
	humanize "github.com/dustin/go-humanize"
)

// CmdQuote ...
func CmdQuote(ctx *system.Context) {
	quote, err := makeQuote(ctx.Ses, ctx.Msg.ChannelID, ctx.Args.After())
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	ctx.ReplyEmbed(quote.SetColor(system.StatusNotify).MessageEmbed)
}

func makeQuote(b *dream.Session, channelID, messageID string) (*dream.Embed, error) {
	msg, err := b.DG.ChannelMessage(channelID, messageID)
	if err != nil {
		if chann, err := b.DG.State.Channel(channelID); err == nil {
			found := false
			for _, v := range chann.Messages {
				if v.ID == messageID {
					msg = v
					found = true
				}
			}
			if !found {
				return nil, errors.New("Message not found")
			}
		} else {
			return nil, errors.New("Channel not found")
		}
	}

	quoteEmbed := dream.NewEmbed().
		SetAuthor(msg.Author.Username, msg.Author.AvatarURL("512")).
		SetDescription(msg.Content)
	if timestamp, err := msg.Timestamp.Parse(); err == nil {
		quoteEmbed.SetFooter(humanize.Time(timestamp))
	}

	return quoteEmbed, nil
}

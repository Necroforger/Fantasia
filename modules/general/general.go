package general

import (
	"time"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/dream"
	humanize "github.com/dustin/go-humanize"
)

// Module ...
type Module struct{}

// Build ...
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.On("ping", Ping).Set("", "responds with the amount of time taken to send and retrieve a message")
	r.On("snowflake", Snowflake).Set("", "gives the creation date of a discord ID")
	r.On("emojify", m.emojifyCommand).Set("", "Emojifies the given text")
	r.On("hex", HexDisplay).Set("", "Returns an image representation of the given hex code. example: `hex ff00ff`")
}

// Ping returns the time taken to send a message and recieve back the discord event
func Ping(ctx *system.Context) {
	embed := dream.NewEmbed().
		SetDescription("pinging").
		SetColor(system.StatusWarning).
		MessageEmbed

	succEmbed := dream.NewEmbed().
		SetColor(system.StatusSuccess).
		MessageEmbed

	start := time.Now()
	m, err := ctx.Ses.DG.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
	if err == nil {
		succEmbed.Description = time.Since(start).String()
		ctx.Ses.DG.ChannelMessageEditEmbed(m.ChannelID, m.ID, succEmbed)
	}
}

// Snowflake gives the creation date of a discord ID
func Snowflake(ctx *system.Context) {
	t, err := dream.CreationTime(ctx.Args.After())
	if err == nil {

		ctx.ReplyEmbed(dream.NewEmbed().
			SetTitle(ctx.Args.After()).
			SetDescription(
				t.Format(
					time.RFC1123) + "\n" +
					humanize.Time(t),
			).
			SetColor(system.StatusNotify).
			MessageEmbed)
	}
}

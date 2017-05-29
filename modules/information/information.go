package information

import (
	"fmt"
	"strings"
	"time"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
	"github.com/Necroforger/ytdl"
)

// Module ...
type Module struct{}

// Build adds this modules commands to the system's router
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.SetCategory("information")

	// Subrouter must be added before adding commands because it sets the current category
	// For the commands to be added to as the parent routers category.
	test, _ := system.NewSubCommandRouter("^"+r.Prefix+"test", "test")
	r.AddSubrouter(test)

	test.Router.On("help", m.Depthmap).Set("", "Lists the available commands")
	test.Router.On("args", m.Argtest).Set("", "Displays your arguments")
	test.Router.On("play", m.Play).Set("", "plays the given youtube URL")
	test.Router.On("ping", func(ctx *system.Context) {
		start := time.Now()
		m, err := ctx.ReplyStatus(system.StatusNotify, "pinging...")
		if err == nil {
			ctx.Ses.DG.ChannelMessageEditEmbed(m.ChannelID, m.ID, dream.NewEmbed().
				SetDescription(fmt.Sprint(time.Since(start))).
				SetTitle("pong").
				SetColor(system.StatusSuccess).
				MessageEmbed)
		}
	})

	sub, _ := system.NewSubCommandRouter("^ meme", "meme")
	test.Router.AddSubrouter(sub)

	sub.Router.On("cachigga", func(ctx *system.Context) { ctx.ReplyStatus(system.StatusWarning, "warning you have been cachiggad") })
}

// Help ...
func (m *Module) Help(ctx *system.Context) {
	ctx.ReplyStatus(system.StatusWarning, "Help command not yet implemented")
}

// Argtest ...
func (m *Module) Argtest(ctx *system.Context) {
	var text string
	for i, v := range ctx.Args {
		text += fmt.Sprintf("%d:\t[%s]\n", i, v)
	}
	ctx.ReplyStatus(system.StatusNotify, text)
}

// Depthmap maps the depth of subrouters and their commands to an embed
func (m *Module) Depthmap(ctx *system.Context) {

	depthString := func(text string, depth int, subrouter bool) string {
		quote := ""
		if subrouter {
			quote = "`"
		}
		return strings.Repeat("  ", depth) + quote + text + quote + "\n"
	}

	var depthcharge func(r *system.CommandRouter, embed *dream.Embed, depth int) *dream.Embed
	depthcharge = func(r *system.CommandRouter, embed *dream.Embed, depth int) *dream.Embed {
		if embed == nil {
			embed = dream.NewEmbed()
		}

		getField := func(name string) *discordgo.MessageEmbedField {
			for _, v := range embed.Fields {
				if v.Name == name {
					return v
				}
			}
			if name == "" {
				name = "undefined"
			}
			field := &discordgo.MessageEmbedField{Name: name}
			embed.Fields = append(embed.Fields, field)
			return field
		}

		for _, v := range r.Routes {
			field := getField(v.Category)
			field.Value += depthString(v.Name, depth, false)
		}

		for _, v := range r.Subrouters {
			field := getField(v.Category())
			field.Value += depthString(v.Name, depth, true)
			embed = depthcharge(v.Router, embed, depth+1)
		}

		return embed
	}

	_, err := ctx.ReplyEmbed(depthcharge(ctx.System.CommandRouter, nil, 0).SetColor(system.StatusNotify).
		SetThumbnail(ctx.Ses.DG.State.User.AvatarURL("2048")).
		InlineAllFields().
		SetDescription("subcommands are represented by indentation.").
		MessageEmbed)
	if err != nil {
		ctx.Reply(fmt.Sprint(err))
	}
}

// Play plays the given song
func (m *Module) Play(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyStatus(system.StatusError, "No arguments provided")
		return
	}

	vs, err := ctx.Ses.UserVoiceState(ctx.Msg.Author.ID)
	if err != nil {
		ctx.ReplyStatus(system.StatusError, fmt.Sprint(err))
		return
	}

	vc, err := ctx.Ses.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
	if err != nil {
		ctx.ReplyStatus(system.StatusError, fmt.Sprint(err))
		return
	}

	info, err := ytdl.GetVideoInfo(ctx.Args.After())
	if err != nil {
		ctx.ReplyStatus(system.StatusError, fmt.Sprint(err))
		return
	}

	stream, err := system.YoutubeDLFromInfo(info)
	if err != nil {
		ctx.ReplyStatus(system.StatusError, fmt.Sprint(err))
		return
	}

	disp := ctx.Ses.PlayStream(vc, stream)

	ctx.ReplyStatus(system.StatusSuccess, "playing `"+info.Title+"`\ntype 'stop' to stop playing the video")

	for {
		var msg *discordgo.MessageCreate
		for msg = ctx.Ses.NextMessageCreate(); msg.Author.ID != ctx.Msg.Author.ID; msg = ctx.Ses.NextMessageCreate() {
		}
		if msg.Content == "stop" {
			disp.Stop()
			ctx.ReplyStatus(system.StatusNotify, "Video stopped")
			return
		}
	}
}

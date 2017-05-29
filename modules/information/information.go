package information

import (
	"fmt"
	"strings"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
	"github.com/rylio/ytdl"
)

// Module ...
type Module struct{}

// Build adds this modules commands to the system's router
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.SetCategory("information")

	// Testing subrouter
	test, _ := system.NewSubCommandRouter("^"+r.Prefix+"test", "test")
	r.AddSubrouter(test)

	test.Router.On("help", m.Help).Set("", "Lists the available commands")
	test.Router.On("args", m.Argtest).Set("", "Displays your arguments")
	test.Router.On("depthcharge", m.Depthmap).Set("", "Displays the depth of the command routers commands")
	test.Router.On("play", m.Play).Set("", "plays the given youtube URL")

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
		text += fmt.Sprintf("%d:\t%s\n", i, v)
	}
	ctx.ReplyStatus(system.StatusNotify, text)
}

// Depthmap maps the depth of subrouters and their commands
func (m *Module) Depthmap(ctx *system.Context) {

	depthString := func(text string, depth int, subrouter bool) string {
		quote := ""
		if subrouter {
			quote = "**"
		}
		return strings.Repeat("\t", depth) + quote + text + quote + "\n"
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
			field := getField(v.Category)
			field.Value += depthString(v.Name, depth, true)
			embed = depthcharge(v.Router, embed, depth+1)
		}

		return embed
	}

	_, err := ctx.ReplyEmbed(depthcharge(ctx.System.CommandRouter, nil, 0).MessageEmbed)
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

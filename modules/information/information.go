package information

import (
	"strings"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
)

// Module ...
type Module struct{}

// Build adds this modules commands to the system's router
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.SetCategory("information")

	r.On("help", m.Help).Set("", "help", "displays a help menu with the available commands")
}

// Help maps a list of available commands and descends into subrouters.
func (m *Module) Help(ctx *system.Context) {

	var depthcharge func(r *system.CommandRouter, embed *dream.Embed, depth int) *dream.Embed

	_, err := ctx.ReplyEmbed(depthcharge(ctx.System.CommandRouter, nil, 0).SetColor(system.StatusNotify).
		SetThumbnail(ctx.Ses.DG.State.User.AvatarURL("2048")).
		InlineAllFields().
		SetDescription("subcommands are represented by indentation.").
		MessageEmbed)
	if err != nil {
		ctx.ReplyError(err)
	}
}

func depthcharge(r *system.CommandRouter, embed *dream.Embed, depth int) *dream.Embed {
	if embed == nil {
		embed = dream.NewEmbed()
	}

	depthString := func(text string, depth int, subrouter bool) string {
		quote := ""
		if subrouter {
			quote = "`"
		}
		return strings.Repeat("  ", depth) + quote + text + quote + "\n"
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

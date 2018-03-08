package information

import (
	"errors"
	"strings"

	"Fantasia/system"
	"github.com/bwmarrin/discordgo"
	"github.com/Necroforger/dream"
)

// Module ...
type Module struct{}

// Build adds this modules commands to the system's router
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.On("help", m.Help).Set("", "Displays a help menu with the available commands")
}

// Help maps a list of available commands and descends into subrouters.
func (m *Module) Help(ctx *system.Context) {

	if cmd := ctx.Args.After(); cmd != "" {
		if route, _ := ctx.System.CommandRouter.FindMatch(cmd); route != nil {
			ctx.ReplyEmbed(dream.NewEmbed().
				SetTitle(route.Name).
				SetDescription(route.Desc).
				SetColor(system.StatusNotify).
				MessageEmbed)
			return
		}
		ctx.ReplyError(errors.New("Command not found"))
		return
	}

	_, err := ctx.ReplyEmbed(depthcharge(ctx.System.CommandRouter, nil, 0).
		SetColor(system.StatusNotify).
		SetThumbnail(ctx.Ses.DG.State.User.AvatarURL("2048")).
		InlineAllFields().
		SetDescription("type `help [command]` for more information\nCommands separated with `|` represent alternative names.\nIndented commands are subroutes of their parent commands").
		MessageEmbed)
	if err != nil {
		ctx.ReplyError(err)
	}
}

// Depthcharge recursively generates a help embed from a CommandRouter and its subrouters
func depthcharge(r *system.CommandRouter, embed *dream.Embed, depth int) *dream.Embed {
	if embed == nil {
		embed = dream.NewEmbed()
	}

	depthString := func(text string, depth int, subrouter bool) string {
		quote := ""
		if subrouter {
			quote = "`"
		}
		return strings.Repeat("\t", depth) + quote + text + quote + "\n"
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

	removeField := func(name string) {
		for i, v := range embed.Fields {
			if v.Name == name {
				embed.Fields = append(embed.Fields[:i], embed.Fields[i+1:]...)
				return
			}
		}
	}

	for _, v := range r.Routes {
		field := getField(v.Category)

		var tag string
		if !v.Disabled {
			field.Value += depthString(tag+v.Name+tag, depth, false)
		}

		if field.Value == "" {
			removeField(v.Category)
		}
	}

	for _, v := range r.Subrouters {
		field := getField(v.Category())
		field.Value += depthString(v.Name, depth, true)
		embed = depthcharge(v.Router, embed, depth+1)
	}

	return embed
}
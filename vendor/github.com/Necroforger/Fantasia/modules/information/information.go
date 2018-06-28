package information

import (
	"errors"
	"sort"
	"strings"

	"github.com/Necroforger/Fantasia/system"

	"github.com/Necroforger/dream"
	"github.com/bwmarrin/discordgo"
)

// Module ...
type Module struct{}

// Build adds this modules commands to the system's router
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.On("help", m.Help).Set("", "Displays a help menu with the available commands")
}

// CategoriesBySize sorts categories by size
type CategoriesBySize []*discordgo.MessageEmbedField

// Size ...
func (c CategoriesBySize) Size() int {
	return len(c)
}

// Len ...
func (c CategoriesBySize) Len() int {
	return len(c)
}

// Less ...
func (c CategoriesBySize) Less(a, b int) bool {
	return len(c[a].Value) < len(c[b].Value)
}

// Swap ...
func (c CategoriesBySize) Swap(a, b int) {
	c[a], c[b] = c[b], c[a]
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

	embed := depthcharge(ctx.System.CommandRouter, nil, 0)

	for _, v := range embed.Fields {
		v.Value = "```\n" + v.Value + "```"
	}

	sort.Sort(sort.Reverse(CategoriesBySize(embed.Fields)))

	_, err := ctx.ReplyEmbed(embed.
		SetColor(system.StatusNotify).
		// SetThumbnail(ctx.Ses.DG.State.User.AvatarURL("512")).
		InlineAllFields().
		SetDescription("`Bot prefix: " + ctx.System.Config.Prefix + "` type `help [command]` for more information\nCommands separated with `|` represent alternative names.\nIndented commands are subroutes of their parent commands").
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
		if v.Disabled || (v.CommandRoute != nil && v.CommandRoute.Disabled) {
			continue
		}
		field := getField(v.Category())
		field.Value += depthString(v.Name, depth, true)
		embed = depthcharge(v.Router, embed, depth+1)
	}

	return embed
}

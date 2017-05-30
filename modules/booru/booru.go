package booru

import (
	"strconv"

	"github.com/Necroforger/Boorudl/extractor"
	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/dream"
)

// Module ...
type Module struct{}

// Build ...
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter

	AddBooru(r, "http://danbooru.donmai.us", "danbooru")
	AddBooru(r, "https://safebooru.org/", "safebooru")

}

// AddBooru adds a booru command to the router
func AddBooru(r *system.CommandRouter, booruURL string, commandName string) {
	r.On(commandName, MakeBooruSearcher(booruURL)).
		Set("", "Returns an image result from ["+commandName+"]("+booruURL+")\n"+
			"Usage: `"+commandName+" [tags] [post index]`\n"+
			"Enclose the tag list in quotes to include multiple tags")
}

// MakeBooruSearcher returns a command that searches the given booru link
func MakeBooruSearcher(booruURL string) func(*system.Context) {
	return func(ctx *system.Context) {
		index := 0
		if n, err := strconv.Atoi(ctx.Args.Get(1)); err == nil {
			index = n
		}

		posts, err := extractor.Search(booruURL, extractor.SearchQuery{
			Limit:  index + 1,
			Page:   0,
			Tags:   ctx.Args.Get(0),
			Random: false,
		})
		if err != nil {
			ctx.ReplyError(err)
		}

		if index < len(posts) {
			post := posts[index]
			ctx.ReplyEmbed(dream.NewEmbed().
				SetColor(system.StatusNotify).
				SetImage(post.ImageURL).
				MessageEmbed)
		}
	}
}

package images

//genmodules:config

import (
	"strconv"

	"github.com/Necroforger/Boorudl/extractor"
	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/dream"
)

// Config ...
type Config struct {
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}

// Module ...
type Module struct {
	Config *Config
}

// Build ...
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter

	AddBooru(r, "http://danbooru.donmai.us", "danbooru")
	AddBooru(r, "https://safebooru.org/", "safebooru")
	AddBooru(r, "http://google.com", "googleimg")

}

// AddBooru adds a booru command to the router
func AddBooru(r *system.CommandRouter, booruURL string, commandName string) {
	r.On(commandName, MakeBooruSearcher(booruURL)).
		Set("", "Returns an image result from ["+commandName+"]("+booruURL+")\n"+
			"Usage: `"+commandName+" [tags] [post index] [to post index]`\n"+
			"Enclose the tag list in quotes to include multiple tags")
}

// MakeBooruSearcher returns a command that searches the given booru link
func MakeBooruSearcher(booruURL string) func(*system.Context) {
	return func(ctx *system.Context) {
		index := 0
		if n, err := strconv.Atoi(ctx.Args.Get(1)); err == nil {
			index = n
		}

		indexTo := index
		if n, err := strconv.Atoi(ctx.Args.Get(2)); err == nil {
			indexTo = n
		}

		if indexTo > index+10 {
			ctx.ReplyError("You cannot bulk view more than 10 images at a time")
		}

		posts, err := extractor.Search(booruURL, extractor.SearchQuery{
			Limit:  indexTo + 1,
			Page:   0,
			Tags:   ctx.Args.Get(0),
			Random: false,
		})
		if err != nil {
			ctx.ReplyError(err)
		}

		for i := index; i < indexTo; i++ {
			if i >= 0 && i < len(posts) {
				post := posts[i]
				ctx.ReplyEmbed(dream.NewEmbed().
					SetColor(system.StatusNotify).
					SetImage(post.ImageURL).
					MessageEmbed)
			}
		}

	}
}

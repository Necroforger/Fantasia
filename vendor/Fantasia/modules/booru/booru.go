package booru

//genmodules:config
import (
	"Fantasia/system"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Necroforger/Boorudl/extractor"
	"github.com/Necroforger/dgwidgets"
	"github.com/Necroforger/dream"
)

// Config ...
type Config struct {
	BooruCommandsCategory string
	BooruCommands         [][]string
}

// NewConfig returns the default configuration
func NewConfig() *Config {
	return &Config{
		// Default booru commands
		BooruCommands: [][]string{
			{"googleimg", "http://google.com"},
			{"safebooru", "https://safebooru.org"},
			{"danbooru", "http://danbooru.donmai.us"},
		},
	}
}

// Module ...
type Module struct {
	Config *Config
}

// Build ...
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	if m.Config.BooruCommandsCategory != "" {
		r.CurrentCategory = m.Config.BooruCommandsCategory
	}

	for _, v := range m.Config.BooruCommands {
		if len(v) < 2 {
			log.Println("error creating booru command " + fmt.Sprint(v) + ", array must be in the form of [command name, booru url]")
			continue
		}
		AddBooru(r, v[0], v[1])
	}
}

// AddBooru adds a booru command to the router
func AddBooru(r *system.CommandRouter, commandName string, booruURL string) {
	r.On(commandName, MakeBooruSearcher(booruURL, true)).
		Set("", "Returns an image result from ["+commandName+"]("+booruURL+")\n"+
			"Usage: `"+commandName+" [tags] ~index[-indexTo]`\n"+
			"Example: `"+commandName+" cirno~0-10` would return a list of posts from index 0-10.\n"+
			"You can omit the 0 to fetch the 10'th post Ex: `"+commandName+" cirno~10`")
}

// MakeBooruSearcher returns a command that searches the given booru link
//     booruURL   : address of booru to search.
//     enforceSFW : enforces that the command will only return SFW results in non-nsfw channels.
func MakeBooruSearcher(booruURL string, enforceSFW bool) func(*system.Context) {
	return func(ctx *system.Context) {
		var (
			index   = 0
			limit   = 1
			indexTo = 1
			tags    string
		)

		if s := strings.Split(ctx.Args.After(), "~"); len(s) > 1 {
			r := strings.Split(s[1], "-")
			if n, err := strconv.Atoi(r[0]); err == nil {
				index = n
				tags = s[0]
			}
			if len(r) > 1 {
				if n, err := strconv.Atoi(r[1]); err == nil {
					indexTo = n
				}
			} else {
				indexTo = index + 1
			}
		} else {
			tags = ctx.Args.After()
		}

		if indexTo-index < 0 {
			ctx.ReplyError("The second index cannot be less than the first")
			return
		}

		if indexTo-index > 1 {
			limit = 30
		}

		page := index / limit
		index %= limit

		func() {
			if enforceSFW {
				channel, err := ctx.Channel()
				if err != nil {
					ctx.ReplyError("Error obtaning channel information ", err)
					return
				}
				if channel.NSFW {
					return
				}
				t := strings.Split(tags, " ")
				for _, v := range t {
					if v == "rating:safe" || v == "rating:s" {
						return
					}
				}
				tags = strings.Join(append(t, "rating:safe"), " ")
			}
		}()

		ctx.ReplyNotify("You are searching with tags\n", tags)

		posts, err := extractor.Search(booruURL, extractor.SearchQuery{
			Limit:  limit,
			Page:   page,
			Tags:   tags,
			Random: true,
		})
		if err != nil {
			ctx.ReplyError("Error searching for posts: ", err)
			return
		}

		if len(posts) == 0 {
			ctx.ReplyError("No posts found")
			return
		}

		if limit == 1 {
			_, err = ctx.ReplyEmbed(dream.NewEmbed().
				SetColor(system.StatusNotify).
				SetImage(posts[0].ImageURL).
				MessageEmbed)
			if err != nil {
				ctx.ReplyError(err)
				ctx.ReplyNotify(posts[0].ImageURL)
			}
		} else {
			paginator := dgwidgets.NewPaginator(ctx.Ses.DG, ctx.Msg.ChannelID)
			for idx, post := range posts {
				if idx < index || idx > indexTo {
					continue
				}
				paginator.Add(dream.NewEmbed().
					SetColor(system.StatusNotify).
					SetImage(post.ImageURL).
					MessageEmbed)
			}
			paginator.SetPageFooters()
			paginator.Widget.Timeout = time.Minute * 3
			paginator.DeleteReactionsWhenDone = true

			paginator.Spawn()
		}
	}
}

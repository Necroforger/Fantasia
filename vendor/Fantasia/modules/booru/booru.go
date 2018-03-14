// package booru

package booru

//genmodules:config
import (
	"Fantasia/system"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Necroforger/Boorudl/extractor"
	"github.com/Necroforger/dgwidgets"
	"github.com/Necroforger/dream"
	"github.com/bwmarrin/discordgo"
)

// Config ...
type Config struct {
	BooruCommandsCategory string
	BooruCommands         [][]string
}

// NewConfig returns the default configuration
func NewConfig() *Config {
	return &Config{
		// BooruCommands
		// [command name] [booru URL] [true/false; enforce rating:safe in non-nsfw channels]
		BooruCommands: [][]string{
			{"googleimg", "https://google.com", "false"},
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

	r.On("replay", CmdOpenSave).Set("", "replays a saved list of posts")

	for _, v := range m.Config.BooruCommands {
		if len(v) < 2 {
			log.Println("error creating booru command " + fmt.Sprint(v) + ", array must be in the form of [command name, booru url]")
			continue
		}

		enforceSFW := true
		if len(v) > 2 {
			enforceSFW = !(v[2] == "false")
		}

		AddBooru(r, v[0], v[1], enforceSFW)
	}
}

// CmdOpenSave opens a previously saved list of posts
func CmdOpenSave(ctx *system.Context) {
	var posts extractor.Posts
	if err := func() error {
		var fileURL string
		switch {
		case len(ctx.Msg.Attachments) > 0:
			fileURL = ctx.Msg.Attachments[0].URL
		case ctx.Args.After() != "":
			fileURL = ctx.Args.After()
		default:
			ctx.ReplyNotify("Upload a saved list of posts or give a file url")
			var nxtmsg *discordgo.MessageCreate
			for nxtmsg = ctx.Ses.NextMessageCreate(); nxtmsg.Author.ID != ctx.Msg.Author.ID; nxtmsg = ctx.Ses.NextMessageCreate() {
			}
			if len(nxtmsg.Attachments) == 0 {
				fileURL = nxtmsg.Content
			} else {
				fileURL = nxtmsg.Attachments[0].URL
			}
		}
		resp, err := http.Get(fileURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(&posts)
	}(); err != nil {
		ctx.ReplyError(err)
		return
	}

	viewPosts(ctx, posts, "", -1, -1)
}

// AddBooru adds a booru command to the router
func AddBooru(r *system.CommandRouter, commandName string, booruURL string, enforceSFW bool) {
	r.On(commandName, MakeBooruSearcher(booruURL, enforceSFW)).
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
			page    = 0
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
			} else {
				ctx.ReplyError("Error converting index to integer: ", err)
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

		// Custom pagination for google
		if booruURL == "http://google.com" || booruURL == "https://google.com" {
			limit = 100
			page = 0
		} else {
			page = index / limit
			index %= limit
		}

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

		posts, err := extractor.Search(booruURL, extractor.SearchQuery{
			Limit:  limit,
			Page:   page,
			Tags:   tags,
			Random: true,
		})
		if err != nil {
			ctx.ReplyError("Error searching for posts: ", err,
				"\nYour tags were: ", tags,
				"\nPage: ", page,
				"\nindex: ", index)
			return
		}

		if len(posts) == 0 {
			ctx.ReplyError("No posts found")
			return
		}

		if (indexTo - index) == 1 {
			var post extractor.Post
			if index >= 0 && index < len(posts) {
				post = posts[index]
			} else {
				ctx.ReplyError("Post out of bounds, returning the first result")
				post = posts[0]
			}
			_, err = ctx.ReplyEmbed(dream.NewEmbed().
				SetColor(system.StatusSuccess).
				SetImage(post.ImageURL).
				SetTitle("source").
				SetURL(post.ImageURL).
				MessageEmbed)
			if err != nil {
				ctx.ReplyError(err)
				ctx.ReplyNotify(post.ImageURL)
			}
		} else {
			viewPosts(ctx, posts, tags, index, indexTo)
		}
	}
}

func viewPosts(ctx *system.Context, posts extractor.Posts, tags string, index, indexTo int) {
	paginator := dgwidgets.NewPaginator(ctx.Ses.DG, ctx.Msg.ChannelID)
	for idx, post := range posts {
		if (indexTo != -1 && idx > indexTo) || (index != -1 && idx < index) {
			continue
		}
		paginator.Add(dream.NewEmbed().
			SetColor(system.StatusSuccess).
			SetImage(post.ImageURL).
			SetURL(post.ImageURL).
			MessageEmbed)
	}
	paginator.Widget.Handle(dgwidgets.NavSave, func(w *dgwidgets.Widget, s *discordgo.MessageReaction) {
		rd, wr := io.Pipe()
		go func() {
			json.NewEncoder(wr).Encode(posts)
			wr.Close()
		}()
		channel, err := ctx.Ses.DG.UserChannelCreate(s.UserID)
		if err != nil {
			ctx.ReplyError(err)
			return
		}
		ctx.Ses.DG.ChannelMessageSend(channel.ID, "You saved: `"+tags+"`")
		ctx.Ses.DG.ChannelFileSend(channel.ID, "images.json", rd)
	})
	paginator.SetPageFooters()
	paginator.Widget.Timeout = time.Minute * 3
	paginator.DeleteReactionsWhenDone = true

	paginator.Spawn()
}

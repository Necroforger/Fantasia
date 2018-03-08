package images

//genmodules:config

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// Used in imageFromContext
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/Necroforger/Boorudl/extractor"
	"Fantasia/system"
	"Fantasia/util"
	"github.com/Necroforger/dgwidgets"
	"github.com/Necroforger/dream"
)

// Config ...
type Config struct {
	// FilterCommands includes image filtering commands
	ImageCommandsCategory string
	ImageCommands         [][]string
	BooruCommandsCategory string
	BooruCommands         [][]string
	EffectCommands        bool
	ImageEffectsCategory  string
}

// ImageCommand ...
type ImageCommand struct {
	Name string
	URL  string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		// FilterCommands includes filter commands
		EffectCommands: true,

		// Default Image commands
		ImageCommands: [][]string{},

		// Default booru commands
		BooruCommands: [][]string{
			{"danbooru", "http://danbooru.donmai.us"},
			{"safebooru", "https://safebooru.org"},
			{"googleimg", "http://google.com"},
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
	maincategory := r.CurrentCategory

	setCategory := func(name string) {
		if name != "" {
			r.SetCategory(name)
		} else {
			r.SetCategory(maincategory)
		}
	}

	/////////////////////////////////
	// Effects
	////////////////////////////////
	if m.Config.EffectCommands {
		setCategory(m.Config.ImageEffectsCategory)
		r.On("textify", m.CmdTextify).Set("", "Converts the supplied image to text")
		r.On("edgedetect", MakeConvolutionFunc(MatrixEdgeDetect, getDivisor(MatrixEdgeDetect), 1)).Set("", "`usage: edge [iteratins]` Detects the edges of the given image")
		r.On("blur", MakeConvolutionFunc(MatrixGaussian, getDivisor(MatrixGaussian), 1)).Set("", "`usage: blur [iterations]` Gaussian blurs the given image")
		r.On("motionblur", MakeConvolutionFunc(MatrixMotionBlur, getDivisor(MatrixMotionBlur), 1)).Set("", "`usage: motionblue [iterations]` Applies a motion blur to the given image")
		r.On("sharpen", MakeConvolutionFunc(MatrixSharpen, getDivisor(MatrixSharpen), 1)).Set("", "`usage: motionblue [iterations]`, sharpens the given image")
		r.On("filter", cmdCustomFilter).Set("", "Provide a custom image filter")
		r.On("rotate", CmdRotate).Set("pooka", "Rotates the supplied image 90 degrees")
	}

	///////////////////////////////
	//   Booru commands
	///////////////////////////////
	setCategory(m.Config.BooruCommandsCategory)
	for _, v := range m.Config.BooruCommands {
		if len(v) < 2 {
			log.Println("error creating booru command " + fmt.Sprint(v) + ", array must be in the form of [command name, booru url]")
			continue
		}
		AddBooru(r, v[0], v[1])
	}

	////////////////////////////////
	//  Custom image commands
	///////////////////////////////
	setCategory(m.Config.ImageCommandsCategory)
	for _, v := range m.Config.ImageCommands {
		AddImageCommand(r, v)
	}

}

// AddImageCommand makes an image command from an array of strings in the format
// [command name, description, urls...]
func AddImageCommand(r *system.CommandRouter, cmd []string) {
	if len(cmd) < 3 {
		return
	}
	cmdName := cmd[0]

	r.On(cmdName, MakeImageCommand(cmd[2:], true)).Set("", cmd[1])
}

// MakeImageCommand makes an image command
func MakeImageCommand(urls []string, openFiles bool) func(*system.Context) {
	return func(ctx *system.Context) {
		index := int(rand.Float64() * float64(len(urls)))
		path := urls[index]

		// If the path is not a URL, it will check the file system for the image.
		if !strings.HasPrefix(path, "http://") &&
			!strings.HasPrefix(path, "https://") &&
			openFiles {

			f, err := os.Open(path)
			if err != nil {
				ctx.ReplyError(err)
				return
			}

			info, err := f.Stat()
			if err != nil {
				ctx.ReplyError(err)
				return
			}

			if info.IsDir() {
				randFile, err := util.RandomFileInDir(path)
				if err != nil {
					ctx.ReplyError(err)
					return
				}
				ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, randFile.Name(), randFile)
			} else {
				ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, info.Name(), f)
			}
			return
		}

		ctx.ReplyEmbed(dream.NewEmbed().
			SetImage(urls[index]).
			SetColor(system.StatusNotify).
			MessageEmbed)

	}
}

// AddBooru adds a booru command to the router
func AddBooru(r *system.CommandRouter, commandName string, booruURL string) {
	r.On(commandName, MakeBooruSearcher(booruURL)).
		Set("", "Returns an image result from ["+commandName+"]("+booruURL+")\n"+
			"Usage: `"+commandName+" [tags] ~index[-indexTo]`\n"+
			"Example: `"+commandName+" cirno~0-10` would return a list of posts from index 0-10.\n"+
			"You can omit the 0 to fetch the 10'th post Ex: `"+commandName+" cirno~10`")
}

// MakeBooruSearcher returns a command that searches the given booru link
func MakeBooruSearcher(booruURL string) func(*system.Context) {
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

func imageFromContext(ctx *system.Context) (image.Image, error) {
	msg := ctx.Msg

	var imgurl string
	if len(msg.Attachments) != 0 {
		imgurl = msg.Attachments[0].URL
	} else if ctx.Args.After() != "" {
		imgurl = ctx.Args.After()
	}

	resp, err := http.Get(imgurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	return img, err
}

package images

//genmodules:config

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/Necroforger/Boorudl/extractor"
	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/dream"
)

// Config ...
type Config struct {
	// FilterCommands includes image filtering commands
	FilterCommands        bool
	ImageCommandsCategory string
	ImageCommands         [][]string
	BooruCommandsCategory string
	BooruCommands         [][]string
	ImageFiltersCategory  string
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
		FilterCommands: true,

		// Default Image commands
		ImageCommands: [][]string{},

		// Default booru commands
		BooruCommands: [][]string{
			{"danbooru", "http://danbooru.donmai.us"},
			{"safebooru", "https://safebooru.org/"},
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
	// Convolution filters
	////////////////////////////////
	if m.Config.FilterCommands {
		setCategory(m.Config.ImageFiltersCategory)
		r.On("edgedetect", MakeConvolutionFunc(MatrixEdgeDetect, getDivisor(MatrixEdgeDetect), 1)).Set("", "`usage: edge [iteratins]` Detects the edges of the given image")
		r.On("blur", MakeConvolutionFunc(MatrixGaussian, getDivisor(MatrixGaussian), 1)).Set("", "`usage: blur [iterations]` Gaussian blurs the given image")
		r.On("motionblur", MakeConvolutionFunc(MatrixMotionBlur, getDivisor(MatrixMotionBlur), 1)).Set("", "`usage: motionblue [iterations]` Applies a motion blur to the given image")
		r.On("sharpen", MakeConvolutionFunc(MatrixSharpen, getDivisor(MatrixSharpen), 1)).Set("", "`usage: motionblue [iterations]`, sharpens the given image")
		r.On("filter", cmdCustomFilter).Set("", "Provide a custom image filter")
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
				randFile, err := system.RandomFileInDir(path)
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

		indexTo := index + 1
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
			return
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

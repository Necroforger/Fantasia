package images_old

//genmodules:config

import (
	"image"
	"math/rand"
	"net/http"
	"os"
	"strings"

	// Used in imageFromContext
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/Fantasia/util"

	"github.com/Necroforger/dream"
)

// Config ...
type Config struct {
	// FilterCommands includes image filtering commands
	ImageCommandsCategory string
	ImageCommands         [][]string
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

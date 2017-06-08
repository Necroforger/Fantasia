package audio

import (
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/ytdl"
)

//genmodules:config

// Config ...
type Config struct {
	SoundClipCommandsCategory string
	SoundclipCommands         [][]string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		SoundclipCommands: [][]string{},
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

	r.On("play", m.playHandler).Set("", "Plays the requested song")
	r.On("stop", func(ctx *system.Context) { ctx.Ses.GuildAudioDispatcherStop(ctx.Msg) }).Set("", "Stops the guilds currently playing audio dispatcher")
	r.On("pause", func(ctx *system.Context) { ctx.Ses.GuildAudioDispatcherPause(ctx.Msg) }).Set("", "Pauses the guild's currently playing audio dispatcher")
	r.On("resume", func(ctx *system.Context) { ctx.Ses.GuildAudioDispatcherResume(ctx.Msg) }).Set("", "resumes the guild's currently playing audio dispatcher")
	r.On("join", func(ctx *system.Context) { ctx.Ses.UserVoiceStateJoin(ctx.Msg.Author.ID, false, true) }).Set("", "Joins the calling user's voice channel")
	r.On("leave", func(ctx *system.Context) { ctx.Ses.GuildVoiceConnectionDisconnect(ctx.Msg) }).Set("", "Disconnects from the current guild voice channel")

	////////////////////////////////////
	//           Sound clips
	///////////////////////////////////
	setCategory(m.Config.SoundClipCommandsCategory)
	// Create custom soundclip commands
	for _, v := range m.Config.SoundclipCommands {
		if len(v) >= 3 {
			r.On(v[0], MakeSoundclipFunc(v[2:], true)).Set("", v[1])
		}
	}
}

// AddSoundclip adds a soundclip to a commandrouter
func AddSoundclip(r *system.CommandRouter, name, description string, urls []string) {

}

// MakeSoundclipFunc ...
func MakeSoundclipFunc(urls []string, openFiles bool) func(*system.Context) {
	return func(ctx *system.Context) {
		var (
			stream io.Reader
			path   string
		)

		vc, err := system.ConnectToVoiceChannel(ctx)
		if err != nil {
			ctx.ReplyError(err)
			return
		}

		if len(urls) == 1 {
			path = urls[0]
		} else {
			path = urls[int(rand.Float64()*float64(len(urls)))]
		}

		// Initialize file stream
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
				stream = randFile
			} else {
				stream = f
			}

			//Initialize youtube stream if not a file
		} else {
			info, err := ytdl.GetVideoInfo(path)
			if err != nil {
				return
			}
			stream, err = system.YoutubeDLFromInfo(info)
			if err != nil {
				ctx.ReplyError(err)
			}
		}

		ctx.Ses.PlayStream(vc, stream)
	}
}

func (m *Module) playHandler(ctx *system.Context) {
	b := ctx.Ses

	vc, err := system.ConnectToVoiceChannel(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	if ctx.Args.After() == "" {
		ctx.ReplyError("No arguments provided")
		return
	}

	info, err := ytdl.GetVideoInfo(ctx.Args.After())
	if err != nil {
		ctx.ReplyError("Error obtaining video information")
		return
	}

	stream, err := system.YoutubeDLFromInfo(info)
	if err != nil {
		ctx.ReplyError("Error downloading youtube video")
		return
	}

	b.PlayStream(vc, stream)
}

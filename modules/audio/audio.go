package audio

import (
	"errors"
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/ytdl"
)

//genmodules:config

// Config ...
type Config struct {
	SoundclipCommands [][]string
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
	r.On("play", m.playHandler).Set("", "Plays the requested song")
	r.On("stop", func(ctx *system.Context) { ctx.Ses.GuildAudioDispatcherStop(ctx.Msg) }).Set("", "Stops the guilds currently playing audio dispatcher")
	r.On("pause", func(ctx *system.Context) { ctx.Ses.GuildAudioDispatcherPause(ctx.Msg) }).Set("", "Pauses the guild's currently playing audio dispatcher")
	r.On("resume", func(ctx *system.Context) { ctx.Ses.GuildAudioDispatcherResume(ctx.Msg) }).Set("", "resumes the guild's currently playing audio dispatcher")
	r.On("join", func(ctx *system.Context) { ctx.Ses.UserVoiceStateJoin(ctx.Msg.Author.ID, false, true) }).Set("", "Joins the calling user's voice channel")
	r.On("leave", func(ctx *system.Context) { ctx.Ses.GuildVoiceConnectionDisconnect(ctx.Msg) }).Set("", "Disconnects from the current guild voice channel")

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

		// Connect to voice channel and play clip
		vc, err := ConnectToVoiceChannel(ctx)
		if err != nil {
			ctx.ReplyError(err)
			return
		}

		ctx.Ses.PlayStream(vc, stream)
	}
}

// ConnectToVoiceChannel finds and connects to a user's voice channel
func ConnectToVoiceChannel(ctx *system.Context) (*discordgo.VoiceConnection, error) {
	b := ctx.Ses
	msg := ctx.Msg

	vc, err := b.GuildVoiceConnection(ctx.Msg)
	if err != nil {

		// If not currently in a channel, attempt to join the voice channel of the calling user.
		vs, err := b.UserVoiceState(msg.Author.ID)
		if err != nil {
			ctx.ReplyError(errors.New("Could not find user's voice state"))
			return nil, err
		}
		vc, err = b.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			ctx.ReplyError(errors.New("Could not join user's voice channel"))
			return nil, err
		}

	}

	// Confirm that the user is in the correct voice channel.
	// If the user is in another voice channel, join them.
	vs, err := b.UserVoiceState(msg.Author.ID)
	if err == nil && vs.ChannelID != vc.ChannelID && vs.GuildID != vc.GuildID {
		vc, err = b.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			ctx.ReplyError("Could not join user's voice channel")
			return nil, err
		}
	}

	return vc, nil
}

func (m *Module) playHandler(ctx *system.Context) {
	b := ctx.Ses

	vc, err := ConnectToVoiceChannel(ctx)
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

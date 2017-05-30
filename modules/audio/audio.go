package audio

import (
	"errors"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/ytdl"
)

// Module ...
type Module struct{}

// Build ...
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.On("play", m.playHandler).Set("", "Plays the requested song")
	r.On("stop", func(ctx *system.Context) { ctx.Ses.GuildAudioDispatcherStop(ctx.Msg) }).Set("", "Stops the guilds currently playing audio dispatcher")
	r.On("pause", func(ctx *system.Context) { ctx.Ses.GuildAudioDispatcherPause(ctx.Msg) }).Set("", "Pauses the guild's currently playing audio dispatcher")
	r.On("resume", func(ctx *system.Context) { ctx.Ses.GuildAudioDispatcherResume(ctx.Msg) }).Set("", "resumes the guild's currently playing audio dispatcher")
}

func (m *Module) playHandler(ctx *system.Context) {
	b := ctx.Ses
	msg := ctx.Msg

	vc, err := b.GuildVoiceConnection(ctx.Msg)
	if err != nil {

		// If not currently in a channel, attempt to join the voice channel of the calling user.
		vs, err := b.UserVoiceState(msg.Author.ID)
		if err != nil {
			ctx.ReplyError(errors.New("Could not find user's voice state"))
			return
		}
		vc, err = b.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			ctx.ReplyError(errors.New("Could not join user's voice channel"))
			return
		}

	}

	// Confirm that the user is in the correct voice channel.
	// If the user is in another voice channel, join them.
	vs, err := b.UserVoiceState(msg.Author.ID)
	if err == nil && vs.ChannelID != vc.ChannelID && vs.GuildID != vc.GuildID {
		vc, err = b.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			ctx.ReplyError("Could not join user's voice channel")
			return
		}
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

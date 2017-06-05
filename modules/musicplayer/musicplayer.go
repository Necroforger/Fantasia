package musicplayer

/*
	Requirements:
		+ youtube-dl installed to path
		+ dca-rs
		+ ffmpeg

	WIP
	Musicplayer is a more advanced audio module with the ability to queue and
	manage playlists using youtube-dl.

	// TODO
		* Create graphical menu using embeds with searchable buttons
		* Bypass the ratelimit on adding reactions to messages
		* Make youtube-dl an optional dependency and fall back to the 'ytdl' go library if it
		  	is not available.
		* Create a method to queue songs.
		* Be able to play a queue while simultaneously adding songs to it.
*/

import (
	"time"

	"github.com/Necroforger/Fantasia/system"
)

// ControlCooldown is the cooldown for control event commands
const ControlCooldown = time.Millisecond * 1500

// Module ...
type Module struct {
	GuildRadios map[string]*Radio
}

// Build ...
func (m *Module) Build(s *system.System) {
	m.GuildRadios = map[string]*Radio{}

	r, _ := system.NewSubCommandRouter(`^m(usicplayer)?(\s|$)`, "m | musicplayer")
	r.Router.Prefix = "^"
	r.Set("", "musicplayer subrouter, controls the various actions related to music playing")
	s.CommandRouter.AddSubrouter(r)
	t := r.Router

	t.On("play", m.CmdPlay).Set("", "Plays the current queue")
	t.On("stop", m.CmdStop).Set("", "stops the currently playing queue")
	t.On("pause", m.CmdPause).Set("", "Pauses the currently playing song")
	t.On("resume", m.CmdResume).Set("", "Resumes the currently playing song")
	t.On("next", m.CmdNext).Set("", "Loads the next song in the queue")
	t.On("prev|previous", m.CmdPrevious).Set("prev | previous", "Loads the previous song in the queue")
}

// CmdClear clears the current queue
func (m *Module) CmdClear(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
	}
	radio := m.getRadio(guildID)
	radio.Queue.Clear()
}

// CmdPlay should handle
// 		+ Playing a song from a URL
//		+ Starting the queue if no argument is provided and nothing is playing.
func (m *Module) CmdPlay(ctx *system.Context) {
	vc, err := system.ConnectToVoiceChannel(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	radio := m.getRadio(vc.GuildID)
	radio.Queue.Goto(0)

	err = radio.PlayQueue(ctx, vc)
	if err != nil {
		ctx.ReplyError(err)
	}
}

// CmdStop should
//		+ Stop the queue from playing
//		+ Stop the current song from playing
func (m *Module) CmdStop(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		return
	}

	err = m.getRadio(guildID).Stop()
	if err == nil {
		ctx.ReplySuccess("Queue stopped")
	}
}

// CmdResume resumes the guild's currently playing audio dispatcher.
func (m *Module) CmdResume(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	ctx.Ses.GuildAudioDispatcherResume(guildID)
}

// CmdPause Pauses the guilds currently playing audio dispatcher
func (m *Module) CmdPause(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	ctx.Ses.GuildAudioDispatcherPause(guildID)
}

// CmdNext loads the next song in the queue
func (m *Module) CmdNext(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		return
	}
	radio := m.getRadio(guildID)

	if t := time.Now().Sub(radio.ControlLastUsed); t < ControlCooldown && radio.IsRunning() {
		ctx.ReplyError("This command is on cooldown for `", (ControlCooldown - t).String(), "`. Please use `Goto` or provide an integer argument to skip multiple songs quickly")
		return
	}
	radio.ControlLastUsed = time.Now()

	err = radio.Next()
	if err != nil {
		ctx.ReplyError(err)
		return
	}
}

// CmdPrevious loads the previous song in the queue
func (m *Module) CmdPrevious(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		return
	}
	radio := m.getRadio(guildID)

	if t := time.Now().Sub(radio.ControlLastUsed); t < ControlCooldown && radio.IsRunning() {
		ctx.ReplyError("This command is on cooldown for `", (ControlCooldown - t).String(), "`. Please use `Goto` or provide an integer argument to skip multiple songs quickly")
		return
	}
	radio.ControlLastUsed = time.Now()

	err = radio.Previous()
	if err != nil {
		ctx.ReplyError(err)
	}
}

func (m *Module) getRadio(guildID string) *Radio {
	if v, ok := m.GuildRadios[guildID]; ok {
		return v
	}
	r := NewRadio(guildID)
	m.GuildRadios[guildID] = r

	r.Queue.Playlist = []*Song{
		&Song{
			URL: "https://www.youtube.com/watch?v=-yEDM-ogtBY",
		},
		&Song{
			URL: "https://youtu.be/-yNATRuEvsE",
		},
		&Song{
			URL: "https://www.youtube.com/watch?v=b8AIGUYscYo",
		},
		&Song{
			URL: "https://www.youtube.com/watch?v=m2eeyey-87U",
		},
		&Song{
			URL: "https://youtu.be/tqMfWwnLtKg",
		},
	}
	r.Queue.Loop = true

	return r
}

func guildIDFromContext(ctx *system.Context) (string, error) {
	var guildID string

	vs, err := ctx.Ses.UserVoiceState(ctx.Msg.Author.ID)
	if err != nil {
		guildID, err = ctx.Ses.GuildID(ctx.Msg)
		if err != nil {
			return "", err
		}
	} else {
		guildID = vs.GuildID
	}

	return guildID, nil
}

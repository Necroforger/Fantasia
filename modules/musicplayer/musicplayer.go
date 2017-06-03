package musicplayer

/*
	Requirements:
		+ youtube-dl installed to path
		+ dca-rs
		+ ffmpeg

	WIP
	Musicplayer is a more advanced audio module with the ability to queue and
	manage playlists using youtube-dl.

*/

import (
	"github.com/Necroforger/Fantasia/system"
)

// Module ...
type Module struct {
	GuildRadios map[string]*Radio
}

// Build ...
func (m *Module) Build(s *system.System) {
	m.GuildRadios = map[string]*Radio{}

	r, _ := system.NewSubCommandRouter("^m|musicplayer", "m | musicplayer")
	r.Set("", "musicplayer subrouter, controls the various actions related to music playing")
	s.CommandRouter.AddSubrouter(r)

	t := r.Router
	t.On("play", m.CmdPlay).Set("", ` CmdPlay should handle
 		+ Playing a song from a URL
		+ Starting the queue if no argument is provided and nothing is playing.
		+ Unpausing the currently playing song if the current song is paused.`)

	t.On("stop", m.CmdStop).Set("", `CmdStop should
		+ Stop the queue from playing
		+ Stop the current song from playing`)

	t.On("pause", m.CmdPause).Set("", ` CmdPause should
		+ Pause the currently playing song`)
}

// CmdPlay should handle
// 		+ Playing a song from a URL
//		+ Starting the queue if no argument is provided and nothing is playing.
//		+ Unpausing the currently playing song if the current song is paused.
func (m *Module) CmdPlay(ctx *system.Context) {
	vc, err := ctx.Ses.UserVoiceStateJoin(ctx.Msg.Author.ID, false, true)
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

// CmdPause should
//	+ Pause the currently playing song
func (m *Module) CmdPause(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
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
	m.getRadio(guildID).Next()
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
	}

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

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
}

var tempradio = NewRadio("")

// Build ...
func (m *Module) Build(s *system.System) {

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

	tempradio.Queue.Playlist = []*Song{
		&Song{
			URL: "https://www.youtube.com/watch?v=l7GObEWKTwA&feature=youtu.be",
		},
		&Song{
			URL: "https://youtu.be/-yNATRuEvsE",
		},
		&Song{
			URL: "https://www.youtube.com/watch?v=b8AIGUYscYo",
		},
	}
	tempradio.Queue.Loop = true
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

	tempradio.Queue.Goto(0)

	err = tempradio.PlayQueue(ctx, vc)
	if err != nil {
		ctx.ReplyError(err)
	}
}

// CmdStop should
//		+ Stop the queue from playing
//		+ Stop the current song from playing
func (m *Module) CmdStop(ctx *system.Context) {
	tempradio.Stop()
}

// CmdPause should
//	+ Pause the currently playing song
func (m *Module) CmdPause(ctx *system.Context) {

}

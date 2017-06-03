package musicplayer

import (
	"errors"
	"sync"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/ytdl"
)

// Control constants
const (
	AudioStop = iota
	AudioPlay
	AudioNext
	AudioPrevious
)

// Radio controls queueing and playing music over a guild
// Voice connection.
type Radio struct {
	sync.Mutex

	// Silent specifies if the radio should send messages over the context or not
	Silent bool

	// If true, automatically play the next song in the playlist when the last one finishes
	AutoPlay bool
	GuildID  string
	Queue    *SongQueue

	running bool
	control chan int
}

// NewRadio returns a pointer to a new radio
func NewRadio(guildID string) *Radio {
	return &Radio{
		GuildID:  guildID,
		Queue:    NewSongQueue(),
		control:  make(chan int),
		AutoPlay: true,
	}
}

// PlayQueue plays the radio queue
func (r *Radio) PlayQueue(ctx *system.Context, vc *discordgo.VoiceConnection) error {

	// Check if the queue is already playing
	r.Lock()
	if r.running {
		r.Unlock()
		return errors.New("Radio queue already playing")
	}
	r.running = true
	r.Unlock()

	b := ctx.Ses

	defer func() {
		r.Lock()
		r.running = false
		r.Unlock()
	}()

	for {

		song, err := r.Queue.Song()
		if err != nil {
			return err
		}
		info, err := ytdl.GetVideoInfo(song.URL)
		if err != nil {
			return err
		}
		stream, err := system.YoutubeDLFromInfo(info)
		if err != nil {
			return err
		}
		disp := b.PlayStream(vc, stream)
		done := make(chan bool)
		go func() {
			disp.Wait()
			done <- true
		}()

		select {
		case ctrl := <-r.control:
			switch ctrl {
			case AudioStop: // Stop the queue from playing
				disp.Stop()
				return nil
			case AudioPrevious: // Attempt to load the previous song in the queue
				r.Queue.Previous()
			case AudioNext: // Load the next song if it has been requested early
				err = r.Queue.Next()
				if err != nil {
					return err
				}
			}
		case <-done: // Attempt to load the next song by default
			err = r.Queue.Next()
			if err != nil {
				return err
			}
		}
	}
}

// Next ...
func (r *Radio) Next() {

}

// Previous ...
func (r *Radio) Previous() {
	r.Lock()
	defer r.Unlock()
	if r.running {
		r.control <- AudioPrevious
	}
}

// Stop stops the playing queue
func (r *Radio) Stop() {
	r.Lock()
	defer r.Unlock()

	if r.running {
		r.control <- AudioStop
	}
}

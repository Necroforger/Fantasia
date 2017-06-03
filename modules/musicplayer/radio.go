package musicplayer

import (
	"errors"
	"sync"
	"time"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
	"github.com/Necroforger/ytdl"
)

// Control constants
const (
	AudioStop = iota
	AudioContinue
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

	// Used to prevent commands from being spammed.
	ControlLastUsed time.Time
}

// NewRadio returns a pointer to a new radio
func NewRadio(guildID string) *Radio {
	return &Radio{
		GuildID:  guildID,
		Queue:    NewSongQueue(),
		control:  make(chan int, 100),
		AutoPlay: true,
		Silent:   false,
	}
}

// PlayQueue plays the radio queue
func (r *Radio) PlayQueue(ctx *system.Context, vc *discordgo.VoiceConnection) error {

	// Check if the queue is already playing
	r.Lock()
	if r.running {
		r.Unlock()
		return errors.New("Queue already playing")
	}
	r.running = true
	r.Unlock()

	defer func() {
		r.Lock()
		r.running = false
		r.Unlock()
	}()

	for {
		ctx.Reply("Creating audio dispatcher")
		disp, err := r.Play(ctx.Ses, vc)
		if err != nil {
			return err
		}

		done := make(chan bool)
		go func() {
			ctx.Reply("Waiting for dispatcher to end")
			disp.Wait()
			done <- true
		}()

		select {
		case ctrl := <-r.control:
			switch ctrl {
			case AudioStop:
				ctx.Reply("Received audio stop command")
				disp.Stop()
				return nil
			case AudioContinue:
				ctx.Reply("Received audio continue command")
				disp.Stop()
				continue
			}

		case <-done:
			// Load the next song by default
			ctx.Reply("Loading next song")
			err = r.Queue.Next()
			if err != nil {
				return err
			}

		}
	}
}

// Play plays a single song in the queue
func (r *Radio) Play(b *dream.Bot, vc *discordgo.VoiceConnection) (*dream.AudioDispatcher, error) {
	song, err := r.Queue.Song()
	if err != nil {
		return nil, err
	}
	info, err := ytdl.GetVideoInfo(song.URL)
	if err != nil {
		return nil, err
	}
	stream, err := system.YoutubeDLFromInfo(info)
	if err != nil {
		return nil, err
	}
	disp := b.PlayStream(vc, stream)
	return disp, nil
}

// Next ...
func (r *Radio) Next() error {
	err := r.Queue.Next()
	if err != nil {
		return err
	}
	if r.IsRunning() {
		r.control <- AudioContinue
	}
	return nil
}

// Previous ...
func (r *Radio) Previous() error {
	err := r.Queue.Previous()
	if err != nil {
		return err
	}
	if r.IsRunning() {
		r.control <- AudioContinue
	}
	return nil
}

// IsRunning returns true if the player is currently running
func (r *Radio) IsRunning() bool {
	r.Lock()
	running := r.running
	r.Unlock()
	return running
}

// Stop stops the playing queue
func (r *Radio) Stop() error {

	r.Lock()
	if r.running {
		r.Unlock()

		r.control <- AudioStop
		return nil
	}
	return errors.New("Audio player not running")
}

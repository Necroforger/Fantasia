package musicplayer

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/Fantasia/util"
	"github.com/Necroforger/Fantasia/youtubeapi"
	"github.com/bwmarrin/discordgo"
	"github.com/Necroforger/dream"
	"github.com/Necroforger/ytdl"
)

// Control constants
const (
	AudioStop = iota
	AudioContinue
)

// Radio controls queueing and playing music over a  guild
// Voice connection.
type Radio struct {
	sync.Mutex

	// Silent specifies if the radio should send messages over the context or not
	Silent bool

	// If true, automatically play the next song in the playlist when the last one finishes
	AutoPlay bool
	GuildID  string
	Queue    *SongQueue

	Dispatcher *dream.AudioDispatcher

	running bool
	control chan int

	// UseYoutubeDL specifies which downloader to use when playing videos.
	// If set to true, videos will be downloaded using youtube-dl rather than the golang lib.
	UseYoutubeDL bool

	// Used to prevent commands from being spammed.
	ControlLastUsed time.Time
}

// NewRadio returns a pointer to a new radio
func NewRadio(guildID string) *Radio {
	return &Radio{
		GuildID:  guildID,
		Queue:    NewSongQueue(),
		control:  make(chan int),
		AutoPlay: true,
		Silent:   false,
	}
}

// PlayQueue plays the radio queue
func (r *Radio) PlayQueue(ctx *system.Context, vc *discordgo.VoiceConnection) error {

	// Don't allow two running PlayQueue methods on the same radio.
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
		disp, err := r.Play(ctx.Ses, vc)
		if err != nil {
			return err
		}

		r.Dispatcher = disp

		//----------------- Print information about the currently playing song ---------------- //
		song, err := r.Queue.Song()
		if err == nil && !r.Silent {
			ctx.ReplyEmbed(dream.NewEmbed().
				SetTitle("Now playing").
				SetDescription(fmt.Sprintf("[%d]: %s\nduration:\t %ds", r.Queue.Index, song.Markdown(), song.Duration)).
				SetFooter("added by " + song.AddedBy).
				SetColor(system.StatusNotify).
				MessageEmbed)
		}
		//-------------------------------------------------------------------------------------//

		done := make(chan bool)
		go func() {
			disp.Wait()
			done <- true
		}()

		select {
		case ctrl := <-r.control:
			switch ctrl {
			case AudioStop:
				disp.Stop()
				return nil
			case AudioContinue:
				continue
			}
			close(done)

		case <-done:
			// I only need to check for a closed voice connection after the done event
			// Is received because the dispatcher will end during a timeout error.
			if !vc.Ready {
				return errors.New("Voice connection closed")
			}
			// Load the next song if AutoPlay is enabled.
			if r.AutoPlay {
				err = r.Queue.Next()
				if err != nil {
					return err
				}
				// Otherwise, stop the queue.
			} else {
				return nil
			}
		}
	}
}

// Play plays a single song in the queue
func (r *Radio) Play(b *dream.Session, vc *discordgo.VoiceConnection) (*dream.AudioDispatcher, error) {
	song, err := r.Queue.Song()
	if err != nil {
		return nil, err
	}

	var stream io.Reader

	if r.UseYoutubeDL {
		yt := exec.Command("youtube-dl", "-f", "bestaudio", "--youtube-skip-dash-manifest", "-o", "-", song.URL)
		stream, err = yt.StdoutPipe()
		if err != nil {
			return nil, err
		}
		err = yt.Start()
		if err != nil {
			return nil, err
		}
	} else {
		stream, err = util.YoutubeDL(song.URL)
		if err != nil {
			return nil, err
		}
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

// Goto ...
func (r *Radio) Goto(index int) error {
	err := r.Queue.Goto(index)
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

	if r.IsRunning() {
		r.control <- AudioStop
		return nil
	}
	return errors.New("Audio player not running")
}

// Duration returns the duration the current song has been playing for in seconds
func (r *Radio) Duration() int {
	r.Lock()
	defer r.Unlock()
	if r.Dispatcher != nil {
		return int(r.Dispatcher.Duration.Seconds())
	}
	return 0
}

// SongInfoEmbed returns an embed with information about the currently playing song
func (r *Radio) SongInfoEmbed(index int) (*dream.Embed, error) {
	if index < 0 {
		index = r.Queue.Index
	}

	song, err := r.Queue.Get(index)
	if err != nil {
		return nil, err
	}

	embed := dream.NewEmbed().
		SetTitle(song.Title).
		SetURL(song.URL).
		SetImage(song.Thumbnail).
		SetDescription("Added by\t" + song.AddedBy + "\nindex\t" + fmt.Sprint(index)).
		SetColor(system.StatusNotify)

	if index == r.Queue.Index {
		embed.SetFooter(ProgressBar(r.Duration(), song.Duration, 100) + fmt.Sprintf("[%d/%d]", r.Duration(), song.Duration))
	} else {
		embed.SetFooter(fmt.Sprintf("Duration: %ds", song.Duration))
	}

	return embed, nil
}

///////////////////////////////////////////////////////////////////
//  Song loading
///////////////////////////////////////////////////

// QueueFromURL Queues a youtube video or playlist to the given song slice.
// Supports returning the currently queued song
func QueueFromURL(URL, addedBy string, queue *SongQueue, progress chan *Song) error {
	// Analyze the URL and use the proper method to obtain
	// Information about it.
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	// If the youtube link does not contain a link, use the golang youtube
	// Extractor to obtain the media information. This is significantly faster
	// Than using youtube-dl.
	if (u.Host == "youtube.com" || u.Host == "youtu.be") && u.Query().Get("list") == "" {
		song, err := SongFromYTDL(URL, addedBy)
		if err != nil {
			return err
		}
		if progress != nil {
			progress <- song
		}
		queue.Add(song)
		return nil
	}

	ytdl := exec.Command("youtube-dl", "-j", "-i", URL)
	ytdlout, err := ytdl.StdoutPipe()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(ytdlout)
	err = ytdl.Start()
	if err != nil {
		return err
	}
	var (
		line      []byte
		isPrefix  bool
		totalLine []byte
	)
	err = nil
	for err == nil {
		// Scan each line for song information
		line, isPrefix, err = reader.ReadLine()
		if err != nil && err != io.EOF {
			fmt.Println("SCANNER ERROR: ", err)
		}
		totalLine = append(totalLine, line...)
		if isPrefix {
			continue
		}
		song := &Song{}
		er := json.Unmarshal(totalLine, song)
		if er != nil {
			continue
		}

		// Add song to queue
		song.AddedBy = addedBy
		queue.Add(song)

		// Track progress
		if progress != nil {
			progress <- song
		}

		totalLine = []byte{}
	}
	if err != nil && err != io.EOF {
		return err
	}

	if progress != nil {
		close(progress)
	}

	return nil
}

// SongFromYTDL Uses ytdl to obtain video information and create a song object
func SongFromYTDL(URL, addedBy string) (*Song, error) {
	info, err := ytdl.GetVideoInfo(URL)
	if err != nil {
		return nil, err
	}

	song := &Song{
		Title:       info.Title,
		AddedBy:     addedBy,
		Description: info.Description,
		Duration:    int(info.Duration.Seconds()),
		ID:          info.ID,
		Thumbnail:   info.GetThumbnailURL(ytdl.ThumbnailQualityHigh).String(),
		Uploader:    info.Author,
		URL:         "https://www.youtube.com/watch?v=" + info.ID,
	}

	return song, nil
}

// QueueFromString queues a song from string
func QueueFromString(q *SongQueue, URL, addedBy, googleAPIKey string, UseYoutubeDL bool) error {
	if !strings.HasPrefix(URL, "https://") && !strings.HasPrefix(URL, "http://") {
		if googleAPIKey != "" {
			results, err := youtubeapi.New(googleAPIKey).Search(URL, 1)
			if err != nil {
				return err
			}
			if len(results.Items) != 0 {
				URL = results.Items[0].ID.VideoID
			}
		} else {
			results, err := youtubeapi.ScrapeSearch(URL, 1)
			if err == nil && len(results) > 0 {
				URL = results[0]
			}
		}
	}

	if UseYoutubeDL {
		return QueueFromURL(URL, addedBy, q, nil)
	}
	song, err := SongFromYTDL(URL, addedBy)
	if err != nil {
		return err
	}
	q.Add(song)
	return nil
}

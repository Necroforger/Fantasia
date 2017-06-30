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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/Fantasia/youtubeapi"
	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
)

//genmodules:config
////////////////////////////////
//         Config
////////////////////////////////

// Config module configuration
type Config struct {
	UseSubrouter bool

	// All radios will start with silent mode on.
	// Minimum notifications will be sent from commands.
	RadioSilent bool

	// RadioLoop sets the default loop value for radios.
	// If enabled, playlists will loop by default.
	RadioLoop bool

	// UseYoutubeDL determines if youtube-dl should be used to download videos.
	// If false, the golang library will be used.
	UseYoutubeDL bool

	// Start all radios with a test queue
	Debug bool
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		UseSubrouter: true,
		RadioSilent:  false,
		RadioLoop:    false,
		UseYoutubeDL: false,
		Debug:        false,
	}
}

// ControlCooldown is the cooldown for control event commands
const ControlCooldown = time.Millisecond * 1500

// Module ...
type Module struct {
	Config      *Config
	GuildRadios map[string]*Radio
}

// Build ...
func (m *Module) Build(s *system.System) {
	m.GuildRadios = map[string]*Radio{}

	var t *system.CommandRouter

	if m.Config.UseSubrouter {
		r, _ := system.NewSubCommandRouter(`^m(usicplayer)?(\s|$)`, "m | musicplayer")
		r.Router.Prefix = "^"
		r.Set("", "musicplayer subrouter, controls the various actions related to music playing")
		s.CommandRouter.AddSubrouter(r)
		t = r.Router
	} else {
		t = s.CommandRouter
	}

	// Queue management
	t.On("queue", m.CmdQueue).Set("", "queue")
	t.On("ytqueue", m.CmdYoutubeSearchQueue).Set("", "Searches youtube for the given query and queues the first video found\n`ytqueue [query]`")
	t.On("star", m.CmdStar).Set("", "Stars the song at the given index. Starring songs is akin to a favourites system and will allow you to sort songs based on their star ratings")
	t.On("loop", m.CmdLoop).Set("", "Controls whether the playlist should loop or not. Call with a boolean argument to change the loop mode.\n`loop [true | false]`")
	t.On("silent", m.CmdSilence).Set("", "Set the silence of the radio. If silent is true, the radio will no longer automatically give updates on the currently playing song\nUsage: `silent [true | false]`")
	t.On("remove|del(ete)?", m.CmdRemove).Set("remove", "Remove an index, or multiple indexes, from the queue.\nProvide multiple integer arguments to remove multiple indexes.")
	t.On("info", m.CmdInfo).Set("", "Gives information about the currently playing song")
	t.On("shuffle", m.CmdShuffle).Set("", "Shuffles the current queue, ignoring the current song index")
	t.On("swap", m.CmdSwap).Set("", "Swaps the song at index 'n' with index 't'\nusage: `swap [int: from] [int: to]`")
	t.On("move", m.CmdMove).Set("", "Moves the song at index 'n' to index 't'\nusage: `move [int: from] [int: to]`")
	t.On("clear", m.CmdClear).Set("", "Clears the current song queue")
	t.On("save", m.CmdSave).Set("", "Saves the current queue state to a json file and uploads it to discord")
	t.On("load", m.CmdLoad).Set("", "Loads a json playlist file. Present a URL, file attachment, or upload your file after calling this command")

	// Control commands
	t.On("go", m.CmdGoto).Set("", "Changes the queues current song index\nusage: `go [int: index]`")
	t.On("play", m.CmdPlay).Set("", "Plays the current queue")
	t.On("stop", m.CmdStop).Set("", "stops the currently playing queue")
	t.On("pause", m.CmdPause).Set("", "Pauses the currently playing song")
	t.On("resume", m.CmdResume).Set("", "Resumes the currently playing song")
	t.On("next", m.CmdNext).Set("", "Loads the next song in the queue")
	t.On("prev|previous", m.CmdPrevious).Set("prev | previous", "Loads the previous song in the queue")

	// Other
	t.On("tutorial", m.CmdTutorial).Set("", "A multipage tutorial for using the musicplayer module.\n Call this command in a DM to prevent other people from changing the pages on you")
}

// CmdSilence should toggle the radio from automatically sending messages when the song changes
func (m *Module) CmdSilence(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	radio := m.getRadio(guildID)

	if err != nil {
		ctx.ReplyError(err)
		return
	}
	if ctx.Args.After() == "true" {
		radio.Silent = true
		ctx.ReplySuccess("silent mode `enabled`")
	} else if ctx.Args.After() == "false" {
		radio.Silent = false
		ctx.ReplySuccess("silent mode `disabled`")
	} else {
		ctx.ReplySuccess(fmt.Sprintf("silent: `%t`", radio.Silent))
	}

}

// CmdClear clears the current queue
func (m *Module) CmdClear(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	radio := m.getRadio(guildID)
	radio.Queue.Clear()
	ctx.ReplyNotify("Cleared queue")
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

// CmdLoop toggles the loop state of the radio
func (m *Module) CmdLoop(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	radio := m.getRadio(guildID)

	if ctx.Args.After() == "false" {
		radio.Queue.Loop = false
	} else if ctx.Args.After() == "true" {
		radio.Queue.Loop = true
	}

	ctx.ReplyNotify(fmt.Sprintf("Loop playlists: `%t`", radio.Queue.Loop))

}

// CmdStar gives a rating to the given song index or the current song
func (m *Module) CmdStar(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	radio := m.getRadio(guildID)

	indexes := []int{radio.Queue.Index}
	if ctx.Args.After() != "" {
		indexes = getIndexes(strings.Split(ctx.Args.After(), " "), radio)
		if err != nil {
			ctx.ReplyError(err)
			return
		}
	}

	embedlog := dream.NewEmbed().SetColor(system.StatusNotify)

	for _, index := range indexes {
		song, err := radio.Queue.Get(index)
		if err != nil {
			ctx.ReplyError(err)
			return
		}

		if song.Rating == 1 {
			song.Rating = 0
			embedlog.Description += "Unstarred " + song.Markdown() + "\n"
		} else {
			song.Rating = 1
			embedlog.Description += "Starred " + song.Markdown() + "\n"
		}
	}

	embedlog.TruncateDescription()
	ctx.ReplyEmbed(embedlog.MessageEmbed)
}

// CmdQueue Queues a song or views the queue list
// Attempts to parse an integer value from the command argument.
// If an integer is parsed without error, the queue will present a
// List of songs around that position in the queue. If there is an
// Error parsing the integer, The queue will attempt to queue a video from the supplied URL.
func (m *Module) CmdQueue(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
	}
	radio := m.getRadio(guildID)

	index := 0
	if index, err = strconv.Atoi(ctx.Args.Get(0)); err != nil && ctx.Args.After() != "" {

		// If the URL is not an http link, search youtube for it.
		if !strings.HasPrefix(ctx.Args.After(), "http://") &&
			!strings.HasPrefix(ctx.Args.After(), "https://") {
			m.CmdYoutubeSearchQueue(ctx)
			return
		}

		// Youtube-dl
		if m.Config.UseYoutubeDL {
			progress := make(chan *Song)
			go func() {
				err := QueueFromURL(ctx.Args.After(), ctx.Msg.Author.Username, radio.Queue, progress)
				if err != nil {
					ctx.ReplyError(err)
				}
			}()

			msg, err := ctx.Ses.SendEmbed(ctx.Msg, dream.NewEmbed().
				SetColor(system.StatusNotify).
				SetDescription("Attempting to add to queue..."))
			if err != nil {
				ctx.ReplyError(err)
				return
			}

			radio.Lock()
			startIndex := len(radio.Queue.Playlist)
			radio.Unlock()

			count := 0
			var lastsong *Song
			for v := range progress {
				count++
				ctx.Ses.DG.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, dream.NewEmbed().
					SetColor(system.StatusSuccess).
					SetTitle(v.Title).
					SetDescription(fmt.Sprintf("Queued %d songs starting at index %d", count, startIndex)).
					SetFooter("index: "+fmt.Sprint(startIndex+count-1)).
					MessageEmbed)
				lastsong = v
			}

			finalEmbed := dream.NewEmbed().SetColor(system.StatusSuccess)
			if count == 1 {
				finalEmbed.SetDescription("queued " + lastsong.Markdown()).SetFooter("index: " + fmt.Sprint(radio.Queue.Index))
			} else {
				finalEmbed.SetDescription(fmt.Sprintf("queued %d songs starting at index %d", count, startIndex))
			}
			ctx.Ses.DG.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, finalEmbed.MessageEmbed)

			// YTDL
		} else {
			song, err := SongFromYTDL(ctx.Args.After(), ctx.Msg.Author.Username)
			if err != nil {
				ctx.ReplyError(err)
				return
			}
			radio.Queue.Add(song)
			ctx.ReplyEmbed(dream.NewEmbed().
				SetColor(system.StatusSuccess).
				SetDescription("queued " + song.Markdown()).
				SetFooter(fmt.Sprintf("index: %d", radio.Queue.Index)).
				MessageEmbed)
		}
		return
	}

	if err != nil {
		index = radio.Queue.Index
	}

	ctx.ReplyEmbed(EmbedQueue(radio.Queue, index, 10, 15).
		SetColor(system.StatusNotify).
		SetFooter("playlist length:\t" + fmt.Sprint(len(radio.Queue.Playlist))).
		MessageEmbed)
}

// CmdYoutubeSearchQueue searches youtube and queues the first video result found
func (m *Module) CmdYoutubeSearchQueue(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	radio := m.getRadio(guildID)

	var videoURL string

	if ctx.System.Config.GoogleAPIKey != "" {
		results, err := youtubeapi.New(ctx.System.Config.GoogleAPIKey).Search(ctx.Args.After(), 1)
		if err != nil {
			ctx.ReplyError(err)
			return
		}
		if len(results.Items) != 0 {
			videoURL = results.Items[0].ID.VideoID
		}
	} else {
		results, err := youtubeapi.ScrapeSearch(ctx.Args.After(), 1)
		if err == nil && len(results) > 0 {
			videoURL = results[0]
		}
	}

	if videoURL == "" {
		ctx.ReplyError("No results found")
		return
	}

	if song, err := SongFromYTDL(videoURL, ctx.Msg.Author.Username); err == nil {
		ctx.ReplySuccess(fmt.Sprintf("Queued [%d]: %s", radio.Queue.Add(song), song.Markdown()))
	} else {
		ctx.ReplyError(err)
	}

	return
}

// CmdRemove removes a song from the queue from its index id
func (m *Module) CmdRemove(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
	}
	radio := m.getRadio(guildID)

	if ctx.Args.After() == "" {
		ctx.ReplyError("Please provide the indexes you want to remove as a space separated list")
		return
	}

	ids := getIndexes(strings.Split(ctx.Args.After(), " "), radio)

	err = radio.Queue.Remove(ids...)
	if err != nil {
		if len(ids) == 1 {
			ctx.ReplyError("The index you provided was out of playlist bounds")
		} else {
			ctx.ReplyError("One of the indexes you provided was out of the playlist bounds")
		}
		return
	}
	ctx.ReplySuccess(fmt.Sprintf("Removed %d indexes", len(ids)))
}

// CmdSave saves the current playlist to a json encoded text file and uploads it to discord.
func (m *Module) CmdSave(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	radio := m.getRadio(guildID)

	playlistName := "playlist"
	if ctx.Args.After() != "" {
		playlistName = ctx.Args.After()
	}

	rd, wr := io.Pipe()
	go func() {
		json.NewEncoder(wr).Encode(radio.Queue.Playlist)
		wr.Close()
	}()
	ctx.Ses.SendFile(ctx.Msg.ChannelID, playlistName+".json", rd)
}

// CmdLoad loads a playlist from a previously saved file.
func (m *Module) CmdLoad(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	radio := m.getRadio(guildID)

	var fileURL string
	switch {
	case len(ctx.Msg.Attachments) > 0:
		fileURL = ctx.Msg.Attachments[0].URL
	case ctx.Args.After() != "":
		fileURL = ctx.Args.After()
	default:
		ctx.ReplyNotify("Upload a saved playlist or give a file url")
		var nxtmsg *discordgo.MessageCreate
		for nxtmsg = ctx.Ses.NextMessageCreate(); nxtmsg.Author.ID != ctx.Msg.Author.ID; nxtmsg = ctx.Ses.NextMessageCreate() {
		}
		if len(nxtmsg.Attachments) == 0 {
			fileURL = nxtmsg.Content
		} else {
			fileURL = nxtmsg.Attachments[0].URL
		}
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	defer resp.Body.Close()

	playlist := []*Song{}
	err = json.NewDecoder(resp.Body).Decode(&playlist)
	if err != nil {
		ctx.ReplyError(err)
	}

	radio.Queue.Add(playlist...)
	ctx.ReplySuccess("Loaded playlist into queue.")
}

// CmdInfo returns various info related to the currently playing song
func (m *Module) CmdInfo(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		return
	}
	radio := m.getRadio(guildID)

	var index = radio.Queue.Index
	if n, err := getIndex(ctx.Args.After(), radio); err == nil {
		index = n
	}

	song, err := radio.Queue.Get(index)
	if err != nil {
		ctx.ReplyError(err)
	}

	embed := dream.NewEmbed().
		SetTitle(song.Title).
		SetURL(song.URL).
		SetImage(song.Thumbnail).
		SetDescription("Added by\t" + song.AddedBy + "\nindex\t" + fmt.Sprint(radio.Queue.Index)).
		SetColor(system.StatusNotify)

	if index == radio.Queue.Index {
		embed.SetFooter(ProgressBar(radio.Duration(), song.Duration, 100) + fmt.Sprintf("[%d/%d]", radio.Duration(), song.Duration))
	} else {
		embed.SetFooter(fmt.Sprintf("Duration: %ds", song.Duration))
	}

	ctx.ReplyEmbed(embed.MessageEmbed)
}

// CmdShuffle Shuffles the current queue
func (m *Module) CmdShuffle(ctx *system.Context) {
	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	radio := m.getRadio(guildID)
	radio.Queue.Shuffle()
	ctx.ReplySuccess("Queue shuffled")
}

// CmdGoto changes the current song index to the specified index
func (m *Module) CmdGoto(ctx *system.Context) {

	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter the song index to go to")
		return
	}

	guildID, err := guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	radio := m.getRadio(guildID)

	index, err := getIndex(ctx.Args.After(), radio)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	err = radio.Goto(index)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	if song, err := radio.Queue.Song(); err == nil {
		ctx.ReplyEmbed(dream.NewEmbed().
			SetColor(system.StatusNotify).
			SetTitle("Selected song").
			SetDescription(fmt.Sprintf("[%d]: %s", radio.Queue.Index, song.Markdown())).
			SetColor(system.StatusSuccess).
			MessageEmbed)
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
		ctx.ReplyError("This command is on cooldown for `", (ControlCooldown - t).String(), "`. Please use `Go` or provide an integer argument to skip multiple songs quickly")
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
		ctx.ReplyError("This command is on cooldown for `", (ControlCooldown - t).String(), "`. Please use `Go` or provide an integer argument to skip multiple songs quickly")
		return
	}
	radio.ControlLastUsed = time.Now()

	err = radio.Previous()
	if err != nil {
		ctx.ReplyError(err)
	}
}

// CmdSwap swaps two queue indexes
func (m *Module) CmdSwap(ctx *system.Context) {
	var (
		guildID string
		from    int
		to      int
		err     error
	)

	guildID, err = guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	radio := m.getRadio(guildID)

	if from, err = getIndex(ctx.Args.Get(0), radio); err != nil {
		ctx.ReplyError(err)
		return
	}

	if to, err = getIndex(ctx.Args.Get(1), radio); err != nil {
		ctx.ReplyError(err)
		return
	}

	err = radio.Queue.Swap(from, to)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	ctx.ReplyNotify(fmt.Sprintf("Song index [%d] swapped with [%d]", from, to))
}

// CmdMove moves the specified index to the given position.
func (m *Module) CmdMove(ctx *system.Context) {
	var (
		guildID string
		from    int
		to      int
		err     error
	)

	guildID, err = guildIDFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	radio := m.getRadio(guildID)

	if from, err = getIndex(ctx.Args.Get(0), radio); err != nil {
		ctx.ReplyError(err)
		return
	}

	if to, err = getIndex(ctx.Args.Get(1), radio); err != nil {
		ctx.ReplyError(err)
		return
	}

	err = radio.Queue.Move(from, to)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	ctx.ReplyNotify(fmt.Sprintf("Song index [%d] moved to [%d]", from, to))
}

func (m *Module) getRadio(guildID string) *Radio {
	if v, ok := m.GuildRadios[guildID]; ok {
		return v
	}
	r := NewRadio(guildID)
	m.GuildRadios[guildID] = r

	if m.Config.RadioSilent {
		r.Silent = true
	}

	if m.Config.RadioLoop {
		r.Queue.Loop = true
	}

	r.UseYoutubeDL = m.Config.UseYoutubeDL

	if m.Config.Debug {
		r.Queue.Playlist = []*Song{
			&Song{
				Title: "Touhou Erhu［東方名曲］土著神醮 ／ 平行世界",
				URL:   "https://www.youtube.com/watch?v=-yEDM-ogtBY",
			},
			&Song{
				Title: "東方 [Piano] Necrofantasia 『5』",
				URL:   "https://youtu.be/-yNATRuEvsE",
			},
			&Song{
				Title: "東方 [Piano] Love-coloured Master Spark 『4』",
				URL:   "https://www.youtube.com/watch?v=b8AIGUYscYo",
			},
			&Song{
				Title: "Touhou Remix E.149 (Horror) Satori Maiden ~ 3rd Eye",
				URL:   "https://www.youtube.com/watch?v=m2eeyey-87U",
			},
			&Song{
				Title: "Elder - Reflections of a Floating World (2017) (New Full Album)",
				URL:   "https://youtu.be/tqMfWwnLtKg",
			},
			&Song{
				Title: "東方 [Piano] Reach for the Moon, Immortal Smoke 『3』",
				URL:   "https://www.youtube.com/watch?v=WUJdZDM8pKk",
			},
			&Song{
				Title: "東方 [Piano] Septette for the Dead Princess 『8』",
				URL:   "https://www.youtube.com/watch?v=c55RF62YZgI",
			},
			&Song{
				Title: "東方 [Piano] Septette for the Dead Princess 『6』",
				URL:   "https://youtu.be/RYTEJ7fdDrg",
			},
		}
	}

	return r
}

package system

import (
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/ytdl"
)

var customHTTPClient = http.Client{
	Timeout: time.Second * 10,
}

////////////////////////////////////////////////
// YOUTUBE DOWNLOADING
////////////////////////////////////////////////

// YoutubeDLFromInfo ...
func YoutubeDLFromInfo(info *ytdl.VideoInfo) (io.Reader, error) {
	if len(info.Formats.Best(ytdl.FormatAudioEncodingKey)) == 0 {
		return nil, errors.New("Error processing video")
	}
	u, err := info.GetDownloadURL(info.Formats.Best(ytdl.FormatAudioEncodingKey)[0])
	if err != nil {
		return nil, err
	}

	resp, err := customHTTPClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New("invalid status code")
	}

	return resp.Body, nil
}

// YoutubeDL ...
func YoutubeDL(URL string) (io.Reader, error) {
	info, err := ytdl.GetVideoInfo(URL)
	if err != nil {
		return nil, err
	}

	return YoutubeDLFromInfo(info)
}

///////////////////////////////////////
//   FILES
//////////////////////////////////////

// RandomFileInDir retrieves a random file from a folder
func RandomFileInDir(path string) (*os.File, error) {
	info, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := []string{}
	for _, v := range info {
		if !v.IsDir() {
			files = append(files, v.Name())
		}
	}

	if len(files) == 0 {
		return nil, errors.New("No files in directory")
	}
	if len(files) == 1 {
		return os.Open(filepath.Join(path, files[0]))
	}

	return os.Open(filepath.Join(path, files[int(rand.Float64()*float64(len(files)))]))
}

/////////////////////////////////////////////////////
// Audio
/////////////////////////////////////////////////////

// ConnectToVoiceChannel finds and connects to a user's voice channel
func ConnectToVoiceChannel(ctx *Context) (*discordgo.VoiceConnection, error) {
	b := ctx.Ses
	msg := ctx.Msg

	vc, err := b.GuildVoiceConnection(ctx.Msg)
	if err != nil {

		// If not currently in a channel, attempt to join the voice channel of the calling user.
		vs, err := b.UserVoiceState(msg.Author.ID)
		if err != nil {
			return nil, err
		}

		vc, err = b.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			return nil, err
		}

		return vc, nil
	}

	// Confirm that the user is in the correct voice channel.
	// If the user is in another voice channel, join them.
	vs, err := b.UserVoiceState(msg.Author.ID)
	if err == nil && vc != nil && vs.ChannelID != vc.ChannelID || vs.GuildID != vc.GuildID {
		vc, err = b.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			return nil, err
		}
	}

	if !vc.Ready {
		vc, err = b.ChannelVoiceJoin(vc.GuildID, vc.ChannelID, false, true)
		if err != nil {
			return nil, err
		}
	}

	return vc, nil
}

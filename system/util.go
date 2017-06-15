package system

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Necroforger/discordgo"
)

var customHTTPClient = http.Client{
	Timeout: time.Second * 10,
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
	if err == nil && (vc != nil && vs != nil && vs.ChannelID != vc.ChannelID || vs.GuildID != vc.GuildID) {
		vc, err = b.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			return nil, err
		}
		return vc, nil
	}

	if !vc.Ready {
		vc, err = b.ChannelVoiceJoin(vc.GuildID, vc.ChannelID, false, true)
		if err != nil {
			return nil, err
		}
	}

	return vc, nil
}

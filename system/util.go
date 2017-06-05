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

	"github.com/Necroforger/ytdl"
)

var customhHttpClient = http.Client{
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

	resp, err := customhHttpClient.Get(u.String())
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

// RandomFileInFolder retrieves a random file from a folder
func RandomFileInFolder(path string) (*os.File, error) {
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

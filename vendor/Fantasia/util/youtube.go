package util

import (
	"errors"
	"io"

	"github.com/Necroforger/ytdl"
)

// YoutubeDLFromInfo ...
func YoutubeDLFromInfo(info *ytdl.VideoInfo) (io.ReadCloser, error) {
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
func YoutubeDL(URL string) (io.ReadCloser, error) {
	info, err := ytdl.GetVideoInfo(URL)
	if err != nil {
		return nil, err
	}

	return YoutubeDLFromInfo(info)
}

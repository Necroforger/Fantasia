package util

import (
	"errors"
	"image"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/Necroforger/dream"
)

var urlRegex = regexp.MustCompile(`(http|ftp|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)

// ReadCloserList is a slice of ReadClosers
type ReadCloserList []io.ReadCloser

// CloseAll closes all the ReadCloser connections
func (r ReadCloserList) CloseAll() {
	for _, v := range r {
		v.Close()
	}
}

// RequestMessage waits for a message from the given user
//    s       : Dream session
//    userID  : userID of the user you wish to request a message from
//    timeout : Duration to wait before returning an error. -1 to never timeout.
func RequestMessage(s *dream.Session, userID string, timeout time.Duration) (msg *discordgo.MessageCreate, err error) {
	for {
		if timeout > 0 {
			select {
			case <-time.After(timeout):
				err = errors.New("err: timed out waiting for response")
				return
			case msg = <-s.NextMessageCreateC():
				if msg.Author.ID != userID {
					continue
				}
				return
			}
		} else {
			msg = s.NextMessageCreate()
			if msg.Author.ID != userID {
				continue
			}
			return
		}
	}
}

// RequestFiles waits for the user to upload multiple files via URLs or attachments
//     s       : Dream session
//     userid  : userID of the user you wish to request the files from
//     timeout : How long the bot should wait before returning an error. -1 to never time out.
func RequestFiles(s *dream.Session, userID string, timeout time.Duration) (readers ReadCloserList, err error) {
	msg, err := RequestMessage(s, userID, timeout)
	if err != nil {
		return nil, err
	}

	readers = FilesFromMessage(msg.Message)
	return
}

// FilesFromMessage obtains a slice of io.ReadClosers from a message's content
// Remember to close ALL the readers when you are done.
//    msg : message to obtain the files from
func FilesFromMessage(msg *discordgo.Message) (readers ReadCloserList) {
	readers = []io.ReadCloser{}
	URLs := []string{}

	for _, a := range msg.Attachments {
		URLs = append(URLs, a.URL)
	}

	for _, u := range urlRegex.FindAllString(msg.Content, -1) {
		URLs = append(URLs, u)
	}

	for _, URL := range URLs {
		resp, err := http.Get(URL)
		if err == nil {
			readers = append(readers, resp.Body)
		}
	}

	return
}

// RequestImages requests images from a user
//     s       : dream session
//     userID  : userID of the user to fetch the images from
//     timeout : How long the bot should wait before returning an error. -1 to never time out.
func RequestImages(s *dream.Session, userID string, timeout time.Duration) (images []image.Image, err error) {
	msg, err := RequestMessage(s, userID, timeout)
	if err != nil {
		return nil, err
	}
	images = ImagesFromMessage(msg.Message)
	return
}

// ImagesFromMessage returns a slice of images from a message
//     msg : discordgo message to obtain the images from.
func ImagesFromMessage(msg *discordgo.Message) (images []image.Image) {
	images = []image.Image{}
	files := FilesFromMessage(msg)

	for _, v := range files {
		img, _, err := image.Decode(v)
		if err == nil {
			images = append(images, img)
		}
		v.Close()
	}

	return
}

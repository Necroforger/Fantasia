package images

import (
	"Fantasia/system"
	"image"
	"image/png"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// util.go : provides various utility functions
var urlRegex = regexp.MustCompile(`(http|ftp|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)

// ImageSuffixes is a list of image suffixes
var ImageSuffixes = []string{
	"png",
	"jpeg",
	"jpg",
	"gif",
}

// ReplyImage replies to the sender with the given image
func ReplyImage(ctx *system.Context, img image.Image) {
	rd, wr := io.Pipe()
	go func() {
		png.Encode(wr, img)
		wr.Close()
	}()
	ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, "image.png", rd)
}

// PullImages requests images from the user, or retrieves them from the cache.
//    amount    : amount of messages to retrieve
//    channelID : channelID to pull images from
//    message   : message user supplied. If images are present, pull them from here first.
//                otherwise fall back to the image cache.
func (m *Module) PullImages(limit int, channelID string, message *discordgo.Message) ([]image.Image, error) {
	images := make([]image.Image, 0, limit)

	// Pull images from message issuing command.
	tmp, err := PullImagesFromMessage(limit, message)
	if err != nil {
		return images, err
	}
	images = append(images, tmp...)

	// If we have satisfied the message limit, return
	if len(images) >= limit {
		return images, nil
	}

	// Else continue searching the cache for images
	tmp, err = PullImagesFromCache(m.ImgCache, limit-len(images), channelID)
	if err != nil {
		return images, err
	}
	images = append(images, tmp...)

	return images, nil
}

// PullImagesFromMessage retrieves images from a message
func PullImagesFromMessage(limit int, msg *discordgo.Message) ([]image.Image, error) {
	URLs := ImageURLsInMessage(msg)
	images := []image.Image{}

	for _, v := range URLs {
		img, err := ImageFromURL(v)
		if err != nil {
			return images, err
		}
		images = append(images, img)

		limit--
		if limit == 0 {
			break
		}
	}

	return images, nil
}

// ImageFromURL gets an image from a URL
//    URL : URL to perform an HTTP get reqeust to
func ImageFromURL(URL string) (image.Image, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	return img, err
}

// ImageURLsInMessage returns a list of image URLs in a message
// Images are obtained from embed images, and thumbnails, and attachments.
//    msg : message to obtain the URLs from
func ImageURLsInMessage(msg *discordgo.Message) []string {
	URLs := []string{}

	for _, a := range msg.Attachments {
		URLs = append(URLs, a.URL)
	}

	for _, embed := range msg.Embeds {
		if embed.Image != nil && embed.Image.ProxyURL != "" {
			URLs = append(URLs, embed.Image.ProxyURL)
		}
		if embed.Thumbnail != nil && embed.Thumbnail.ProxyURL != "" {
			URLs = append(URLs, embed.Thumbnail.ProxyURL)
		}
	}

	return URLs
}

// PullImagesFromCache ...
func PullImagesFromCache(cache *MessageCache, limit int, channelID string) ([]image.Image, error) {
	images := make([]image.Image, 0, limit)

	messages, err := cache.Messages(channelID)
	if err != nil {
		return nil, err
	}

	for i := len(messages) - 1; i >= 0; i-- {
		v := messages[i]
		tmp, err := PullImagesFromMessage(limit, v)
		if err != nil {
			return nil, err
		}
		images = append(images, tmp...)

		// Decrement the image limit
		limit -= len(images)
		if limit <= 0 {
			break
		}
	}

	return images, nil
}

// HasImage returns if a message has an image
func HasImage(msg *discordgo.Message) bool {

	for _, embed := range msg.Embeds {
		if embed.Image != nil && embed.Image.URL != "" {
			return true
		}
		if embed.Thumbnail != nil && embed.Thumbnail.URL != "" {
			return true
		}
	}

	for _, v := range msg.Attachments {
		if HasImageSuffix(v.Filename) {
			return true
		}
	}

	return false
}

// HasImageSuffix returns true if an image suffix is present
func HasImageSuffix(filename string) bool {
	for _, v := range ImageSuffixes {
		if strings.HasSuffix(filename, v) {
			return true
		}
	}
	return false
}

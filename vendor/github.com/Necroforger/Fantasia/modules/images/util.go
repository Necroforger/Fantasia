package images

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/gif"
	_ "image/jpeg" // Needed to decode jpegs
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Necroforger/Fantasia/system"

	"github.com/bwmarrin/discordgo"
)

// ImageMaxDimensions is the maximum size the image decoder is allowed to decode
const ImageMaxDimensions = 3840 * 2160

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
		err := png.Encode(wr, img)
		if err != nil {
			ctx.ReplyError("Error encoding png: ", err)
		}
		wr.Close()
	}()
	_, err := ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, "image.png", rd)
	if err != nil {
		ctx.ReplyError("Error sending file to channel: ", err)
	}
}

// ReplyGif replies to the sender with the given gif
func ReplyGif(ctx *system.Context, g *gif.GIF) {
	rd, wr := io.Pipe()
	go func() {
		err := gif.EncodeAll(wr, g)
		if err != nil {
			ctx.ReplyError("error encoding gif: ", err)
		}
		wr.Close()
	}()
	_, err := ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, "image.gif", rd)
	if err != nil {
		ctx.ReplyError("Error uploading gif: ", err)
	}
}

// CompressGif attempts to compress a Gif's images
//    src      : source gif to compress
//    distance : allowed color distance for colors to be considered equal
func CompressGif(src *gif.GIF, distance uint32) *gif.GIF {
	dst := &gif.GIF{
		Image:    make([]*image.Paletted, len(src.Image)),
		Disposal: make([]byte, len(src.Disposal)),
		Delay:    src.Delay,
	}
	// Set disposal to none.
	for i := 0; i < len(dst.Disposal); i++ {
		dst.Disposal[i] = 1
	}

	diff := func(a, b uint32) bool {
		if a >= b-distance && a <= b+distance {
			return false
		}
		return true
	}

	var frame = src.Image[0]
	dst.Image[0] = frame
	for i := 1; i < len(src.Image); i++ {
		cp := src.Image[i]
		np := image.NewPaletted(src.Image[i].Bounds(), cp.Palette)
		np.Palette[0] = color.RGBA{0, 0, 0, 0} // Insert an alpha pixel into the palette

		w, h := cp.Bounds().Dx(), cp.Bounds().Dy()
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				pidx := y*cp.Stride + x*1
				// Compare the colours of frame(last image) and cp(current paletted),
				// if they are different, add the pixel to new paletted (np)
				r, g, b, a := cp.Palette[cp.Pix[pidx]].RGBA()
				r1, g1, b1, a1 := frame.Palette[frame.Pix[pidx]].RGBA()

				if diff(r, r1) || diff(g, g1) || diff(b, b1) || diff(a, a1) {
					np.Pix[pidx] = cp.Pix[pidx]
				} else {
					np.Pix[pidx] = uint8(np.Palette.Index(color.RGBA{0, 0, 0, 0}))
				}
			}
		}

		dst.Image[i] = np
		frame = cp // Update last frame
	}

	return dst
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

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 100000000)) // read a maximum of 100mb
	if err != nil {
		return nil, err
	}
	bodyBuf := bytes.NewReader(body)

	// Read from first reader
	config, _, err := image.DecodeConfig(bodyBuf)
	if err != nil {
		return nil, err
	}

	if config.Width*config.Height >= ImageMaxDimensions {
		return nil, errors.New("Image is too big: maximum dimensions: " + strconv.Itoa(ImageMaxDimensions))
	}

	// restore to the original dimensions
	bodyBuf.Seek(0, io.SeekStart)

	// Read from second reader
	img, _, err := image.Decode(bodyBuf)
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

	if len(msg.Embeds) == 0 {
		for _, u := range urlRegex.FindAllString(msg.Content, -1) {
			if IsImageURL(u) {
				URLs = append(URLs, u)
			}
		}
	}

	return URLs
}

// IsImageURL guesses if a URL is an image
func IsImageURL(path string) bool {
	t, err := url.Parse(path)
	if err != nil {
		log.Println("error parsing URL")
	} else {
		if HasImageSuffix(strings.ToLower(t.Path)) {
			return true
		}
	}
	return false
}

// PullImagesFromCache pulls images from the cache
//    cache     : message cache to pull from
//    limit     : maximum number of images to retrieve
//    channelID : channelID of messages to retrieve
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

	for _, u := range urlRegex.FindAllString(msg.Content, -1) {
		if IsImageURL(u) {
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

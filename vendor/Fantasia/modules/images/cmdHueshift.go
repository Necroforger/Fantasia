package images

import (
	"Fantasia/modules/images/animate"
	"Fantasia/system"
	"image"

	"github.com/nfnt/resize"

	"github.com/anthonynsimon/bild/adjust"
)

// CmdHueshift animates an image's hue
func (m *Module) CmdHueshift(ctx *system.Context) {
	images, err := m.PullImages(1, ctx.Msg.ChannelID, ctx.Msg)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	if len(images) == 0 {
		ctx.ReplyError("No images found")
		return
	}

	images[0] = resize.Thumbnail(300, 300, images[0], resize.NearestNeighbor)
	ReplyGif(ctx, animate.Animate(images[0], effecthue, 0, 360, 10, 10))
}

func effecthue(img image.Image, amount float64) *image.RGBA {
	return adjust.Hue(img, int(amount))
}

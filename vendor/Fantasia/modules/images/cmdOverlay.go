package images

import (
	"Fantasia/system"
	"image"
	"image/draw"

	"github.com/anthonynsimon/bild/clone"
	"github.com/nfnt/resize"
)

// CmdOverlay ...
func (m *Module) CmdOverlay(ctx *system.Context) {
	images, err := m.PullImages(2, ctx.Msg.ChannelID, ctx.Msg)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	if len(images) < 2 {
		ctx.ReplyError("You need atleast two images in the cache to use this command")
		return
	}

	dst := clone.AsRGBA(images[1])
	images[0] = resize.Resize(0, uint(images[1].Bounds().Dy()), images[0], resize.NearestNeighbor)
	draw.Draw(dst, dst.Bounds(), images[0], image.ZP, draw.Over)

	ReplyImage(ctx, dst)
}

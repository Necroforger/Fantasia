package images

import (
	"Fantasia/system"
	"strconv"

	"github.com/anthonynsimon/bild/adjust"
)

// CmdSaturation adjusts the saturation of an image
func (m *Module) CmdSaturation(ctx *system.Context) {
	if ctx.Args.Get(0) == "" {
		ctx.ReplyError("Please supply a float indicating the change in saturation")
		return
	}
	images, err := m.PullImages(1, ctx.Msg.ChannelID, ctx.Msg)
	if err != nil {
		ctx.ReplyError("Error fetching images: ", err)
		return
	}
	if len(images) == 0 {
		ctx.ReplyError(ErrNoImagesFound)
		return
	}
	amount, err := strconv.ParseFloat(ctx.Args.Get(0), 64)
	if err != nil {
		ctx.ReplyError("Error parsing float")
		return
	}
	ReplyImage(ctx, adjust.Saturation(images[0], amount))
}

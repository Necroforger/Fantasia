package images

import (
	"Fantasia/system"
	"strconv"

	"github.com/anthonynsimon/bild/adjust"
)

// CmdHue adjusts the hue of an image
func (m *Module) CmdHue(ctx *system.Context) {
	images, err := m.PullImages(1, ctx.Msg.ChannelID, ctx.Msg)
	if err != nil {
		ctx.ReplyError("Error fetching images: ", err)
		return
	}
	if len(images) == 0 {
		ctx.ReplyError(ErrNoImagesFound)
		return
	}

	amount, err := strconv.Atoi(ctx.Args.Get(0))
	if err != nil {
		amount = 90
	}

	ReplyImage(ctx, adjust.Hue(images[0], amount))
}

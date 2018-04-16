package images

import (
	"Fantasia/system"

	"github.com/anthonynsimon/bild/effect"
)

// CmdInvert inverts an image
func (m *Module) CmdInvert(ctx *system.Context) {
	images, err := m.PullImages(1, ctx.Msg.ChannelID, ctx.Msg)
	if err != nil {
		ctx.ReplyError("Error fetching images: ", err)
		return
	}
	if len(images) == 0 {
		ctx.ReplyError(ErrNoImagesFound)
		return
	}

	ReplyImage(ctx, effect.Invert(images[0]))
}

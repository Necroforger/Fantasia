package images

import (
	"Fantasia/system"

	"github.com/Necroforger/textify"
)

// CmdTextify converts an image to text
func (m *Module) CmdTextify(ctx *system.Context) {
	img, err := imageFromContext(ctx)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	opts := textify.NewOptions()
	opts.Thumbnail = true
	opts.Resize = true
	opts.Palette = textify.PaletteReverse
	opts.Width = 50
	opts.Height = 50
	text, err := textify.Encode(img, opts)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	ctx.Reply("```" + text + "```")
}

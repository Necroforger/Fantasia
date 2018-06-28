package images

import (
	"github.com/Necroforger/Fantasia/system"

	"github.com/Necroforger/textify"
)

// CmdTextify converts an image to text
func (m *Module) CmdTextify(ctx *system.Context) {
	images, err := m.PullImages(1, ctx.Msg.ChannelID, ctx.Msg)
	if err != nil {
		ctx.ReplyError(err)
		return
	}
	if len(images) == 0 {
		ctx.ReplyError("No images found")
		return
	}

	// opts := textify.NewOptions()
	// opts.Palette = []string{".", ".", " "}
	// opts.StrideW = 0.5
	// opts.StrideH = 2.3
	// opts.Thumbnail = true
	// opts.Resize = true
	// opts.Height = 100
	// opts.Width = 100
	// txt, err := textify.Encode(images[0], opts)

	opts := textify.NewOptions()
	opts.Palette = textify.PaletteReverse
	opts.Resize = true
	opts.Thumbnail = true
	opts.Height = 43
	opts.Width = 43
	txt, err := textify.Encode(images[0], opts)
	if err != nil {
		ctx.ReplyError(err)
	}

	ctx.Reply("```\n" + txt + "\n```")

	// buffer := ""
	// for i, v := range strings.Split(txt, "\n") {
	// 	buffer += v + "\n"
	// 	if i%10 == 0 {
	// 		ctx.Reply(buffer)
	// 		buffer = ""
	// 	}
	// }
	// if buffer != "" {
	// 	ctx.Reply(buffer)
	// }
}

package themeify

import (
	"image/png"
	"io"

	"github.com/Necroforger/Fantasia/system"
	"github.com/fogleman/gg"
)

// CmdDuoText ...
func CmdDuoText(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter your text enclosed with quotes")
		return
	}

	fnt, _ := gg.LoadFontFace("fonts/swanse.ttf", 50)

	img := mergeImages(
		createTextImage(ctx.Args.Get(0), LightThemeClr, fnt, 300),
		createTextImage(ctx.Args.Get(1), DarkThemeClr, fnt, 300),
	)

	rd, wr := io.Pipe()
	go func() {
		png.Encode(wr, img)
		wr.Close()
	}()
	ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, "duotext.png", rd)
}

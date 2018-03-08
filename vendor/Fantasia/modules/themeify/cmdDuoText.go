package themeify

import (
	"image/png"
	"io"

	"Fantasia/fonts"
	"Fantasia/system"
	"github.com/golang/freetype/truetype"
)

// CmdDuoText ...
func CmdDuoText(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter your text enclosed with quotes")
		return
	}

	fnt := truetype.NewFace(fonts.Swanse, &truetype.Options{
		Size: 50,
	})
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

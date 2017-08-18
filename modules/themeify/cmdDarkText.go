package themeify

import (
	"image/png"
	"io"

	"github.com/golang/freetype/truetype"

	"github.com/Necroforger/Fantasia/fonts"
	"github.com/Necroforger/Fantasia/system"
)

// CmdDarkText creates text that only a dark theme user can see
func CmdDarkText(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter some text")
		return
	}

	fnt := truetype.NewFace(fonts.Swanse, &truetype.Options{
		Size: 25,
	})
	img := createTextImage(ctx.Args.After(), DarkThemeClr, fnt, 300)

	rd, wr := io.Pipe()
	go func() {
		png.Encode(wr, img)
		wr.Close()
	}()
	ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, "text.png", rd)
}

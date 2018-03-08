package general

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"image/png"
	"io"

	"Fantasia/fonts"
	"Fantasia/system"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

// CmdHexDisplay ...
func CmdHexDisplay(ctx *system.Context) {
	hexb, err := hex.DecodeString(ctx.Args.After())
	if err != nil {
		return
	}

	hexclr := color.RGBA{0, 0, 0, 255}
	switch l := len(hexb); {
	case l > 2:
		hexclr.B = hexb[2]
		fallthrough
	case l > 1:
		hexclr.G = hexb[1]
		fallthrough
	case l > 0:
		hexclr.R = hexb[0]
	}

	gc := gg.NewContext(300, 200)

	gc.SetColor(hexclr)
	gc.DrawRectangle(0, 0, float64(gc.Width()), float64(gc.Height()))
	gc.Fill()

	gc.SetFontFace(truetype.NewFace(fonts.MonospaceTypewriter, &truetype.Options{
		Size: 50,
	}))

	gc.SetColor(color.RGBA{hexclr.R ^ 0xff, hexclr.G ^ 0xff, hexclr.B ^ 0xff, 255})
	gc.DrawString(ctx.Args.After(), 10, 50)

	gc.SetFontFace(truetype.NewFace(fonts.MonospaceTypewriter, &truetype.Options{
		Size: 15,
	}))

	gc.DrawString("RGB"+fmt.Sprint(hexb), 10, float64(gc.Height())-10)

	rd, wr := io.Pipe()
	go func() {
		png.Encode(wr, gc.Image())
		wr.Close()
	}()

	ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, "hex.png", rd)
}

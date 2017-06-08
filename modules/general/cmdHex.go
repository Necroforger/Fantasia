package general

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"image/png"
	"io"

	"github.com/Necroforger/Fantasia/system"
	"github.com/fogleman/gg"
)

// HexDisplay ...
func HexDisplay(ctx *system.Context) {
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

	// gc.LoadFontFace(monospaceFont, 50)

	gc.SetColor(color.RGBA{255 - hexclr.R, 255 - hexclr.G, 255 - hexclr.B, 255})
	gc.DrawString(ctx.Args.After(), 10, 50)
	// gc.LoadFontFace(monospaceFont, 15)
	gc.DrawString("RGB"+fmt.Sprint(hexb), 10, float64(gc.Height())-10)

	rd, wr := io.Pipe()
	go func() {
		png.Encode(wr, gc.Image())
		wr.Close()
	}()

	ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, "hex.png", rd)
}

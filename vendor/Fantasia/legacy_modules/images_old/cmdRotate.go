package images_old

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"strconv"
	"time"

	"github.com/disintegration/imaging"

	"Fantasia/system"
	"Fantasia/util"
)

// CmdRotate rotates an image 90 degrees
func CmdRotate(ctx *system.Context) {
	var img image.Image
	if imgs := util.ImagesFromMessage(ctx.Msg); len(imgs) > 0 {
		img = imgs[0]
	} else {
		ctx.ReplyNotify("Please upload an image or enter an image URL")
		imgs, err := util.RequestImages(ctx.Ses, ctx.Msg.Author.ID, time.Minute*5)
		if err != nil {
			ctx.ReplyError(ctx.Msg.Author.Mention() + ": timed out while waiting for image")
			return
		}
		if len(imgs) == 0 {
			return
		}
		img = imgs[0]
	}

	var angle float64
	if ctx.Args.After() != "" {
		n, err := strconv.ParseFloat(ctx.Args.After(), 64)
		if err != nil {
			ctx.ReplyError("invalid angle")
			return
		}
		angle = n
	} else {
		angle = 90
	}

	img = imaging.Rotate(img, angle, color.Black)

	rd, wr := io.Pipe()
	go func() {
		png.Encode(wr, img)
		wr.Close()
	}()
	ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, "image.png", rd)
}

package themeify

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"time"

	"github.com/nfnt/resize"

	"Fantasia/system"
	"Fantasia/util"
)

// CmdDuoImage ...
func CmdDuoImage(ctx *system.Context) {
	const (
		width  = uint(300)
		height = uint(300)
	)

	var (
		img1 image.Image
		img2 image.Image
	)

	images := util.ImagesFromMessage(ctx.Msg)
	if len(images) >= 2 {
		img1 = images[0]
		img2 = images[1]
	} else {
		ctx.ReplyNotify("Upload an image as an attachment or enter a list of image urls: ")
		imgs, err := util.RequestImages(ctx.Ses, ctx.Msg.Author.ID, time.Minute*5)
		if err != nil {
			ctx.ReplyError("Timed out waiting for images...")
			return
		}
		images = append(images, imgs...)
		if len(images) < 2 {
			ctx.ReplyNotify("Upload another image or enter an image url: ")
			imgs, err := util.RequestImages(ctx.Ses, ctx.Msg.Author.ID, time.Minute*5)
			if err != nil {
				ctx.ReplyError("Timed out waiting for images...")
				return
			}
			images = append(images, imgs...)
		}
		if len(images) >= 2 {
			img1 = images[0]
			img2 = images[1]
		} else {
			ctx.ReplyError("Two images are required to merge")
			return
		}
	}

	img1 = resize.Thumbnail(width, height, img1, resize.Lanczos3)
	img2 = resize.Thumbnail(width, height, img2, resize.Lanczos3)

	gray1 := convertToGrayscale(img1, []color.Color{color.Transparent, color.RGBA{0xff, 0xff, 0xff, 0xff}})
	gray2 := convertToGrayscale(img2, []color.Color{color.RGBA{0x36, 0x39, 0x3E, 0xff}, color.Transparent})

	rd, wr := io.Pipe()
	go func() {
		png.Encode(wr, mergeImages(gray1, gray2))
		wr.Close()
	}()
	ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, "mono.png", rd)
}

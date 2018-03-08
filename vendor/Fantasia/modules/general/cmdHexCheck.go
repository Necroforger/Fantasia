package general

import (
	"fmt"
	"time"

	"Fantasia/system"
	"Fantasia/util"
)

// CmdHexCheck returns the hex value of the center pixel of the supplied image
func CmdHexCheck(ctx *system.Context) {
	images := util.ImagesFromMessage(ctx.Msg)
	if len(images) == 0 {
		ctx.ReplyNotify("Upload an images or enter the urls ")
		imgs, err := util.RequestImages(ctx.Ses, ctx.Msg.Author.ID, time.Minute*5)
		if err != nil {
			ctx.ReplyError(ctx.Msg.Author.Mention() + "Timed out while waiting for images ")
		}
		images = append(images, imgs...)
	}

	if len(images) == 0 {
		ctx.ReplyError("No images found")
	}

	var results string
	for _, v := range images {
		if dx, dy := v.Bounds().Dx(), v.Bounds().Dy(); dx > 0 && dy > 0 {
			r, g, b, _ := v.At(dx/2, dy/2).RGBA()
			results += fmt.Sprintf("#%02X%02X%02X\n", r>>8, g>>8, b>>8)
		}
	}

	ctx.ReplyNotify(results)
}

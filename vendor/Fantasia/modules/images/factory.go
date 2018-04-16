package images

import (
	"Fantasia/system"
	"image"
	"strconv"
)

// NewEffectCmdSingle ...
func (m *Module) NewEffectCmdSingle(fn func(image.Image) *image.RGBA) func(ctx *system.Context) {
	return func(ctx *system.Context) {
		images, err := m.PullImages(1, ctx.Msg.ChannelID, ctx.Msg)
		if err != nil {
			ctx.ReplyError("Error fetching images: ", err)
			return
		}
		if len(images) == 0 {
			ctx.ReplyError(ErrNoImagesFound)
			return
		}

		ReplyImage(ctx, fn(images[0]))
	}
}

// NewEffectCommandFloat produces an effect command that accepts an image and a float
func (m *Module) NewEffectCommandFloat(fn func(img image.Image, amount float64) *image.RGBA) func(ctx *system.Context) {
	return func(ctx *system.Context) {
		if ctx.Args.Get(0) == "" {
			ctx.ReplyError("Please supply a float")
			return
		}
		images, err := m.PullImages(1, ctx.Msg.ChannelID, ctx.Msg)
		if err != nil {
			ctx.ReplyError("Error fetching images: ", err)
			return
		}
		if len(images) == 0 {
			ctx.ReplyError(ErrNoImagesFound)
			return
		}
		amount, err := strconv.ParseFloat(ctx.Args.Get(0), 64)
		if err != nil {
			ctx.ReplyError("Error parsing float")
			return
		}
		ctx.Reply(amount)
		ReplyImage(ctx, fn(images[0], amount))
	}
}

// NewEffectCommandInt ...
func (m *Module) NewEffectCommandInt(fn func(img image.Image, amount int) *image.RGBA) func(ctx *system.Context) {
	return func(ctx *system.Context) {
		if ctx.Args.Get(0) == "" {
			ctx.ReplyError("Please supply an integer")
			return
		}

		images, err := m.PullImages(1, ctx.Msg.ChannelID, ctx.Msg)
		if err != nil {
			ctx.ReplyError("Error fetching images: ", err)
			return
		}
		if len(images) == 0 {
			ctx.ReplyError(ErrNoImagesFound)
			return
		}

		amount, err := strconv.Atoi(ctx.Args.Get(0))
		if err != nil {
			ctx.ReplyError("error parsing integer", err)
		}

		ReplyImage(ctx, fn(images[0], amount))
	}
}

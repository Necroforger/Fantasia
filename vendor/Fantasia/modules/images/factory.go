package images

import (
	"Fantasia/modules/images/animate"
	"Fantasia/system"
	"image"
	"strconv"

	"github.com/nfnt/resize"
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

// EffectOptions are Effect options to be passed to the
// Effect creation commands.
type EffectOptions struct {
	Max, Min                   float64
	ConstrainMax, ConstrainMin bool
	UseDefault                 bool
	Default                    float64
}

// NewEffectCommandFloat produces an effect command that accepts an image and a float
func (m *Module) NewEffectCommandFloat(fn func(img image.Image, amount float64) *image.RGBA, opts ...EffectOptions) func(ctx *system.Context) {
	return func(ctx *system.Context) {
		if ctx.Args.Get(0) == "" && (len(opts) != 0 && !opts[0].UseDefault) {
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

		var amount float64
		if ctx.Args.Get(0) != "" {
			amount, err = strconv.ParseFloat(ctx.Args.Get(0), 64)
			if err != nil {
				ctx.ReplyError("Error parsing float")
				return
			}
		} else if len(opts) != 0 && opts[0].UseDefault {
			amount = opts[0].Default
		}

		if len(opts) != 0 {
			o := opts[0]
			if o.ConstrainMax && amount > o.Max {
				amount = o.Max
			}
			if o.ConstrainMin && amount < o.Min {
				amount = o.Min
			}
		}

		ctx.Reply(amount)
		ReplyImage(ctx, fn(images[0], amount))
	}
}

// NewGifCommand creates an animated effect command
func (m *Module) NewGifCommand(fn animate.Effect, opts *animate.Options) func(ctx *system.Context) {
	return func(ctx *system.Context) {
		images, err := m.PullImages(1, ctx.Msg.ChannelID, ctx.Msg)
		if err != nil {
			ctx.ReplyError(err)
			return
		}
		if len(images) == 0 {
			ctx.ReplyError("No images found")
			return
		}

		// Resize image to something small
		images[0] = resize.Thumbnail(300, 300, images[0], resize.NearestNeighbor)
		ReplyGif(ctx, animate.Animate(images[0], fn, opts))
	}
}

// NewEffectCommandInt ...
func (m *Module) NewEffectCommandInt(fn func(img image.Image, amount int) *image.RGBA, opts ...EffectOptions) func(ctx *system.Context) {
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

		if len(opts) != 0 {
			o := opts[0]
			if o.ConstrainMax && amount > int(o.Max) {
				amount = int(o.Max)
			}
			if o.ConstrainMin && amount < int(o.Min) {
				amount = int(o.Min)
			}
		}

		ReplyImage(ctx, fn(images[0], amount))
	}
}

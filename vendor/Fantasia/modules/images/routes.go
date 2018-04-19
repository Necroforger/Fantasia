package images

import (
	"Fantasia/modules/images/animate"
	"Fantasia/modules/images/exeffects"
	"image"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/transform"

	"github.com/anthonynsimon/bild/adjust"

	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/effect"
)

// CreateCommands adds the image commands
func (m *Module) CreateCommands() {
	r := m.Sys.CommandRouter

	// =================== Adjustments ========================
	// !______________________________________________________!
	r.On("hue", m.NewEffectCommandInt(adjust.Hue)).Set("", "adjusts the hue of the supplied image;\nex: `hue [degree]`")
	r.On("animatehue", m.NewGifCommand(exeffects.Hue, &animate.Options{From: 0, To: 360, Increment: 10, Delay: 10})).Set("", "Creates an image with an animated hue")

	r.On("saturation", m.NewEffectCommandFloat(adjust.Saturation)).Set("", "Adjusts the saturation of an image;\nex: `saturation [value]`")
	r.On("contrast", m.NewEffectCommandFloat(adjust.Contrast)).Set("", "Adjusts the contrast of an image;\nex: `contrast [value]`")
	r.On("gamma", m.NewEffectCommandFloat(adjust.Gamma)).Set("", "Adjusts the gamma of an image;\nex: `gamma [value]`")
	r.On("brightness", m.NewEffectCommandFloat(adjust.Brightness)).Set("", "Adjusts the brightness of an image;\nex: `brightness [value]`")

	// =================== Effects ============================
	// !______________________________________________________!
	r.On("pixelate", m.NewEffectCommandFloat(exeffects.Pixelate, constraints(1, 0, ODefault(0.1)))).Set("", "Piexelates an image\nUsage: `pixelate [scale 0-1.0]")
	r.On("jpegify", m.NewEffectCommandFloat(exeffects.Jpegify, constraints(100, 0, ODefault(1)))).Set("", "Almost as good as lossy audio\nUsage: `jpegify [quality 0-100]`")
	r.On("animatejpegify", m.NewGifCommand(exeffects.Jpegify, &animate.Options{From: 50, To: 1, Increment: 5, Delay: 10})).Set("", "Animates the jpegification of an image")
	r.On("textify", m.CmdTextify).Set("", "Converts an image to text")
	r.On("overlay", m.CmdOverlay).Set("", "Overlay one image ontop of another")

	r.On("sharpen", m.NewEffectCmdSingle(effect.Sharpen)).Set("", "Applies a sharpen effect to an image")
	r.On("invert", m.NewEffectCmdSingle(effect.Invert)).Set("", "inverts an image")
	r.On("emboss", m.NewEffectCmdSingle(effect.Emboss)).Set("", "applies an emboss effect to an image")
	r.On("sepia", m.NewEffectCmdSingle(effect.Sepia)).Set("", "applies a sepia effect to an image")
	r.On("sobel", m.NewEffectCmdSingle(effect.Sobel)).Set("", "applies a sobel effect to an image")
	r.On("grayscale", m.NewEffectCmdSingle(func(img image.Image) *image.RGBA { return clone.AsRGBA(effect.Grayscale(img)) })).Set("", "applies a grayscale effect to an image")
	r.On("edgedetect", m.NewEffectCmdSingle(exeffects.EdgeDetect)).Set("", "Perform an edge detection")

	r.On("erode", m.NewEffectCommandFloat(effect.Erode, constraints(3))).Set("", "applies an erode effect to an image\nUsage: `erode [radius]`")
	r.On("dilate", m.NewEffectCommandFloat(effect.Dilate, constraints(5, 0))).Set("", "Dilate the image.\nUsage: `dilate [radius]`")

	// ================== Blur ============================
	// !__________________________________________________!
	r.On("blur", m.NewEffectCommandFloat(blur.Gaussian, constraints(10, 0))).Set("", "creates a gaussian blur:\nUsage: `blur [radius]`")
	r.On("boxblur", m.NewEffectCommandFloat(blur.Box, constraints(10, 0))).Set("", "creates a box blue:\nUsage: `boxblur [radius]")

	// ================= Transform =======================
	// !_________________________________________________!
	r.On("rotate", m.NewEffectCommandFloat(exeffects.Rotate, constraints(360, -360))).Set("", "rotate an image [n] degrees\nUsage: `rotate [degrees]`")
	r.On("animaterotate", m.NewGifCommand(exeffects.Rotate, &animate.Options{From: 0, To: 360, Increment: 10, Delay: 10})).Set("", "Animates the rotation of an image")
	r.On("shearh", m.NewEffectCommandFloat(transform.ShearH, constraints(360, -360))).Set("", "shear horizontal\nUsage: `shearh [amount]`")
	r.On("shearv", m.NewEffectCommandFloat(transform.ShearV, constraints(360, -360))).Set("", "shear vertical\nUsage: `shearh [amount]`")
	r.On("fliph", m.NewEffectCmdSingle(transform.FlipH)).Set("", "flip an image over the horizontal axis")
	r.On("flipv", m.NewEffectCmdSingle(transform.FlipV)).Set("", "flip an image over the vertical axis")
}

// constraints [max] [min]
// construncts a constrain
func constraints(args ...interface{}) EffectOptions {
	c := EffectOptions{}
	if len(args) > 0 {
		if n, ok := args[0].(float64); ok {
			c.ConstrainMax = true
			c.Max = n
		}
	}
	if len(args) > 1 {
		if n, ok := args[1].(float64); ok {
			c.ConstrainMin = true
			c.Min = n
		}
	}

	for _, v := range args {
		switch t := v.(type) {
		case TOptionDefault:
			c.UseDefault = true
			c.Default = float64(t)
		}
	}

	return c
}

// TOptionDefault ...
type TOptionDefault float64

// ODefault ...
func ODefault(data float64) TOptionDefault {
	return TOptionDefault(data)
}

package images

import (
	"github.com/Necroforger/Fantasia/modules/images/animate"
	"github.com/Necroforger/Fantasia/modules/images/exeffects"
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
	r.On("hue", m.NewEffectCommandFloat(exeffects.Hue)).Set("", "adjusts the hue of the supplied image;\nex: `hue [degree]`")
	r.On("animatehue", m.NewGifCommand(exeffects.Hue, &animate.Options{From: 0, To: 360, Increment: 10, Delay: 10})).Set("", "Creates an image with an animated hue")

	r.On("saturation", m.NewEffectCommandFloat(adjust.Saturation)).Set("", "Adjusts the saturation of an image;\nex: `saturation [value]`")
	r.On("contrast", m.NewEffectCommandFloat(adjust.Contrast)).Set("", "Adjusts the contrast of an image;\nex: `contrast [value]`")
	r.On("gamma", m.NewEffectCommandFloat(adjust.Gamma)).Set("", "Adjusts the gamma of an image;\nex: `gamma [value]`")
	r.On("brightness", m.NewEffectCommandFloat(adjust.Brightness)).Set("", "Adjusts the brightness of an image;\nex: `brightness [value]`")

	// =================== Effects ============================
	// !______________________________________________________!
	r.On("pixelate", m.NewEffectCommandFloat(exeffects.Pixelate, constraints(oMax(1), oMin(0), oDefault(0.1)))).Set("", "Piexelates an image\nUsage: `pixelate [scale 0-1.0]")
	r.On("jpegify", m.NewEffectCommandFloat(exeffects.Jpegify, constraints(oMax(100), oMin(0), oDefault(1)))).Set("", "Almost as good as lossy audio\nUsage: `jpegify [quality 0-100]`")
	r.On("animatejpegify", m.NewGifCommand(exeffects.Jpegify, &animate.Options{From: 100, To: 1, Increment: 5, Delay: 10})).Set("", "Animates the jpegification of an image")
	r.On("textify", m.CmdTextify).Set("", "Converts an image to text")
	r.On("overlay", m.NewBlendCommand(exeffects.Overlay)).Set("", "Overlays the last sent image over the image sent before it")
	r.On("duoimage", m.NewBlendCommand(exeffects.DuoImage)).Set("", "Merge two images so that one is visible only on discord light theme, "+
		"and the other only visible on discord dark theme")

	r.On("sharpen", m.NewEffectCmdSingle(effect.Sharpen)).Set("", "Applies a sharpen effect to an image")
	r.On("invert", m.NewEffectCmdSingle(effect.Invert)).Set("", "inverts an image")
	r.On("emboss", m.NewEffectCmdSingle(effect.Emboss)).Set("", "applies an emboss effect to an image")
	r.On("sepia", m.NewEffectCmdSingle(effect.Sepia)).Set("", "applies a sepia effect to an image")
	r.On("sobel", m.NewEffectCmdSingle(effect.Sobel)).Set("", "applies a sobel effect to an image")
	r.On("grayscale", m.NewEffectCmdSingle(func(img image.Image) *image.RGBA { return clone.AsRGBA(effect.Grayscale(img)) })).Set("", "applies a grayscale effect to an image")
	r.On("edgedetect", m.NewEffectCmdSingle(exeffects.EdgeDetect)).Set("", "Perform an edge detection")

	r.On("erode", m.NewEffectCommandFloat(effect.Erode, constraints(oMax(5)))).Set("", "applies an erode effect to an image\nUsage: `erode [radius]`")
	r.On("dilate", m.NewEffectCommandFloat(effect.Dilate, constraints(oMax(5), oMin(0)))).Set("", "Dilate the image.\nUsage: `dilate [radius]`")
	r.On("animatedilate", m.NewGifCommand(effect.Dilate, &animate.Options{From: 0, To: 6, Increment: 0.5, Delay: 10, StopFor: 30, LoopBackwards: true}))

	// ================== Blur ============================
	// !__________________________________________________!
	r.On("blur", m.NewEffectCommandFloat(blur.Gaussian, constraints(oMax(10), oMin(0)))).Set("", "creates a gaussian blur:\nUsage: `blur [radius]`")
	r.On("boxblur", m.NewEffectCommandFloat(blur.Box, constraints(oMax(10), oMin(0)))).Set("", "creates a box blue:\nUsage: `boxblur [radius]")

	// ================= Transform =======================
	// !_________________________________________________!
	r.On("rotate", m.NewEffectCommandFloat(exeffects.Rotate, constraints(oMax(360), oMin(-360)))).Set("", "rotate an image [n] degrees\nUsage: `rotate [degrees]`")
	r.On("animaterotate", m.NewGifCommand(exeffects.Rotate, &animate.Options{From: 0, To: 360, Increment: 10, Delay: 10})).Set("", "Animates the rotation of an image")
	r.On("shearh", m.NewEffectCommandFloat(transform.ShearH, constraints(oMax(360), oMin(-360)))).Set("", "shear horizontal\nUsage: `shearh [amount]`")
	r.On("shearv", m.NewEffectCommandFloat(transform.ShearV, constraints(oMax(360), oMin(-360)))).Set("", "shear vertical\nUsage: `shearh [amount]`")
	r.On("fliph", m.NewEffectCmdSingle(transform.FlipH)).Set("", "flip an image over the horizontal axis")
	r.On("flipv", m.NewEffectCmdSingle(transform.FlipV)).Set("", "flip an image over the vertical axis")
}

// constraints simplifies adding constraints a to a command
func constraints(args ...interface{}) EffectOptions {
	c := EffectOptions{}

	for _, v := range args {
		switch t := v.(type) {
		case oDefault:
			c.UseDefault = true
			c.Default = float64(t)
		case oMax:
			c.ConstrainMax = true
			c.Max = float64(t)
		case oMin:
			c.ConstrainMin = true
			c.Min = float64(t)
		}
	}

	return c
}

type oDefault float64

type oMax float64

type oMin float64

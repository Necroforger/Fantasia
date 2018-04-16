package images

import (
	"image"

	"github.com/anthonynsimon/bild/adjust"

	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/effect"
)

// CreateCommands adds the image commands
func (m *Module) CreateCommands() {
	r := m.Sys.CommandRouter

	// Adjustments
	r.On("hue", m.NewEffectCommandInt(adjust.Hue)).Set("", "adjusts the hue of the supplied image;\nex: `hue [degree]`")
	r.On("saturate", m.NewEffectCommandFloat(adjust.Saturation)).Set("", "Adjusts the saturation of an image;\nex: `saturation [value]`")
	r.On("contrast", m.NewEffectCommandFloat(adjust.Contrast)).Set("", "Adjusts the contrast of an image;\nex: `contrast [value]`")
	r.On("gamma", m.NewEffectCommandFloat(adjust.Gamma)).Set("", "Adjusts the gamma of an image;\nex: `gamma [value]`")
	r.On("brightness", m.NewEffectCommandFloat(adjust.Brightness)).Set("", "Adjusts the brightness of an image;\nex: `brightness [value]`")

	// Effects
	r.On("sharpen", m.NewEffectCmdSingle(effect.Sharpen)).Set("", "Applies a sharpen effect to an image")
	r.On("invert", m.NewEffectCmdSingle(effect.Invert)).Set("", "inverts an image")
	r.On("emboss", m.NewEffectCmdSingle(effect.Emboss)).Set("", "applies an emboss effect to an image")
	r.On("sepia", m.NewEffectCmdSingle(effect.Sepia)).Set("", "applies a sepia effect to an image")
	r.On("sobel", m.NewEffectCmdSingle(effect.Sobel)).Set("", "applies a sobel effect to an image")
	r.On("erode", m.NewEffectCommandFloat(effect.Erode)).Set("", "applies an erode effect to an image\nUsage: `erode [radius]`")
	r.On("dilate", m.NewEffectCommandFloat(effect.Dilate)).Set("", "Dilate the image.\nUsage: `dilate [radius]`")
	r.On("grayscale", m.NewEffectCmdSingle(func(img image.Image) *image.RGBA { return clone.AsRGBA(effect.Grayscale(img)) })).Set("", "applies a grayscale effect to an image")
	r.On("edgedetect", m.NewEffectCmdSingle(edgeDetect)).Set("", "Perform an edge detection")
}

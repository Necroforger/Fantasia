package exeffects

import (
	"image"

	"github.com/anthonynsimon/bild/adjust"
)

// Hue wraps adjust.Hue with a float64
func Hue(img image.Image, amount float64) *image.RGBA {
	return adjust.Hue(img, int(amount))
}

package images

import (
	"image"

	"github.com/anthonynsimon/bild/transform"
)

// Rotate rotates an image
func Rotate(src image.Image, amount float64) *image.RGBA {
	return transform.Rotate(src, amount, &transform.RotationOptions{ResizeBounds: true})
}

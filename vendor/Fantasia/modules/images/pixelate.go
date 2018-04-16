package images

import (
	"image"

	"github.com/anthonynsimon/bild/clone"

	"github.com/nfnt/resize"
)

// Pixelate pixelates an image
func Pixelate(src image.Image, amount float64) *image.RGBA {
	dx, dy := src.Bounds().Dx(), src.Bounds().Dy()
	ndx, ndy := int(amount*float64(dx)), int(amount*float64(dy))

	return clone.AsRGBA(
		resize.Resize(
			uint(dx),
			uint(dy),
			resize.Resize(
				uint(ndx),
				uint(ndy),
				src,
				resize.NearestNeighbor,
			),
			resize.NearestNeighbor,
		),
	)
}

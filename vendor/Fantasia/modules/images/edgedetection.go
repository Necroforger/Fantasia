package images

import (
	"image"

	"github.com/anthonynsimon/bild/convolution"
)

// EdgeDetect detects edges in an image
func EdgeDetect(img image.Image) *image.RGBA {
	k := convolution.Kernel{
		Matrix: []float64{
			-1, -1, -1,
			-1, 8, -1,
			-1, -1, -1,
		},
		Width:  3,
		Height: 3,
	}
	return convolution.Convolve(img, &k, &convolution.Options{Bias: 0, Wrap: false, KeepAlpha: true})
}

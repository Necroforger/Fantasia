package images

import (
	"image"
	"image/jpeg"
	"io"

	"github.com/anthonynsimon/bild/clone"
)

// Jpegify converts an image to a low quality jpeg
func Jpegify(src image.Image, quality float64) *image.RGBA {
	rd, wr := io.Pipe()
	go func() {
		jpeg.Encode(wr, src, &jpeg.Options{Quality: int(quality)})
		wr.Close()
	}()
	dst, _ := jpeg.Decode(rd)
	return clone.AsRGBA(dst)
}

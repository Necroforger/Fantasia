package exeffects

import (
	"image"
	"image/draw"

	"github.com/anthonynsimon/bild/clone"
	"github.com/nfnt/resize"
)

// Overlay overlays two images on to each other
func Overlay(srca, srcb image.Image) *image.RGBA {
	dst := clone.AsRGBA(srcb)
	srca = resize.Resize(0, uint(srcb.Bounds().Dy()), srca, resize.NearestNeighbor)
	draw.Draw(dst, dst.Bounds(), srca, image.ZP, draw.Over)
	return dst
}

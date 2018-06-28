package exeffects

import (
	"github.com/Necroforger/Fantasia/modules/images/duoimage"
	"image"
)

// DuoImage merges two images together so that one is visible on dark theme
// And the other is visible on discord light theme
func DuoImage(srca, srcb image.Image) *image.RGBA {
	return duoimage.Merge(srca, srcb, &duoimage.Options{
		Mask: duoimage.MaskHorizontal,
	})
}

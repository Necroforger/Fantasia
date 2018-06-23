// Package duoimage helps create images that differ depending on what theme you are using
// on Discord.
package duoimage

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/nfnt/resize"

	"github.com/anthonynsimon/bild/parallel"

	"github.com/anthonynsimon/bild/clone"
)

// Color constants
const (
	DarkThemeClr  = int(0x36393Eff)
	LightThemeClr = int(0xffffffff)
)

// MaskOption is an enum of mask types
type MaskOption int

// Mask constants
const (
	MaskHorizontal MaskOption = iota
	MaskVertical
	MaskDiagonal
)

// Options are the options
type Options struct {
	// Mask is one of the MaskOption types
	// EG. MaskVertical, MaskHorizontal, MaskDiagonal
	Mask MaskOption
}

// Merge merges two images with the specified mask option
func Merge(srca, srcb image.Image, o *Options) *image.RGBA {
	if o == nil {
		o = &Options{
			Mask: MaskHorizontal,
		}
	}

	// Resize the images to be the same width
	srca = resize.Resize(0, 300, srca, resize.NearestNeighbor)
	srcb = resize.Resize(0, 300, srcb, resize.NearestNeighbor)

	var (
		xmod = 1
		ymod = 1
	)

	// Set the mask options depending on the specified
	// Mask type
	switch o.Mask {
	case MaskHorizontal:
		xmod = 1
		ymod = 2
	case MaskVertical:
		xmod = 2
		ymod = 1
	case MaskDiagonal:
		xmod = 2
		ymod = 2
	}

	log.Println(xmod, ymod)

	// The destination image is the same size as srca.
	dst := image.NewRGBA(srca.Bounds())
	ca, cb := clone.AsRGBA(srca), clone.AsRGBA(srcb)

	// Reduce the colour schemes of the image to
	// A transparent or colour of the given theme
	ca = Reduce(ca, []int{0x0, LightThemeClr})
	cb = Reduce(cb, []int{DarkThemeClr, 0x00})

	// debugWriteImage(ca, "tempCA.png")
	// debugWriteImage(cb, "tempCB.png")

	// Pad the pixels of CB to be the same dimensions as CA
	// cb = clone.Pad(cb, ca.Rect.Dx(), ca.Rect.Dy(), clone.NoFill)
	cb = Pad(cb, ca.Rect.Dx(), cb.Rect.Dy())

	width := ca.Rect.Dx()

	parallel.Line(dst.Rect.Dy(), func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				idx := y*ca.Stride + x*4
				if x%xmod != 0 || y%ymod != 0 {
					// if y%2 != 0 {
					dst.Pix[idx+0] = ca.Pix[idx+0]
					dst.Pix[idx+1] = ca.Pix[idx+1]
					dst.Pix[idx+2] = ca.Pix[idx+2]
					dst.Pix[idx+3] = ca.Pix[idx+3]
				} else {
					dst.Pix[idx+0] = cb.Pix[idx+0]
					dst.Pix[idx+1] = cb.Pix[idx+1]
					dst.Pix[idx+2] = cb.Pix[idx+2]
					dst.Pix[idx+3] = cb.Pix[idx+3]
				}
			}
		}
	})

	return dst
}

// Pad pads an image to the dimensions width and height
//     src    : image to pad
//     width  : width of the new dimensions
//     height : height of the new dimensions
func Pad(source image.Image, width, height int) *image.RGBA {
	src := clone.AsRGBA(source)
	var w, h int

	if src.Rect.Dx() > width {
		w = src.Rect.Dx()
	} else {
		w = width
	}

	if src.Rect.Dy() > height {
		h = src.Rect.Dy()
	} else {
		h = height
	}

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Over)

	return dst
}

func debugWriteImage(img image.Image, dst string) {
	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	png.Encode(f, img)
}

// Reduce reduces the colours of an image into the given palette
//     src     : Source image.
//     palette : Palette of colours.
func Reduce(src image.Image, palette []int) *image.RGBA {
	ca := clone.AsRGBA(src)
	dst := image.NewRGBA(ca.Rect)

	// ab := averageBrightness(ca)

	width := ca.Rect.Dx()
	parallel.Line(ca.Rect.Dy(), func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				idx := y*ca.Stride + x*4
				r := ca.Pix[idx+0]
				g := ca.Pix[idx+1]
				b := ca.Pix[idx+2]

				brightness := float64(int(r)+int(g)+int(b)) / 3
				pidx := int(brightness / 256 * float64(len(palette)))
				dr, dg, db, da := intToRgba(palette[pidx])

				dst.Pix[idx+0] = dr
				dst.Pix[idx+1] = dg
				dst.Pix[idx+2] = db
				dst.Pix[idx+3] = da
			}
		}
	})

	return dst
}

// nearestColor returns the nearest color in a palette.
func nearestColor(src int, palette []int) int {
	var nearestColor int
	var lowestDistance = math.MaxFloat64

	for _, v := range palette {
		r, g, b, a := intToRgba(v)
		r1, g1, b1, a1 := intToRgba(v)

		dist := colorDistance(
			r, g, b, a,
			r1, g1, b1, a1,
		)

		if dist < lowestDistance {
			lowestDistance = dist
			nearestColor = v
		}
	}

	return nearestColor
}

// colorDistance returns the euclidean distance between two colours
func colorDistance(r1, b1, g1, a1, r2, b2, g2, a2 uint8) float64 {
	return math.Sqrt(
		math.Pow(float64(r2)-float64(r1), 2) +
			math.Pow(float64(g2)-float64(g1), 2) +
			math.Pow(float64(b2)-float64(b1), 2) +
			math.Pow(float64(a2)-float64(a1), 2),
	)
}

// rgbaToInt converts RGBA to an integer
func rgbaToInt(r, g, b, a uint8) int {
	return (int(r) << 24) | ((int(g)) << 16) | (int(b) << 8) | (int(a))
}

// intToRGBA converts an integer to RGBA
func intToRgba(c int) (r, g, b, a uint8) {
	r = uint8((c >> 24) & 0xFF)
	g = uint8((c >> 16) & 0xFF)
	b = uint8((c >> 8) & 0xFF)
	a = uint8(c & 0xFF)
	return
}

// averageBrightness computes the average brightness of an image
func averageBrightness(src image.Image) float64 {
	var sum uint64
	ca := clone.AsRGBA(src)

	width := ca.Rect.Dy()
	parallel.Line(ca.Rect.Dy(), func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				idx := y*ca.Stride + x*4
				sum += uint64((ca.Pix[idx+0] + ca.Pix[idx+1] + ca.Pix[idx+2])) / 3
			}
		}
	})

	return float64(sum) / float64(ca.Rect.Dx()*ca.Rect.Dy())
}

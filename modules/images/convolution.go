package images

import (
	"encoding/json"
	"image"
	"image/color"
	"io"
	"math"
	"net/http"
	"strconv"

	// Decoding libs
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"github.com/Necroforger/Fantasia/system"
	"github.com/nfnt/resize"
)

// Filter variables
var (
	MatrixGaussian = [][]float64{
		{1, 4, 6, 4, 1},
		{4, 16, 24, 16, 4},
		{6, 24, 36, 24, 6},
		{4, 16, 24, 16, 4},
		{1, 4, 6, 4, 1},
	}

	MatrixEdgeDetect = [][]float64{
		{-1, -1, -1},
		{-1, 8, -1},
		{-1, -1, -1},
	}

	MatrixMotionBlur = [][]float64{
		{1, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0},
		{0, 0, 0, 0, 1},
	}

	MatrixSharpen = [][]float64{
		{0, 0, 0, 0, 0},
		{0, 0, -1, 0, 0},
		{0, -1, 5, -1, 0},
		{0, 0, -1, 1, 0},
		{0, 0, 0, 0, 0},
	}
)

func getDivisor(matrix [][]float64) int {
	var divisor int
	for _, v := range matrix {
		for _, j := range v {
			divisor += int(j)
		}
	}
	if divisor == 0 {
		divisor = 1
	}
	return divisor
}

// MakeConvolutionFunc returns a convolution command for the supplied matrix
func MakeConvolutionFunc(matrix [][]float64, divisor int, iterations int) func(*system.Context) {
	return func(ctx *system.Context) {
		if len(ctx.Msg.Attachments) == 0 {
			ctx.ReplyError("Please upload an image with this command")
			return
		}

		// Copy the iterations variable so that it is not modified the next time the
		// Command is used.
		it := iterations
		// iterations bypass
		if n, err := strconv.Atoi(ctx.Args.After()); err == nil {
			if n > 20 {
				n = 20
			}
			it = n
		}

		resp, err := http.Get(ctx.Msg.Attachments[0].URL)
		if err != nil {
			ctx.ReplyError("Error fetching attachment")
			return
		}

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			ctx.ReplyError(err)
			return
		}

		img = resize.Thumbnail(500, 500, img, resize.Lanczos3)

		imgresult := img
		for i := 0; i < it; i++ {
			imgresult = Convolute(imgresult, matrix, divisor)
		}

		rd, wr := io.Pipe()
		go func() {
			png.Encode(wr, imgresult)
			wr.Close()
		}()

		ctx.Ses.SendFile(ctx.Msg, "convolution.png", rd)
	}
}

func cmdCustomFilter(ctx *system.Context) {
	if len(ctx.Msg.Attachments) == 0 {
		ctx.ReplyError("Please upload an image with this command")
		return
	}

	resp, err := http.Get(ctx.Msg.Attachments[0].URL)
	if err != nil {
		ctx.ReplyError("Error fetching attachment")
		return
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	var filters [][][]float64
	err = json.Unmarshal([]byte(ctx.Args.After()), &filters)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	var matrixR [][]float64
	var matrixG [][]float64
	var matrixB [][]float64

	if len(filters) == 0 {
		ctx.ReplyError("Please give a three dimensional array of convolution filters to use. The first filter will be used for Red, second: green, third: blue")
		return
	}

	matrixR = filters[0]
	matrixG = matrixR
	matrixB = matrixG

	if len(filters) > 1 {
		matrixG = filters[1]
	}
	if len(filters) > 2 {
		matrixB = filters[2]
	}

	if !validFilter(matrixR) || !validFilter(matrixG) || !validFilter(matrixG) {
		ctx.ReplyError("One of your filters presents an invalid form")
		return
	}

	img = ConvoluteRGB(img, matrixR, matrixG, matrixB, getDivisor(matrixR), getDivisor(matrixG), getDivisor(matrixB))

	rd, wr := io.Pipe()
	go func() {
		png.Encode(wr, img)
		wr.Close()
	}()

	ctx.Ses.SendFile(ctx.Msg, "convolution.png", rd)
}

// Convolute applies a convolution matrix to an image
func Convolute(img image.Image, matrix [][]float64, divisor int) *image.NRGBA {
	round := func(n float64) int {
		return int(math.Floor(n + 0.5))
	}

	var (
		w      = img.Bounds().Dx()
		h      = img.Bounds().Dy()
		result = image.NewNRGBA(img.Bounds())
		mid    = round(float64(len(matrix[0]) / 2.0))
	)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {

			var sumR float64
			var sumG float64
			var sumB float64

			var (
				r uint32
				g uint32
				b uint32
				a uint32
			)

			offset := mid
			for i := -(offset); i <= (offset); i++ {
				for j := -(offset); j <= (offset); j++ {
					r, g, b, a = img.At(x+i, y+j).RGBA()
					sumR += (float64(r) * matrix[i+offset][j+offset])
					sumG += (float64(g) * matrix[i+offset][j+offset])
					sumB += (float64(b) * matrix[i+offset][j+offset])
				}
			}

			getColor := func(n float64) uint8 {
				round := int(n) >> 8
				if round > 255 {
					round = 255
				}
				if round < 0 {
					round = 0
				}
				return uint8(round)
			}

			div := float64(divisor)
			result.Set(x, y, color.NRGBA{
				R: getColor(sumR / div),
				G: getColor(sumG / div),
				B: getColor(sumB / div),

				// G: uint8(g >> 8),
				// B: uint8(b >> 8),
				// R: uint8(r >> 8),
				A: uint8(a >> 8),
			})

		}
	}

	return result
}

// ConvoluteRGB applies individual RGB filters to an image
func ConvoluteRGB(img image.Image, matrixR, matrixG, matrixB [][]float64, divisorR, divisorG, divisorB int) *image.NRGBA {
	round := func(n float64) int {
		return int(math.Floor(n + 0.5))
	}

	var (
		w      = img.Bounds().Dx()
		h      = img.Bounds().Dy()
		result = image.NewNRGBA(img.Bounds())
		midR   = round(float64(len(matrixR[0]) / 2.0))
		midG   = round(float64(len(matrixG[0]) / 2.0))
		midB   = round(float64(len(matrixB[0]) / 2.0))
	)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {

			var sumR float64
			var sumG float64
			var sumB float64

			var (
				r uint32
				g uint32
				b uint32
			)

			// Red filter
			offset := midR
			for i := -(offset); i <= (offset); i++ {
				for j := -(offset); j <= (offset); j++ {
					r, _, _, _ = img.At(x+i, y+j).RGBA()
					sumR += (float64(r) * matrixR[i+offset][j+offset])
				}
			}

			// Green filter
			offset = midG
			for i := -(offset); i <= (offset); i++ {
				for j := -(offset); j <= (offset); j++ {
					_, g, _, _ = img.At(x+i, y+j).RGBA()
					sumG += (float64(g) * matrixR[i+offset][j+offset])
				}
			}

			// Blue filter
			offset = midB
			for i := -(offset); i <= (offset); i++ {
				for j := -(offset); j <= (offset); j++ {
					_, _, b, _ = img.At(x+i, y+j).RGBA()
					sumB += (float64(b) * matrixR[i+offset][j+offset])
				}
			}

			getColor := func(n float64) uint8 {
				round := int(n) >> 8
				if round > 255 {
					round = 255
				}
				if round < 0 {
					round = 0
				}
				return uint8(round)
			}

			result.Set(x, y, color.NRGBA{
				R: getColor(sumR / float64(divisorR)),
				G: getColor(sumG / float64(divisorG)),
				B: getColor(sumB / float64(divisorB)),

				// G: uint8(g >> 8),
				// B: uint8(b >> 8),
				// R: uint8(r >> 8),
				// A: uint8(a >> 8),

				A: 255,
			})

		}
	}

	return result
}

func validFilter(matrix [][]float64) bool {
	if len(matrix) == 0 {
		return false
	}

	rows := len(matrix[0])
	for _, v := range matrix {
		if len(v) != rows {
			return false
		}
	}

	return true
}

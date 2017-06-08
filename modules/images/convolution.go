package images

import (
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

		// iterations bypass
		if n, err := strconv.Atoi(ctx.Args.After()); err == nil {
			if n > 10 {
				n = 10
			}
			iterations = n
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
		for i := 0; i < iterations; i++ {
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

			offset := mid
			for i := -(offset); i <= (offset); i++ {
				for j := -(offset); j <= (offset); j++ {
					r, g, b, _ := img.At(x+i, y+j).RGBA()
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
				A: 255,
			})

		}
	}

	return result
}

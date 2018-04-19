package animate

import (
	"image"
	"image/gif"
	"runtime"
	"sync"

	"github.com/esimov/colorquant"
)

// Effect is an image effect that accepts an amount as a parameter
type Effect func(src image.Image, amount float64) *image.RGBA

// Animate creates a GIF animation from an Effect
func Animate(src image.Image, fn Effect, from, to, increment float64, delay int) *gif.GIF {
	// Set the minimum delay for gifs at 10 ms
	if delay < 0 {
		delay = 0
	}

	var reverse bool

	if from > to {
		from, to = to, from
		reverse = true
	}
	if increment < 0 {
		increment = -increment
	}

	gsize := int((to - from) / increment)
	g := gif.GIF{
		Image:    make([]*image.Paletted, gsize),
		Delay:    make([]int, gsize),
		Disposal: make([]byte, gsize),
	}

	// Set default delays
	for i := 0; i < len(g.Delay); i++ {
		g.Delay[i] = delay
	}

	// Set disposal
	for i := 0; i < len(g.Disposal); i++ {
		g.Disposal[i] = gif.DisposalNone
	}

	// Have the number of running goroutines limited by the number of cpus
	var wg sync.WaitGroup
	ncpu := runtime.GOMAXPROCS(0)
	tokens := make(chan struct{}, ncpu)
	for i := 0; i < ncpu; i++ {
		tokens <- struct{}{}
	}

	wg.Add(gsize)
	for i := from; i < to; i += increment {
		<-tokens
		go func(i float64) {
			dst := image.NewPaletted(src.Bounds(), nil)
			mod := fn(src, i)
			dst.Palette = colorquant.Quant{}.Quantize(mod, 256).(*image.Paletted).Palette

			colorquant.Dither{Filter: [][]float32{ // Floyd steinburg dithering
				[]float32{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
				[]float32{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
				[]float32{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
			}}.Quantize(mod, dst, 256, true, true)

			g.Image[int(i/increment)] = dst
			tokens <- struct{}{}
			wg.Done()
		}(i)
	}

	wg.Wait()

	if reverse {
		for i, j := 0, len(g.Image)-1; i < j; i, j = i+1, j-1 {
			g.Image[i], g.Image[j] = g.Image[j], g.Image[i]
		}
	}

	return &g
}

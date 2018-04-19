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

// FloydSteinberg is a dithering filter
var FloydSteinberg = [][]float32{
	[]float32{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
	[]float32{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
	[]float32{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
}

// Options are animation options
type Options struct {
	From, To, Increment float64

	// Delay of each frame
	Delay int

	// Use dithering
	Dither bool

	// Loop backwards at the end of the loop
	LoopBackwards bool

	// StopFor is the delay to stop for before beginning the next loop
	StopFor int
}

// Animate creates a GIF animation from an Effect
func Animate(src image.Image, fn Effect, opts *Options) *gif.GIF {
	if opts == nil {
		opts = &Options{}
	}

	// Set the minimum delay for gifs at 10 ms
	if opts.Delay < 0 {
		opts.Delay = 0
	}

	var reverse bool

	if opts.From > opts.To {
		opts.From, opts.To = opts.To, opts.From
		reverse = true
	}
	if opts.Increment < 0 {
		opts.Increment = -opts.Increment
	}

	gsize := int(((opts.To - opts.From) / opts.Increment) + 0.5) // Round up
	g := gif.GIF{
		Image:    make([]*image.Paletted, gsize),
		Delay:    make([]int, gsize),
		Disposal: make([]byte, gsize),
	}

	// Set default delays
	for i := 0; i < len(g.Delay); i++ {
		g.Delay[i] = opts.Delay
	}
	if opts.StopFor > 0 {
		g.Delay[gsize-1] = opts.StopFor
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
	for i := opts.From; i < opts.To; i += opts.Increment {
		<-tokens
		go func(i float64) {
			mod := fn(src, i)

			if opts.Dither {
				dst := image.NewPaletted(src.Bounds(), nil)
				dst.Palette = colorquant.Quant{}.Quantize(mod, 256).(*image.Paletted).Palette
				colorquant.Dither{Filter: FloydSteinberg}.Quantize(mod, dst, 256, true, true)
				g.Image[int(i/opts.Increment)] = dst
			} else {
				g.Image[int(i/opts.Increment)] = colorquant.Quant{}.Quantize(mod, 256).(*image.Paletted)
			}

			tokens <- struct{}{}
			wg.Done()
		}(i)
	}
	wg.Wait()

	if opts.LoopBackwards {
		tmp := make([]*image.Paletted, gsize)
		copy(tmp, g.Image)
		for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
			tmp[i], tmp[j] = tmp[j], tmp[i]
		}
		g.Image = append(tmp, g.Image...)
		g.Delay = append(g.Delay, g.Delay...)
		g.Disposal = append(g.Disposal, g.Disposal...)
	}

	if reverse {
		for i, j := 0, len(g.Image)-1; i < j; i, j = i+1, j-1 {
			g.Image[i], g.Image[j] = g.Image[j], g.Image[i]
		}
	}

	return &g
}

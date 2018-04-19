package animate

import (
	"fmt"
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
			fmt.Println(i/increment, " IDX : ", gsize, " I : ", i)
			g.Image[int(i/increment)] = colorquant.Quant{}.Quantize(fn(src, i), 256).(*image.Paletted)
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

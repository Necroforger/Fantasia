package fonts

/*
Package fonts contains a collection of fonts for use in image generation.
Example:
	fnt := truetype.NewFace(fonts.Swanse, &truetype.Options{
		Size: 50,
	})
*/

import (
	"log"

	"github.com/golang/freetype/truetype"
)

// Fonts
var (
	Swanse              *truetype.Font
	Trench100           *truetype.Font
	MonospaceTypewriter *truetype.Font
)

func init() {
	Swanse = loadFont("data/swanse.ttf")
	Trench100 = loadFont("data/trench100.ttf")
	MonospaceTypewriter = loadFont("data/MonospaceTypewriter.ttf")
}

func loadFont(path string) *truetype.Font {
	b, err := Asset(path)
	if err != nil {
		log.Println("package font: Error loading asset: ", path, " ERR: ", err)
		return nil
	}
	fnt, err := truetype.Parse(b)
	if err != nil {
		log.Println("package font: Error loading font: ", path, " ERR: ", err)
		return nil
	}
	return fnt
}

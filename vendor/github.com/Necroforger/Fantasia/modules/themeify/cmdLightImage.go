package themeify

import (
	"image/color"

	"github.com/Necroforger/Fantasia/system"
)

// CmdLightImage produces an image coloured with the background of the light theme
func CmdLightImage(ctx *system.Context) {
	cmdThemeImage(ctx, []color.Color{color.Transparent, color.RGBA{0xff, 0xff, 0xff, 0xff}})
}

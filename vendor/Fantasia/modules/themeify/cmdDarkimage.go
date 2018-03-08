package themeify

import (
	"image/color"

	"Fantasia/system"
)

// CmdDarkImage returns colours an image with the background colour of the dark theme
func CmdDarkImage(ctx *system.Context) {
	cmdThemeImage(ctx, []color.Color{color.RGBA{0x36, 0x39, 0x3E, 0xff}, color.Transparent})
}

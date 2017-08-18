package themeify

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/nfnt/resize"

	"golang.org/x/image/font"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/Fantasia/util"
	"github.com/fogleman/gg"
)

// Constants
const (
	DarkThemeClr  = int(0x36393E)
	LightThemeClr = int(0xffffff)
)

// Module ...
type Module struct{}

// Build ...
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter

	r.On("duoimage", CmdDuoImage).Set("", "Merge two images together so that they may be observed differently from different themes."+
		"Upload images as attachments or a list of URLs. Requires two images to work.")
	r.On("duotext", CmdDuoText).Set("", "Creates text that differes based on the theme you are using."+
		"Pass two strings enclosed by quotes. The first one will be visible in dark theme, the second will be visible in light theme.\n`duotext \"Goodnight world\", \"Hello World\"")
	r.On("darktext", CmdDarkText).Set("", "Creates text coloured with the background of the dark theme")
	r.On("lighttext", CmdLightText).Set("", "Creates text coloured with the background of the light theme")
	r.On("darkimage", CmdDarkImage).Set("", "Create an image coloured with the background color of the dark theme")
	r.On("lightimage", CmdLightImage).Set("", "Create an image with the background colour of the light theme")
}

func mergeImages(img1, img2 image.Image) image.Image {
	var highestX int
	var highestY int
	if img1.Bounds().Dx() > img2.Bounds().Dx() {
		highestX = img1.Bounds().Dx()
	} else {
		highestX = img2.Bounds().Dx()
	}
	if img1.Bounds().Dy() > img2.Bounds().Dy() {
		highestY = img1.Bounds().Dy()
	} else {
		highestY = img2.Bounds().Dy()
	}

	var wg sync.WaitGroup

	maskDark := image.NewRGBA(image.Rect(0, 0, highestX, highestY))
	wg.Add(1)
	go func() {
		for y := 0; y < highestY; y++ {
			for x := 0; x < highestX; x++ {
				if y%2 == 0 {
					maskDark.Set(x, y, color.RGBA{0, 0, 0, 255})
				}
			}

		}
		wg.Done()
	}()

	maskLight := image.NewRGBA(image.Rect(0, 0, highestX, highestY))
	wg.Add(1)
	go func() {
		for y := 0; y < highestY; y++ {
			for x := 0; x < highestX; x++ {
				if y%2 != 0 {
					maskLight.Set(x, y, color.RGBA{0, 0, 0, 255})
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()

	combined := image.NewRGBA(maskDark.Bounds())
	wg.Add(1)
	go func() {
		draw.DrawMask(combined, combined.Bounds(), img2, image.ZP, maskLight, image.ZP, draw.Over)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		draw.DrawMask(combined, combined.Bounds(), img1, image.ZP, maskDark, image.ZP, draw.Over)
		wg.Done()
	}()
	wg.Wait()

	return combined
}

func convertToGrayscale(img image.Image, palette []color.Color) image.Image {
	imgmono := image.NewRGBA(img.Bounds())
	min, max := imgmono.Bounds().Min, imgmono.Bounds().Max
	for y := min.Y; y < max.Y; y++ {
		for x := min.X; x < max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			imgmono.Set(x, y, palette[int(((float32((r+g+b)/3)/65536.0)*float32(len(palette))))])
		}
	}
	return imgmono
}

func createDoubleTextImage(text, text2 string, clr1, clr2 int, font1, font2 font.Face, maxwidth float64) image.Image {
	darktheme := createTextImage(text, clr1, font1, maxwidth)
	whitetheme := createTextImage(text2, clr2, font2, maxwidth)

	return mergeImages(darktheme, whitetheme)
}

func createTextImage(text string, clr int, textfont font.Face, width float64) image.Image {
	r := (clr >> 16) & 0xff
	g := (clr >> 8) & 0xff
	b := clr & 0xff

	c := gg.NewContext(0, 0)
	if textfont != nil {
		c.SetFontFace(textfont)
	}

	var (
		totalHeight float64
		totalWidth  float64
		heights     = []float64{}
		lines       []string
	)

	if width < 0 {
		lines = strings.Split(text, "\n")
	} else {
		lines = c.WordWrap(text, width)
	}

	for _, l := range lines {
		w, h := c.MeasureString(l)
		w *= 1.3
		h *= 1.3
		totalHeight += h
		if w > totalWidth {
			totalWidth = w
		}
		heights = append(heights, h)
	}

	c = gg.NewContext(int(totalWidth), int(totalHeight))
	if textfont != nil {
		c.SetFontFace(textfont)
	}
	c.SetColor(color.RGBA{uint8(r), uint8(g), uint8(b), 255})
	for i := 0; i < len(lines); i++ {
		c.DrawStringAnchored(lines[i], 0, float64(i)*heights[i], 0, 1)
	}

	return c.Image()
}

func cmdThemeImage(ctx *system.Context, clrs []color.Color) {
	var img image.Image

	if images := util.ImagesFromMessage(ctx.Msg); len(images) != 0 {
		img = images[0]
	} else {
		ctx.ReplyNotify("Upload an image or enter an image url: ")
		imgs, err := util.RequestImages(ctx.Ses, ctx.Msg.Author.ID, time.Minute*5)
		if err != nil {
			ctx.ReplyError(ctx.Msg.Author.Mention() + ": Timed out while waiting for images")
			return
		}
		if len(imgs) != 0 {
			img = imgs[0]
		} else {
			ctx.ReplyError("No image was supplied")
			return
		}
	}

	img = convertToGrayscale(img, clrs)
	img = resize.Thumbnail(300, 300, img, resize.Lanczos3)
	rd, wr := io.Pipe()
	go func() {
		png.Encode(wr, img)
		wr.Close()
	}()
	ctx.Ses.DG.ChannelFileSend(ctx.Msg.ChannelID, "image.png", rd)
}

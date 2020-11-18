//go:generate go-bindata -pkg gif ./font
package gif

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"log"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var defaultFace *truetype.Font
var palette []color.Color
var fg, bg image.Image

const dpi float64 = 200
const size float64 = 10
const hinting = font.HintingFull

func init() {
	var err error
	defaultFace, err = loadFace("font/DejaVuSansMono.ttf")
	if err != nil {
		log.Fatalf("Failed to load face: %v", err)
	}

	fg, bg = image.Black, image.White

	var n, i uint16
	n = 64

	step := (color.White.Y - color.Black.Y) / n

	palette = make([]color.Color, n+1)

	for i = 0; i < n; i++ {
		palette[i] = color.Gray16{color.White.Y - i*step}
	}
	palette[n] = color.Black
}

func DefaultFace() *truetype.Font {
	return defaultFace
}

func drawFrame(f *truetype.Font, img draw.Image, lineHeight, heightShift int, frame2d []string) error {
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)
	c.SetHinting(hinting)

	for i, line := range frame2d {
		point := fixed.Point26_6{X: fixed.I(0), Y: fixed.I(((i + 1) * lineHeight) + heightShift)}
		_, err := c.DrawString(line, point)
		if err != nil {
			return err
		}
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func lineBounds(face font.Face, line string) (int, int) {
	height := 0
	width := 0
	for _, r := range []rune(line) {
		bounds, adv, ok := face.GlyphBounds(r)
		if !ok {
			log.Fatalf("No glyph for rune %v", r)
		}
		width += adv.Round() //bounds.Max.X.Round() - bounds.Min.X.Round()
		sh := bounds.Max.Y.Round() - bounds.Min.Y.Round()

		if sh > height {
			height = sh
		}
	}
	return width, height
}

func DrawGif(f *truetype.Font, frames []string, delays []int, out io.Writer) error {
	if len(frames) != len(delays) {
		return errors.New("Number of shots does not match number of delays")
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: hinting,
	})

	frames2d := make([][]string, len(frames))

	maxWidth := 0
	maxHeight := 0
	maxLineHeight := 0
	for i, frame := range frames {
		frames2d[i] = strings.Split(frame, "\n")

		frameHeight := 0
		for _, line := range frames2d[i] {
			lineWidth, lineHeight := lineBounds(face, line)

			log.Printf("Line bounds %d, %d", lineWidth, lineHeight)
			
			maxWidth = max(lineWidth, maxWidth)
			frameHeight += lineHeight
			maxLineHeight = max(maxLineHeight, lineHeight)
		}

		maxHeight = max(maxHeight, frameHeight)
	}

	imageSize := max(maxWidth, maxHeight)
	heightShift := imageSize - maxHeight

	var images []*image.Paletted
	for i, frame := range frames2d {
		img := image.NewPaletted(image.Rect(0, 0, imageSize, imageSize), palette)
		images = append(images, img)

		err := drawFrame(f, img, maxLineHeight, heightShift/2, frame)
		if err != nil {
			return fmt.Errorf("Failed to draw label %d=%v due: %v", i, frame, err)
		}
	}

	return gif.EncodeAll(out, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

func loadFace(asset string) (*truetype.Font, error) {
	data, err := Asset(asset)
	if err != nil {
		return nil, fmt.Errorf("No such font: %s, %v", asset, err)
	}

	return freetype.ParseFont(data)
}

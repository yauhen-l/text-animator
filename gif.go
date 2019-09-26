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
	"runtime"

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
	//defaultFace, err = loadFace("font/Roboto-Regular.ttf")
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

func addLabel(f *truetype.Font, img draw.Image, x, y int, label string) error {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)
	c.SetHinting(hinting)

	_, err := c.DrawString(label, point)
	return err
}

func DrawGif(f *truetype.Font, shots []string, delays []int, out io.Writer) error {
	if len(shots) != len(delays) {
		return errors.New("Number of shots does not match number of delays")
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: hinting,
	})

	//m := f.Bounds(fixed.Int26_6(f.FUnitsPerEm()))

	//sw := m.Max.X.Round() - m.Min.X.Round()
	//	sh := m.Max.Y.Round() - m.Min.Y.Round()
	//	fmt.Printf("%+v, sw=%v, sh=%v", m, sw, sh)

	maxWidth := 0
	maxHeight := 0
	for _, shot := range shots {
		curWidth := 0
		for _, r := range []rune(shot) {
			bounds, adv, ok := face.GlyphBounds(r)
			if !ok {
				return fmt.Errorf("No glyph for rune %v", r)
			}
			sw := adv.Round() //bounds.Max.X.Round() - bounds.Min.X.Round()
			sh := bounds.Max.Y.Round() - bounds.Min.Y.Round()

			curWidth += sw

			if sh > maxHeight {
				maxHeight = sh
			}

			//log.Printf("curWidth=%v, maxHeight=%v, adv=%v", curWidth, maxHeight, adv.Round())
		}

		if curWidth > maxWidth {
			maxWidth = curWidth
		}
	}

	height := maxHeight
	if maxWidth > maxHeight {
		height = maxWidth
	}

	var images []*image.Paletted
	for i, shot := range shots {
		img := image.NewPaletted(image.Rect(0, 0, maxWidth, height), palette)
		images = append(images, img)

		err := addLabel(f, img, 0, height/2+maxHeight/3, shot)
		if err != nil {
			return fmt.Errorf("Failed to draw label %d=%s due: %v", i, shot, err)
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

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Printf("Alloc = %v mb, \tTotalAlloc = %v mb\tSys = %v mb\tNumGC = %v", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

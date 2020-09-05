package processor

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
)

var ErrCoordinates = errors.New("invalid coordinates")

type service struct {
	lights []Light
}

type Service interface {
	AddLight(id string, x float64, y float64) error
	Exact(frame image.Image) []Light
	Area(frame image.Image, size int) []Light
}

// New creates a new service instance
func New() *service {
	return &service{}
}

// AddLight adds a x,y axis for a light location.  We follow the
// standards created by the hue entertainment api.
// The X-axis and Y-axis positions for the location go from -1 to 1.
//          1
//          |
//          |
//  -1 -----+----- 1
//          |
//          |
//          -1
func (s *service) AddLight(id string, x float64, y float64) error {
	if x > 1 || x < -1 {
		return fmt.Errorf("%q: %w", "x out of bounds", ErrCoordinates)
	}

	if y > 1 || y < -1 {
		return fmt.Errorf("%q: %w", "y out of bounds", ErrCoordinates)
	}

	light := Light{
		ID: id,
		Coords: LightCoordinates{
			x, y,
		},
	}
	s.lights = append(s.lights, light)

	log.WithFields(log.Fields{
		"id":          id,
		"coordinates": fmt.Sprintf("%f,%f", x, y),
	}).Debug("added light")
	return nil
}

// Exact takes an image and returns the light id and the
// color that corresponds to it its specific pixel location.
func (s service) Exact(frame image.Image) []Light {
	bounds := frame.Bounds()

	width := bounds.Max.X
	height := bounds.Max.Y

	lights := []Light{}
	for _, l := range s.lights {
		// .At() starts at 0, but the width and height start at 1
		// to keep from going out of bounds we need to reduce it by one.
		coords := l.pixel(width-1, height-1)
		lc, _ := colorful.MakeColor(frame.At(coords.X, coords.Y))
		r, g, b, _ := lc.RGBA()
		log.WithFields(log.Fields{
			"id":    l.ID,
			"pixel": fmt.Sprintf("%d,%d", coords.X, coords.Y),
			"rgb":   fmt.Sprintf("(%d,%d,%d)", r>>8, g>>8, b>>8),
		}).Debug("found color in frame")
		//fmt.Printf("Light ID: %s\n  Location:%d, %d\n  RGB: %d,%d,%d\n", l.ID, x, y, r>>8, g>>8, b>>8)

		lights = append(lights, Light{ID: l.ID, Coords: l.Coords, Color: lc})
	}
	return lights
}

// Area creates a bounding box around the light location
// if a bounding box will break out of the image, its moved
// over until it fits, if it doesn't the location of the light
// will be the center of the box.
// The color returned is the most prominate color in the box.
func (s service) Area(frame image.Image, size int) []Light {
	bounds := frame.Bounds()

	width := bounds.Max.X
	height := bounds.Max.Y

	lights := []Light{}

	for _, l := range s.lights {
		box := l.area(width-1, height-1, size)

		img := image.NewRGBA(box)
		draw.Draw(img, box, frame, box.Min, draw.Src)

		var r, g, b uint32
		for x := box.Min.X; x < box.Max.X; x++ {
			for y := box.Min.Y; y < box.Max.Y; y++ {
				pr, pg, pb, _ := img.At(x, y).RGBA()
				r += pr
				g += pg
				b += pb
			}
		}
		d := uint32(box.Dy() * box.Dx())

		r /= d
		g /= d
		b /= d

		avgColor, _ := colorful.MakeColor(color.NRGBA{uint8(r / 0x101), uint8(g / 0x101), uint8(b / 0x101), 255})
		lights = append(lights, Light{ID: l.ID, Coords: l.Coords, Color: avgColor})
	}

	return lights
}

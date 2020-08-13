package processor

import (
	"image"
	"image/color"
)

// Light is a representation of a bulb with its current location
// and color
type Light struct {
	ID     string
	Coords LightCoordinates
	Color  color.Color
}

// LightCoordinates sets the location of the light based on the
// the display source.
// The X-axis and Y-axis positions for the location go from -1 to 1.
type LightCoordinates struct {
	X float64
	Y float64
}

// pixel converts the x, y axis to a matching pixel location
func (l Light) pixel(width int, height int) image.Point {
	x := int((l.Coords.X + 1) * (float64(width) / 2))
	y := int(((-1 * l.Coords.Y) + 1) * (float64(height) / 2))
	return image.Point{x, y}
}

// area returns the bounding box around the light location.
// if a bounding box will break out of the image, its moved
// over until it fits, if it doesn't the location of the light
// will be the center of the box.
func (l Light) area(width int, height int, area int) image.Rectangle {
	bw := int(float64(width) * (float64(area) / float64(100)))
	bh := int(float64(height) * (float64(area) / float64(100)))

	coords := l.pixel(width, height)
	tl := image.Point{}
	br := image.Point{}
	if coords.X-int(bw/2) < 0 {
		tl.X = 0
		br.X = bw
	} else if coords.X+int(bw/2) > width {
		tl.X = width - bw
		br.X = width
	} else {
		tl.X = coords.X - int(bw/2)
		br.X = coords.X + int(bw/2)
	}

	if coords.Y-int(bh/2) < 0 {
		tl.Y = 0
		br.Y = bh
	} else if coords.Y+int(bh/2) > height {
		tl.Y = height - bh
		br.Y = height
	} else {
		tl.Y = coords.Y - int(bh/2)
		br.Y = coords.Y + int(bh/2)
	}

	rect := image.Rectangle{tl, br}
	return rect
}

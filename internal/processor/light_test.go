package processor

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationPixel(t *testing.T) {
	var lights = []struct {
		light    Light
		width    int
		height   int
		expected image.Point
	}{
		{Light{ID: "middle", Coords: LightCoordinates{0, 0}}, 200, 100, image.Point{100, 50}},
		{Light{ID: "top_right", Coords: LightCoordinates{-1, 1}}, 200, 100, image.Point{0, 0}},
		{Light{ID: "bottom_right", Coords: LightCoordinates{1, -1}}, 200, 100, image.Point{200, 100}},
		{Light{ID: "offset", Coords: LightCoordinates{-.5, .5}}, 200, 100, image.Point{50, 25}},
	}

	for _, tt := range lights {
		coords := tt.light.pixel(tt.width, tt.height)
		assert.Equal(t, tt.expected, coords)
	}
}

func TestLocationArea(t *testing.T) {
	var lights = []struct {
		light        Light
		width        int
		height       int
		expectedRect image.Rectangle
	}{
		{Light{ID: "middle", Coords: LightCoordinates{0, 0}}, 200, 100, image.Rectangle{image.Point{95, 48}, image.Point{105, 52}}},
		{Light{ID: "top_left", Coords: LightCoordinates{-1, 1}}, 200, 100, image.Rectangle{image.Point{0, 0}, image.Point{10, 5}}},
		{Light{ID: "top_right", Coords: LightCoordinates{1, 1}}, 200, 100, image.Rectangle{image.Point{190, 0}, image.Point{200, 5}}},
		{Light{ID: "bottom_left", Coords: LightCoordinates{-1, -1}}, 200, 100, image.Rectangle{image.Point{0, 95}, image.Point{10, 100}}},
		{Light{ID: "bottom_right", Coords: LightCoordinates{1, -1}}, 200, 100, image.Rectangle{image.Point{190, 95}, image.Point{200, 100}}},
		{Light{ID: "offset", Coords: LightCoordinates{-.5, .5}}, 200, 100, image.Rectangle{image.Point{45, 23}, image.Point{55, 27}}},
	}

	for _, tt := range lights {
		r := tt.light.area(tt.width, tt.height, 5)
		assert.Equal(t, tt.expectedRect, r)
	}
}

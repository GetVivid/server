package processor

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAddLight(t *testing.T) {
	var cTests = []struct {
		id       string
		x        float64
		y        float64
		expected error
	}{
		{"middle", 0, 0, nil},
		{"top_left", -1, 1, nil},
		{"top_right", 1, 1, nil},
		{"bottom_left", -1, -1, nil},
		{"bottom_right", 1, -1, nil},
		{"invalid_1", -2, 1, ErrCoordinates},
		{"invalid_2", 0, -2, ErrCoordinates},
	}

	p := New()
	for _, tt := range cTests {
		err := p.AddLight(tt.id, tt.x, tt.y)
		e := errors.Is(err, tt.expected)
		assert.Equal(t, e, true)
	}
}

func TestExactColor(t *testing.T) {
	var lTests = []struct {
		id       string
		x        float64
		y        float64
		expected [3]uint32
	}{
		{"middle", 0, 0, [3]uint32{33, 54, 105}},
		{"top_left", -1, 1, [3]uint32{92, 125, 178}},
		{"top_right", 1, 1, [3]uint32{250, 250, 250}},
		{"bottom_left", -1, -1, [3]uint32{19, 36, 78}},
		{"bottom_right", 1, -1, [3]uint32{25, 38, 43}},
	}

	tstPath, err := filepath.Abs("../../testdata")
	assert.NoError(t, err)

	reader, err := os.Open(filepath.Join(tstPath, "john_wick.jpg"))
	assert.NoError(t, err)

	m, _, err := image.Decode(reader)
	assert.NoError(t, err)
	log.SetLevel(log.InfoLevel)
	p := New()
	for _, tt := range lTests {
		p.AddLight(tt.id, tt.x, tt.y)
	}
	for i, l := range p.Exact(m) {
		assert.Equal(t, lTests[i].id, l.ID)

		r, g, b, _ := l.Color.RGBA()
		assert.Equal(t, lTests[i].expected, [3]uint32{r >> 8, g >> 8, b >> 8})

		assert.Equal(t, lTests[i].x, l.Coords.X)
		assert.Equal(t, lTests[i].y, l.Coords.Y)

	}
}

func TestAreaColor(t *testing.T) {
	var lTests = []struct {
		id       string
		x        float64
		y        float64
		expected [3]uint32
	}{
		{"top_left", -1, 1, [3]uint32{106, 140, 188}},
		{"top_right", 1, 1, [3]uint32{250, 250, 250}},
		{"bottom_left", -1, -1, [3]uint32{24, 41, 88}},
		{"bottom_right", 1, -1, [3]uint32{29, 43, 49}},
		{"middle", 0, 0, [3]uint32{35, 55, 109}},
	}

	tstPath, err := filepath.Abs("../../testdata")
	assert.NoError(t, err)

	reader, err := os.Open(filepath.Join(tstPath, "john_wick.jpg"))
	assert.NoError(t, err)

	m, _, err := image.Decode(reader)
	assert.NoError(t, err)
	log.SetLevel(log.InfoLevel)
	p := New()
	for _, tt := range lTests {
		p.AddLight(tt.id, tt.x, tt.y)
	}
	for i, l := range p.Area(m, 2) {
		assert.Equal(t, lTests[i].id, l.ID)

		r, g, b, _ := l.Color.RGBA()
		assert.Equal(t, lTests[i].expected, [3]uint32{r >> 8, g >> 8, b >> 8})

		assert.Equal(t, lTests[i].x, l.Coords.X)
		assert.Equal(t, lTests[i].y, l.Coords.Y)
	}
}

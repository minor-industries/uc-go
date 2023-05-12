package noise

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tinygo/noise/f32"
	"tinygo/noise/f32b"
)

func TestNoise(t *testing.T) {
	for x := -5.0; x < 5; x += 0.1 {
		for y := -5.0; y < 5; y += 0.1 {
			a := f32.Noise2(x, y)
			b := f32b.Noise2(x, y)
			assert.Equal(t, a, b)
		}
	}

	for x := -5.0; x < 5; x += 0.1 {
		for y := -5.0; y < 5; y += 0.1 {
			a := f32.Noise2(x, y)
			b := f32b.Noise2(x, y)
			assert.Equal(t, a, b)
		}
	}
}

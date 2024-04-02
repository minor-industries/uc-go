package noise

import (
	"fmt"
	"github.com/minor-industries/uc-go/pkg/noise/f32"
	"github.com/minor-industries/uc-go/pkg/noise/f32b"
	"github.com/minor-industries/uc-go/pkg/noise/fixed"
	"github.com/stretchr/testify/assert"
	"testing"
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
			t.Run(fmt.Sprintf("%0.02f %0.02f", x, y), func(t *testing.T) {
				a := f32.Noise2(x, y)
				b := fixed.Noise2(x, y)
				assert.InDelta(t, a, b, 0.025)
			})
		}
	}
}

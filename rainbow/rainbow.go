package rainbow

import (
	"math"
	"tinygo/noise/fixed"
	"tinygo/strip"
)

const (
	meter = 1.0
	inch  = 0.0254 * meter

	ledRingRadius = (15.0/2 - 1) * inch
)

type callback func(t, dt float64)

type App struct {
	Strip *strip.Strip
}

func Rainbow1(app *App, cfg *FaderConfig) callback {
	fade := newFader(app, cfg)

	return func(t, dt float64) {
		fade.fade(0, t)
	}
}

func rainbow2(app *App, cfg *FaderConfig) callback {
	fade := newFader(app, cfg)

	return func(t, dt float64) {
		fade.fade(8*math.Sin(t/4.0), t)
	}
}

type FaderConfig struct {
	TimeScale float64
}

type fader struct {
	positions []complex64
	cfg       *FaderConfig
	app       *App
}

func newFader(app *App, cfg *FaderConfig) *fader {
	return &fader{
		app:       app,
		cfg:       cfg,
		positions: make([]complex64, app.Strip.NumLeds()),
	}
}

const (
	rangeR = 0.5
	rangeG = 0.5
	rangeB = 0.5
)

func (f *fader) fade(
	theta float64,
	t_ float64,
) {
	t := float32(t_)
	scale := float32(f.cfg.TimeScale)
	f.calculatePositions(theta)

	f.app.Strip.Each(func(i int, led *strip.Led) {
		pos := f.positions[i]

		r := rangeR * (0.5 + 0.5*fixed.Noise2(
			real(pos)+000+t*scale,
			imag(pos)+000,
		))

		g := rangeG * (0.5 + 0.5*fixed.Noise2(
			real(pos)+100+t*scale,
			imag(pos)+100,
		))

		b := rangeB * (0.5 + 0.5*fixed.Noise2(
			real(pos)+200+t*scale,
			imag(pos)+200,
		))

		led.R, led.G, led.B = float64(r), float64(g), float64(b)
	})
}

func (f *fader) calculatePositions(theta float64) {
	n := len(f.positions)

	c := complex64(complex(ledRingRadius*3.333, 0))
	c *= complex64(complex(math.Cos(theta), math.Sin(theta)))

	dPhi := (2 * math.Pi) * (1.0 / float64(n))
	incr := complex64(complex(math.Cos(dPhi), math.Sin(dPhi)))

	// Calculate real-world approximate position of LEDS, rotated by theta
	for i := 0; i < n; i++ {
		f.positions[i] = c
		c *= incr
	}
}

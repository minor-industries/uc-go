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
	fade.calculatePositions(0)

	return func(t, dt float64) {
		fade.fade(0, t)
	}
}

func rainbow2(app *App, cfg *FaderConfig) callback {
	fade := newFader(app, cfg)
	fade.calculatePositions(0)

	return func(t, dt float64) {
		fade.fade(8*math.Sin(t/4.0), t)
	}
}

type FaderConfig struct {
	TimeScale float64
}

type fader struct {
	positions []Vec2
	cfg       *FaderConfig
	app       *App
}

func newFader(app *App, cfg *FaderConfig) *fader {
	return &fader{
		app:       app,
		cfg:       cfg,
		positions: make([]Vec2, app.Strip.NumLeds()),
	}
}

const (
	rangeR = 0.5
	rangeG = 0.5
	rangeB = 0.5
)

func (f *fader) fade(
	theta float64,
	t float64,
) {
	//f.calculatePositions(theta)

	f.app.Strip.Each(func(i int, led *strip.Led) {
		pos := f.positions[i]
		led.R = rangeR * (0.5 + 0.5*fixed.Noise2(
			pos.x+000+t*f.cfg.TimeScale,
			0,
		))

		led.G = rangeG * (0.5 + 0.5*fixed.Noise2(
			pos.x+100+t*f.cfg.TimeScale,
			0,
		))

		led.B = rangeB * (0.5 + 0.5*fixed.Noise2(
			pos.x+200+t*f.cfg.TimeScale,
			0,
		))
	})
}

func (f *fader) calculatePositions(theta float64) {
	n := len(f.positions)

	// Calculate real-world approximate position of LEDS, rotated by theta
	for i := 0; i < n; i++ {
		phi := (2 * math.Pi) * (float64(i) / float64(n))
		phi += theta // rotate
		u := Vec2{math.Cos(phi), math.Sin(phi)}
		u = u.Scale(ledRingRadius * 3.333)
		f.positions[i] = u
	}
}

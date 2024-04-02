package bikelights

import (
	"github.com/minor-industries/uc-go/app/bikelights/cfg"
	"github.com/minor-industries/uc-go/pkg/leds"
	"github.com/minor-industries/uc-go/pkg/leds/animations/bounce"
	"github.com/minor-industries/uc-go/pkg/leds/animations/rainbow"
	"github.com/minor-industries/uc-go/pkg/leds/strip"
	"github.com/minor-industries/uc-go/pkg/util"
	"image/color"
	"sync/atomic"
	"time"
)

const (
	ledMaxLevel = 0.5 // brightness level of NeoPxels (0~1)
)

func runLeds(
	config *util.SyncConfig[cfg.Config],
	driver any,
) {
	pixels := make([]color.RGBA, 150)

	ledStrip := strip.NewStrip(config.SnapShot())

	tickDuration := 30 * time.Millisecond

	count := uint32(0)

	animations := map[string]func(t, dt float64){
		"rainbow1": rainbow.Rainbow1(
			&rainbow.App{Strip: ledStrip},
			&rainbow.FaderConfig{TimeScale: 0.3},
		),
		"rainbow2": rainbow.Rainbow2(
			&rainbow.App{Strip: ledStrip},
			&rainbow.FaderConfig{TimeScale: 0.03},
		),
		"bounce": bounce.Bounce(
			&bounce.App{Strip: ledStrip},
		).Tick,
		"white": func(t, dt float64) {
			ledStrip.Each(func(i int, led *strip.Led) {
				led.R = 1.0
				led.G = 1.0
				led.B = 1.0
			})
		},
	}

	f := func() {
		curCfg := config.SnapShot()
		atomic.AddUint32(&count, 1)

		cb := animations[curCfg.CurrentAnimation]
		t := float64(time.Now().UnixNano()) / 1e9
		cb(t, tickDuration.Seconds())
		writeColors(driver, curCfg.Scale, pixels, ledStrip)
	}

	ticker := time.NewTicker(tickDuration)
	for range ticker.C {
		f()
	}
}

func writeColors(
	driver any,
	scale float32,
	pixels []color.RGBA,
	st *strip.Strip,
) {
	convert := func(x float32) uint8 {
		val := x * scale
		return uint8(util.Clamp(0, val, 1.0) * ledMaxLevel * 255.0)
	}

	st.Each(func(i int, led *strip.Led) {
		pixels[i].R = convert(led.R)
		pixels[i].G = convert(led.G)
		pixels[i].B = convert(led.B)
	})

	leds.Write(driver, pixels)
}

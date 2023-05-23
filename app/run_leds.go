package app

import (
	"fmt"
	"image/color"
	"sync/atomic"
	"time"
	"uc-go/bounce"
	"uc-go/cfg"
	"uc-go/leds"
	"uc-go/pio"
	"uc-go/rainbow"
	"uc-go/strip"
	"uc-go/util"
)

const (
	ledMaxLevel = 0.5 // brightness level of NeoPxels (0~1)
)

func RunLeds(
	config *cfg.SyncConfig,
	sm *pio.PIOStateMachine,
) {
	pixels := make([]color.RGBA, 150)

	strip := strip.NewStrip(config.SnapShot())

	tickDuration := 30 * time.Millisecond

	count := uint32(0)
	t0 := time.Now()

	go func() {
		for range time.NewTicker(time.Second).C {
			count := atomic.LoadUint32(&count)
			dt := time.Now().Sub(t0)
			line := fmt.Sprintf(
				"count = %d, t=%s, fps=%0.02f, txfull=%d",
				count,
				time.Now().String(),
				float64(count)/dt.Seconds(),
				atomic.LoadInt64(&leds.TxFullCounter),
			)
			log(line)
		}
	}()

	animations := map[string]func(t, dt float64){
		"rainbow1": rainbow.Rainbow1(
			&rainbow.App{Strip: strip},
			&rainbow.FaderConfig{TimeScale: 0.3},
		),
		"rainbow2": rainbow.Rainbow2(
			&rainbow.App{Strip: strip},
			&rainbow.FaderConfig{TimeScale: 0.03},
		),
		"bounce": bounce.Bounce(
			&bounce.App{Strip: strip},
		).Tick,
	}

	f := func() {
		curCfg := config.SnapShot()
		atomic.AddUint32(&count, 1)

		cb := animations[curCfg.CurrentAnimation]
		t := float64(time.Now().UnixNano()) / 1e9
		cb(t, tickDuration.Seconds())
		writeColors(sm, curCfg.Scale, pixels, strip)
	}

	ticker := time.NewTicker(tickDuration)
	for {
		select {
		case <-ticker.C:
			f()
		}
	}
}

func writeColors(
	sm *pio.PIOStateMachine,
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

	leds.Write(sm, pixels)
}

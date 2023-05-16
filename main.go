//go:build rp2040

package main

import (
	"bytes"
	"fmt"
	"image/color"
	"os"
	"sync/atomic"
	"time"
	"tinygo.org/x/drivers/irremote"
	"tinygo/bounce"
	"tinygo/cfg"
	"tinygo/exe/ir"
	"tinygo/framing"
	"tinygo/leds"
	"tinygo/pio"
	"tinygo/rainbow"
	"tinygo/strip"
	"tinygo/util"
)

const (
	ledMaxLevel = 0.5 // brightness level of NeoPxels (0~1)
)

func main() {
	config := cfg.NewSyncConfig(cfg.Config{
		CurrentAnimation: "rainbow1",
		NumLeds:          150,
		StartIndex:       0,
		Length:           5.0,
		Scale:            0.5,
		MinScale:         0.3,
		ScaleIncr:        0.02,
	})

	irMsg := make(chan irremote.Data, 10)
	ir.Main(func(data irremote.Data) {
		irMsg <- data
	})

	go handleIR(config, irMsg)

	sm := leds.Setup()
	runLeds(config, sm)
}

func handleIR(
	config *cfg.SyncConfig,
	msgs chan irremote.Data,
) {
	for msg := range msgs {
		line := fmt.Sprintf(
			"0x%02x, 0x%02x, 0x%02x 0x%02x",
			msg.Code,
			msg.Flags,
			msg.Command,
			msg.Address,
		)

		frame := framing.Encode([]byte(line))

		framing.Decode(bytes.NewBuffer(nil), func(i []byte) {

		})

		_, _ = os.Stdout.Write(frame)

		switch msg.Command {
		case 0x00: // vol-
			config.ScaleDown()
		case 0x02: // vol+
			config.ScaleUp()
		case 0x10: // 1
			config.SetAnimation("rainbow1")
		case 0x11: // 2
			config.SetAnimation("rainbow2")
		case 0x12: // 2
			config.SetAnimation("bounce")
		}
	}
}

func runLeds(
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
				"count = %d, t=%s, fps=%0.02f",
				count,
				time.Now().String(),
				float64(count)/dt.Seconds(),
			)
			fmt.Printf(line + "\r\n")
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

//go:build rp2040

package main

import (
	"tinygo.org/x/drivers/irremote"
	"uc-go/app"
	"uc-go/cfg"
	"uc-go/exe/ir"
	"uc-go/leds"
	"uc-go/storage"
	"uc-go/util"
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

	storedLogs := util.NewStoredLogs(100)
	storage.Setup(storedLogs)

	go app.DecodeFrames(storedLogs)

	irMsg := make(chan irremote.Data, 10)
	ir.Main(func(data irremote.Data) {
		irMsg <- data
	})

	go app.HandleIR(config, irMsg)

	sm := leds.Setup()
	go app.RunLeds(config, sm)

	select {}
}

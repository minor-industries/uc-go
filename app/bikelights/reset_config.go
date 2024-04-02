package bikelights

import (
	"github.com/minor-industries/uc-go/app/bikelights/cfg"
	"github.com/minor-industries/uc-go/pkg/logger"
	"github.com/minor-industries/uc-go/pkg/storage"
)

func resetConfig(logger logger.Logger) {
	lfs, err := storage.Setup(logger)
	if err != nil {
		panic("no")
	}

	if err := storage.WriteMsgp(logger, lfs, &cfg.DefaultConfig, "/cfg.msgp"); err != nil {
		panic("no")
	}

	select {}
}

func setConfig(logger logger.Logger) {
	lfs, err := storage.Setup(logger)
	if err != nil {
		panic("no")
	}

	c := cfg.Config{
		CurrentAnimation: "rainbow1",
		NumLeds:          150,
		StartIndex:       10,
		Length:           5,
		Scale:            0.5,
		MinScale:         0.04,
		ScaleIncr:        0.02,
	}

	if err := storage.WriteMsgp(logger, lfs, &c, "/cfg.msgp"); err != nil {
		panic("no")
	}
}

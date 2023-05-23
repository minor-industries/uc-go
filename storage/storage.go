package storage

import (
	"fmt"
	"machine"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/util"
)

var (
	blockDevice = machine.Flash
	filesystem  = littlefs.New(blockDevice)
)

func Setup(storedLogs *util.StoredLogs) {
	lfs := filesystem.Configure(&littlefs.Config{
		CacheSize:     512,
		LookaheadSize: 512,
		BlockCycles:   100,
	})

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	storedLogs.Log(fmt.Sprintf(
		"lsblk start=0x%x, end=0x%x",
		machine.FlashDataStart(),
		machine.FlashDataEnd(),
	))

	_ = lfs // TODO: remove

	//n, err := lfs.Size()
	//if err != nil {
	//	storedLogs.Log(fmt.Sprintf("size=%d", n))
	//} else {
	//	storedLogs.Log(errors.Wrap(err, "error: size").Error())
	//}
}

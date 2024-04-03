package main

import (
	"fmt"
	"github.com/minor-industries/uc-go/app/simple-ir/cfg"
	"tinygo.org/x/drivers/irremote"
)

func main() {
	ir := irremote.NewReceiver(cfg.IrPin)
	ir.Configure()

	ch := make(chan irremote.Data, 10)

	ir.SetCommandHandler(func(data irremote.Data) {
		ch <- data
	})

	irCount := 0

	for {
		select {
		case data := <-ch:
			if data.Flags&irremote.DataFlagIsRepeat != 0 {
				continue
			}
			switch data.Command {
			case 16:
				irCount++
			case 17:
				irCount--
			}

			fmt.Println(irCount)
		}
	}
}

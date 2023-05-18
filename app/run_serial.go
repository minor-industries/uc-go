package app

import (
	"fmt"
	"os"
	"uc-go/protocol/framing"
)

func DecodeFrames() {
	ch := make(chan []byte, 10)

	go func() {
		framing.Decode(os.Stdin, func(msg []byte) {
			ch <- msg
		})
	}()

	for msg := range ch {
		reply := fmt.Sprintf("got frame: [%s]", msg)
		log(reply)
	}
}

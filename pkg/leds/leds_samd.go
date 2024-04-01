//go:build itsybitsy_m4

package leds

import "image/color"

var TxFullCounter int64 // TODO

type X struct{}

func Setup() *X {
	return &X{}
}

func Write(driver any, pixels []color.RGBA) {

}

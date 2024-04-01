//go:build pico

package ir

import "machine"

const pinIRIn = machine.GP21 // TODO: should be config?
const powerPin = machine.GP22

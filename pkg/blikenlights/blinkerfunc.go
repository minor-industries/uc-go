package blikenlights

type blinkerFunc struct {
	fcn func(on bool)
}

func BlinkerFunc(fcn func(on bool)) Blinker {
	return &blinkerFunc{fcn: fcn}
}

func (b *blinkerFunc) Set(on bool) {
	b.fcn(on)
}

package cfg

//go:generate msgp

type Config struct {
	CurrentAnimation string
	NumLeds          int
	StartIndex       int
	Length           float32

	Scale     float32
	MinScale  float32
	ScaleIncr float32
}

var DefaultConfig = Config{
	CurrentAnimation: "rainbow1",
	NumLeds:          150,
	StartIndex:       0,
	Length:           5.0,
	Scale:            0.5,
	MinScale:         0.04,
	ScaleIncr:        0.02,
}

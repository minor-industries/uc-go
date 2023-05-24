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

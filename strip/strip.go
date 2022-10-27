package strip

import "tinygo/cfg"

type Led struct {
	R, G, B float64
}

type Strip struct {
	leds []Led
	env  *cfg.Cfg
}

func NewStrip(env *cfg.Cfg) *Strip {
	return &Strip{
		leds: make([]Led, env.NumLeds),
		env:  env,
	}

}

func (s *Strip) Fill(x0, x1 int, color Led) {
	lastIndex := len(s.leds) - 1

	if x0 > x1 {
		x0, x1 = x1, x0
	}

	if x0 < s.env.StartIndex {
		x0 = s.env.StartIndex
	}

	if x1 > lastIndex {
		x1 = lastIndex
	}

	for i := x0; i < x1; i++ {
		s.leds[i] = color
	}
}

func (s *Strip) Each(cb func(i int, led *Led)) {
	for pos := 0; pos < s.env.StartIndex; pos++ {
		s.leds[pos].R = 0
		s.leds[pos].G = 0
		s.leds[pos].B = 0
	}

	i := 0
	for pos := s.env.StartIndex; pos < len(s.leds); pos++ {
		cb(i, &s.leds[pos])
		i++
	}
}

func (s *Strip) Tx(x float64) int {
	// TODO: very strip specific (needs config, etc)
	numLeds := float64(len(s.leds) - s.env.StartIndex)
	return int(numLeds*(x)/s.env.Length) + s.env.StartIndex
}

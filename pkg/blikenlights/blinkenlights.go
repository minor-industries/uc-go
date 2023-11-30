package blikenlights

import (
	"math"
	"time"
)

const (
	Long = math.MaxInt32
)

type Blinker func(on bool)

type Light struct {
	led    Blinker
	ctrl   chan []int
	seq    []int
	pos    int
	remain int
	ticker *time.Ticker
}

func NewLight(led Blinker) *Light {
	return &Light{
		led:    led,
		ctrl:   make(chan []int),
		seq:    []int{32, 32},
		pos:    0,
		remain: 0,
		ticker: time.NewTicker(25 * time.Millisecond),
	}
}

func (li *Light) Run() {
	for {
		select {
		case <-li.ticker.C:
			li.tick()
		case li.seq = <-li.ctrl:
			li.reset()
		}
	}
}

func (li *Light) reset() {
	li.pos = 0
	li.remain = li.seq[li.pos]
}

func (li *Light) tick() {
	if li.remain == 0 {
		li.pos++
		if li.pos >= len(li.seq) {
			li.pos = 0
		}
		li.remain = li.seq[li.pos]
	}
	li.remain--

	on := li.pos%2 == 0
	li.led(on)
}

func (li *Light) Seq(seq []int) {
	li.ctrl <- seq
}

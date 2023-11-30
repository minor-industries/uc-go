package blikenlights

import (
	"math"
	"time"
)

const (
	Long = math.MaxInt32
)

type Blinker interface {
	Set(on bool)
}

type Light struct {
	led  Blinker
	ctrl chan []int

	ticker *time.Ticker

	seq []int

	pos    int
	remain int
}

func NewLight(led Blinker) *Light {
	return &Light{
		led:    led,
		ctrl:   make(chan []int),
		seq:    nil,
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
	if len(li.seq) == 0 {
		li.led.Set(false)
		return
	}

	li.pos = 0
	li.remain = li.seq[li.pos]
}

func (li *Light) tick() {
	if len(li.seq) == 0 {
		return
	}

	for li.remain == 0 {
		li.pos++
		if li.pos >= len(li.seq) {
			li.pos = 0
		}
		li.remain = li.seq[li.pos]
	}
	li.remain--

	on := li.pos%2 == 0
	li.led.Set(on)
}

func (li *Light) Seq(seq []int) {
	li.ctrl <- seq
}

func (li *Light) Off() {
	li.ctrl <- nil
}

func (li *Light) On() {
	li.ctrl <- []int{Long, 0}
}

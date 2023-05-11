package fixed

import (
	"golang.org/x/image/math/fixed"
)

type FloatT struct {
	x fixed.Int26_6
}

func (f FloatT) Add(x FloatT) FloatT {
	return FloatT{f.x + x.x}
}

func (f FloatT) Sub(x FloatT) FloatT {
	return FloatT{f.x - x.x}
}

func (f FloatT) Mul(x FloatT) FloatT {
	return FloatT{f.x.Mul(x.x)}
}

func (f FloatT) Neg() FloatT {
	return FloatT{-f.x}
}

func (f FloatT) Gt(x FloatT) bool {
	return f.x > x.x
}

func (f FloatT) Lt(x FloatT) bool {
	return f.x < x.x
}

func (f FloatT) Int() int {
	return f.x.Floor()
}

func (f FloatT) Float64() float64 {
	return 0.0 // TODO
}

func New(x float32) FloatT {
	//math.Floor(x) // TODO
	return FloatT{} // TODO
}

func INew(i int) FloatT {
	return FloatT{fixed.I(i)}
}

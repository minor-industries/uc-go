package fixed

type FloatT struct {
	float32
}

func (f FloatT) Add(x FloatT) FloatT {
	return FloatT{f.float32 + x.float32}
}

func (f FloatT) Sub(x FloatT) FloatT {
	return FloatT{f.float32 - x.float32}
}

func (f FloatT) Mul(x FloatT) FloatT {
	return FloatT{f.float32 * x.float32}
}

func (f FloatT) Div(x FloatT) FloatT {
	return FloatT{f.float32 / x.float32}
}

func (f FloatT) Neg() FloatT {
	return FloatT{-f.float32}
}

func (f FloatT) Gt(x FloatT) bool {
	return f.float32 > x.float32
}

func (f FloatT) Lt(x FloatT) bool {
	return f.float32 < x.float32
}

func (f FloatT) Int() int {
	return int(f.float32)
}

func New(x float32) FloatT {
	return FloatT{x}
}

func INew(i int) FloatT {
	return FloatT{float32(i)}
}

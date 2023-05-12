package fixed

const (
	precision = 6
	scale     = float32(1 << 6)
	invScale  = 1.0 / float64(1<<6)
)

type FloatT struct {
	x int32
}

func (f FloatT) Add(x FloatT) FloatT {
	return FloatT{f.x + x.x}
}

func (f FloatT) Sub(x FloatT) FloatT {
	return FloatT{f.x - x.x}
}

func (f FloatT) Mul(x FloatT) FloatT {
	var result = int64(f.x) * int64(x.x)
	result >>= precision * 2

	return FloatT{int32(result)}
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

func (f FloatT) Int() int32 {
	return f.x >> precision
}

func (f FloatT) Float64() float64 {
	return float64(f.x) * invScale
}

func New(x float32) FloatT {
	return FloatT{x: int32(x * scale)}
}

func INew(i int32) FloatT {
	return FloatT{x: int32(i << precision)}
}

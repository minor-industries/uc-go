package fixed

const (
	precision = 13
	scale     = 1 << precision
	invScale  = 1.0 / (1 << precision)
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
	result >>= precision

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
	if f.x < 0 {
		return -(-f.x >> precision)
	}
	return f.x >> precision
}

func (f FloatT) Float64() float64 {
	return float64(f.x) * invScale
}

func (f FloatT) Float32() float32 {
	return float32(f.x) * invScale
}

func New(x float32) FloatT {
	return FloatT{x: int32(x * scale)}
}

func INew(i int32) FloatT {
	return FloatT{x: i << precision}
}

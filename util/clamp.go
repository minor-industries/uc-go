package util

func Clamp[T float32 | float64 | int](a, x, b T) T {
	if x < a {
		return a
	}

	if x > b {
		return b
	}

	return x
}

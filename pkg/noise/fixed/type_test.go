package fixed

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloatT_Add(t *testing.T) {
	assert.InDelta(
		t,
		5.5,
		New(2.2).Add(New(3.3)).Float64(),
		0.001,
	)

	assert.InDelta(
		t,
		2.2*3.3,
		New(2.2).Mul(New(3.3)).Float64(),
		0.001,
	)

	assert.InDelta(
		t,
		-1.0*3.3,
		New(-1.0).Mul(New(3.3)).Float64(),
		0.002,
	)

	assert.InDelta(
		t,
		-1.0*-3.3,
		New(-1.0).Mul(New(-3.3)).Float64(),
		0.002,
	)

	a := -3.3
	assert.Equal(t, int32(a), New(-3.3).Int())
}

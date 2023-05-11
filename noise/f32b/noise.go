package f32b

type Float32 float32

func (f Float32) Add(x Float32) Float32 {
	return f + x
}

func (f Float32) Sub(x Float32) Float32 {
	return f - x
}

func (f Float32) Mul(x Float32) Float32 {
	return f * x
}

func (f Float32) Div(x Float32) Float32 {
	return f / x
}

var perm = [512]uint8{
	151, 160, 137, 91, 90, 15,
	131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
	190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
	88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
	102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
	135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
	5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
	223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
	251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
	138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
	151, 160, 137, 91, 90, 15,
	131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
	190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
	88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
	102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
	135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
	5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
	223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
	251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
	138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
}

func Q(cond bool, v1 Float32, v2 Float32) Float32 {
	if cond {
		return v1
	}
	return v2
}

func FASTFLOOR(x Float32) int {
	if x > 0 {
		return int(x)
	}
	return int(x) - 1
}

func grad2(hash uint8, x Float32, y Float32) Float32 {
	h := hash & 7       // Convert low 3 bits of hash code
	u := Q(h < 4, x, y) // into 8 simple gradient directions,
	v := Q(h < 4, y, x) // and compute the dot product with (x,y).
	return Q(h&1 != 0, -u, u).Add(Q(h&2 != 0, -2*v, 2*v))
}

// 2D simplex noise
func Noise2(x, y float64) float64 {
	return float64(noise2(Float32(x), Float32(y)))
}

func noise2(x, y Float32) Float32 {

	const F2 = 0.366025403 // F2 = 0.5*(sqrt(3.0)-1.0)
	const G2 = 0.211324865 // G2 = (3.0-Math.sqrt(3.0))/6.0

	var n0, n1, n2 Float32 // Noise contributions from the three corners

	// Skew the input space to determine which simplex cell we're in
	//s := (x + y) * F2 // Hairy factor for 2D
	s := x.Add(y).Mul(F2)
	xs := x.Add(s)
	ys := y.Add(s)
	i := FASTFLOOR(xs)
	j := FASTFLOOR(ys)

	t := Float32(i+j) * G2
	X0 := Float32(i).Sub(t) // Unskew the cell origin back to (x,y) space
	Y0 := Float32(j).Sub(t)
	x0 := x.Sub(X0) // The x,y distances from the cell origin
	y0 := y.Sub(Y0)

	// For the 2D case, the simplex shape is an equilateral triangle.
	// Determine which simplex we are in.
	var i1, j1 int // Offsets for second (middle) corner of simplex in (i,j) coords
	if x0 > y0 {
		i1 = 1
		j1 = 0 // lower triangle, XY order: (0,0)->(1,0)->(1,1)
	} else {
		i1 = 0
		j1 = 1
	} // upper triangle, YX order: (0,0)->(0,1)->(1,1)

	// A step of (1,0) in (i,j) means a step of (1-c,-c) in (x,y), and
	// a step of (0,1) in (i,j) means a step of (-c,1-c) in (x,y), where
	// c = (3-sqrt(3))/6

	x1 := x0.Sub(Float32(i1)).Add(G2) // Offsets for middle corner in (x,y) unskewed coords
	y1 := y0.Sub(Float32(j1)).Add(G2)
	x2 := x0.Sub(1).Add(2 * G2) // Offsets for last corner in (x,y) unskewed coords
	y2 := y0.Sub(1).Add(2 * G2)

	// Wrap the integer indices at 256, to avoid indexing perm[] out of bounds
	ii := i & 0xff
	jj := j & 0xff

	// Calculate the contribution from the three corners
	t0 := Float32(0.5).Sub(x0.Mul(x0)).Sub(y0.Mul(y0))
	if t0 < 0 {
		n0 = 0
	} else {
		t0 = t0.Mul(t0)
		n0 = t0.Mul(t0).Mul(grad2(perm[ii+int(perm[jj])], x0, y0))
	}

	t1 := Float32(0.5).Sub(x1.Mul(x1)).Sub(y1.Mul(y1))
	if t1 < 0 {
		n1 = 0
	} else {
		t1 = t1.Mul(t1)
		n1 = t1.Mul(t1).Mul(grad2(perm[ii+i1+int(perm[jj+j1])], x1, y1))
	}

	t2 := Float32(0.5).Sub(x2.Mul(x2)).Sub(y2.Mul(y2))
	if t2 < 0 {
		n2 = 0
	} else {
		t2 = t2.Mul(t2)
		n2 = t2.Mul(t2).Mul(grad2(perm[ii+1+int(perm[jj+1])], x2, y2))
	}

	// Add contributions from each corner to get the final noise value.
	// The result is scaled to return values in the interval [-1,1].
	return (n0.Add(n1).Add(n2)).Div(0.022108854818853867)
}

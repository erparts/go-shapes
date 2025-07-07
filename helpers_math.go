package shapes

import (
	"math"
)

const MaxFloat16 = 65504

var roFloat32Inf = float32(math.Inf(1))

func Float32Inf() float32 {
	return roFloat32Inf
}

type GoldenRatioGen struct {
	n float64
}

func (gen *GoldenRatioGen) Reset() {
	gen.n = 0
}

func (gen *GoldenRatioGen) Float64() float64 {
	const phi = 1.618033988749895 // golden ratio
	gen.n += 1.0
	if gen.n == 514230.0 {
		gen.n = 1.0
	}
	v := math.Mod(gen.n/phi, 1.0)
	return v
}

func abs[Float float32 | float64](a Float) Float {
	if a < 0 {
		return -a
	}
	return a
}

func lerp[Float float32 | float64](a, b, t Float) Float {
	return a + t*(b-a)
}

// Notice: geometry code is derived from etxt@v0.0.8 emask/helper_funcs.go

// Given two points of a line, it returns its A, B and C
// coefficients from the form "Ax + By + C = 0".
func toLinearFormABC(ox, oy, fx, fy float64) (float64, float64, float64) {
	a, b, c := fy-oy, -(fx - ox), (fx-ox)*oy-(fy-oy)*ox
	return a, b, c
}

// If we had two line equations like this:
// >> a1*x + b1*y = c1
// >> a2*x + b2*y = c2
// We would apply cramer's rule to solve the system:
// >> x = (b2*c1 - b1*c2)/(b2*a1 - b1*a2)
// This function solves this system, but assuming c1 and c2 have
// a negative sign (ax + by + c = 0).
func shortCramer(a1, b1, c1, a2, b2, c2 float64) (float64, float64) {
	xdiv := b2*a1 - b1*a2
	if xdiv == 0 {
		panic("parallel lines")
	}

	// actual application of cramer's rule
	x := (b2*-c1 - b1*-c2) / xdiv
	if b1 != 0 {
		return x, (-c1 - a1*x) / b1
	}
	return x, (-c2 - a2*x) / b2
}

// given a line equation in the form Ax + By + C = 0, it returns
// C1 and C2 such that two new line equations can be created that
// are parallel to the original line, but at distance 'dist' from it
func parallelsAtDist(a, b, c float64, dist float64) (float64, float64) {
	norm := math.Hypot(a, b)
	if norm == 0 {
		return c, c // degenerate case
	}
	shift := dist * norm
	return c - shift, c + shift
}

// gaussian elimination 8x8 homogeneous linear system solver
func gaussSolver8x8(sys [8][8]float32, weights [8]float32) [8]float32 {
	var x [8]float32
	for i := range 8 {
		// find pivot
		maxRow := i
		for k := i + 1; k < 8; k++ {
			if abs(sys[k][i]) > abs(sys[maxRow][i]) {
				maxRow = k
			}
		}

		// swap rows
		sys[i], sys[maxRow] = sys[maxRow], sys[i]
		weights[i], weights[maxRow] = weights[maxRow], weights[i]

		// eliminate
		for k := i + 1; k < 8; k++ {
			f := sys[k][i] / sys[i][i]
			for j := i; j < 8; j++ {
				sys[k][j] -= f * sys[i][j]
			}
			weights[k] -= f * weights[i]
		}
	}

	// substitution
	for i := 7; i >= 0; i-- {
		sum := float32(0)
		for j := i + 1; j < 8; j++ {
			sum += sys[i][j] * x[j]
		}
		x[i] = (weights[i] - sum) / sys[i][i]
	}

	return x
}

// points are given in clockwise order, from top-left.
// returned matrix is row-major order
func computeHomography(fromQuad, toQuad [4]PointF32) [9]float32 {
	var system [8][8]float32
	var weights [8]float32

	var i int
	for j, pt := range fromQuad {
		u, v := toQuad[j].X, toQuad[j].Y
		system[i+0] = [8]float32{pt.X, pt.Y, 1, 0, 0, 0, -u * pt.X, -u * pt.Y}
		system[i+1] = [8]float32{0, 0, 0, pt.X, pt.Y, 1, -v * pt.X, -v * pt.Y}
		weights[i+0] = u
		weights[i+1] = v
		i += 2
	}

	solutions := gaussSolver8x8(system, weights)
	var homography [9]float32
	_ = copy(homography[:], solutions[:])
	homography[8] = 1.0
	return homography
}

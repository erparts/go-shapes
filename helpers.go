package shapes

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

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

func ColorToF32(clr color.Color) [4]float32 {
	r, g, b, a := clr.RGBA()
	return [4]float32{float32(r) / 65535.0, float32(g) / 65535.0, float32(b) / 65535.0, float32(a) / 65535.0}
}

func colorToF64(clr color.Color) [4]float64 {
	r, g, b, a := clr.RGBA()
	return [4]float64{float64(r) / 65535.0, float64(g) / 65535.0, float64(b) / 65535.0, float64(a) / 65535.0}
}

func f32ToRGBA64(r, g, b, a float32) color.RGBA64 {
	return color.RGBA64{
		R: uint16(r * 65535.0),
		G: uint16(g * 65535.0),
		B: uint16(b * 65535.0),
		A: uint16(a * 65535.0),
	}
}

func lerp[Float float32 | float64](a, b, t Float) Float {
	return a + t*(b-a)
}

func interpColor(ox, oy, fx, fy float32, tlClr, trClr, blClr, brClr [4]float32, x, y float32) [4]float32 {
	u := min(max((x-ox)/(fx-ox), 0), 1)
	v := min(max((y-oy)/(fy-oy), 0), 1)

	var result [4]float32
	for i := range 4 {
		topClr := tlClr[i]*(1-u) + trClr[i]*u
		bottomClr := blClr[i]*(1-u) + brClr[i]*u
		result[i] = topClr*(1-v) + bottomClr*v
	}
	return result
}

func interpVertexColor(a, b ebiten.Vertex, t float32) (cr, cg, cb, ca float32) {
	return lerp(a.ColorR, b.ColorR, t), lerp(a.ColorG, b.ColorG, t), lerp(a.ColorB, b.ColorB, t), lerp(a.ColorA, b.ColorA, t)
}

// Common blend modes not directly exposed on Ebitengine.
var (
	BlendSubtract = ebiten.Blend{
		BlendFactorSourceRGB:        ebiten.BlendFactorOne,
		BlendFactorSourceAlpha:      ebiten.BlendFactorOne,
		BlendFactorDestinationRGB:   ebiten.BlendFactorOne,
		BlendFactorDestinationAlpha: ebiten.BlendFactorOne,
		BlendOperationRGB:           ebiten.BlendOperationReverseSubtract,
		BlendOperationAlpha:         ebiten.BlendOperationAdd,
	}
	BlendMultiply = ebiten.Blend{
		BlendFactorSourceRGB:        ebiten.BlendFactorDestinationColor,
		BlendFactorSourceAlpha:      ebiten.BlendFactorDestinationColor,
		BlendFactorDestinationRGB:   ebiten.BlendFactorOneMinusSourceAlpha,
		BlendFactorDestinationAlpha: ebiten.BlendFactorOneMinusSourceAlpha,
		BlendOperationRGB:           ebiten.BlendOperationAdd,
		BlendOperationAlpha:         ebiten.BlendOperationAdd,
	}
)

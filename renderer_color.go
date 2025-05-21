package shapes

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// FlatPaint draws the mask onto the given target using the renderer vertex colors.
func (r *Renderer) FlatPaint(target, mask *ebiten.Image, ox, oy float32) {
	ensureShaderFlatPaintLoaded()
	r.DrawShaderAt(target, mask, ox, oy, 0, 0, shaderFlatPaint)
}

// SimpleGradient paints a high quality gradient over the given target.
// Common gradient directions are:
//   - Left to right: 0
//   - Right to left: math.Pi
//   - Top to bottom: math.Pi/2
//   - Bottom to top: -math.Pi/2
//   - Top-left to bottom-right: math.Pi/4
//   - Bottom-left to top-right: -math.Pi/4
//   - Top-right to bottom-left: 3*math.Pi/4
//   - Bottom-right to top-left: -3*math.Pi/4
func (r *Renderer) SimpleGradient(target *ebiten.Image, from, to color.RGBA, dirRadians float32) {
	r.Gradient(target, nil, 0, 0, from, to, -1, dirRadians, 1.0)
}

// Gradient paints a high quality gradient over the given target. If mask is nil,
// the target will have the gradient applied starting from (ox, oy) throughout
// the entire image.
//
// CurveFactor allows making the gradient linear (1.0), or ease it towards an
// early start (e.g. 0.5) or late start (e.g. 2.0). Reasonable CurveFactor values
// typically fall in the ~[0.2...4.0] range.
//
// See also [Renderer.SimpleGradient]().
func (r *Renderer) Gradient(target, mask *ebiten.Image, ox, oy float32, from, to color.RGBA, numSteps int, dirRadians, curveFactor float32) {
	if curveFactor < 0.001 {
		panic("curveFactor must be positive above 0.001")
	}

	var srcBounds image.Rectangle
	dstBounds := target.Bounds()
	if mask == nil {
		srcBounds = dstBounds
	} else {
		srcBounds = mask.Bounds()
	}

	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	minX, minY := dstMinX+ox, dstMinY+oy
	maxX, maxY := minX+srcWidth, minY+srcHeight
	r.setDstRectCoords(minX, minY, maxX, maxY)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX, srcMinY, srcMaxX, srcMaxY)
	fromF64, toF64 := colorToF64(from), colorToF64(to)
	memo := r.GetColorF32()
	fromOklab := rgbToOklab([3]float64(fromF64[:3]))
	toOklab := rgbToOklab([3]float64(toF64[:3]))
	r.SetColorF32(float32(toOklab[0]), float32(toOklab[1]), float32(toOklab[2]), float32(toF64[3]))
	r.setFlatCustomVAs(float32(fromOklab[0]), float32(fromOklab[1]), float32(fromOklab[2]), float32(fromF64[3]))

	clear(r.opts.Uniforms)
	r.opts.Uniforms["Area"] = [4]float32{ox, oy, srcWidth, srcHeight}
	r.opts.Uniforms["DirRadians"] = dirRadians
	r.opts.Uniforms["NumSteps"] = numSteps
	r.opts.Uniforms["CurveFactor"] = curveFactor
	if mask != nil { //
		r.opts.Uniforms["UseMask"] = 1
	} else {
		r.opts.Uniforms["UseMask"] = 0
	}

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderGradientLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderGradient, &r.opts)
	r.opts.Images[0] = nil
	r.SetColorF32(memo[0], memo[1], memo[2], memo[3])
}

func rgbToOklab(rgb [3]float64) [3]float64 {
	linR, linG, linB := linearize(rgb[0]), linearize(rgb[1]), linearize(rgb[2])
	x := math.Pow(0.4122214708*linR+0.5363325363*linG+0.0514459929*linB, 1.0/3.0)
	y := math.Pow(0.2119034982*linR+0.6806995451*linG+0.1073969566*linB, 1.0/3.0)
	z := math.Pow(0.0883024619*linR+0.2817188376*linG+0.6299787005*linB, 1.0/3.0)

	l := 0.2104542553*x + 0.7936177850*y - 0.0040720468*z
	a := 1.9779984951*x - 2.4285922050*y + 0.4505937099*z
	b := 0.0259040371*x + 0.7827717662*y - 0.8086757660*z
	return [3]float64{l, a, b}
}

func linearize(colorChan float64) float64 {
	if colorChan >= 0.04045 {
		return math.Pow((colorChan+0.055)/1.055, 2.4)
	} else {
		return colorChan / 12.92
	}
}

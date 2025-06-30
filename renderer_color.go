package shapes

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Direction constants for use with gradient generation functions.
const (
	DirRadsLTR float32 = 0.0          // left to right
	DirRadsRTL float32 = math.Pi      // right to left
	DirRadsTTB float32 = math.Pi / 2  // top to bottom
	DirRadsBTT float32 = -math.Pi / 2 // bottom to top

	DirRadsTLBR float32 = math.Pi / 4      // top-left to bottom-right
	DirRadsBLTR float32 = -math.Pi / 4     // bottmo-left to top-right
	DirRadsTRBL float32 = 3 * math.Pi / 4  // top-right to bottom-left
	DirRadsBRTL float32 = -3 * math.Pi / 4 // bottom-right to top-left
)

func (r *Renderer) NewSimpleGradient(w, h int, from, to color.RGBA, dirRadians float32) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	r.SimpleGradient(img, from, to, dirRadians)
	return img
}

// FlatPaint draws the mask onto the given target using the renderer vertex colors.
func (r *Renderer) FlatPaint(target, mask *ebiten.Image, ox, oy float32) {
	ensureShaderFlatPaintLoaded()
	r.DrawShaderAt(target, mask, ox, oy, 0, 0, shaderFlatPaint)
}

// SimpleGradient paints a high quality gradient over the given target.
// See [DirRadsLTR] and similar constants for common gradient directions.
func (r *Renderer) SimpleGradient(target *ebiten.Image, from, to color.RGBA, dirRadians float32) {
	r.Gradient(target, nil, 0, 0, from, to, -1, dirRadians, 1.0)
}

// Gradient paints a high quality gradient over the given target. If mask is nil,
// the target will have the gradient applied starting from (ox, oy) throughout
// the entire image. See [DirRadsLTR] and similar constants for common gradient
// directions.
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

	r.opts.Uniforms["Area"] = [4]float32{ox, oy, srcWidth, srcHeight}
	r.opts.Uniforms["DirRadians"] = dirRadians
	r.opts.Uniforms["NumSteps"] = numSteps
	r.opts.Uniforms["CurveFactor"] = curveFactor
	if mask != nil {
		r.opts.Uniforms["UseMask"] = 1
	} else {
		r.opts.Uniforms["UseMask"] = 0
	}

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderGradientLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderGradient, &r.opts)
	r.opts.Images[0] = nil
	clear(r.opts.Uniforms)
	r.SetColorF32(memo[0], memo[1], memo[2], memo[3])
}

// GradientRadial paints a high quality radial gradient over the given target.
//
// CurveFactor allows making the gradient linear (1.0), or ease it towards an
// early start (e.g. 0.5) or late start (e.g. 2.0). Reasonable CurveFactor values
// typically fall in the ~[0.2...4.0] range.
//
// Three radiuses are necessary:
//   - fromRadius: distances below this threshold take 'from' color. Use 0.0 if
//     you don't need a solid central area.
//   - transRadius: distances below this threshold but above fromRadius interpolate
//     colors between 'from' and 'to'.
//   - toRadius: distances below this threshold but above transRadius take 'to' color.
//     Distances above this threshold are not painted. Use toRadius = transRadius for
//     a gradient that ends at the given radius, or Float32Inf() if you want 'to'
//     color to extend beyond the gradient radius.
//
// To mask the gradient over an existing image, consider [Renderer.SetBlend](ebiten.BlendSourceIn)
// and similar tricks.
func (r *Renderer) GradientRadial(target *ebiten.Image, cx, cy float32, from, to color.RGBA, fromRadius, transRadius, toRadius float32, numSteps int, curveFactor float32) {
	if curveFactor < 0.001 {
		panic("curveFactor must be positive above 0.001")
	}
	if transRadius < fromRadius || toRadius < transRadius {
		panic("invalid radius values (radiuses must be equal or increasing)")
	}

	dstBounds := target.Bounds()
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	dstWidthF64, dstHeightF64 := float64(dstBounds.Dx()), float64(dstBounds.Dy())
	cxF64, cyF64, toRadiusF64 := float64(cx), float64(cy), float64(toRadius)
	ox, oy := float32(max(math.Floor(cxF64-toRadiusF64), 0)), float32(max(math.Floor(cyF64-toRadiusF64), 0))
	fx, fy := float32(min(math.Ceil(cxF64+toRadiusF64), dstWidthF64)), float32(min(math.Ceil(cyF64+toRadiusF64), dstHeightF64))
	minX, minY := dstMinX+ox, dstMinY+oy
	maxX, maxY := dstMinX+fx, dstMinY+fy
	r.setDstRectCoords(minX, minY, maxX, maxY)

	fromF64, toF64 := colorToF64(from), colorToF64(to)
	memo := r.GetColorF32()
	fromOklab := rgbToOklab([3]float64(fromF64[:3]))
	toOklab := rgbToOklab([3]float64(toF64[:3]))
	r.SetColorF32(float32(toOklab[0]), float32(toOklab[1]), float32(toOklab[2]), float32(toF64[3]))
	r.setFlatCustomVAs(float32(fromOklab[0]), float32(fromOklab[1]), float32(fromOklab[2]), float32(fromF64[3]))

	r.opts.Uniforms["Radius"] = [3]float32{fromRadius, transRadius, toRadius}
	r.opts.Uniforms["Origin"] = [2]float32{cx, cy}
	r.opts.Uniforms["NumSteps"] = numSteps
	r.opts.Uniforms["CurveFactor"] = curveFactor

	// draw shader
	ensureShaderGradientRadialLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderGradientRadial, &r.opts)
	clear(r.opts.Uniforms)
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

func (r *Renderer) OklabShiftChroma(target, source *ebiten.Image, x, y, chromaShift float32) {
	ensureShaderOklabShiftChromaLoaded()
	r.setFlatCustomVA0(chromaShift)
	r.DrawShaderAt(target, source, x, y, 0, 0, shaderOklabShiftChroma)
}

// ColorMix draws 'base' and 'over' to 'target' using the mix() function
// for color mixing instead of BlendSourceOver or other standard composition
// operations. This is useful to interpolate color transitions or other image
// changes when the images have translucent areas.
//
// The bounds of 'base' and 'over' must match.
func (r *Renderer) ColorMix(target, base, over *ebiten.Image, x, y int, alpha, mixLevel float32) {
	srcBounds := base.Bounds()
	if !srcBounds.Eq(over.Bounds()) {
		panic("'base' and 'over' bounds must match")
	}

	srcWidth, srcHeight := srcBounds.Dx(), srcBounds.Dy()
	srcWidthF32, srcHeightF32 := float32(srcWidth), float32(srcHeight)
	dstBounds := target.Bounds()
	minX := float32(dstBounds.Min.X) + float32(x)
	minY := float32(dstBounds.Min.Y) + float32(y)
	r.setDstRectCoords(minX, minY, minX+srcWidthF32, minY+srcHeightF32)
	minX = float32(srcBounds.Min.X)
	minY = float32(srcBounds.Min.Y)
	r.setSrcRectCoords(minX, minY, minX+srcWidthF32, minY+srcHeightF32)

	r.setFlatCustomVAs01(alpha, mixLevel)
	r.opts.Images[0] = base
	r.opts.Images[1] = over
	ensureShaderColorMixLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderColorMix, &r.opts)
	r.opts.Images[0] = nil
	r.opts.Images[1] = nil
}

// AlphaMask draws 'source' over 'target', but using 'mask' as an alpha mask.
func (r *Renderer) AlphaMask(target, source, mask *ebiten.Image, x, y, xMask, yMask float32) {
	srcBounds := source.Bounds()
	srcWidth, srcHeight := srcBounds.Dx(), srcBounds.Dy()
	srcWidthF32, srcHeightF32 := float32(srcWidth), float32(srcHeight)
	dstBounds := target.Bounds()
	minX := float32(dstBounds.Min.X) + x
	minY := float32(dstBounds.Min.Y) + y
	r.setDstRectCoords(minX, minY, minX+srcWidthF32, minY+srcHeightF32)
	minX = float32(srcBounds.Min.X)
	minY = float32(srcBounds.Min.Y)
	r.setSrcRectCoords(minX, minY, minX+srcWidthF32, minY+srcHeightF32)

	r.setFlatCustomVAs01(x-xMask, y-yMask)
	r.opts.Images[0] = source
	r.opts.Images[1] = mask
	ensureShaderAlphaMaskLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderAlphaMask, &r.opts)
	r.opts.Images[0] = nil
	r.opts.Images[1] = nil
}

// AlphaHorzFade draws 'source' over 'target' but with an horizontal alpha fade between
// the given points.
func (r *Renderer) AlphaHorzFade(target, source *ebiten.Image, x, y, inX, outX float32) {
	ensureShaderAlphaHorzFadeLoaded()
	r.setFlatCustomVAs01(inX, outX)
	r.DrawShaderAt(target, source, x, y, 0, 0, shaderAlphaHorzFade)
}

var DitherBayes [16]float32 = [16]float32{
	0.0 / 16.0, 12.0 / 16.0, 3.0 / 16.0, 15.0 / 16.0,
	8.0 / 16.0, 4.0 / 16.0, 11.0 / 16.0, 7.0 / 16.0,
	2.0 / 16.0, 14.0 / 16.0, 1.0 / 16.0, 13.0 / 16.0,
	10.0 / 16.0, 6.0 / 16.0, 9.0 / 16.0, 5.0 / 16.0,
}
var DitherDots [16]float32 = [16]float32{
	12.0 / 16.0, 4.0 / 16.0, 11.0 / 16.0, 15.0 / 16.0,
	5.0 / 16.0, 0.0 / 16.0, 3.0 / 16.0, 10.0 / 16.0,
	6.0 / 16.0, 1.0 / 16.0, 2.0 / 16.0, 9.0 / 16.0,
	13.0 / 16.0, 7.0 / 16.0, 8.0 / 16.0, 14.0 / 16.0,
}
var DitherSerp [16]float32 = [16]float32{
	0.0 / 16.0, 12.0 / 16.0, 13.0 / 16.0, 1.0 / 16.0,
	3.0 / 16.0, 7.0 / 16.0, 6.0 / 16.0, 2.0 / 16.0,
	4.0 / 16.0, 8.0 / 16.0, 9.0 / 16.0, 5.0 / 16.0,
	11.0 / 16.0, 15.0 / 16.0, 14.0 / 16.0, 10.0 / 16.0,
}
var DitherGlitch [16]float32 = [16]float32{
	0.0 / 16.0, 1.0 / 16.0, 2.0 / 16.0, 3.0 / 16.0,
	4.0 / 16.0, 5.0 / 16.0, 6.0 / 16.0, 7.0 / 16.0,
	8.0 / 16.0, 9.0 / 16.0, 10.0 / 16.0, 11.0 / 16.0,
	12.0 / 16.0, 13.0 / 16.0, 14.0 / 16.0, 15.0 / 16.0,
}

var DitherCrumbs [16]float32 = [16]float32{
	0.0 / 16.0, 4.0 / 16.0, 8.0 / 16.0, 1.0 / 16.0,
	11.0 / 16.0, 14.0 / 16.0, 12.0 / 16.0, 5.0 / 16.0,
	7.0 / 16.0, 13.0 / 16.0, 15.0 / 16.0, 9.0 / 16.0,
	3.0 / 16.0, 10.0 / 16.0, 6.0 / 16.0, 2.0 / 16.0,
}

var DitherBW []float32 = []float32{
	0.0, 0.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0,
}
var DitherBW4 []float32 = []float32{
	0.0, 0.0, 0.0, 1.0,
	0.333, 0.333, 0.333, 1.0,
	0.666, 0.666, 0.666, 1.0,
	1.0, 1.0, 1.0, 1.0,
}

var DitherAlpha8 []float32 = []float32{
	0.0, 0.0, 0.0, 0.0,
	1.0 / 7.0, 1.0 / 7.0, 1.0 / 7.0, 1.0 / 7.0,
	2.0 / 7.0, 2.0 / 7.0, 2.0 / 7.0, 2.0 / 7.0,
	3.0 / 7.0, 3.0 / 7.0, 3.0 / 7.0, 3.0 / 7.0,
	4.0 / 7.0, 4.0 / 7.0, 4.0 / 7.0, 4.0 / 7.0,
	5.0 / 7.0, 5.0 / 7.0, 5.0 / 7.0, 5.0 / 7.0,
	6.0 / 7.0, 6.0 / 7.0, 6.0 / 7.0, 6.0 / 7.0,
	1.0, 1.0, 1.0, 1.0,
}
var DitherBRG []float32 = []float32{
	0.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
}

// DitherMat4 draws the given mask to the target applying a static 4x4 dithering pattern to select colors from
// rgbaColors. The rgbaColors argument can contain up to 8 colors, flattened as RGBA quadruplets in [0...1] range.
// You can test with DitherBW4. The ditherMatrix argument is a 4x4 dithering matrix in column major order (like
// GLSL), where the values indicate the thresholds of the pattern in 0...1 range. You can test with [DitherBayes].
func (r *Renderer) DitherMat4(target, mask *ebiten.Image, ox, oy float32, xOffset, yOffset int, rgbaColors []float32, ditherMatrix [16]float32, rendererClrMix, maskColorMix float32) {
	if len(rgbaColors)%4 != 0 {
		panic("rgbaColors must have length multiple of 4")
	}
	numColors := len(rgbaColors) / 4
	if numColors > 8 {
		panic("DitherMat4 currently only supports up to 8 colors")
	} else if numColors <= 1 {
		panic("DitherMat4 expects at least 2 colors (as 8 float32 values)")
	}
	var palette [4 * 8]float32
	copy(palette[:], rgbaColors)

	r.setFlatCustomVAs(float32(xOffset), float32(yOffset), rendererClrMix, maskColorMix)
	r.opts.Uniforms["Matrix"] = ditherMatrix
	r.opts.Uniforms["NumColors"] = numColors
	r.opts.Uniforms["Colors"] = palette
	ensureShaderDitherMat4Loaded()
	r.DrawShaderAt(target, mask, ox, oy, 0, 0, shaderDitherMat4)
}

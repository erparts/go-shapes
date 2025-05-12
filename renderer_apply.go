package shapes

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Precondition: thickness can't exceed 32.
func (r *Renderer) ApplyExpansion(target *ebiten.Image, mask *ebiten.Image, ox, oy, thickness float32) {
	if thickness > 32 {
		panic("thickness can't exceed 32")
	}

	srcBounds := mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	ht32 := thickness / 2.0
	dstBounds := target.Bounds()
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	minX, minY := dstMinX+ox-ht32, dstMinY+oy-ht32
	maxX, maxY := dstMinX+ox+srcWidth+ht32, dstMinY+oy+srcHeight+ht32
	r.setDstRectCoords(minX-1, minY-1, maxX+1, maxY+1)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX-ht32-1, srcMinY-ht32-1, srcMaxX+ht32+1, srcMaxY+ht32+1)
	r.setFlatCustomVA0(thickness)

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderExpansionLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderExpansion, &r.opts)
	r.opts.Images[0] = nil
}

// Precondition: thickness can't exceed 32.
func (r *Renderer) ApplyErosion(target *ebiten.Image, mask *ebiten.Image, ox, oy, thickness float32) {
	if thickness > 32 {
		panic("thickness can't exceed 32")
	}

	srcBounds := mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	dstBounds := target.Bounds()
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	minX, minY := dstMinX+ox, dstMinY+oy
	maxX, maxY := dstMinX+ox+srcWidth, dstMinY+oy+srcHeight
	r.setDstRectCoords(minX-1, minY-1, maxX+1.0, maxY+1.0)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX-1, srcMinY-1, srcMaxX+1.0, srcMaxY+1.0)
	r.setFlatCustomVA0(thickness)

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderErosionLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderErosion, &r.opts)
	r.opts.Images[0] = nil
}

// Precondition: thickness can't exceed 32.
func (r *Renderer) ApplyOutline(target *ebiten.Image, mask *ebiten.Image, ox, oy, thickness float32) {
	if thickness > 32 {
		panic("thickness can't exceed 32")
	}

	srcBounds := mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	ht32 := thickness / 2.0
	dstBounds := target.Bounds()
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	minX, minY := dstMinX+ox-ht32, dstMinY+oy-ht32
	maxX, maxY := dstMinX+ox+srcWidth+ht32, dstMinY+oy+srcHeight+ht32
	r.setDstRectCoords(minX-1, minY-1, maxX+1, maxY+1)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX-ht32-1, srcMinY-ht32-1, srcMaxX+ht32+1.0, srcMaxY+ht32+1.0)
	r.setFlatCustomVA0(thickness)

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderOutlineLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderOutline, &r.opts)
	r.opts.Images[0] = nil
}

// ApplyBlur applies a gaussian blur to the given mask and draws it onto the given target.
// colorMix = 0 will use the renderer's vertex colors; colorMix = 1 will use the original mask colors.
//
// For radiuses above 4, you typically will prefer using ApplyBlur2.
func (r *Renderer) ApplyBlur(target *ebiten.Image, mask *ebiten.Image, ox, oy, radius, colorMix float32) {
	if radius > 32 {
		panic("radius can't exceed 32")
	}

	srcBounds := mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	hr32 := radius / 2.0
	dstBounds := target.Bounds()
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	minX, minY := dstMinX+ox-hr32, dstMinY+oy-hr32
	maxX, maxY := dstMinX+ox+srcWidth+hr32, dstMinY+oy+srcHeight+hr32
	r.setDstRectCoords(minX-1, minY-1, maxX+1, maxY+1)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX-hr32-1, srcMinY-hr32-1, srcMaxX+hr32+1.0, srcMaxY+hr32+1.0)
	r.setFlatCustomVAs01(radius, colorMix)

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderBlurLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderBlur, &r.opts)
	r.opts.Images[0] = nil
}

// ApplyBlur2 is similar to ApplyBlur, but uses two 1D passes instead of a single 2D pass.
// This greatly reduces the amount of sampled pixels for the shader, and despite breaking
// batching, tends to be much more efficient than ApplyBlur.
func (r *Renderer) ApplyBlur2(target *ebiten.Image, mask *ebiten.Image, ox, oy, radius, colorMix float32) {
	if radius > 32 {
		panic("radius can't exceed 32")
	}

	srcBounds := mask.Bounds()
	w32, h32 := float32(srcBounds.Dx()), float32(srcBounds.Dy())+radius
	w, h := int(w32), int(math.Ceil(float64(h32)))
	tmp := r.getTemp(0, w, h)
	preBlend := r.opts.Blend
	r.opts.Blend = ebiten.BlendCopy
	hr32 := radius / 2.0
	r.ApplyVertBlur(tmp, mask, 0, hr32+1.0, radius, 1.0)
	r.opts.Blend = preBlend
	r.ApplyHorzBlur(target, tmp, ox, oy-hr32-1.0, radius, colorMix)
}

func (r *Renderer) ApplyVertBlur(target *ebiten.Image, mask *ebiten.Image, ox, oy, radius, colorMix float32) {
	if radius > 32 {
		panic("radius can't exceed 32")
	}

	srcBounds := mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	hr32 := radius / 2.0
	dstBounds := target.Bounds()
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	minX, minY := dstMinX+ox, dstMinY+oy-hr32
	maxX, maxY := dstMinX+ox+srcWidth, dstMinY+oy+srcHeight+hr32
	r.setDstRectCoords(minX, minY-1, maxX, maxY+1)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX, srcMinY-hr32-1, srcMaxX, srcMaxY+hr32+1.0)
	r.setFlatCustomVAs01(radius, colorMix)

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderVertBlurLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderVertBlur, &r.opts)
	r.opts.Images[0] = nil
}

func (r *Renderer) ApplyHorzBlur(target *ebiten.Image, mask *ebiten.Image, ox, oy, radius, colorMix float32) {
	if radius > 32 {
		panic("radius can't exceed 32")
	}

	srcBounds := mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	hr32 := radius / 2.0
	dstBounds := target.Bounds()
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	minX, minY := dstMinX+ox-hr32, dstMinY+oy
	maxX, maxY := dstMinX+ox+srcWidth+hr32, dstMinY+oy+srcHeight
	r.setDstRectCoords(minX-1, minY, maxX+1, maxY)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX-hr32-1, srcMinY, srcMaxX+hr32+1.0, srcMaxY)
	r.setFlatCustomVAs01(radius, colorMix)

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderHorzBlurLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderHorzBlur, &r.opts)
	r.opts.Images[0] = nil
}

type Clamping uint8

const (
	ClampNone   Clamping = 0b0000
	ClampTop    Clamping = 0b1000
	ClampBottom Clamping = 0b0100
	ClampLeft   Clamping = 0b0010
	ClampRight  Clamping = 0b0001

	ClampTopLeft     Clamping = ClampTop | ClampLeft
	ClampTopRight    Clamping = ClampTop | ClampRight
	ClampBottomLeft  Clamping = ClampBottom | ClampLeft
	ClampBottomRight Clamping = ClampBottom | ClampRight
)

func (r *Renderer) ApplyHardShadow(target *ebiten.Image, mask *ebiten.Image, ox, oy, xOffset, yOffset float32, clamping Clamping) {
	dstBounds, srcBounds := target.Bounds(), mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	leftOff, topOff := min(0, xOffset), min(0, yOffset)
	rightOff, bottomOff := max(0, xOffset), max(0, yOffset)
	minX, minY := dstMinX+ox+leftOff, dstMinY+oy+topOff
	maxX, maxY := dstMinX+srcWidth+ox+rightOff, dstMinY+srcHeight+oy+bottomOff
	r.setDstRectCoords(minX, minY, maxX, maxY)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX+leftOff, srcMinY+topOff, srcMaxX+rightOff, srcMaxY+bottomOff)
	r.setFlatCustomVAs(xOffset, yOffset, float32(clamping), 0)

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderHardShadowLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderHardShadow, &r.opts)
	r.opts.Images[0] = nil
}

func (r *Renderer) ApplyShadow(target *ebiten.Image, mask *ebiten.Image, ox, oy, xOffset, yOffset, radius float32, clamping Clamping) {
	if radius > 32 {
		panic("radius can't exceed 32")
	}

	dstBounds, srcBounds := target.Bounds(), mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	leftOff, topOff := min(0, xOffset), min(0, yOffset)
	rightOff, bottomOff := max(0, xOffset), max(0, yOffset)
	hr32 := radius / 2.0
	topHR32, bottomHR32, leftHR32, rightHR32 := hr32, hr32, hr32, hr32
	if clamping&ClampBottom != 0 {
		bottomHR32 = 0
	}
	if clamping&ClampTop != 0 {
		topHR32 = 0
	}
	if clamping&ClampLeft != 0 {
		leftHR32 = 0
	}
	if clamping&ClampRight != 0 {
		rightHR32 = 0
	}

	minX, minY := dstMinX+ox+leftOff-leftHR32, dstMinY+oy+topOff-topHR32
	maxX, maxY := dstMinX+srcWidth+ox+rightOff+rightHR32, dstMinY+srcHeight+oy+bottomOff+bottomHR32
	r.setDstRectCoords(minX, minY, maxX, maxY)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX+leftOff-leftHR32, srcMinY+topOff-topHR32, srcMaxX+rightOff+rightHR32, srcMaxY+bottomOff+bottomHR32)
	r.setFlatCustomVAs(xOffset, yOffset, radius, float32(clamping))

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderShadowLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderShadow, &r.opts)
	r.opts.Images[0] = nil
}

func (r *Renderer) ApplyZoomShadow(target *ebiten.Image, mask *ebiten.Image, ox, oy, xOffset, yOffset, zoom float32, clamping Clamping) {
	if zoom < 1.0 || zoom > 16.0 {
		panic("zoom must be in [1, 16] range")
	}

	dstBounds, srcBounds := target.Bounds(), mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	dstMinX, dstMinY := float32(dstBounds.Min.X), float32(dstBounds.Min.Y)
	leftOff, topOff := min(0, xOffset), min(0, yOffset)
	rightOff, bottomOff := max(0, xOffset), max(0, yOffset)
	var topCut, leftCut, bottomCut, rightCut float32
	topPad, leftPad := srcHeight*0.5*(zoom-1.0), srcWidth*0.5*(zoom-1.0)
	bottomPad, rightPad := topPad, leftPad
	if clamping&ClampBottom != 0 {
		bottomCut = bottomPad / zoom
		bottomPad = 0
	}
	if clamping&ClampTop != 0 {
		topCut = topPad / zoom
		topPad = 0
	}
	if clamping&ClampLeft != 0 {
		leftCut = leftPad / zoom
		leftPad = 0
	}
	if clamping&ClampRight != 0 {
		rightCut = rightPad / zoom
		rightPad = 0
	}

	minX, minY := dstMinX+ox+leftOff-leftPad, dstMinY+oy+topOff-topPad
	maxX, maxY := dstMinX+srcWidth+ox+rightOff+rightPad, dstMinY+srcHeight+oy+bottomOff+bottomPad
	r.setDstRectCoords(minX, minY, maxX, maxY)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX+leftOff/zoom-leftCut, srcMinY+topOff/zoom-topCut, srcMaxX+rightOff/zoom-rightCut, srcMaxY+bottomOff/zoom-bottomCut)
	r.setFlatCustomVAs(xOffset/zoom, yOffset/zoom, float32(clamping), zoom)

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderZoomShadowLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderZoomShadow, &r.opts)
	r.opts.Images[0] = nil
}

// ApplySimpleGlow draws the given mask into the target, at the given coordinates, with
// an glow effect added. The effect mix intensity is determined by the renderer's color
// alphas. For finer control, see also [Renderer.ApplyGlow]().
func (r *Renderer) ApplySimpleGlow(target *ebiten.Image, mask *ebiten.Image, ox, oy, radius float32) {
	r.ApplyGlow(target, mask, ox, oy, radius, radius, 0.4, 0.7, 1.0)
}

// ApplyGlow draws a horizontal glow effect for the given mask into the target, at the
// given coordinates. The effect mix intensity is determined by the renderer's color alphas.
//
// Regarding the advanced control parameters:
//   - threshStart and threshEnd indicate the start luminosity threshold at which the glow
//     effect kicks in and the point at which it's fully active. threshStart must be <=
//     threshEnd, and the values must be in [0, 1] range.
//   - colorMix controls the glow's color. If 0, the glow color will be determined fully
//     by the renderer's vertex colors. If 1, the glow color will be determined by the original
//     mask colors. Any values in between will lead to linear interpolation.
//
// Notice that this effect uses an internal offscreen and two passes, which means it will
// always break batching.
func (r *Renderer) ApplyGlow(target *ebiten.Image, mask *ebiten.Image, ox, oy, horzRadius, vertRadius, threshStart, threshEnd, colorMix float32) {
	if threshStart > threshEnd {
		panic("threshStart > threshEnd")
	}
	if horzRadius > 32 || vertRadius > 32 {
		panic("radius can't exceed 32")
	}

	srcBounds := mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())
	w32, h32 := float32(srcWidth), float32(srcHeight)+vertRadius
	w, h := int(w32), int(math.Ceil(float64(h32)))
	tmp := r.getTemp(0, w, h)

	hr32 := vertRadius / 2.0
	r.setDstRectCoords(0, 0, w32, h32+2)

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX, srcMinY-hr32-1, srcMaxX, srcMaxY+hr32+1.0)
	r.setFlatCustomVAs(vertRadius, threshStart, threshEnd, 1.0)

	// first pass (threshold + vertical blur)
	r.opts.Images[0] = mask
	preBlend := r.opts.Blend
	r.opts.Blend = ebiten.BlendCopy
	ensureShaderGlowFirstPassLoaded()
	tmp.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderGlowFirstPass, &r.opts)
	r.opts.Images[0] = nil

	// second pass
	r.opts.Blend = ebiten.BlendLighter
	r.ApplyHorzBlur(target, tmp, ox, oy-hr32-1.0, horzRadius, colorMix)
	r.opts.Blend = preBlend
}

// ApplyHorzGlow draws a horizontal glow effect for the given mask into the target, at the
// given coordinates. See [Renderer.ApplyGlow]() for additional documentation. Comparedto
// Renderer.ApplyGlow, this effect only applies the glow horizontally and it's much cheaper,
// requiring no offscreen and a single pass.
func (r *Renderer) ApplyHorzGlow(target *ebiten.Image, mask *ebiten.Image, ox, oy, horzRadius, threshStart, threshEnd, colorMix float32) {
	if threshStart > threshEnd {
		panic("threshStart > threshEnd")
	}
	if horzRadius > 32 {
		panic("radius can't exceed 32")
	}

	srcBounds := mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())

	hr32 := horzRadius / 2.0
	r.setDstRectCoords(ox-hr32-1.0, oy, ox+float32(srcWidth)+hr32+1.0, oy+float32(srcHeight))

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX-hr32-1, srcMinY, srcMaxX+hr32+1, srcMaxY)
	r.setFlatCustomVAs(horzRadius, threshStart, threshEnd, colorMix)

	r.opts.Images[0] = mask
	ensureShaderHorzGlowLoaded()
	preBlend := r.opts.Blend
	r.opts.Blend = ebiten.BlendLighter
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderHorzGlow, &r.opts)
	r.opts.Blend = preBlend
	r.opts.Images[0] = nil
}

// ApplyDarkHorzGlow is the "negative" version of [Renderer.ApplyHorzGlow](). Instead of
// using an additive blending effect around high luminosity areas, it uses multiplicative
// blending around dark areas.
//
// Notice that unlike regular glow effects, dark glows expects threshStart >= threshEnd.
func (r *Renderer) ApplyDarkHorzGlow(target *ebiten.Image, mask *ebiten.Image, ox, oy, horzRadius, threshStart, threshEnd, colorMix float32) {
	if threshStart < threshEnd {
		panic("threshStart < threshEnd")
	}
	if horzRadius > 32 {
		panic("radius can't exceed 32")
	}

	srcBounds := mask.Bounds()
	srcWidth, srcHeight := float32(srcBounds.Dx()), float32(srcBounds.Dy())

	hr32 := horzRadius / 2.0
	r.setDstRectCoords(ox-hr32-1.0, oy, ox+float32(srcWidth)+hr32+1.0, oy+float32(srcHeight))

	srcMinX, srcMinY := float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	srcMaxX, srcMaxY := float32(srcBounds.Max.X), float32(srcBounds.Max.Y)
	r.setSrcRectCoords(srcMinX-hr32-1, srcMinY, srcMaxX+hr32+1, srcMaxY)
	r.setFlatCustomVAs(horzRadius, threshStart, threshEnd, colorMix)

	r.opts.Images[0] = mask
	ensureShaderDarkHorzGlowLoaded()
	preBlend := r.opts.Blend
	r.opts.Blend = BlendMultiply
	//r.opts.Blend = BlendSubtract // also possible with a shader flag, but multiply feels more natural
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderDarkHorzGlow, &r.opts)
	r.opts.Blend = preBlend
	r.opts.Images[0] = nil
}

type GaussKern uint8

const (
	GaussKern3 GaussKern = iota
	GaussKern5
	GaussKern7
	GaussKern9
	GaussKern11
	GaussKern13
	GaussKern15
	GaussKern17
)

func (k GaussKern) Radius() int {
	ik := int(k)
	return 1 + ik + ik
}

func (k GaussKern) Size() int {
	ik := int(k)
	return 3 + ik + ik
}

var gaussKerns = [][9]float32{
	{0.5000, 0.2500, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000},
	{0.3750, 0.2500, 0.0625, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000},
	{0.3990, 0.2420, 0.0540, 0.0040, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000},
	{0.2220, 0.1460, 0.0370, 0.0050, 0.0002, 0.0000, 0.0000, 0.0000, 0.0000},
	{0.2005, 0.1644, 0.1038, 0.0506, 0.0176, 0.0045, 0.0000, 0.0000, 0.0000},
	{0.1995, 0.1690, 0.1172, 0.0650, 0.0285, 0.0093, 0.0024, 0.0000, 0.0000},
	{0.1966, 0.1712, 0.1268, 0.0803, 0.0434, 0.0192, 0.0064, 0.0017, 0.0000},
	{0.1920, 0.1716, 0.1346, 0.0937, 0.0571, 0.0291, 0.0126, 0.0044, 0.0013},
}

// ApplyBlurD4 is a less flexible form of blur, similar to [Renderer.ApplyBlur2](),
// that downscales the source x4 before applying a gaussian kernel. This blur
// implementation tends to be more efficient than ApplyBlur2 when it comes to less
// powerful hardware and large blur areas (it uses less memory and compute at the
// cost of more steps). When enough resources are available, ApplyBlur2 tends to be
// slightly more efficient than ApplyBlurD4 when it comes to medium-sized or small
// blurs.
func (r *Renderer) ApplyBlurD4(target *ebiten.Image, mask *ebiten.Image, ox, oy float32, kernel GaussKern, colorMix float32) {
	const downscaling = 4
	maskBounds := mask.Bounds()
	maskWidth, maskHeight := maskBounds.Dx(), maskBounds.Dy()
	maskW64, maskH64 := float64(maskWidth), float64(maskHeight)
	downW64, downH64 := maskW64/downscaling, maskH64/downscaling
	downImgWidth, downImgHeight := math.Ceil(downW64)+2, math.Ceil(downH64)+2
	down := r.getTemp(0, int(downImgWidth), int(downImgHeight))
	down.Clear()

	var opts ebiten.DrawImageOptions
	opts.Filter = ebiten.FilterLinear
	opts.GeoM.Scale(1.0/downscaling, 1.0/downscaling)
	opts.GeoM.Translate(1, 1)
	opts.Blend = ebiten.BlendCopy
	down.DrawImage(mask, &opts)

	// apply kern horz blur
	halfMargin := float64(kernel.Radius())
	margin := halfMargin * 2.0
	dblurW64, dblurH64 := downW64+margin, downH64+margin
	dblurImgWidth, dblurImgHeight := math.Ceil(dblurW64)+2, math.Ceil(dblurH64)+2
	dblurHorz := r.getTemp(1, int(dblurImgWidth), int(downImgHeight))

	r.setDstRectCoords(0, 0, float32(dblurW64)+2, float32(downH64)+2)
	r.setSrcRectCoords(float32(-halfMargin), float32(0), float32(downW64+halfMargin)+2, float32(downH64)+2)
	preBlend := r.opts.Blend
	r.opts.Blend = ebiten.BlendCopy
	r.setFlatCustomVA0(colorMix)
	r.opts.Images[0] = down
	r.opts.Uniforms["KernelLen"] = kernel.Size()
	r.opts.Uniforms["Kernel"] = gaussKerns[kernel]
	ensureShaderHorzBlurKernLoaded()
	dblurHorz.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderHorzBlurKern, &r.opts)

	// apply kern vert blur
	dblur := r.getTemp(0, int(dblurImgWidth), int(dblurImgHeight))
	r.setDstRectCoords(0, 0, float32(dblurW64)+2, float32(dblurH64)+2)
	r.setSrcRectCoords(0, float32(-halfMargin), float32(dblurW64)+2, float32(downH64+halfMargin)+2)
	r.opts.Images[0] = dblurHorz
	ensureShaderVertBlurKernLoaded()
	dblur.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderVertBlurKern, &r.opts)

	// upscale
	r.opts.Blend = preBlend
	fx, fy := ox+float32(-downscaling-halfMargin*downscaling), oy+float32(-downscaling-halfMargin*downscaling)
	r.Upscale(target, dblur, fx, fy, downscaling, false)
}

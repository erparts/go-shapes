package shapes

import (
	"fmt"
	"image"
	"math"
	"math/bits"

	"github.com/hajimehoshi/ebiten/v2"
)

type JFMInitMode uint8

const (
	// JFMBoundary initialization mode places the seeds on the
	// boundary pixels of the image.
	JFMBoundary JFMInitMode = iota

	// JFMPixel initialization mode places the seeds on any pixel
	// within the specified alpha range.
	JFMPixel
)

// AAMargin is the standard antialias margin or soft edge value
// recommended for operations that accept it explicitly.
const AAMargin = 1.333

// JFMCompute computes a jumping flood map of the given source and stores it
// into jfmap. source and jfmap must have the same dimensions. jfmap doesn't
// need to be cleared before computation.
//
// A jumping flood map encodes offsets to nearest seeds, which are determined
// by the given [JFMInitMode], selected within the given [minAlpha, maxAlpha]
// range (inclusive). Jumping flood maps can be used to speed up morphological
// operations on solid shapes like outlining, expansion and erosion. Internal
// encoding details are documented in jfm_pass.kage.
//
// The function panics if maxDistance <= 0, maxDistance > 32k, source size != jfmap
// size, minAlpha > maxAlpha, minAlpha < 0, maxAlpha > 1 or an invalid initMode is
// given.
//
// For minAlpha, if you want to make the value exclusive, a typically safe epsilon
// is 0.001. E.g., to select (0.5, 1.0], use minAlpha = 0.5 + 0.001, maxAlpha = 1.0.
//
// This function uses one internal offscreen (#0), and jfmap and source can be on
// the same internal atlas.
func (r *Renderer) JFMCompute(jfmap, source *ebiten.Image, initMode JFMInitMode, maxDistance int, minAlpha, maxAlpha float32) {
	// safety assertions
	if maxDistance <= 0 {
		panic("maxDistance <= 0")
	}
	if maxDistance > 32000 { // up to 32766 should be technically distinguishable
		panic("maxDistance > 32000")
	}
	if minAlpha < 0 {
		panic("minAlpha < 0")
	}
	if maxAlpha > 1 {
		panic("maxAlpha > 1")
	}
	if minAlpha > maxAlpha {
		panic("minAlpha > maxAlpha")
	}

	sbounds := source.Bounds()
	tbounds := jfmap.Bounds()
	sw, sh := sbounds.Dx(), sbounds.Dy()
	tw, th := tbounds.Dx(), tbounds.Dy()
	if sw != tw || sh != th {
		panic(fmt.Sprintf("source size != jfmap size (%dx%d != %dx%d)", sw, sh, tw, th))
	}

	ensureShaderJFMPassLoaded()
	var initShader *ebiten.Shader
	switch initMode {
	case JFMBoundary:
		ensureShaderJFMInitBoundaryLoaded()
		initShader = shaderJFMInitBoundary
	case JFMPixel:
		ensureShaderJFMInitFillLoaded()
		initShader = shaderJFMInitFill
	default:
		panic(initMode) // invalid JFMInitMode
	}
	r.setFlatCustomVAs01(minAlpha, maxAlpha)

	// init
	memoBlend := r.opts.Blend
	r.opts.Blend = ebiten.BlendCopy
	temp := r.getTemp(0, sw, sh)

	dstOX, dstOY := float32(tbounds.Min.X), float32(tbounds.Min.Y)
	w, h := float32(sw), float32(sh)
	mapCoords := [2][4]float32{{dstOX, dstOY, dstOX + w, dstOY + h}, {0, 0, w, h}}
	r.setDstRectCoords(0, 0, w, h) // dst is 'temp' for the initialization
	srcOX, srcOY := float32(sbounds.Min.X), float32(sbounds.Min.Y)
	r.setSrcRectCoords(srcOX, srcOY, srcOX+w, srcOY+h)
	r.opts.Images[0] = source
	temp.DrawTrianglesShader(r.vertices[:], r.indices[:], initShader, &r.opts)

	// we use 1+JFA, so the first pass uses jump size = 1
	r.setFlatCustomVAs01(1.0, float32(maxDistance))
	r.opts.Images[0] = temp
	r.setDstRectCoords(mapCoords[0][0], mapCoords[0][1], mapCoords[0][2], mapCoords[0][3])
	r.setSrcRectCoords(mapCoords[1][0], mapCoords[1][1], mapCoords[1][2], mapCoords[1][3])
	jfmap.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderJFMPass, &r.opts)

	// - main JFA loop -
	// jump size starts at the base power of 2 of the current number
	jumpSize := 1 << (15 - bits.LeadingZeros16(uint16(maxDistance)))
	maps := [2]*ebiten.Image{jfmap, temp}
	mapIndex := 1
	for jumpSize > 0 {
		r.setFlatCustomVA0(float32(jumpSize)) // set only jump size, maxDistance is already ok
		r.setDstRectCoords(mapCoords[mapIndex][0], mapCoords[mapIndex][1], mapCoords[mapIndex][2], mapCoords[mapIndex][3])
		newIndex := 1 - mapIndex
		r.setSrcRectCoords(mapCoords[newIndex][0], mapCoords[newIndex][1], mapCoords[newIndex][2], mapCoords[newIndex][3])
		r.opts.Images[0] = maps[newIndex]
		maps[mapIndex].DrawTrianglesShader(r.vertices[:], r.indices[:], shaderJFMPass, &r.opts)
		mapIndex = newIndex
		jumpSize /= 2
	}

	// copy to jfmap if last step was done on temp
	if mapIndex == 0 {
		var opts ebiten.DrawImageOptions
		opts.Blend = ebiten.BlendCopy
		jfmap.DrawImage(temp, &opts)
	}

	// cleanup
	r.opts.Blend = memoBlend
	r.opts.Images[0] = nil
}

// JFMHeat is a debug and utility method to draw a heatmap for jfmap into the given target,
// using 0 and maxDistance as reference distances for "hot" and "cold".
func (r *Renderer) JFMHeat(target, jfmap *ebiten.Image, ox, oy float32, maxDistance int) {
	ensureShaderJFMHeatLoaded()
	r.setFlatCustomVA0(float32(maxDistance))
	r.DrawShaderAt(target, jfmap, ox, oy, 0, 0, shaderJFMHeat)
}

// JFMComputeUnsafeTemp is a utility method that puts both source and a newly generated jumping flood map into
// the given internal offscreen, returning references to the new images. This uses [Renderer.JFMCompute]()
// internally, so the internal offscreen index can't be #0 and all the derived parameter conditions apply.
//
// For details on internal renderer offscreens, please see [Renderer.UnsafeTemp]().
func (r *Renderer) JFMComputeUnsafeTemp(offscreenIndex int, source *ebiten.Image, initMode JFMInitMode, maxDistance int, minAlpha, maxAlpha float32) (sourceTemp, jfmapTemp *ebiten.Image) {
	if offscreenIndex == 0 {
		panic("JFMComputeTemp expects an offscreenIndex > 0; #0 is already used by JFMCompute")
	}
	_, _, w, h := rectOriginSize(source.Bounds())
	ox, oy := 0, 0
	if h <= w {
		oy = h
	} else {
		ox = w
	}
	temp := r.UnsafeTemp(offscreenIndex, ox+w, oy+h)
	temp.Clear()
	var opts ebiten.DrawImageOptions
	opts.Blend = ebiten.BlendCopy
	temp.DrawImage(source, &opts)

	sourceTemp = temp.SubImage(image.Rect(0, 0, w, h)).(*ebiten.Image)
	jfmapTemp = temp.SubImage(image.Rect(ox, oy, ox+w, oy+h)).(*ebiten.Image)
	r.JFMCompute(jfmapTemp, sourceTemp, initMode, maxDistance, minAlpha, maxAlpha)
	return sourceTemp, jfmapTemp
}

// JFMExpand performs morphological expansion. Thickness must be in [0, 32k].
//
//   - jfmap can be nil, in which case it will be automatically generated for only this operation
//     using [JFMPixel] mode with [0.001, 1.0] alpha interval (all not fully transparent pixels
//     are seeds).
//   - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
//   - aaMargin is the antialias margin. [AAMargin] can be used for a reasonable default.
func (r *Renderer) JFMExpand(target, source, jfmap *ebiten.Image, ox, oy, thickness, aaMargin float32) {
	if thickness < 0 {
		panic("thickness < 0")
	}
	if thickness > 32000 {
		panic("thickness > 32k")
	}

	if jfmap == nil {
		jfmapMaxDist := max(int(math.Ceil(float64(thickness))), 1)
		source, jfmap = r.JFMComputeUnsafeTemp(1, source, JFMPixel, jfmapMaxDist, 0.001, 1.0)
	}

	ensureShaderJFMExpansionLoaded()
	r.opts.Images[1] = jfmap
	r.setFlatCustomVAs01(thickness, aaMargin)
	r.DrawShaderAt(target, source, ox, oy, 0, 0, shaderJFMExpansion)
	r.opts.Images[1] = nil
}

// TODO: unimplemented
//
// JFMErode performs morphological erosion.
//
//   - colorMix controls the outline color (0 = use vertex colors, 1 = use source colors)
//   - jfmap can be nil, in which case it will be automatically generated for only this operation
//     using [JFMPixel] mode with [0.0, 0.0] alpha interval (transparent pixels are seeds).
//   - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
//   - aaMargin is the antialias margin. [AAMargin] can be used for a reasonable default.
func (r *Renderer) JFMErode(target, source, jfmap *ebiten.Image, ox, oy, thickness, aaMargin, colorMix float32) {
	panic("unimplemented")
}

// TODO: unimplemented
//
// JFMOutline performs morphological outlining.
//
//   - colorMix controls the outline color (0 = use vertex colors, 1 = use source colors)
//   - jfmap can be nil, in which case it will be automatically generated for only this operation
//     using [JFMBoundary] mode.
//   - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
func (r *Renderer) JFMOutline(target, source, jfmap *ebiten.Image, ox, oy, inThickness, outThickness, inOpacity, colorMix float32) {
	panic("unimplemented")
}

// TODO: unimplemented
//
// JFMInsetContour is a specific effect designed mainly for text animations. It creates an
// internal outline, which includes the image borders where the target clips the source, while
// also allowing to control the inner fill opacity.
//
//   - colorMix controls the outline color (0 = use vertex colors, 1 = use source colors)
//   - jfmap can be nil, in which case it will be automatically generated for only this operation
//     using [JFMBoundary] mode.
//   - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
func (r *Renderer) JFMInsetContour(target, source, jfmap *ebiten.Image, ox, oy, inThickness, inOpacity, colorMix float32) {
	panic("unimplemented")
}

//func (r *Renderr) JFMFeather(target, source, jfmap *ebiten.Image, ox, oy, radius, curve float32) {}

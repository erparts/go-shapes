package shapes

import (
	"fmt"
	"image"
	"math/bits"

	"github.com/hajimehoshi/ebiten/v2"
)

type JFMInitMode uint8

const (
	JFMBoundary JFMInitMode = iota
	JFMInside
	JFMOutside
)

// JFMCompute computes a jumping flood map of the given source and stores it
// into jfmap. source and jfmap must have the same dimensions. jfmap doesn't
// need to be cleared before computation.
//
// A jumping flood map encodes offsets to nearest seeds, which are determined
// by the given [JFMInitMode]. These maps can be used to speed up morphological
// operations like outlining, expansion and erosion. Internal encoding details
// are documented in jfm_pass.kage
//
// The maxDistance can't exceed 32k due to Ebitengine texture format limitations.
//
// The function panics if maxDistance <= 0, maxDistance > 32k, source size != jfmap
// size or an invalid initMode is given.
//
// This function uses one internal offscreen (#0), and jfmap and source can be on
// the same internal atlas.
func (r *Renderer) JFMCompute(jfmap, source *ebiten.Image, initMode JFMInitMode, maxDistance int) {
	// safety assertions
	if maxDistance <= 0 {
		panic("maxDistance <= 0")
	}
	if maxDistance > 32000 { // 32511 technically
		panic("maxDistance > 32000")
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
	case JFMInside:
		ensureShaderJFMInitFillLoaded()
		initShader = shaderJFMInitFill
		r.setFlatCustomVAs01(0.001, 1.0)
	case JFMOutside:
		ensureShaderJFMInitFillLoaded()
		initShader = shaderJFMInitFill
		r.setFlatCustomVAs01(0.0, 0.001)
	default:
		panic(initMode) // invalid JFMInitMode
	}

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

// JFMComputeUnsafeTemp is a utility method that puts both source and a newly generated jumping flood map into
// the given internal offscreen, returning references to the new images. This uses [Renderer.JFMCompute]()
// internally, so the internal offscreen index can't be #0 and all the derived parameter conditions apply.
//
// For details on internal renderer offscreens, please see [Renderer.UnsafeTemp]().
func (r *Renderer) JFMComputeUnsafeTemp(offscreenIndex int, source *ebiten.Image, initMode JFMInitMode, maxDistance int) (sourceTemp, jfmapTemp *ebiten.Image) {
	if offscreenIndex == 0 {
		panic("JFMComputeTemp expects an offscreenIndex > 0, as 0 is already used by JFMCompute")
	}
	_, _, w, h := rectOriginSize(source.Bounds())
	ox, oy := 0, 0
	if h <= w {
		oy = h
	} else {
		ox = w
	}
	temp := r.UnsafeTemp(offscreenIndex, ox+w, oy+h)
	var opts ebiten.DrawImageOptions
	opts.Blend = ebiten.BlendCopy
	temp.DrawImage(source, &opts)

	sourceTemp = temp.SubImage(image.Rect(ox, oy, ox+w, oy+h)).(*ebiten.Image)
	jfmapTemp = temp.SubImage(image.Rect(ox, oy, ox+w, oy+h)).(*ebiten.Image)
	r.JFMCompute(jfmapTemp, sourceTemp, initMode, maxDistance)
	return sourceTemp, jfmapTemp
}

// TODO: unimplemented
//
// JFMExpand performs morphological expansion.
//
//   - colorMix controls the outline color (0 = use vertex colors, 1 = use source colors)
//   - jfmap can be nil, in which case it will be automatically generated for only this operation
//     using [JFMInitInside] mode.
//   - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
func (r *Renderer) JFMExpand(target, source, jfmap *ebiten.Image, ox, oy, thickness, colorMix float32) {
	panic("unimplemented")
}

// TODO: unimplemented
//
// JFMErode performs morphological erosion.
//
//   - colorMix controls the outline color (0 = use vertex colors, 1 = use source colors)
//   - jfmap can be nil, in which case it will be automatically generated for only this operation
//     using [JFMInitInside]
//   - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
func (r *Renderer) JFMErode(target, source, jfmap *ebiten.Image, ox, oy, thickness, colorMix float32) {
	panic("unimplemented")
}

// TODO: unimplemented
//
// JFMOutline performs morphological outlining.
//
//   - colorMix controls the outline color (0 = use vertex colors, 1 = use source colors)
//   - jfmap can be nil, in which case it will be automatically generated for only this operation
//     using [JFMInitBoundary] mode.
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
//     using [JFMInitBoundary] mode.
//   - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
func (r *Renderer) JFMInsetContour(target, source, jfmap *ebiten.Image, ox, oy, inThickness, inOpacity, colorMix float32) {
	panic("unimplemented")
}

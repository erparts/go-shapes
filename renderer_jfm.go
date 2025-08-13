package shapes

import (
	"fmt"
	"math/bits"

	"github.com/hajimehoshi/ebiten/v2"
)

type JFMInitMode uint8

const (
	JFMBoundary JFMInitMode = iota
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
// The maxDistance can't exceed 380 due to Ebitengine texture format limitations.
//
// The function panics if maxDistance <= 0, maxDistance > 380, source size != jfmap
// size or an invalid initMode is given.
//
// This function uses one internal offscreen (#0), and jfmap and source can be on
// the same internal atlas.
func (r *Renderer) JFMCompute(jfmap, source *ebiten.Image, initMode JFMInitMode, maxDistance int) {
	// safety assertions
	if maxDistance <= 0 {
		panic("maxDistance <= 0")
	}
	if maxDistance > 380 {
		panic("maxDistance > 380")
	}
	sbounds := source.Bounds()
	tbounds := jfmap.Bounds()
	sw, sh := sbounds.Dx(), sbounds.Dy()
	tw, th := tbounds.Dx(), tbounds.Dy()
	if sw != tw || sh != th {
		panic(fmt.Sprintf("source size != jfmap size (%dx%d != %dx%d)", sw, sh, tw, th))
	}

	ensureShaderJFMInitBoundaryLoaded()
	ensureShaderJFMPassLoaded()
	var initShader *ebiten.Shader
	switch initMode {
	case JFMBoundary:
		initShader = shaderJFMInitBoundary
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

// TODO: unimplemented
//
// JFMExpand performs morphological expansion.
//
// - colorMix controls the outline color (0 = use vertex colors, 1 = use source colors)
// - jfmap can be nil, in which case it will be automatically generated for only this operation
// - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
func (r *Renderer) JFMExpand(target, source, jfmap *ebiten.Image, ox, oy, thickness, colorMix float32) {
	panic("unimplemented")
}

// TODO: unimplemented
//
// JFMErode performs morphological erosion.
//
// - colorMix controls the outline color (0 = use vertex colors, 1 = use source colors)
// - jfmap can be nil, in which case it will be automatically generated for only this operation
// - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
func (r *Renderer) JFMErode(target, source, jfmap *ebiten.Image, ox, oy, thickness, colorMix float32) {
	panic("unimplemented")
}

// TODO: unimplemented
//
// JFMOutline performs morphological outlining.
//
// - colorMix controls the outline color (0 = use vertex colors, 1 = use source colors)
// - jfmap can be nil, in which case it will be automatically generated for only this operation
// - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
func (r *Renderer) JFMOutline(target, source, jfmap *ebiten.Image, ox, oy, inThickness, outThickness, colorMix float32) {
	panic("unimplemented")
}

// TODO: unimplemented
//
// JFMInsetContour is a specific effect designed mainly for text animations. It creates an
// internal outline, which includes the image borders where the target clips the source, while
// also allowing to control the inner fill opacity.
//
// - colorMix controls the outline color (0 = use vertex colors, 1 = use source colors)
// - jfmap can be nil, in which case it will be automatically generated for only this operation
// - source and jfmap should be in the same atlas to avoid automatic atlasing issues.
func (r *Renderer) JFMInsetContour(target, source, jfmap *ebiten.Image, ox, oy, inThickness, inOpacity, colorMix float32) {
	panic("unimplemented")
}

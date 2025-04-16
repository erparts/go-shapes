package shapes

import (
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
	r.setFlatCustomVAs(thickness, 0, 0, 0)

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
	r.setFlatCustomVAs(thickness, 0, 0, 0)

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
	r.setFlatCustomVAs(thickness, 0, 0, 0)

	// draw shader
	r.opts.Images[0] = mask
	ensureShaderOutlineLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderOutline, &r.opts)
	r.opts.Images[0] = nil
}

// colorMix = 0 will use the renderer's vertex colors; colorMix = 1 will use the original mask colors.
func (r *Renderer) ApplyBlur(target *ebiten.Image, mask *ebiten.Image, ox, oy, radius, colorMix float32) {
	// ...
}

func (r *Renderer) ApplyHardShadow(target *ebiten.Image, mask *ebiten.Image, ox, oy, thickness float32) {
	// ...
}

func (r *Renderer) ApplyGlow(target *ebiten.Image, mask *ebiten.Image, radius, threshold, colorMix float32) {
	// ...
}

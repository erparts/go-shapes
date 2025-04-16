package shapes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// TileDotsHex draws dots of the given radius distributed in a hexagonal
// lattice. HorzSpacing should always be at least twice the radius.
// TODO: offsets x/y
func (r *Renderer) TileDotsHex(target *ebiten.Image, radius, horzSpacing, xOffset, yOffset float32) {
	bounds := target.Bounds()
	minX, minY := float32(bounds.Min.X), float32(bounds.Min.Y)
	maxX, maxY := float32(bounds.Max.X), float32(bounds.Max.Y)
	r.setDstRectCoords(minX, minY, maxX, maxY)
	r.setFlatCustomVAs(radius, horzSpacing, xOffset, yOffset)

	// draw shader
	ensureShaderPatternDotsLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderPatternDots, &r.opts)
}

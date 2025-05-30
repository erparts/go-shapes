package shapes

import "github.com/hajimehoshi/ebiten/v2"

// WarpBarrel draws the given image with a simple, CRT-like barrel warp.
// Intensity should be in ~[0.2, 1.5], with 0.5 being a good starting value
// to play with.
//
// The size of the output image will always be equal or smaller than the
// input source, as the corner vertices are warped towards the interior.
func (r *Renderer) WarpBarrel(target, source *ebiten.Image, ox, oy float32, horzWarp, vertWarp float32) {
	if horzWarp < 0 || vertWarp < 0 {
		panic("horzWarp < 0 || vertWarp < 0")
	}
	r.setFlatCustomVAs01(horzWarp, vertWarp)
	ensureShaderWarpBarrelLoaded()
	r.DrawShaderAt(target, source, ox, oy, 0, 0, shaderWarpBarrel)
}

package shapes

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestTileDotsHex ./... -count 1
func TestTileDotsHex(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)
		ctx.Renderer.SetColor(color.RGBA{0, 192, 0, 255})
		ctx.Renderer.TileDotsHex(canvas, 8, 24, 0, 0)
		ctx.Renderer.SetColor(color.RGBA{192, 0, 192, 255})
		xOffset := float32(ctx.DistAnim(6, 1.0))
		yOffset := float32(ctx.DistAnim(6, 0.5))
		ctx.Renderer.TileDotsHex(canvas, 4, 12, xOffset, yOffset)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestTileDotsGrid ./... -count 1
func TestTileDotsGrid(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)
		ctx.Renderer.SetColor(color.RGBA{220, 64, 32, 255})
		ctx.Renderer.TileDotsGrid(canvas, 8, 24, 0, 0)

		ctx.Renderer.SetColor(color.RGBA{32, 64, 250, 255})
		ctx.Renderer.TileDotsGrid(canvas, 5, 24, 12, 12)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

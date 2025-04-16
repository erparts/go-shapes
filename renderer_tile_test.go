package shapes

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestTileDots ./... -count 1
func TestTileDots(t *testing.T) {
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

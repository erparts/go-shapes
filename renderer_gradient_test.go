package shapes

import (
	"image/color"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestGradient ./... -count 1

func TestGradient(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		var opts ebiten.DrawImageOptions
		lx, ly := ctx.LeftClickF64()
		opts.GeoM.Translate(lx, ly)
		canvas.DrawImage(ctx.Images[0], &opts)

		rx, ry := ctx.RightClickF64()
		opts.GeoM.Reset()
		opts.GeoM.Translate(rx, ry)
		canvas.DrawImage(ctx.Images[1], &opts)
	})

	rectA := app.Renderer.NewRect(120, 80)
	circ := app.Renderer.NewCircle(64.0)
	rectB := ebiten.NewImageWithOptions(circ.Bounds(), nil)
	app.Renderer.Gradient(rectA, nil, 0, 0, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 255, 0, 255}, 4, math.Pi/7, 1.0)
	app.Renderer.Gradient(rectB, circ, 0, 0, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 255, 0, 255}, -1, math.Pi, 0.2)
	app.Images = append(app.Images, rectA, rectB)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

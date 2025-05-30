package shapes

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestWarpBarrel . -count 1
func TestWarpBarrel(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)
		d := float32(ctx.DistAnim(1.5, 0.5))
		ctx.Renderer.WarpBarrel(canvas, ctx.Images[0], 64, 32, d/2.0, d)
	})
	w, h := 640/2, 480/2
	img := ebiten.NewImage(w, h)
	app.Renderer.SetColor(color.RGBA{255, 0, 0, 255}, 0)
	app.Renderer.SetColor(color.RGBA{255, 255, 0, 255}, 1)
	app.Renderer.SetColor(color.RGBA{0, 255, 0, 255}, 2)
	app.Renderer.SetColor(color.RGBA{0, 255, 255, 255}, 3)
	app.Renderer.DrawIntArea(img, 0, 0, w, h)

	app.Images = append(app.Images, img)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

package shapes

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestApplyExpansion ./... -count 1
// go test -run ^TestApplyErosion ./... -count 1
// go test -run ^TestApplyOutline ./... -count 1

func TestApplyExpansion(t *testing.T) {
	radius := float32(64.0)
	expansion := float32(16.0)
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{0, 255, 0, 255})
		ctx.Renderer.DrawCircle(canvas, lx, ly, radius+expansion/2.0+1.0)
		ctx.Renderer.SetColor(color.RGBA{255, 0, 0, 255})
		ctx.Renderer.ApplyExpansion(canvas, ctx.Images[0], lx-radius, ly-radius, expansion)
	})
	app.Images = append(app.Images, app.Renderer.NewCircle(float64(radius)))
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestApplyErosion(t *testing.T) {
	radius := float32(64.0)
	erosion := float32(16.0)
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{255, 255, 255, 255})
		ctx.Renderer.DrawCircle(canvas, lx, ly, radius)

		ctx.Renderer.SetColor(color.RGBA{255, 0, 0, 255})
		ctx.Renderer.DrawCircle(canvas, lx, ly, radius-erosion/2+1.0)

		ctx.Renderer.SetColor(color.RGBA{0, 0, 255, 255})
		ctx.Renderer.ApplyErosion(canvas, ctx.Images[0], lx-radius, ly-radius, erosion)
	})
	app.Images = append(app.Images, app.Renderer.NewCircle(float64(radius)))
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestApplyOutline(t *testing.T) {
	radius := float32(64.0)
	thick := float32(8.0)
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{0, 255, 0, 255})
		ctx.Renderer.DrawCircle(canvas, lx, ly, radius+thick/2+1.0)
		ctx.Renderer.SetColor(color.RGBA{255, 0, 0, 255})
		ctx.Renderer.DrawCircle(canvas, lx, ly, radius-thick/2-1.0)

		ctx.Renderer.SetColor(color.RGBA{0, 0, 255, 255})
		ctx.Renderer.ApplyOutline(canvas, ctx.Images[0], lx-radius, ly-radius, thick)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.SetColor(color.RGBA{255, 255, 255, 255})
		ctx.Renderer.ApplyOutline(canvas, ctx.Images[0], rx-radius, ry-radius, thick)
	})
	app.Images = append(app.Images, app.Renderer.NewCircle(float64(radius)))
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

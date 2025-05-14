package shapes

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestDrawShapes ./... -count 1
// go test -run ^TestDrawEllipse ./... -count 1
// ...

func TestDrawShapes(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)
		ctx.Renderer.SetColor(color.RGBA{255, 255, 255, 255})
		lx, ly := ctx.LeftClickF64()
		rx, ry := ctx.RightClickF64()
		ctx.Renderer.DrawLine(canvas, lx, ly, rx, ry, 6.0)
		ctx.Renderer.DrawCircle(canvas, 540, 80, 60)
		x, y := float64(160), float64(40)
		ctx.Renderer.DrawTriangle(canvas, x, y, x+30, y+10, x+16, y+50, 0)
		x, y = float64(80), float64(260)

		ctx.Renderer.SetColor(color.RGBA{240, 48, 48, 255})
		ctx.Renderer.DrawTriangle(canvas, x, y, x+70, y-20, x+114, y+80, 0)
		ctx.Renderer.SetColor(color.RGBA{255, 255, 255, 255})
		ctx.Renderer.DrawTriangle(canvas, x, y, x+70, y-20, x+114, y+80, 8)
		x, y = float64(200), float64(300)
		ctx.Renderer.StrokeTriangle(canvas, x+70, y-20, x, y, x+114, y+80, 4, 8)
		v := uint8(32 + ctx.DistAnim(196-32, 1.0))
		ctx.Renderer.SetColor(color.RGBA{v, 0, 0, v})
		ctx.Renderer.StrokeTriangle(canvas, x+70, y-20, x, y, x+114, y+80, -4, 0)

		rads := ctx.RadsAnim(1.0)
		ctx.Renderer.DrawHexagon(canvas, 80, 400, 60, 0, float32(rads))
		ctx.Renderer.DrawHexagon(canvas, 80, 400, 60, 0, float32(rads))
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestDrawArea(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		lx, ly := ctx.LeftClickF32()
		rx, ry := ctx.RightClickF32()

		ctx.Renderer.SetColorF32(1.0, 1.0, 1.0, 1.0)
		w1, h1 := float32(128), float32(48)
		w2, h2 := float32(48), float32(128)
		ctx.Renderer.DrawArea(canvas, lx-w1/2, ly-h1/2, w1, h1, float32(ctx.DistAnim(float64(min(w1, h1))/2.0, 1.0)))
		ctx.Renderer.DrawArea(canvas, rx-w2/2, ry-h2/2, w2, h2, 0)

		ctx.Renderer.SetColorF32(0.2, 0.0, 0.2, 0.2)
		ctx.Renderer.DrawCircle(canvas, lx, ly, max(w1, h1)/2.0)
		ctx.Renderer.DrawCircle(canvas, rx, ry, max(w2, h2)/2.0)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestStrokeIntArea(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		lx, ly := ctx.LeftClick.X, ctx.LeftClick.Y
		ctx.Renderer.SetColor(color.RGBA{255, 255, 255, 255})
		ctx.Renderer.DrawIntArea(canvas, lx, ly, 200, 50)
		ctx.Renderer.SetColor(color.RGBA{0, 255, 0, 255})
		ctx.Renderer.StrokeIntArea(canvas, lx-1, ly-1, 200+2, 50+2, 1, 0)

		ctx.Renderer.SetColor(color.RGBA{0, 128, 0, 128})
		ctx.Renderer.StrokeIntArea(canvas, lx, ly, 200, 50, 0, 1)

		rx, ry := ctx.RightClick.X, ctx.RightClick.Y
		ctx.Renderer.SetColor(color.RGBA{240, 0, 240, 255}, 0, 2)
		ctx.Renderer.StrokeIntArea(canvas, rx, ry, 100, 50, 4, 4)

		ctx.Renderer.SetColor(color.RGBA{64, 128, 64, 128})
		ctx.Renderer.DrawIntArea(canvas, rx, ry, 100, 50)

		ctx.Renderer.SetColor(color.RGBA{255, 0, 0, 255}, 0)
		ctx.Renderer.SetColor(color.RGBA{0, 255, 0, 255}, 1)
		ctx.Renderer.SetColor(color.RGBA{0, 0, 255, 255}, 2)
		ctx.Renderer.SetColor(color.RGBA{0, 255, 255, 255}, 3)
		ctx.Renderer.StrokeIntArea(canvas, lx, ry, 80, 50, 8, 8)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestDrawEllipse(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		lx, ly := ctx.LeftClickF32()
		rx, ry := ctx.RightClickF32()

		ctx.Renderer.SetColorF32(0.5, 0.5, 0.5, 0.5)
		ctx.Renderer.DrawCircle(canvas, lx, ly, 64.0)
		ctx.Renderer.DrawCircle(canvas, rx, ry, 32.0)

		ctx.Renderer.SetColorF32(1.0, 1.0, 1.0, 1.0)
		ctx.Renderer.DrawEllipse(canvas, lx, ly, 24.0, 64.0, ctx.RadsAnim(1.0))
		ctx.Renderer.DrawEllipse(canvas, rx, ry, 32.0, 16.0, 0)

		ctx.Renderer.SetColorF32(0.0, 0.5, 0.5, 0.5)
		ctx.Renderer.DrawCircle(canvas, lx, ly, 24.0)
		ctx.Renderer.DrawCircle(canvas, rx, ry, 16.0)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestDrawRing(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		lx, ly := ctx.LeftClickF32()
		rx, ry := ctx.RightClickF32()

		ctx.Renderer.SetColorF32(1.0, 1.0, 1.0, 1.0)
		ctx.Renderer.DrawCircle(canvas, lx, ly, 64.0)
		ctx.Renderer.DrawCircle(canvas, rx, ry, 48.0)

		ctx.Renderer.SetColorF32(1.0, 0.0, 1.0, 1.0, 0, 1)
		ctx.Renderer.SetColorF32(0.5, 1.0, 0.5, 1.0, 2, 3)
		ctx.Renderer.DrawRing(canvas, lx, ly, 65.0, 67.0)
		ctx.Renderer.DrawRing(canvas, rx, ry, 48.0-4, 48+0)

		ctx.Renderer.SetColorF32(0.5, 0.5, 0.5, 0.5)
		ctx.Renderer.DrawCircle(canvas, rx, ry, 48.0)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

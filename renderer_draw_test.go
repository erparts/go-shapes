package shapes

import (
	"image/color"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestDrawShapes . -count 1
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

// go test -run ^TestDrawArea . -count 1
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

// go test -run ^TestDrawIntArea . -count 1
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

// go test -run ^TestStrokeArea . -count 1
func TestStrokeArea(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{255, 255, 255, 255})
		ctx.Renderer.DrawArea(canvas, lx, ly, 200, 50, 16)
		ctx.Renderer.SetColor(color.RGBA{0, 255, 0, 255})
		ctx.Renderer.StrokeArea(canvas, lx, ly, 200, 50, 2, 0, 16)

		ctx.Renderer.SetColor(color.RGBA{128, 0, 0, 128})
		ctx.Renderer.StrokeArea(canvas, lx, ly, 200, 50, 0, 2, 16)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.SetColor(color.RGBA{240, 0, 240, 255}, 0, 2)
		ctx.Renderer.StrokeArea(canvas, rx, ry, 100, 50, 4, 4, 25)

		a := uint8(ctx.DistAnim(144.0, 1.0))
		ctx.Renderer.SetColor(color.RGBA{a, a, a, a})
		ctx.Renderer.DrawIntArea(canvas, int(rx), int(ry), 100, 50)

		ctx.Renderer.SetColor(color.RGBA{255, 0, 0, 255}, 0)
		ctx.Renderer.SetColor(color.RGBA{0, 255, 0, 255}, 1)
		ctx.Renderer.SetColor(color.RGBA{0, 0, 255, 255}, 2)
		ctx.Renderer.SetColor(color.RGBA{0, 255, 255, 255}, 3)
		extra := float32(ctx.DistAnim(16, 1.0))
		subRounding := float32(ctx.DistAnim(20, 1.0))
		ctx.Renderer.StrokeArea(canvas, lx, ry, 80+extra, 50, 8, 8, 25-subRounding)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestDrawEllipse . -count 1
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

// go test -run ^TestDrawRing . -count 1
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

// go test -run ^TestDrawRingSector$ . -count 1
func TestDrawRingSector(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		w, h := rectSizeF32(canvas.Bounds())
		cx, cy := w/2.0, h/2.0
		startRads := ctx.ModAnim(2*math.Pi, 1.0)
		endRads := startRads + 0.4 + ctx.DistAnim(1.6, 1.0)
		ctx.Renderer.SetColorF32(1.0, 1.0, 1.0, 1.0)
		ctx.Renderer.DrawRingSector(canvas, cx, cy, 48, 128, startRads, endRads, 0.0)
		ctx.Renderer.SetColorF32(0.0, 0.5, 0.5, 0.5)
		ctx.Renderer.DrawRingSector(canvas, cx, cy, 48, 128, startRads, endRads, 8.0)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestDrawPie$ . -count 1
func TestDrawPie(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		w, h := rectSizeF32(canvas.Bounds())
		cx, cy := w/2.0, h/2.0
		rate := -0.01 + ctx.DistAnim(1.02, 1.0)

		ctx.Renderer.SetColorF32(1.0, 1.0, 1.0, 1.0)
		ctx.Renderer.DrawPieRate(canvas, cx, cy, 96.0, RadsRight, rate, 6.0)

		ctx.Renderer.SetColorF32(0.0, 1.0, 0.0, 1.0)
		ctx.Renderer.DrawPie(canvas, cx, cy, 64.0, RadsRight+rate, RadsBottom, 3.0)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestDrawQuad$ . -count 1
func TestDrawQuad(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		lx, ly := ctx.LeftClickF32()
		rx, ry := ctx.RightClickF32()

		w, h := rectSizeF32(canvas.Bounds())
		quad := [4]PointF32{
			{X: lx, Y: ly},
			{X: w/2.0 + w/4.0, Y: h/2.0 - h/4.0},
			{X: rx, Y: ry},
			{X: w/2.0 - w/4.0, Y: h/2.0 + h/4.0},
		}
		thickening := float32(ctx.DistAnim(48.0, 1.0))
		ctx.Renderer.DrawQuad(canvas, quad, thickening)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestDrawQuadSoft$ . -count 1
func TestDrawQuadSoft(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		lx, ly := ctx.LeftClickF32()
		rx, ry := ctx.RightClickF32()

		w, h := rectSizeF32(canvas.Bounds())
		quad := [4]PointF32{
			{X: lx, Y: ly},
			{X: w/2.0 + w/4.0, Y: h/2.0 - h/4.0},
			{X: rx, Y: ry},
			{X: w/2.0 - w/4.0, Y: h/2.0 + h/4.0},
		}
		thickening := float32(ctx.DistAnim(48.0, 1.0))
		ctx.Renderer.DrawQuadSoft(canvas, quad, thickening, 64.0)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

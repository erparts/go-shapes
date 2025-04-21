package shapes

import (
	"image"
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

func TestApplyBlur(t *testing.T) {
	radius := float32(64.0)
	fxRadius := float32(32.0)
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{0, 0, 255, 255})
		ctx.Renderer.ApplyBlur(canvas, ctx.Images[0], lx-radius, ly-radius, fxRadius, 1.0)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.SetColor(color.RGBA{255, 0, 0, 255})
		ctx.Renderer.ApplyBlur(canvas, ctx.Images[0], rx-radius, ry-radius, fxRadius, 0.0)
	})
	app.Images = append(app.Images, app.Renderer.NewCircle(float64(radius)))
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestApplyBlur2(t *testing.T) {
	radius := float32(64.0)
	fxRadius := float32(32.0)
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{0, 0, 255, 255})
		ctx.Renderer.ApplyBlur2(canvas, ctx.Images[0], lx-radius, ly-radius, fxRadius, 1.0)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.SetColor(color.RGBA{255, 0, 0, 255})
		ctx.Renderer.ApplyBlur2(canvas, ctx.Images[0], rx-radius, ry-radius, fxRadius, 0.0)
	})
	app.Images = append(app.Images, app.Renderer.NewCircle(float64(radius)))
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestApplyDirBlur(t *testing.T) {
	radius := float32(64.0)
	fxRadius := float32(32.0)
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{0, 0, 255, 255})
		ctx.Renderer.ApplyVertBlur(canvas, ctx.Images[0], lx-radius, ly-radius, fxRadius, 1.0)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.SetColor(color.RGBA{255, 0, 0, 255})
		ctx.Renderer.ApplyHorzBlur(canvas, ctx.Images[0], rx-radius, ry-radius, fxRadius, 0.0)
	})
	app.Images = append(app.Images, app.Renderer.NewCircle(float64(radius)))
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestApplyHardShadow(t *testing.T) {
	radius := float32(64.0)
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{0, 128, 128, 128})
		ctx.Renderer.ApplyHardShadow(canvas, ctx.Images[0], lx-radius, ly-radius, 0, 8.0, ClampLeft)
		ctx.Renderer.SetColor(color.RGBA{255, 255, 255, 255})
		ctx.Renderer.DrawCircle(canvas, lx, ly, radius)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.SetColor(color.RGBA{128, 128, 128, 128})
		ctx.Renderer.ApplyHardShadow(canvas, ctx.Images[0], rx-radius, ry-radius, 8.0, 0.0, ClampNone)
		ctx.Renderer.SetColor(color.RGBA{255, 255, 255, 255})
		ctx.Renderer.DrawCircle(canvas, rx, ry, radius)
	})
	app.Images = append(app.Images, app.Renderer.NewCircle(float64(radius)))
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestApplyShadow(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{0, 128, 128, 128})
		ctx.Renderer.ApplyShadow(canvas, ctx.Images[0], lx, ly, 0, -16.0, 8.0, ClampBottom)
		var opts ebiten.DrawImageOptions
		opts.GeoM.Translate(float64(lx), float64(ly))
		canvas.DrawImage(ctx.Images[0], &opts)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.SetColor(color.RGBA{128, 128, 128, 128})
		ctx.Renderer.ApplyShadow(canvas, ctx.Images[0], rx, ry, -12.0, -12.0, 8.0, ClampBottom)
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(rx), float64(ry))
		canvas.DrawImage(ctx.Images[0], &opts)
	})
	circle := app.Renderer.NewCircle(64.0)
	circBounds := circle.Bounds()
	halfCircle := circle.SubImage(image.Rect(0, 0, circBounds.Dx(), circBounds.Dy()/2)).(*ebiten.Image)
	app.Images = append(app.Images, halfCircle)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

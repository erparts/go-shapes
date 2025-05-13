package shapes

import (
	"image"
	"image/color"
	"math"
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

		ctx.Renderer.SetColor(color.RGBA{0, 255, 255, 255})
		ctx.Renderer.SetBlend(ebiten.BlendLighter)
		ctx.Renderer.ApplyBlur(canvas, ctx.Images[0], rx-radius, ry-radius, fxRadius, 0.0)
		ctx.Renderer.SetBlend(ebiten.BlendSourceOver)
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

func TestApplyBlurKern(t *testing.T) {
	radius := float32(64.0)
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{255, 0, 255, 255})
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			ctx.Renderer.ApplyBlur2(canvas, ctx.Images[0], lx-radius, ly-radius, 32.0, 1.0)
		} else {
			ctx.Renderer.ApplyBlurD4(canvas, ctx.Images[0], lx-radius, ly-radius, GaussKern17, GaussKern3, 1.0)
		}
		ctx.Renderer.SetColor(color.RGBA{128, 0, 128, 128})
		ctx.Renderer.DrawCircle(canvas, lx, ly, radius)

		rx, ry := ctx.RightClickF32()
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			ctx.Renderer.ApplyBlur2(canvas, ctx.Images[0], rx-radius, ry-radius, 32.0, 1.0)
		} else {
			ctx.Renderer.ApplyBlurD4(canvas, ctx.Images[0], rx-radius, ry-radius, GaussKern17, GaussKern3, 1.0)
		}
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
		ctx.Renderer.ApplyShadow(canvas, ctx.Images[0], lx, ly, 0, -16.0, 4.0, ClampBottom)
		var opts ebiten.DrawImageOptions
		opts.GeoM.Translate(float64(lx), float64(ly))
		canvas.DrawImage(ctx.Images[0], &opts)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.SetColor(color.RGBA{128, 128, 128, 128})
		ctx.Renderer.ApplyShadow(canvas, ctx.Images[0], rx, ry, -12.0, -12.0, 9.0, ClampBottom)
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(rx), float64(ry))
		canvas.DrawImage(ctx.Images[0], &opts)

		mx, my := min(lx, rx), max(ly, ry)
		ctx.Renderer.SetColor(color.RGBA{0, 196, 196, 196})
		ctx.Renderer.ApplyShadow(canvas, ctx.Images[1], mx, my, 0, 0, 32.0, ClampNone)
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(mx), float64(my))
		canvas.DrawImage(ctx.Images[1], &opts)
	})
	circle := app.Renderer.NewCircle(64.0)
	circBounds := circle.Bounds()
	halfCircle := circle.SubImage(image.Rect(0, 0, circBounds.Dx(), circBounds.Dy()/2)).(*ebiten.Image)
	app.Images = append(app.Images, halfCircle, circle)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestApplyZoomShadow(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{0, 128, 128, 128})
		ctx.Renderer.ApplyZoomShadow(canvas, ctx.Images[0], lx, ly, 0, 0, 2.0, ClampNone)
		var opts ebiten.DrawImageOptions
		opts.GeoM.Translate(float64(lx), float64(ly))
		canvas.DrawImage(ctx.Images[0], &opts)

		rx, ry := ctx.RightClickF32()

		ctx.Renderer.SetColor(color.RGBA{128, 0, 128, 128})
		ctx.Renderer.ApplyZoomShadow(canvas, ctx.Images[1], rx, ry, 0, 0, 1.2, ClampBottom)

		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(rx), float64(ry))
		canvas.DrawImage(ctx.Images[1], &opts)
	})
	circle := app.Renderer.NewCircle(64.0)
	circBounds := circle.Bounds()
	halfCircle := circle.SubImage(image.Rect(0, 0, circBounds.Dx(), circBounds.Dy()/2)).(*ebiten.Image)
	app.Images = append(app.Images, circle, halfCircle)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestApplyGlow(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		var opts ebiten.DrawImageOptions
		opts.GeoM.Translate(float64(lx), float64(ly))
		canvas.DrawImage(ctx.Images[0], &opts)
		ctx.Renderer.ApplySimpleGlow(canvas, ctx.Images[0], lx, ly, 16)

		rx, ry := ctx.RightClickF32()
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(rx), float64(ry))
		canvas.DrawImage(ctx.Images[0], &opts)
		ctx.Renderer.SetColor(color.RGBA{255, 192, 192, 255})
		dynRadius := float32(ctx.DistAnim(6, 2.0))
		ctx.Renderer.ApplyGlow(canvas, ctx.Images[0], rx, ry, 24+dynRadius, 18, 0.5, 0.6, 0.0)
	})
	const s, m = 96, 16
	cross := ebiten.NewImage(s, s)
	app.Renderer.SetColor(color.RGBA{96, 240, 240, 255})
	app.Renderer.DrawLine(cross, m, m, s-m, s-m, m/2)
	app.Renderer.DrawLine(cross, s-m, m, m, s-m, m/2)
	app.Images = append(app.Images, cross)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestApplyHorzGlow(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		var opts ebiten.DrawImageOptions
		opts.GeoM.Translate(float64(lx), float64(ly))
		canvas.DrawImage(ctx.Images[0], &opts)
		ctx.Renderer.ApplyHorzGlow(canvas, ctx.Images[0], lx, ly, 16, 0.4, 0.5, 1.0)

		rx, ry := ctx.RightClickF32()
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(rx), float64(ry))
		canvas.DrawImage(ctx.Images[0], &opts)
		ctx.Renderer.SetColor(color.RGBA{255, 192, 192, 255})
		dynRadius := float32(ctx.DistAnim(6, 2.0))
		ctx.Renderer.ApplyHorzGlow(canvas, ctx.Images[0], rx, ry, 24+dynRadius, 0.5, 0.6, 0.0)
	})
	const s, m = 96, 16
	cross := ebiten.NewImage(s, s)
	app.Renderer.SetColor(color.RGBA{96, 240, 240, 255})
	app.Renderer.DrawLine(cross, m, m, s-m, s-m, m/2)
	app.Renderer.DrawLine(cross, s-m, m, m, s-m, m/2)
	app.Images = append(app.Images, cross)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestApplyDarkHorzGlow(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.White)

		lx, ly := ctx.LeftClickF32()
		var opts ebiten.DrawImageOptions
		opts.GeoM.Translate(float64(lx), float64(ly))
		canvas.DrawImage(ctx.Images[0], &opts)
		ctx.Renderer.ApplyDarkHorzGlow(canvas, ctx.Images[0], lx, ly, 16, 0.5, 0.01, 1.0)
		opts.GeoM.Translate(0, float64(120))
		canvas.DrawImage(ctx.Images[0], &opts)

		rx, ry := ctx.RightClickF32()
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(rx), float64(ry))
		canvas.DrawImage(ctx.Images[1], &opts)
		dynRadius := float32(ctx.DistAnim(6, 2.0))
		ctx.Renderer.SetColor(color.RGBA{64, 0, 0, 255})
		ctx.Renderer.ApplyDarkHorzGlow(canvas, ctx.Images[1], rx, ry, 24+dynRadius, 1, 0.5, 0.0)
	})
	const s, m = 96, 16
	cross := ebiten.NewImage(s, s)
	app.Renderer.SetColor(color.RGBA{0, 0, 128, 255})
	app.Renderer.DrawLine(cross, m, m, s-m, s-m, m/2)
	app.Renderer.DrawLine(cross, s-m, m, m, s-m, m/2)
	img := ebiten.NewImage(s, s)
	app.Renderer.SetColor(color.RGBA{0, 0, 0, 255})
	app.Renderer.SimpleGradient(img, color.RGBA{255, 255, 255, 255}, color.RGBA{128, 0, 0, 255}, math.Pi/2)
	app.Images = append(app.Images, img, cross)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

package shapes

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestApplyExpansion . -count 1
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

// go test -run ^TestApplyErosion . -count 1
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

// go test -run ^TestApplyOutline . -count 1
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

// go test -run ^TestApplyBlur . -count 1
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
		clrMix := float32(ctx.DistAnim(1.0, 1.0))
		ctx.Renderer.ApplyBlur(canvas, ctx.Images[0], rx-radius, ry-radius, fxRadius, clrMix)
	})
	app.Images = append(app.Images, app.Renderer.NewCircle(float64(radius)))
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestApplyBlur2 . -count 1
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
		clrMix := float32(ctx.DistAnim(1.0, 1.0))
		ctx.Renderer.ApplyBlur2(canvas, ctx.Images[0], rx-radius, ry-radius, fxRadius, clrMix)

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

// go test -run ^TestApplyDirBlur . -count 1
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

// go test -run ^TestApplyBlurKern . -count 1
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

// go test -run ^TestApplyHardShadow . -count 1
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

// go test -run ^TestApplyShadow . -count 1
func TestApplyShadow(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{0, 128, 128, 128})
		ctx.Renderer.ApplyShadow(canvas, ctx.Images[0], lx, ly, 0, -16.0, 4.0, ClampBottom)
		ctx.DrawAtF32(canvas, ctx.Images[0], lx, ly)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.SetColor(color.RGBA{128, 128, 128, 128})
		ctx.Renderer.ApplyShadow(canvas, ctx.Images[0], rx, ry, -12.0, -12.0, 9.0, ClampBottom)
		ctx.DrawAtF32(canvas, ctx.Images[0], rx, ry)

		mx, my := min(lx, rx), max(ly, ry)
		ctx.Renderer.SetColor(color.RGBA{0, 196, 196, 196})
		ctx.Renderer.ApplyShadow(canvas, ctx.Images[1], mx, my, 0, 0, 32.0, ClampNone)
		ctx.DrawAtF32(canvas, ctx.Images[1], mx, my)
	})
	circle := app.Renderer.NewCircle(64.0)
	circBounds := circle.Bounds()
	halfCircle := circle.SubImage(image.Rect(0, 0, circBounds.Dx(), circBounds.Dy()/2)).(*ebiten.Image)
	app.Images = append(app.Images, halfCircle, circle)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestApplyZoomShadow . -count 1
func TestApplyZoomShadow(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColor(color.RGBA{0, 128, 128, 128})
		ctx.Renderer.ApplyZoomShadow(canvas, ctx.Images[0], lx, ly, 0, 0, 2.0, ClampNone)
		ctx.DrawAtF32(canvas, ctx.Images[0], lx, ly)

		rx, ry := ctx.RightClickF32()

		ctx.Renderer.SetColor(color.RGBA{128, 0, 128, 128})
		ctx.Renderer.ApplyZoomShadow(canvas, ctx.Images[1], rx, ry, 0, 0, 1.2, ClampBottom)
		ctx.DrawAtF32(canvas, ctx.Images[1], rx, ry)
	})
	circle := app.Renderer.NewCircle(64.0)
	circBounds := circle.Bounds()
	halfCircle := circle.SubImage(image.Rect(0, 0, circBounds.Dx(), circBounds.Dy()/2)).(*ebiten.Image)
	app.Images = append(app.Images, circle, halfCircle)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestApplyGlow . -count 1
func TestApplyGlow(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[0], lx, ly)
		ctx.Renderer.ApplySimpleGlow(canvas, ctx.Images[0], lx, ly, 16)

		rx, ry := ctx.RightClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[0], rx, ry)
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

// go test -run ^TestApplyHorzGlow . -count 1
func TestApplyHorzGlow(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[0], lx, ly)
		ctx.Renderer.ApplyHorzGlow(canvas, ctx.Images[0], lx, ly, 16, 0.4, 0.5, 1.0)

		rx, ry := ctx.RightClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[0], rx, ry)
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

// go test -run ^TestApplyDarkHorzGlow . -count 1
func TestApplyDarkHorzGlow(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.White)

		lx, ly := ctx.LeftClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[0], lx, ly)
		ctx.Renderer.ApplyDarkHorzGlow(canvas, ctx.Images[0], lx, ly, 16, 0.5, 0.01, 1.0)
		ctx.DrawAtF32(canvas, ctx.Images[0], lx, ly+120)

		rx, ry := ctx.RightClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[1], rx, ry)
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

// go test -run ^TestApplyGlowD4 . -count 1
func TestApplyGlowD4(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[0], lx, ly)
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			ctx.Renderer.ApplyGlow(canvas, ctx.Images[0], lx, ly, 12, 12, 0.2, 0.8, 1.0)
		} else {
			ctx.Renderer.ApplyGlowD4(canvas, ctx.Images[0], lx, ly, GaussKern5, GaussKern5, 0.2, 0.8, 1.0)
		}

		rx, ry := ctx.RightClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[0], rx, ry)
		ctx.Renderer.SetColor(color.RGBA{255, 192, 192, 255})
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			ctx.Renderer.ApplyGlow(canvas, ctx.Images[0], rx, ry, 16, 32, 0.5, 0.6, 0.0)
		} else {
			ctx.Renderer.ApplyGlowD4(canvas, ctx.Images[0], rx, ry, GaussKern5, GaussKern15, 0.5, 0.6, 0.0)
		}
	})
	const s, m = 96, 16
	tri := ebiten.NewImage(s, s)
	app.Renderer.SetColor(color.RGBA{96, 240, 240, 255})
	app.Renderer.DrawTriangle(tri, m, s-m, s/2, m, s-m, s-m, 0)
	app.Images = append(app.Images, tri)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestApplyBlurKernBleed . -count 1
func TestApplyBlurKernBleed(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		i1, i2, i3 := 0, 1, 2
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			i1, i2, i3 = i2, i3, i1
		}
		ctx.Renderer.ApplyBlurD4(canvas, ctx.Images[i1], 16, 16, GaussKern9, GaussKern9, 1.0)
		ctx.Renderer.ApplyBlurD4(canvas, ctx.Images[i2], 16+96*1, 16, GaussKern17, GaussKern17, 1.0)
		ctx.Renderer.ApplyBlurD4(canvas, ctx.Images[i3], 16+96*2, 16, GaussKern11, GaussKern11, 1.0)
		ctx.Renderer.ApplyBlurD4(canvas, ctx.Images[i3], 16, 16+96*1, GaussKern5, GaussKern5, 1.0)
		ctx.Renderer.ApplyBlurD4(canvas, ctx.Images[i2], 16+96*1, 16+96*1, GaussKern13, GaussKern13, 1.0)
		ctx.Renderer.ApplyBlurD4(canvas, ctx.Images[i1], 16+96*2, 16+96*1, GaussKern9, GaussKern9, 1.0)
	})
	app.Renderer.SetColorF32(1, 0, 1, 1)
	img1 := app.Renderer.NewRect(33, 33)
	app.Renderer.SetColorF32(0, 1, 1, 1)
	img2 := app.Renderer.NewRect(50, 50)
	app.Renderer.SetColorF32(1, 1, 0, 1)
	img3 := app.Renderer.NewRect(67, 67)
	app.Renderer.SetColorF32(1, 1, 1, 1)
	app.Images = append(app.Images, img1, img2, img3)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestApplyGlowKernBleed . -count 1
// Notice: some bleeding edge cases are quite difficult to reproduce
// and haven't been able to catch them through tests yet, only live
// code in more complex projects.
func TestApplyGlowKernBleed(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		i1, i2, i3 := 0, 1, 2
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			i1, i2, i3 = i2, i3, i1
		}
		const st, et = 0.0, 0.5
		ctx.Renderer.ApplyGlowD4(canvas, ctx.Images[i1], 16, 16, GaussKern9, GaussKern9, st, et, 1.0)
		ctx.Renderer.ApplyGlowD4(canvas, ctx.Images[i2], 16+96*1, 16, GaussKern17, GaussKern17, st, et, 1.0)
		ctx.Renderer.ApplyGlowD4(canvas, ctx.Images[i3], 16+96*2, 16, GaussKern11, GaussKern11, st, et, 1.0)
		ctx.Renderer.ApplyGlowD4(canvas, ctx.Images[i3], 16, 16+96*1, GaussKern5, GaussKern5, st, et, 1.0)
		ctx.Renderer.ApplyGlowD4(canvas, ctx.Images[i2], 16+96*1, 16+96*1, GaussKern13, GaussKern13, st, et, 1.0)
		ctx.Renderer.ApplyGlowD4(canvas, ctx.Images[i1], 16+96*2, 16+96*1, GaussKern9, GaussKern9, st, et, 1.0)
	})
	app.Renderer.SetColorF32(1, 0, 1, 1)
	img1 := app.Renderer.NewRect(33, 33)
	app.Renderer.SetColorF32(0, 1, 1, 1)
	img2 := app.Renderer.NewRect(50, 50)
	app.Renderer.SetColorF32(1, 1, 0, 1)
	img3 := app.Renderer.NewRect(67, 67)
	app.Renderer.SetColorF32(1, 1, 1, 1)
	app.Images = append(app.Images, img1, img2, img3)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestApplyColorGlowD4 . -count 1
func TestApplyColorGlowD4(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[0], lx, ly)
		if !ebiten.IsKeyPressed(ebiten.KeySpace) {
			loThresh := 0.1 + ctx.DistAnim(0.4, 1.0)
			ctx.Renderer.ApplyColorGlowD4(canvas, ctx.Images[0], lx, ly, GaussKern7, GaussKern7, RGBF32(color.RGBA{255, 255, 0, 255}), float32(loThresh), 1.0, 1.0)
		}
	})

	circ := app.Renderer.NewCircle(96.0)
	img := ebiten.NewImage(circ.Bounds().Dx(), circ.Bounds().Dy())
	app.Renderer.Gradient(img, circ, 0, 0, color.RGBA{255, 255, 0, 255}, color.RGBA{255, 0, 255, 255}, -1, 0, 1.0)
	app.Images = append(app.Images, img)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestScanlinesSharp . -count 1
func TestScanlinesSharp(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.White)
		const darkThick, clearThick = 3, 1
		offset := float32(ctx.ModAnim(darkThick+clearThick, 1.0))
		ctx.Renderer.ApplyScanlinesSharp(canvas, darkThick, clearThick, 0.05, offset)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

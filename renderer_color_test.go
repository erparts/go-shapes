package shapes

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestGradient . -count 1
func TestGradient(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[0], lx, ly)

		rx, ry := ctx.RightClickF32()
		ctx.DrawAtF32(canvas, ctx.Images[1], rx, ry)

		ox, oy := 50, 400
		sub := canvas.SubImage(image.Rect(ox, oy, ox+80, oy+60)).(*ebiten.Image)
		ctx.Renderer.SimpleGradient(sub, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255}, 0)
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

// go test -run ^TestFlatPaint . -count 1
func TestFlatPaint(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColorF32(1.0, 0.0, 0.0, 1.0, 0, 1)
		ctx.Renderer.SetColorF32(1.0, 0.0, 1.0, 1.0, 2, 3)
		ctx.Renderer.FlatPaint(canvas, ctx.Images[0], lx, ly)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.SetColorF32(0.0, 1.0, 0.0, 1.0, 0, 3)
		ctx.Renderer.SetColorF32(0.0, 1.0, 1.0, 1.0, 1, 2)
		ctx.Renderer.FlatPaint(canvas, ctx.Images[1], rx, ry)
	})

	rect := app.Renderer.NewRect(120, 80)
	circ := app.Renderer.NewCircle(64.0)
	app.Renderer.Options().Blend = ebiten.BlendDestinationOut
	app.Renderer.SetColorF32(0.8, 0.8, 0.8, 0.8, 0, 1)
	app.Renderer.SetColorF32(0.3, 0.3, 0.3, 0.3, 2, 3)
	app.Renderer.DrawCircle(circ, 64.0, 64.0, 42.0)
	app.Renderer.Options().Blend = ebiten.BlendSourceOver
	app.Images = append(app.Images, rect, circ)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

// go test -run ^TestDitherMat4 . -count 1
func TestDitherMat4(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.RGBA{128, 0, 128, 255})

		var mat [16]float32
		switch int(ctx.ModAnim(5.0, 0.25)) {
		case 0:
			mat = DitherBayes
		case 1:
			mat = DitherDots
		case 2:
			mat = DitherGlitch
		case 3:
			mat = DitherSerp
		case 4:
			mat = DitherCrumbs
		}

		lx, ly := ctx.LeftClickF32()
		ctx.Renderer.SetColorF32(1.0, 0.0, 0.0, 1.0, 0, 1)
		ctx.Renderer.SetColorF32(1.0, 0.0, 1.0, 1.0, 2, 3)
		anim := float32(ctx.DistAnim(1.0, 1.0))
		yOffset := int(ctx.ModAnim(4.0, 1.0))
		xOffset := 8 - int(ctx.DistAnim(16.0, 1.0))
		ctx.Renderer.DitherMat4(canvas, ctx.Images[0], lx, ly, xOffset, yOffset, DitherBRG, mat, anim, 0.0)

		rx, ry := ctx.RightClickF32()
		ctx.Renderer.DitherMat4(canvas, ctx.Images[1], rx, ry, 0, 0, DitherAlpha8, mat, 0, anim)
	})

	from, to := color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}
	gradient := app.Renderer.NewSimpleGradient(160, 160, from, to, DirRadsLTR)
	app.Images = append(app.Images, gradient)
	from, to = color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}
	gradient = app.Renderer.NewSimpleGradient(160, 160, from, to, DirRadsTLBR)
	app.Images = append(app.Images, gradient)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

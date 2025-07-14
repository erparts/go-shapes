package shapes

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestMap ./... -count 1
func TestMap(t *testing.T) {
	const CardWidth, CardHeight = 128, 164
	card := newCardWaver(CardWidth, CardHeight)

	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		card.Update()

		canvas.Fill(color.Black)

		ox, oy, w, h := rectOriginSize(canvas.Bounds())
		c1x, c1y := float32(ox+w/4), float32(oy+h/4)
		c2x, c2y := float32(ox+3*w/4), float32(oy+h/4)
		c3x, c3y := float32(ox+w/4), float32(oy+3*h/4)
		c4x, c4y := float32(ox+3*w/4), float32(oy+3*h/4)

		ctx.Renderer.mapQuad2(canvas, ctx.Images[0], card.Quad(c1x, c1y))
		ctx.Renderer.MapQuad4(canvas, ctx.Images[0], card.Quad(c2x, c2y))
		ctx.Renderer.SetColorF32(1, 1, 1, 1)
		if (ctx.Ticks/180)&1 == 1 {
			ctx.Renderer.SetColorF32(0, 0, 0, 0, 2, 3)
		}
		ctx.Renderer.MapProjective(canvas, ctx.Images[0], card.Quad(c3x, c3y))
		ctx.Renderer.MapQuad4(canvas, ctx.Images[0], card.Quad(c4x, c4y))
	})

	img := ebiten.NewImage(CardWidth+2, CardHeight+2)
	app.Renderer.DrawRect(img, image.Rect(1, 1, CardWidth+1, CardHeight+1), 6.0)
	app.Renderer.SetBlend(ebiten.BlendSourceAtop)
	app.Renderer.SetColorF32(0, 0, 0, 0.1)
	app.Renderer.TileDotsHex(img, 4.0, 12.0, 0, 0)

	app.Renderer.SetBlend(ebiten.BlendSourceOver)
	app.Renderer.SetColorF32(1, 1, 1, 1)

	app.Images = append(app.Images, img)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

func TestMapProjectiveStress(t *testing.T) {
	const CardWidth, CardHeight = 32, 48
	const NumRows = 9
	const NumCols = 18

	ebiten.SetVsyncEnabled(false)

	cardWavers := make([]*cardWaver, NumRows*NumCols)
	for i := range cardWavers {
		cardWavers[i] = newCardWaver(CardWidth, CardHeight)
	}

	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		for _, cardWaver := range cardWavers {
			cardWaver.Update()
		}
		canvas.Fill(color.Black)

		ox, oy, w, h := rectOriginSize(canvas.Bounds())
		dx, dy := float32(w)/(NumCols+1), float32(h)/(NumRows+1)
		waverIdx := 0
		cy := float32(oy) + dy
		for range NumRows {
			cx := float32(ox) + dx
			for range NumCols {
				ctx.Renderer.MapProjective(canvas, ctx.Images[waverIdx&1], cardWavers[waverIdx].Quad(cx, cy))
				cx += dx
				waverIdx += 1
			}
			cy += dy
		}
	})

	img := ebiten.NewImage(CardWidth+2, CardHeight+2)
	img2 := ebiten.NewImage(CardWidth+2, CardHeight+2)
	app.Renderer.DrawRect(img, image.Rect(1, 1, CardWidth+1, CardHeight+1), 6.0)
	app.Renderer.SetColorF32(1, 0.3, 1, 1)
	app.Renderer.DrawRect(img2, image.Rect(1, 1, CardWidth+1, CardHeight+1), 6.0)

	app.Renderer.SetBlend(ebiten.BlendSourceAtop)
	app.Renderer.SetColorF32(0, 0, 0, 0.1)
	app.Renderer.TileDotsHex(img, CardWidth/18, CardWidth/8, 0, 0)
	app.Renderer.SetColorF32(1, 0, 0, 0.1)
	app.Renderer.TileDotsHex(img2, CardWidth/18, CardWidth/8, 0, 0)

	app.Renderer.SetBlend(ebiten.BlendSourceOver)
	app.Renderer.SetColorF32(1, 1, 1, 1)

	app.Images = append(app.Images, img, img2)
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

type cardWaver struct {
	w, h  float64
	xTilt float64
	yTilt float64
	rads  float64
}

func newCardWaver(w, h float64) *cardWaver {
	var waver cardWaver
	waver.w, waver.h = w, h
	return &waver
}

func (cw *cardWaver) Update() {
	cw.rads += 0.02
	if cw.rads > 2*math.Pi {
		cw.rads -= 2 * math.Pi
	}
	cw.xTilt = math.Cos(cw.rads) * -0.15
	cw.yTilt = math.Sin(cw.rads) * 0.11
}

func (cw *cardWaver) Quad(cx, cy float32) [4]PointF32 {
	const DepthStrength = 1.0 / 3.0

	xTiltSin, xTiltCos := math.Sincos(cw.xTilt)
	yTiltSin, yTiltCos := math.Sincos(cw.yTilt)

	xCosSigns := [4]float64{-1, 1, 1, -1}
	yCosSigns := [4]float64{-1, -1, 1, 1}
	xzSinSigns := [4]float64{-1, 1, -1, 1}
	yzSinSigns := [4]float64{-1, 1, -1, 1}
	var pts [4]PointF32
	for i := range pts {
		xOffset := xTiltCos * xCosSigns[i] * cw.w / 2.0
		yOffset := yTiltCos * yCosSigns[i] * cw.h / 2.0

		xzOffset := xTiltSin * xzSinSigns[i] * (cw.w / 2.0) * DepthStrength
		yzOffset := yTiltSin * yzSinSigns[i] * (cw.h / 2.0) * DepthStrength

		pts[i].X = cx + float32(xOffset+yzOffset)
		pts[i].Y = cy + float32(yOffset+xzOffset)
	}

	return pts
}

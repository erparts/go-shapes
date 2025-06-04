package shapes

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestGoldenRatioGen . -count 1
func TestGoldenRatioGen(t *testing.T) {
	const BarWidth, BarHeight = 32, 96
	const Intersp = 4

	var gen GoldenRatioGen
	gen.n = 514230.0 - 16.0
	var values [16]float64
	for i := range len(values) {
		values[i] = gen.Float64() * BarHeight
	}

	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)

		if ctx.Ticks%60 == 0 {
			copy(values[:], values[1:])
			values[len(values)-1] = gen.Float64() * BarHeight
		}

		oxF64 := float64(canvas.Bounds().Dx())/2.0 - (float64(len(values))*(BarWidth+Intersp)-Intersp)/2.0
		oyF64 := float64(canvas.Bounds().Dy())/2.0 + BarHeight/2.0
		for _, v := range values {
			ctx.Renderer.DrawArea(canvas, float32(oxF64), float32(oyF64), BarWidth, -float32(v), 2.0)
			oxF64 += BarWidth + Intersp
		}
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

package shapes

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// go test -run ^TestStudyWaveFuncs . -count 1
func TestStudyWaveFuncs(t *testing.T) {
	app := NewTestApp(func(canvas *ebiten.Image, ctx TestAppCtx) {
		canvas.Fill(color.Black)
		ctx.Renderer.studyWaveFuncs(canvas, 1.0, 8.0)
	})
	if err := ebiten.RunGame(app); err != nil {
		t.Fatal(err)
	}
}

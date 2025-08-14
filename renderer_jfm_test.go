package shapes

import (
	"fmt"
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func jfmDebugPrint(t *testing.T, out *image.RGBA) {
	const DisplayMode string = "coords" // "coords", "rgba" or "dual"
	const DebugPrint bool = false

	decodeAxisOffsetToSeed := func(a, b int) int {
		magnitude := ((a >> 1) << 8) + b
		sign := (1 - ((a & 1) << 1))
		return sign * magnitude
	}
	fmtRGBA := func(rgba color.RGBA) string {
		x := decodeAxisOffsetToSeed(int(rgba.R), int(rgba.G))
		y := decodeAxisOffsetToSeed(int(rgba.B), int(rgba.A))
		switch DisplayMode {
		case "coords":
			return fmt.Sprintf("[%+04d %+04d]", x, y)
		case "rgba":
			return fmt.Sprintf("[%03d %03d %03d %03d]", rgba.R, rgba.G, rgba.B, rgba.A)
		case "dual":
			return fmt.Sprintf("[%+04d %+04d](%03d %03d %03d %03d)", x, y, rgba.R, rgba.G, rgba.B, rgba.A)
		default:
			panic("invalid display mode '" + DisplayMode + "'")
		}
	}
	if DebugPrint {
		for y := range out.Rect.Max.Y {
			fmt.Printf("row#%d ", y)
			for x := range out.Rect.Max.X {
				fmt.Printf("%s ", fmtRGBA(out.RGBAAt(x, y)))
			}
			fmt.Printf("\n")
		}
		t.Fatalf("debug print")
	}
}

// go test -run ^TestJFMCompute$ . -count 1
func TestJFMCompute(t *testing.T) {
	r := NewRenderer()
	src := ebiten.NewImage(9, 9)
	// top-left hollow rectangle
	src.Set(0, 0, color.White)
	src.Set(1, 0, color.White)
	src.Set(2, 0, color.White)
	src.Set(0, 1, color.White)
	src.Set(2, 1, color.White)
	src.Set(0, 2, color.White)
	src.Set(1, 2, color.White)
	src.Set(2, 2, color.White)

	dst := ebiten.NewImage(9, 9)
	r.JFMCompute(dst, src, JFMBoundary, 4)

	out := image.NewRGBA(image.Rect(0, 0, 9, 9))
	if err := ebiten.RunGame(&testOutputWriter{subject: dst, out: out.Pix}); err != nil {
		t.Fatal(err)
	}
	jfmDebugPrint(t, out)

	expectedOut := image.NewRGBA(image.Rect(0, 0, 9, 9))
	for i := range expectedOut.Pix {
		expectedOut.Pix[i] = 255
	}
	expectedOut.SetRGBA(0, 0, color.RGBA{0, 0, 0, 0})
	expectedOut.SetRGBA(1, 0, color.RGBA{0, 0, 0, 0})
	expectedOut.SetRGBA(2, 0, color.RGBA{0, 0, 0, 0})
	expectedOut.SetRGBA(0, 1, color.RGBA{0, 0, 0, 0})
	expectedOut.SetRGBA(2, 1, color.RGBA{0, 0, 0, 0})
	expectedOut.SetRGBA(0, 2, color.RGBA{0, 0, 0, 0})
	expectedOut.SetRGBA(1, 2, color.RGBA{0, 0, 0, 0})
	expectedOut.SetRGBA(2, 2, color.RGBA{0, 0, 0, 0})

	expectedOut.SetRGBA(1, 1, color.RGBA{1, 1, 0, 0})
	for c := range 3 {
		for i := range 4 {
			expectedOut.SetRGBA(c, 3+i, color.RGBA{0, 0, 1, uint8(i + 1)})
			expectedOut.SetRGBA(3+i, c, color.RGBA{1, uint8(i + 1), 0, 0})
		}
	}
	low := [][]color.RGBA{
		{{1, 1, 1, 1}, {1, 2, 1, 1}, {1, 3, 1, 1}},
		{{1, 1, 1, 2}, {1, 2, 1, 2}, {1, 3, 1, 2}},
		{{1, 1, 1, 3}, {1, 2, 1, 3}},
	}
	for y, row := range low {
		for x, value := range row {
			expectedOut.SetRGBA(3+x, 3+y, value)
		}
	}

	if !slices.Equal(out.Pix, expectedOut.Pix) {
		t.Fatalf("expected slices.Equal(expected, out)\nexpected:\n%v\nout:\n%v", expectedOut.Pix, out.Pix)
	}
}

// go test -run ^TestJFMCompute2$ . -count 1
func TestJFMCompute2(t *testing.T) {
	r := NewRenderer()
	src := ebiten.NewImage(1, 260)
	dst := ebiten.NewImage(1, 260)
	src.Set(0, 258, color.White)
	r.JFMCompute(dst, src, JFMBoundary, 257)

	out := image.NewRGBA(image.Rect(0, 0, 1, 260))
	if err := ebiten.RunGame(&testOutputWriter{subject: dst, out: out.Pix}); err != nil {
		t.Fatal(err)
	}
	jfmDebugPrint(t, out)

	expectedOut := image.NewRGBA(image.Rect(0, 0, 1, 260))
	for i := range 260 {
		n := 258 - i
		if n < 256 {
			expectedOut.SetRGBA(0, i, color.RGBA{0, 0, 0, uint8(n)})
		} else {
			expectedOut.SetRGBA(0, i, color.RGBA{0, 0, 2, uint8(n - 256)})
		}
	}
	expectedOut.SetRGBA(0, 258, color.RGBA{0, 0, 0, 0})       // seed
	expectedOut.SetRGBA(0, 259, color.RGBA{0, 0, 1, 1})       // after seed
	expectedOut.SetRGBA(0, 0, color.RGBA{255, 255, 255, 255}) // outside range

	if !slices.Equal(out.Pix, expectedOut.Pix) {
		t.Fatalf("expected slices.Equal(expected, out)\nexpected:\n%v\nout:\n%v", expectedOut.Pix, out.Pix)
	}
}

type testOutputWriter struct {
	ticks   int
	subject *ebiten.Image
	out     []byte
}

func (t *testOutputWriter) Draw(*ebiten.Image) {}
func (t *testOutputWriter) Layout(w, h int) (int, int) {
	return w, h
}
func (t *testOutputWriter) Update() error {
	t.ticks += 1
	if t.ticks == 32 {
		t.subject.ReadPixels(t.out)
	}
	if t.ticks >= 64 {
		return ebiten.Termination
	}
	return nil
}

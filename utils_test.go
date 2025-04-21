package shapes

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type BaseTestApp struct{}

func (*BaseTestApp) Layout(w, h int) (int, int) {
	return w, h
}
func (*BaseTestApp) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	return nil
}
func (*BaseTestApp) Draw(*ebiten.Image) {}

type TestAppCtx struct {
	Renderer   *Renderer
	Images     []*ebiten.Image
	Ticks      uint64
	LeftClick  image.Point
	RightClick image.Point
}

func (ctx *TestAppCtx) LeftClickF32() (x, y float32) {
	return float32(ctx.LeftClick.X), float32(ctx.LeftClick.Y)
}
func (ctx *TestAppCtx) RightClickF32() (x, y float32) {
	return float32(ctx.RightClick.X), float32(ctx.RightClick.Y)
}
func (ctx *TestAppCtx) LeftClickF64() (x, y float64) {
	return float64(ctx.LeftClick.X), float64(ctx.LeftClick.Y)
}
func (ctx *TestAppCtx) RightClickF64() (x, y float64) {
	return float64(ctx.RightClick.X), float64(ctx.RightClick.Y)
}
func (ctx *TestAppCtx) RadsAnim(speedFactor float64) float64 {
	return math.Pi * math.Sin(float64(ctx.Ticks)*0.01*speedFactor)
}
func (ctx *TestAppCtx) DistAnim(maxDist, speedFactor float64) float64 {
	return maxDist * (math.Sin(float64(ctx.Ticks)*0.02*speedFactor) + 1.0) / 2.0
}

type TestApp struct {
	BaseTestApp
	TestAppCtx
	drawer func(canvas *ebiten.Image, ctx TestAppCtx)
}

func NewTestApp(drawer func(canvas *ebiten.Image, ctx TestAppCtx), images ...*ebiten.Image) *TestApp {
	//ebiten.SetVsyncEnabled(false)
	var app TestApp
	app.Images = images
	app.LeftClick = image.Pt((128*4)/3, 128)
	app.RightClick = image.Pt((320*4)/3, 320)
	app.Renderer = NewRenderer()
	app.drawer = drawer
	return &app
}

func (app *TestApp) Update() error {
	app.Ticks += 1
	left := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	right := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	if left || right {
		x, y := ebiten.CursorPosition()
		if left {
			app.LeftClick = image.Pt(x, y)
		} else {
			app.RightClick = image.Pt(x, y)
		}
	}
	ebiten.SetWindowTitle(fmt.Sprintf("LeftClick (%d, %d), RightClick (%d, %d) [%.02f FPS]", app.LeftClick.X, app.LeftClick.Y, app.RightClick.X, app.RightClick.Y, ebiten.ActualFPS()))
	return app.BaseTestApp.Update()
}

func (app *TestApp) Draw(canvas *ebiten.Image) {
	app.drawer(canvas, app.TestAppCtx)
}

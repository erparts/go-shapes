package shapes

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shaders/rect.kage
var shaderRectSrc []byte

//go:embed shaders/line.kage
var shaderLineSrc []byte

//go:embed shaders/circle.kage
var shaderCircleSrc []byte

//go:embed shaders/ellipse.kage
var shaderEllipseSrc []byte

//go:embed shaders/triangle.kage
var shaderTriangleSrc []byte

//go:embed shaders/hexagon.kage
var shaderHexagonSrc []byte

//go:embed shaders/expansion.kage
var shaderExpansionSrc []byte

//go:embed shaders/erosion.kage
var shaderErosionSrc []byte

//go:embed shaders/outline.kage
var shaderOutlineSrc []byte

//go:embed shaders/pattern_dots.kage
var shaderPatternDotsSrc []byte

var shaderRect *ebiten.Shader
var shaderLine *ebiten.Shader
var shaderCircle *ebiten.Shader
var shaderEllipse *ebiten.Shader
var shaderTriangle *ebiten.Shader
var shaderHexagon *ebiten.Shader
var shaderExpansion *ebiten.Shader
var shaderErosion *ebiten.Shader
var shaderOutline *ebiten.Shader
var shaderPatternDots *ebiten.Shader

func mustCompile(src []byte) *ebiten.Shader {
	shader, err := ebiten.NewShader(src)
	if err != nil {
		panic(err)
	}
	return shader
}

func ensureShaderLineLoaded() {
	if shaderLine == nil {
		shaderLine = mustCompile(shaderLineSrc)
	}
}

func ensureShaderCircleLoaded() {
	if shaderCircle == nil {
		shaderCircle = mustCompile(shaderCircleSrc)
	}
}

func ensureShaderRectLoaded() {
	if shaderRect == nil {
		shaderRect = mustCompile(shaderRectSrc)
	}
}

func ensureShaderEllipseLoaded() {
	if shaderEllipse == nil {
		shaderEllipse = mustCompile(shaderEllipseSrc)
	}
}

func ensureShaderTriangleLoaded() {
	if shaderTriangle == nil {
		shaderTriangle = mustCompile(shaderTriangleSrc)
	}
}

func ensureShaderHexagonLoaded() {
	if shaderHexagon == nil {
		shaderHexagon = mustCompile(shaderHexagonSrc)
	}
}

func ensureShaderPatternDotsLoaded() {
	if shaderPatternDots == nil {
		shaderPatternDots = mustCompile(shaderPatternDotsSrc)
	}
}

func ensureShaderOutlineLoaded() {
	if shaderOutline == nil {
		shaderOutline = mustCompile(shaderOutlineSrc)
	}
}

func ensureShaderExpansionLoaded() {
	if shaderExpansion == nil {
		shaderExpansion = mustCompile(shaderExpansionSrc)
	}
}

func ensureShaderErosionLoaded() {
	if shaderErosion == nil {
		shaderErosion = mustCompile(shaderErosionSrc)
	}
}

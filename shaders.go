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

//go:embed shaders/blur.kage
var shaderBlurSrc []byte

//go:embed shaders/horz_blur.kage
var shaderHorzBlurSrc []byte

//go:embed shaders/vert_blur.kage
var shaderVertBlurSrc []byte

//go:embed shaders/glow.kage
var shaderGlowSrc []byte

//go:embed shaders/shadow.kage
var shaderShadowSrc []byte

//go:embed shaders/hard_shadow.kage
var shaderHardShadowSrc []byte

//go:embed shaders/tile_dots_hex.kage
var shaderTileDotsHexSrc []byte

var shaderRect *ebiten.Shader
var shaderLine *ebiten.Shader
var shaderCircle *ebiten.Shader
var shaderEllipse *ebiten.Shader
var shaderTriangle *ebiten.Shader
var shaderHexagon *ebiten.Shader
var shaderExpansion *ebiten.Shader
var shaderErosion *ebiten.Shader
var shaderOutline *ebiten.Shader
var shaderBlur *ebiten.Shader
var shaderHorzBlur *ebiten.Shader
var shaderVertBlur *ebiten.Shader
var shaderGlow *ebiten.Shader
var shaderShadow *ebiten.Shader
var shaderHardShadow *ebiten.Shader
var shaderTileDotsHex *ebiten.Shader

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

func ensureShaderOutlineLoaded() {
	if shaderOutline == nil {
		shaderOutline = mustCompile(shaderOutlineSrc)
	}
}

func ensureShaderBlurLoaded() {
	if shaderBlur == nil {
		shaderBlur = mustCompile(shaderBlurSrc)
	}
}

func ensureShaderHorzBlurLoaded() {
	if shaderHorzBlur == nil {
		shaderHorzBlur = mustCompile(shaderHorzBlurSrc)
	}
}

func ensureShaderVertBlurLoaded() {
	if shaderVertBlur == nil {
		shaderVertBlur = mustCompile(shaderVertBlurSrc)
	}
}

func ensureShaderGlowLoaded() {
	if shaderGlow == nil {
		shaderGlow = mustCompile(shaderGlowSrc)
	}
}

func ensureShaderShadowLoaded() {
	if shaderShadow == nil {
		shaderShadow = mustCompile(shaderShadowSrc)
	}
}

func ensureShaderHardShadowLoaded() {
	if shaderHardShadow == nil {
		shaderHardShadow = mustCompile(shaderHardShadowSrc)
	}
}

func ensureShaderTileDotsHexLoaded() {
	if shaderTileDotsHex == nil {
		shaderTileDotsHex = mustCompile(shaderTileDotsHexSrc)
	}
}

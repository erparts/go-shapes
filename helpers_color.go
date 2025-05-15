package shapes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Common blend modes not directly exposed on Ebitengine.
var (
	BlendSubtract = ebiten.Blend{
		BlendFactorSourceRGB:        ebiten.BlendFactorOne,
		BlendFactorSourceAlpha:      ebiten.BlendFactorOne,
		BlendFactorDestinationRGB:   ebiten.BlendFactorOne,
		BlendFactorDestinationAlpha: ebiten.BlendFactorOne,
		BlendOperationRGB:           ebiten.BlendOperationReverseSubtract,
		BlendOperationAlpha:         ebiten.BlendOperationReverseSubtract,
	}
	BlendMultiply = ebiten.Blend{
		BlendFactorSourceRGB:        ebiten.BlendFactorDestinationColor,
		BlendFactorSourceAlpha:      ebiten.BlendFactorDestinationColor,
		BlendFactorDestinationRGB:   ebiten.BlendFactorOneMinusSourceAlpha,
		BlendFactorDestinationAlpha: ebiten.BlendFactorOneMinusSourceAlpha,
		BlendOperationRGB:           ebiten.BlendOperationAdd,
		BlendOperationAlpha:         ebiten.BlendOperationAdd,
	}
)

func ColorToF32(clr color.Color) [4]float32 {
	r, g, b, a := clr.RGBA()
	return [4]float32{float32(r) / 65535.0, float32(g) / 65535.0, float32(b) / 65535.0, float32(a) / 65535.0}
}

func RGBF32(clr color.Color) [3]float32 {
	r, g, b, _ := clr.RGBA()
	return [3]float32{float32(r) / 65535.0, float32(g) / 65535.0, float32(b) / 65535.0}
}

func colorToF64(clr color.Color) [4]float64 {
	r, g, b, a := clr.RGBA()
	return [4]float64{float64(r) / 65535.0, float64(g) / 65535.0, float64(b) / 65535.0, float64(a) / 65535.0}
}

func f32ToRGBA64(r, g, b, a float32) color.RGBA64 {
	return color.RGBA64{
		R: uint16(r * 65535.0),
		G: uint16(g * 65535.0),
		B: uint16(b * 65535.0),
		A: uint16(a * 65535.0),
	}
}

func interpColor(ox, oy, fx, fy float32, tlClr, trClr, blClr, brClr [4]float32, x, y float32) [4]float32 {
	u := min(max((x-ox)/(fx-ox), 0), 1)
	v := min(max((y-oy)/(fy-oy), 0), 1)

	var result [4]float32
	for i := range 4 {
		topClr := tlClr[i]*(1-u) + trClr[i]*u
		bottomClr := blClr[i]*(1-u) + brClr[i]*u
		result[i] = topClr*(1-v) + bottomClr*v
	}
	return result
}

func interpVertexColor(a, b ebiten.Vertex, t float32) (cr, cg, cb, ca float32) {
	return lerp(a.ColorR, b.ColorR, t), lerp(a.ColorG, b.ColorG, t), lerp(a.ColorB, b.ColorB, t), lerp(a.ColorA, b.ColorA, t)
}

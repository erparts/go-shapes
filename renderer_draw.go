package shapes

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func (r *Renderer) NewRect(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	rgba := r.GetColorF32()
	img.Fill(f32ToRGBA64(rgba[0], rgba[1], rgba[2], rgba[3]))
	return img
}

func (r *Renderer) NewCircle(radius float64) *ebiten.Image {
	side := float32(math.Ceil(radius * 2))
	img := ebiten.NewImage(int(side), int(side))
	r.DrawCircle(img, side/2, side/2, float32(radius))
	return img
}

func (r *Renderer) DrawRect(target *ebiten.Image, ox, oy, w, h, rounding float32) {
	r.setDstRectCoords(ox, oy, ox+w, oy+h)
	ensureShaderRectLoaded()
	r.setFlatCustomVAs(ox, oy, w, h)
	r.opts.Uniforms["Rounding"] = rounding
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderRect, &r.opts)
}

// DrawLine draws a smooth line between the given two points, with rounded ends.
func (r *Renderer) DrawLine(target *ebiten.Image, ox, oy, fx, fy float64, thickness float64) {
	vdx, vdy := fx-ox, fy-oy // non-normalized vector
	vpx, vpy := -vdy, vdx    // perpendicular vector
	length := math.Hypot(vdx, vdy)
	if length == 0 {
		length = 1
	}
	// scale for vector normalization
	scale := (thickness / 2) / length

	// adjust bounding ends to include thickness rounding
	box, boy := ox-vdx*scale, oy-vdy*scale
	bfx, bfy := fx+vdx*scale, fy+vdy*scale

	// compute bounding vertices applying the perpendicular offset
	svpx, svpy := vpx*scale, vpy*scale
	r.vertices[0].DstX = float32(box + svpx)
	r.vertices[0].DstY = float32(boy + svpy)
	r.vertices[1].DstX = float32(bfx + svpx)
	r.vertices[1].DstY = float32(bfy + svpy)
	r.vertices[2].DstX = float32(bfx - svpx)
	r.vertices[2].DstY = float32(bfy - svpy)
	r.vertices[3].DstX = float32(box - svpx)
	r.vertices[3].DstY = float32(boy - svpy)

	r.setFlatCustomVAs(float32(ox), float32(oy), float32(fx), float32(fy))
	r.opts.Uniforms["Thickness"] = float32(thickness)

	// draw shader
	ensureShaderLineLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderLine, &r.opts)
}

func (r *Renderer) DrawCircle(target *ebiten.Image, ox, oy, radius float32) {
	r.setDstRectCoords(ox-radius, oy-radius, ox+radius, oy+radius)
	ensureShaderCircleLoaded()
	r.setFlatCustomVAs(ox, oy, radius, 0.0)
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderCircle, &r.opts)
}

// Notice: ellipses don't have a perfect SDF, so approximations can be very slightly
// bigger or smaller than the requested radiuses.
func (r *Renderer) DrawEllipse(target *ebiten.Image, ox, oy, horzRadius, vertRadius float32, rads float64) {
	if rads == 0 {
		r.setDstRectCoords(ox-horzRadius, oy-vertRadius, ox+horzRadius, oy+vertRadius)
		r.opts.Uniforms["Radians"] = 0
	} else {
		hRadiusF64, vRadiusF64 := float64(horzRadius), float64(vertRadius)
		rc, rs := math.Cos(rads), math.Sin(rads)
		halfWidth := float32(math.Hypot(hRadiusF64*rc, vRadiusF64*rs))
		halfHeight := float32(math.Hypot(hRadiusF64*rs, vRadiusF64*rc))
		r.DrawRect(target, ox-halfWidth, oy-halfHeight, halfWidth, halfHeight, 0)
		r.setDstRectCoords(ox-halfWidth, oy-halfHeight, ox+halfWidth, oy+halfHeight)
		r.opts.Uniforms["Radians"] = rads
	}
	r.setFlatCustomVAs(ox, oy, horzRadius, vertRadius)
	ensureShaderEllipseLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderEllipse, &r.opts)
}

// DrawTriangle draws a smooth triangle using the given vertices, and an optional rounding factor.
// Notice that, if provided, handling the rounding is relatively expensive (two dozen f64 products
// and 3 square roots)
func (r *Renderer) DrawTriangle(target *ebiten.Image, ox1, oy1, ox2, oy2, ox3, oy3, rounding float64) {
	area := math.Abs((ox1*(oy2-oy3) + ox2*(oy3-oy1) + ox3*(oy1-oy2)) / 2)
	if area < 1e-6 {
		return // empty triangle
	}

	var iox1, ioy1, iox2, ioy2, iox3, ioy3 float64 = ox1, oy1, ox2, oy2, ox3, oy3
	if rounding != 0 {
		a12, b12, c12 := toLinearFormABC(ox1, oy1, ox2, oy2)
		a23, b23, c23 := toLinearFormABC(ox2, oy2, ox3, oy3)
		a31, b31, c31 := toLinearFormABC(ox3, oy3, ox1, oy1)
		c1_12, c2_12 := parallelsAtDist(a12, b12, c12, rounding)
		c1_23, c2_23 := parallelsAtDist(a23, b23, c23, rounding)
		c1_31, c2_31 := parallelsAtDist(a31, b31, c31, rounding)
		if a12*ox3+b12*oy3+c12 > 0 { // fancy winding order test
			c12, c23, c31 = c1_12, c1_23, c1_31
		} else {
			c12, c23, c31 = c2_12, c2_23, c2_31
		}
		iox1, ioy1 = shortCramer(a31, b31, c31, a12, b12, c12)
		iox2, ioy2 = shortCramer(a12, b12, c12, a23, b23, c23)
		iox3, ioy3 = shortCramer(a23, b23, c23, a31, b31, c31)
	}

	minX, maxX := min(ox1, ox2, ox3), max(ox1, ox2, ox3)
	minY, maxY := min(oy1, oy2, oy3), max(oy1, oy2, oy3)
	r.setDstRectCoords(float32(minX), float32(minY), float32(maxX), float32(maxY))

	// draw shader
	ensureShaderTriangleLoaded()
	r.opts.Uniforms["P0"] = []float32{float32(iox1), float32(ioy1)}
	r.opts.Uniforms["P1"] = []float32{float32(iox2), float32(ioy2)}
	r.opts.Uniforms["P2"] = []float32{float32(iox3), float32(ioy3)}
	r.opts.Uniforms["Rounding"] = float32(rounding)
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderTriangle, &r.opts)
}

// DrawHexagon renders an hexagon that can be fully contained within the given radius.
// Roundness can be used to round the corners. Rads can be used to rotate the hexagon,
// in radians.
func (r *Renderer) DrawHexagon(target *ebiten.Image, ox, oy, radius, roundness, rads float32) {
	r.setDstRectCoords(ox-radius, oy-radius, ox+radius, oy+radius)

	// draw shader
	const apothemToRadiusFactor = 0.866025404 // math.Sqrt(3)/2
	apothem := (radius - roundness) * apothemToRadiusFactor
	r.setFlatCustomVAs(ox, oy, apothem, rads)
	r.opts.Uniforms["Roundness"] = roundness
	ensureShaderHexagonLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderHexagon, &r.opts)
}

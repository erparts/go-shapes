package shapes

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func (r *Renderer) NewRect(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	r.DrawIntArea(img, 0, 0, width, height)
	return img
}

func (r *Renderer) NewCircle(radius float64) *ebiten.Image {
	side := float32(math.Ceil(radius * 2))
	img := ebiten.NewImage(int(side), int(side))
	r.DrawCircle(img, side/2, side/2, float32(radius))
	return img
}

func (r *Renderer) NewRing(inRadius, outRadius float64) *ebiten.Image {
	side := float32(math.Ceil(outRadius * 2))
	img := ebiten.NewImage(int(side), int(side))
	r.DrawRing(img, side/2, side/2, float32(inRadius), float32(outRadius))
	return img
}

// DrawRect is the image.Rectangle compatible equivalent of [Renderer.DrawArea]().
// When no rounding is required, prefer [Renderer.DrawIntArea]() instead.
func (r *Renderer) DrawRect(target *ebiten.Image, rect image.Rectangle, rounding float32) {
	r.DrawArea(target, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), rounding)
}

func (r *Renderer) DrawArea(target *ebiten.Image, ox, oy, w, h, rounding float32) {
	ensureShaderRectLoaded()
	r.setFlatCustomVAs(ox, oy, w, h)
	r.opts.Uniforms["Rounding"] = rounding
	r.DrawRectShader(target, ox, oy, w, h, 0, 0, shaderRect)
	clear(r.opts.Uniforms)
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

func (r *Renderer) DrawCircle(target *ebiten.Image, cx, cy, radius float32) {
	r.setDstRectCoords(cx-radius, cy-radius, cx+radius, cy+radius)
	ensureShaderCircleLoaded()
	r.setFlatCustomVAs(cx, cy, radius, 0.0)
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderCircle, &r.opts)
}

func (r *Renderer) DrawRing(target *ebiten.Image, cx, cy, inRadius, outRadius float32) {
	r.setDstRectCoords(cx-outRadius, cy-outRadius, cx+outRadius, cy+outRadius)
	ensureShaderRingLoaded()
	r.setFlatCustomVAs(cx, cy, outRadius, inRadius)
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderRing, &r.opts)
}

// Notice: ellipses don't have a perfect SDF, so approximations can be very slightly
// bigger or smaller than the requested radiuses.
func (r *Renderer) DrawEllipse(target *ebiten.Image, cx, cy, horzRadius, vertRadius float32, rads float64) {
	if rads == 0 {
		r.setDstRectCoords(cx-horzRadius, cy-vertRadius, cx+horzRadius, cy+vertRadius)
		r.opts.Uniforms["Radians"] = 0
	} else {
		hRadiusF64, vRadiusF64 := float64(horzRadius), float64(vertRadius)
		rc, rs := math.Cos(rads), math.Sin(rads)
		halfWidth := float32(math.Hypot(hRadiusF64*rc, vRadiusF64*rs))
		halfHeight := float32(math.Hypot(hRadiusF64*rs, vRadiusF64*rc))
		r.setDstRectCoords(cx-halfWidth, cy-halfHeight, cx+halfWidth, cy+halfHeight)
		r.opts.Uniforms["Radians"] = rads
	}
	r.setFlatCustomVAs(cx, cy, horzRadius, vertRadius)
	ensureShaderEllipseLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderEllipse, &r.opts)
}

// DrawIntRect is the image.Rectangle compatible equivalent of [Renderer.DrawIntArea]().
func (r *Renderer) DrawIntRect(target *ebiten.Image, rect image.Rectangle) {
	r.DrawIntArea(target, rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy())
}

func (r *Renderer) DrawIntArea(target *ebiten.Image, ox, oy, w, h int) {
	bounds := target.Bounds()
	minX, minY := bounds.Min.X, bounds.Min.Y
	r.setDstRectCoords(float32(minX+ox), float32(minY+oy), float32(minX+ox+w), float32(minY+oy+h))
	ensureShaderDefaultLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderDefault, &r.opts)
}

// StrokeIntRect is the image.Rectangle compatible equivalent of [Renderer.StrokeIntArea]().
func (r *Renderer) StrokeIntRect(target *ebiten.Image, area image.Rectangle, outThickness, inThickness int) {
	r.StrokeIntArea(target, area.Min.X, area.Min.Y, area.Dx(), area.Dy(), outThickness, inThickness)
}

func (r *Renderer) StrokeIntArea(target *ebiten.Image, ox, oy, w, h, outThickness, inThickness int) {
	if outThickness < 0 || inThickness < 0 {
		panic("outThickness < 0 || inThickness < 0")
	}

	if outThickness == 0 {
		if inThickness != 0 {
			r.strokeIntInnerArea(target, ox, oy, w, h, inThickness)
		}
	} else {
		r.strokeIntInnerArea(target, ox-outThickness, oy-outThickness, w+outThickness*2, h+outThickness*2, outThickness+inThickness)
	}
}

func (r *Renderer) strokeIntInnerArea(target *ebiten.Image, ox, oy, w, h, thickness int) {
	bounds := target.Bounds()
	minX, minY := bounds.Min.X, bounds.Min.Y
	oox, ooy := float32(minX+ox), float32(minY+oy)
	ofx, ofy := float32(minX+ox+w), float32(minY+oy+h)
	r.setDstRectCoords(oox, ooy, ofx, ofy)

	// add inner points
	thickF32 := float32(thickness)
	iox, ioy := oox+thickF32, ooy+thickF32
	ifx, ify := ofx-thickF32, ofy-thickF32
	if r.singleClr || r.opts.Blend == ebiten.BlendClear {
		r.vertices = append(r.vertices,
			ebiten.Vertex{DstX: iox, DstY: ioy},
			ebiten.Vertex{DstX: ifx, DstY: ioy},
			ebiten.Vertex{DstX: ifx, DstY: ify},
			ebiten.Vertex{DstX: iox, DstY: ify},
		)
		for i := range 4 {
			r.vertices[4+i].ColorR = r.vertices[i].ColorR
			r.vertices[4+i].ColorG = r.vertices[i].ColorG
			r.vertices[4+i].ColorB = r.vertices[i].ColorB
			r.vertices[4+i].ColorA = r.vertices[i].ColorA
		}
	} else {
		// we need to interpolate colors. this code takes advantage of
		// the heavy symmetries in the geometry to reduce the number of
		// operations, but as a downside, it's a bit tricky to understand

		// compute uv coords for inner points
		iou := min(max((iox-oox)/(ofx-oox), 0), 1)
		iov := min(max((ioy-ooy)/(ofy-ooy), 0), 1)

		// compute top and bottom left colors
		tR, tG, tB, tA := interpVertexColor(r.vertices[0], r.vertices[1], iou)
		bR, bG, bB, bA := interpVertexColor(r.vertices[3], r.vertices[2], iou)

		// append all vertices with left side colors set
		r.vertices = append(r.vertices,
			ebiten.Vertex{DstX: iox, DstY: ioy},
			ebiten.Vertex{DstX: ifx, DstY: ioy},
			ebiten.Vertex{DstX: ifx, DstY: ify},
			ebiten.Vertex{DstX: iox, DstY: ify},
		)

		tli, tri, bli, bri := 4, 5, 7, 6 // NOTE: use other orders for cool effects
		r.vertices[tli].ColorR = lerp(tR, bR, iov)
		r.vertices[tli].ColorG = lerp(tG, bG, iov)
		r.vertices[tli].ColorB = lerp(tB, bB, iov)
		r.vertices[tli].ColorA = lerp(tA, bA, iov)

		r.vertices[bli].ColorR = lerp(bR, tR, iov)
		r.vertices[bli].ColorG = lerp(bG, tG, iov)
		r.vertices[bli].ColorB = lerp(bB, tB, iov)
		r.vertices[bli].ColorA = lerp(bA, tA, iov)

		// compute right side colors by symmetry
		tR = r.vertices[1].ColorR - (tR - r.vertices[0].ColorR)
		tG = r.vertices[1].ColorG - (tG - r.vertices[0].ColorG)
		tB = r.vertices[1].ColorB - (tB - r.vertices[0].ColorB)
		tA = r.vertices[1].ColorA - (tA - r.vertices[0].ColorA)
		bR = r.vertices[2].ColorR - (bR - r.vertices[3].ColorR)
		bG = r.vertices[2].ColorG - (bG - r.vertices[3].ColorG)
		bB = r.vertices[2].ColorB - (bB - r.vertices[3].ColorB)
		bA = r.vertices[2].ColorA - (bA - r.vertices[3].ColorA)

		// set right vertex colors
		r.vertices[tri].ColorR = lerp(tR, bR, iov)
		r.vertices[tri].ColorG = lerp(tG, bG, iov)
		r.vertices[tri].ColorB = lerp(tB, bB, iov)
		r.vertices[tri].ColorA = lerp(tA, bA, iov)

		r.vertices[bri].ColorR = lerp(bR, tR, iov)
		r.vertices[bri].ColorG = lerp(bG, tG, iov)
		r.vertices[bri].ColorB = lerp(bB, tB, iov)
		r.vertices[bri].ColorA = lerp(bA, tA, iov)
	}

	ensureShaderDefaultLoaded()
	target.DrawTrianglesShader(r.vertices[:], r.strokeIndices[:], shaderDefault, &r.opts)
	r.vertices = r.vertices[:4]
}

// StrokeRect is the image.Rectangle compatible equivalent of [Renderer.DrawArea]().
// When no rounding is required, prefer [Renderer.DrawIntArea]() instead.
func (r *Renderer) StrokeRect(target *ebiten.Image, rect image.Rectangle, outThickness, inThickness, rounding float32) {
	r.StrokeArea(target, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), outThickness, inThickness, rounding)
}

func (r *Renderer) StrokeArea(target *ebiten.Image, ox, oy, w, h, outThickness, inThickness, rounding float32) {
	if outThickness < 0 || inThickness < 0 {
		panic("outThickness < 0 || inThickness < 0")
	}

	if outThickness == 0 {
		if inThickness != 0 {
			r.strokeInnerArea(target, ox, oy, w, h, inThickness, rounding)
		}
	} else {
		r.strokeInnerArea(target, ox-outThickness, oy-outThickness, w+outThickness*2, h+outThickness*2, outThickness+inThickness, rounding)
	}
}

func (r *Renderer) strokeInnerArea(target *ebiten.Image, ox, oy, w, h, inThickness, rounding float32) {
	ensureShaderStrokeRectLoaded()
	r.setFlatCustomVAs(ox, oy, w, h)
	r.opts.Uniforms["InnerThickness"] = inThickness
	r.opts.Uniforms["Rounding"] = rounding
	r.DrawRectShader(target, ox, oy, w, h, 0, 0, shaderStrokeRect)
	clear(r.opts.Uniforms)
}

// DrawTriangle draws a smooth triangle using the given vertices and an optional rounding factor.
// Notice that, if provided, handling the rounding is relatively expensive (two dozen f64 products
// and 3 square roots).
func (r *Renderer) DrawTriangle(target *ebiten.Image, ox1, oy1, ox2, oy2, ox3, oy3, rounding float64) {
	r.drawTriangle(target, ox1, oy1, ox2, oy2, ox3, oy3, 0.0, rounding)
}

// StrokeTriangle draws an unfilled triangle. The outline will expand [-thickness/2, +thickness/2] around
// the given points, unless the passed thickness is negative, in which case the outline will be interior
// only, going from [-thickness, 0].
//
// For more details on rounding, see [Renderer.DrawTriangle]().
func (r *Renderer) StrokeTriangle(target *ebiten.Image, ox1, oy1, ox2, oy2, ox3, oy3, thickness, rounding float64) {
	r.drawTriangle(target, ox1, oy1, ox2, oy2, ox3, oy3, thickness, rounding)
}

func (r *Renderer) drawTriangle(target *ebiten.Image, ox1, oy1, ox2, oy2, ox3, oy3, thickness, rounding float64) {
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
	hthick := max(thickness/2.0, 0)
	r.setDstRectCoords(float32(minX-hthick), float32(minY-hthick), float32(maxX+hthick), float32(maxY+hthick))

	// draw shader
	ensureShaderTriangleLoaded()
	r.opts.Uniforms["P0"] = []float32{float32(iox1), float32(ioy1)}
	r.opts.Uniforms["P1"] = []float32{float32(iox2), float32(ioy2)}
	r.opts.Uniforms["P2"] = []float32{float32(iox3), float32(ioy3)}
	r.opts.Uniforms["Rounding"] = float32(rounding)
	r.opts.Uniforms["Thickness"] = float32(thickness)
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

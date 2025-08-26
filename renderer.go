package shapes

import (
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

// Renderer is a helper type for basic shape rendering which
// reuses vertices and options for slightly reduced memory usage.
type Renderer struct {
	vertices []ebiten.Vertex
	indices  []uint16
	opts     ebiten.DrawTrianglesShaderOptions

	singleClr     bool
	strokeIndices []uint16

	temps []offscreen
}

func NewRenderer() *Renderer {
	var renderer Renderer
	renderer.vertices = make([]ebiten.Vertex, 4)
	renderer.SetColor(color.RGBA{255, 255, 255, 255})
	renderer.indices = []uint16{0, 1, 2, 0, 2, 3}
	renderer.opts.Uniforms = make(map[string]any, 8)
	renderer.strokeIndices = []uint16{
		0, 1, 4,
		4, 1, 5,
		5, 1, 2,
		5, 2, 6,
		6, 2, 3,
		6, 3, 7,
		7, 3, 0,
		0, 4, 7,
	}
	return &renderer
}

func (r *Renderer) GetColorF32() [4]float32 {
	return [4]float32{r.vertices[0].ColorR, r.vertices[0].ColorG, r.vertices[0].ColorB, r.vertices[0].ColorA}
}

// SetColor sets the color of all vertices, unless vertexIndices are specifically provided, in
// which case only the given indices will be set. In general, most shaders use vertex 0 as top-left,
// vertex 1 as top-right, vertex 2 as bottom-right, vertex 3 as bottom-left, but this is shader
// dependent (or even variable in some cases).
func (r *Renderer) SetColor(clr color.Color, vertexIndices ...int) {
	clrF32 := ColorToF32(clr)
	r.SetColorF32(clrF32[0], clrF32[1], clrF32[2], clrF32[3], vertexIndices...)
}

func (r *Renderer) SetColorF32(red, green, blue, alpha float32, vertexIndices ...int) {
	if len(vertexIndices) == 0 {
		r.singleClr = true
		vertexIndices = []int{0, 1, 2, 3}
	} else {
		r.singleClr = false
	}
	for _, i := range vertexIndices {
		r.vertices[i].ColorR = red
		r.vertices[i].ColorG = green
		r.vertices[i].ColorB = blue
		r.vertices[i].ColorA = alpha
	}
}

func (r *Renderer) ScaleAlphaBy(alphaFactor float32) {
	for i := range r.vertices {
		r.vertices[i].ColorR *= alphaFactor
		r.vertices[i].ColorG *= alphaFactor
		r.vertices[i].ColorB *= alphaFactor
		r.vertices[i].ColorA *= alphaFactor
	}
}

func (r *Renderer) SetBlend(blend ebiten.Blend) {
	r.opts.Blend = blend
}

func (r *Renderer) Options() *ebiten.DrawTrianglesShaderOptions {
	return &r.opts
}

func (r *Renderer) DrawShaderAt(target, source *ebiten.Image, ox, oy, horzMargin, vertMargin float32, shader *ebiten.Shader) {
	srcOX, srcOY, srcWidthF32, srcHeightF32 := rectOriginSizeF32(source.Bounds())
	dstOX, dstOY := rectOriginF32(target.Bounds())
	dstOX, dstOY = dstOX+ox, dstOY+oy
	r.setDstRectCoords(dstOX-horzMargin, dstOY-vertMargin, dstOX+srcWidthF32+horzMargin, dstOY+srcHeightF32+vertMargin)
	r.setSrcRectCoords(srcOX-horzMargin, srcOY-vertMargin, srcOX+srcWidthF32+horzMargin, srcOY+srcHeightF32+vertMargin)

	r.opts.Images[0] = source
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shader, &r.opts)
	r.opts.Images[0] = nil
}

func (r *Renderer) DrawRectShader(target *ebiten.Image, ox, oy, w, h, horzMargin, vertMargin float32, shader *ebiten.Shader) {
	dstOX, dstOY := rectOriginF32(target.Bounds())
	dstOX, dstOY = dstOX+ox, dstOY+oy
	r.setDstRectCoords(dstOX-horzMargin, dstOY-vertMargin, dstOX+w+horzMargin, dstOY+h+vertMargin)
	r.setSrcRectCoords(-horzMargin, -vertMargin, w+horzMargin, h+vertMargin)
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shader, &r.opts)
}

func (r *Renderer) DrawShader(target *ebiten.Image, horzMargin, vertMargin float32, shader *ebiten.Shader) {
	bounds := target.Bounds()
	r.DrawRectShader(target, 0, 0, float32(bounds.Dx()), float32(bounds.Dy()), horzMargin, vertMargin, shader)
}

// Scale draws the source into the given target with two differences from Ebitengine's scaling:
//   - scaledSampling can be set to true to mimic Ebitengine's v2.9.0 FilterPixelated.
//   - Subimages can be scaled without bleeding edges, as the shader uses clamping.
func (r *Renderer) Scale(target, source *ebiten.Image, ox, oy, scale float32, scaledSampling bool) {
	srcBounds := source.Bounds()
	srcWidth, srcHeight := srcBounds.Dx(), srcBounds.Dy()
	srcWidthF32, srcHeightF32 := float32(srcWidth), float32(srcHeight)

	dstBounds := target.Bounds()
	minX := float32(dstBounds.Min.X) + ox
	minY := float32(dstBounds.Min.Y) + oy
	r.setDstRectCoords(minX, minY, minX+srcWidthF32*scale, minY+srcHeightF32*scale)

	minX, minY = float32(srcBounds.Min.X), float32(srcBounds.Min.Y)
	r.setSrcRectCoords(minX, minY, minX+srcWidthF32, minY+srcHeightF32)
	r.opts.Images[0] = source

	ensureShaderBilinearLoaded()
	if scaledSampling {
		r.setFlatCustomVAs(1.0/scale, 1.0/scale, 0, 0)
	} else {
		r.setFlatCustomVAs(1.0, 1.0, 0, 0)
	}
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shaderBilinear, &r.opts)
	r.opts.Images[0] = nil
}

// UnsafeTemp allows requesting offscreens to the renderer. These offscreens might have
// already been created while the renderer was doing complex operations, so reusing them
// can prevent the creation of additional offscreens.
//
// The offscreens returned by this function should only be used for local operations, and
// the offscreen must not be stored. Any renderer function documented to use an internal
// offscreen can panic or fail in any other way if an offscreen returned by this function
// if passed as an input parameter.
func (r *Renderer) UnsafeTemp(offscreenIndex int, w, h int) *ebiten.Image {
	return r.getTemp(offscreenIndex, w, h, false)
}

// UnsafeTempClear behaves like UnsafeTemp but returns a cleared image, including 1
// extra pixel of clear padding to prevent problems with bleeding edges.
func (r *Renderer) UnsafeTempClear(offscreenIndex int, w, h int) *ebiten.Image {
	return r.getTemp(offscreenIndex, w, h, true)
}

// UnsafeTempCopy calls [Renderer.UnsafeTemp]() and copies the contents of source into
// the returned offscreen. See safety warnings and docs for UnsafeTemp. The 'clear'
// argument allows specifying whether a 1 pixel clear margin is required or not.
func (r *Renderer) UnsafeTempCopy(offscreenIndex int, source *ebiten.Image, clear bool) *ebiten.Image {
	bounds := source.Bounds()
	temp := r.getTemp(offscreenIndex, bounds.Dx(), bounds.Dy(), clear)
	var opts ebiten.DrawImageOptions
	opts.Blend = ebiten.BlendCopy
	temp.DrawImage(source, &opts)
	return temp
}

func (r *Renderer) setDstRectCoords(minX, minY, maxX, maxY float32) {
	r.vertices[0].DstX = minX
	r.vertices[0].DstY = minY
	r.vertices[1].DstX = maxX
	r.vertices[1].DstY = minY
	r.vertices[2].DstX = maxX
	r.vertices[2].DstY = maxY
	r.vertices[3].DstX = minX
	r.vertices[3].DstY = maxY
}

func (r *Renderer) setSrcRectCoords(minX, minY, maxX, maxY float32) {
	r.vertices[0].SrcX = minX
	r.vertices[0].SrcY = minY
	r.vertices[1].SrcX = maxX
	r.vertices[1].SrcY = minY
	r.vertices[2].SrcX = maxX
	r.vertices[2].SrcY = maxY
	r.vertices[3].SrcX = minX
	r.vertices[3].SrcY = maxY
}

func (r *Renderer) setFlatCustomVAs(cva0, cva1, cva2, cva3 float32) {
	for i := range len(r.vertices) {
		r.vertices[i].Custom0 = cva0
		r.vertices[i].Custom1 = cva1
		r.vertices[i].Custom2 = cva2
		r.vertices[i].Custom3 = cva3
	}
}

func (r *Renderer) setFlatCustomVA0(cva0 float32) {
	for i := range len(r.vertices) {
		r.vertices[i].Custom0 = cva0
	}
}

func (r *Renderer) setFlatCustomVAs01(cva0, cva1 float32) {
	for i := range len(r.vertices) {
		r.vertices[i].Custom0 = cva0
		r.vertices[i].Custom1 = cva1
	}
}

// SetCustomVAs configures up to 4 custom vertex attributes.
func (r *Renderer) SetCustomVAs(vas ...float32) {
	switch len(vas) {
	case 0:
		// nothing
	case 1:
		r.setFlatCustomVA0(vas[0])
	case 2:
		r.setFlatCustomVAs01(vas[0], vas[1])
	case 3:
		r.setFlatCustomVAs(vas[0], vas[1], vas[2], 0.0)
	case 4:
		r.setFlatCustomVAs(vas[0], vas[1], vas[2], vas[3])
	default:
		panic("only up to 4 custom VAs allowed")
	}
}

func (r *Renderer) getTemp(offscreenIndex int, w, h int, clear bool) *ebiten.Image {
	if offscreenIndex >= len(r.temps) {
		growth := offscreenIndex + 1 - len(r.temps)
		r.temps = slices.Grow(r.temps, growth)
		r.temps = r.temps[:offscreenIndex+1]
		r.temps[offscreenIndex] = newOffscreen(0, 0, 64)
		clear = false // we have already cleared, skip requirement
	}
	return r.temps[offscreenIndex].WithSize(w, h, clear)
}

package shapes

import (
	"image"
	"image/color"

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

	tmp       *ebiten.Image
	tmpParent *ebiten.Image
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
	srcBounds := source.Bounds()
	srcWidth, srcHeight := srcBounds.Dx(), srcBounds.Dy()
	srcWidthF32, srcHeightF32 := float32(srcWidth), float32(srcHeight)
	dstBounds := target.Bounds()
	minX := float32(dstBounds.Min.X) + ox - horzMargin
	minY := float32(dstBounds.Min.Y) + oy - vertMargin
	r.setDstRectCoords(minX, minY, minX+srcWidthF32+horzMargin*2, minY+srcHeightF32+vertMargin*2)
	minX = float32(srcBounds.Min.X) - horzMargin
	minY = float32(srcBounds.Min.Y) - vertMargin
	r.setSrcRectCoords(minX, minY, minX+srcWidthF32+horzMargin*2, minY+srcHeightF32+vertMargin*2)
	r.opts.Images[0] = source
	target.DrawTrianglesShader(r.vertices[:], r.indices[:], shader, &r.opts)
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

func (r *Renderer) ensureOffscreenSize(w, h int) {
	const ExtraMargin = 64
	if r.tmp == nil {
		r.tmpParent = ebiten.NewImage(w+ExtraMargin, h+ExtraMargin)
		r.tmp = r.tmpParent.SubImage(image.Rect(0, 0, w, h)).(*ebiten.Image)
		return
	}

	bounds := r.tmp.Bounds()
	if bounds.Dx() == w && bounds.Dy() == h {
		return
	}

	bounds = r.tmpParent.Bounds()
	currWidth, currHeight := bounds.Dx(), bounds.Dy()
	if currWidth >= w && currHeight >= h {
		r.tmp = r.tmpParent.SubImage(image.Rect(0, 0, w, h)).(*ebiten.Image)
	} else {
		r.tmpParent = ebiten.NewImage(w+ExtraMargin, h+ExtraMargin)
		r.tmp = r.tmpParent.SubImage(image.Rect(0, 0, w, h)).(*ebiten.Image)
	}
}

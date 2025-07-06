package previewer

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/Advik-B/Axon/pkg/axon"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	colorBg         = color.RGBA{R: 24, G: 25, B: 26, A: 255}
	colorNodeBorder = color.RGBA{R: 80, G: 85, B: 90, A: 255}
	colorText       = color.White
	colorTextDim    = color.Gray{Y: 150}
	colorTextImpl   = color.RGBA{R: 156, G: 163, B: 175, A: 255}
	colorExecEdge   = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	colorDataEdge   = color.RGBA{R: 76, G: 151, B: 255, A: 255}
	colorPort       = color.RGBA{R: 150, G: 155, B: 160, A: 255}
	nodeColors      = map[axon.NodeType]color.Color{
		axon.NodeType_START:      color.RGBA{R: 16, G: 185, B: 129, A: 255},
		axon.NodeType_END:        color.RGBA{R: 239, G: 68, B: 68, A: 255},
		axon.NodeType_RETURN:     color.RGBA{R: 217, G: 70, B: 239, A: 255},
		axon.NodeType_CONSTANT:   color.RGBA{R: 59, G: 130, B: 246, A: 255},
		axon.NodeType_FUNCTION:   color.RGBA{R: 99, G: 102, B: 241, A: 255},
		axon.NodeType_OPERATOR:   color.RGBA{R: 249, G: 115, B: 22, A: 255},
		axon.NodeType_IGNORE:     color.RGBA{R: 236, G: 72, B: 153, A: 255},
		axon.NodeType_STRUCT_DEF: color.RGBA{R: 14, G: 165, B: 233, A: 255},
		axon.NodeType_FUNC_DEF:   color.RGBA{R: 34, G: 197, B: 94, A: 255},
	}
	defaultNodeColor = color.RGBA{R: 45, G: 48, B: 51, A: 255}
)

var (
	whitePixelOnce sync.Once
	whitePixel     *ebiten.Image
)

// getWhitePixel provides a standard 1x1 white image, which is required as a texture source
// for drawing solid-colored vector graphics.
func getWhitePixel() *ebiten.Image {
	whitePixelOnce.Do(func() {
		whitePixel = ebiten.NewImage(1, 1)
		whitePixel.Fill(color.White)
	})
	return whitePixel
}

func drawNode(screen *ebiten.Image, node *LayoutNode, face text.Face, op *ebiten.DrawImageOptions) {
	tx, ty := op.GeoM.Apply(float64(node.Rect.Min.X), float64(node.Rect.Min.Y))
	x, y := float32(tx), float32(ty)
	zoom := float32(op.GeoM.Element(0, 0))
	w, h := float32(node.Rect.Dx())*zoom, float32(node.Rect.Dy())*zoom

	nodeColor, ok := nodeColors[node.Type]
	if !ok {
		nodeColor = defaultNodeColor
	}
	vector.DrawFilledRect(screen, x, y, w, h, nodeColor, false)
	vector.StrokeRect(screen, x, y, w, h, 1, colorNodeBorder, false)

	textOp := &text.DrawOptions{}
	textOp.ColorScale.ScaleWithColor(colorText)
	textOp.GeoM.Translate(float64(x+10), float64(y+20))
	text.Draw(screen, node.Label, face, textOp)

	textOp.GeoM.Reset()
	textOp.ColorScale.Reset()
	textOp.ColorScale.ScaleWithColor(colorTextDim)
	textOp.GeoM.Translate(float64(x+10), float64(y+40))
	text.Draw(screen, node.Type.String(), face, textOp)

	if node.Type == axon.NodeType_FUNCTION && node.ImplReference != "" {
		textOp.GeoM.Reset()
		textOp.ColorScale.Reset()
		textOp.ColorScale.ScaleWithColor(colorTextImpl)
		textOp.GeoM.Translate(float64(x+10), float64(y+65))
		text.Draw(screen, node.ImplReference, face, textOp)
	}

	for _, p := range node.InputPorts {
		ptx, pty := op.GeoM.Apply(float64(p.X), float64(p.Y))
		vector.DrawFilledCircle(screen, float32(ptx), float32(pty), portRadius*zoom, colorPort, false)
	}
	for _, p := range node.OutputPorts {
		ptx, pty := op.GeoM.Apply(float64(p.X), float64(p.Y))
		vector.DrawFilledCircle(screen, float32(ptx), float32(pty), portRadius*zoom, colorPort, false)
	}
}

func drawBezierCurve(screen *ebiten.Image, p0, p3 image.Point, clr color.Color, op *ebiten.DrawImageOptions) {
	dx := p3.X - p0.X
	dy := p3.Y - p0.Y
	var p1, p2 image.Point
	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		p1 = image.Pt(p0.X+dx/3, p0.Y)
		p2 = image.Pt(p3.X-dx/3, p3.Y)
	} else {
		p1 = image.Pt(p0.X, p0.Y+dy/3)
		p2 = image.Pt(p3.X, p3.Y-dy/3)
	}

	var path vector.Path
	v0x, v0y := op.GeoM.Apply(float64(p0.X), float64(p0.Y))
	v1x, v1y := op.GeoM.Apply(float64(p1.X), float64(p1.Y))
	v2x, v2y := op.GeoM.Apply(float64(p2.X), float64(p2.Y))
	v3x, v3y := op.GeoM.Apply(float64(p3.X), float64(p3.Y))

	path.MoveTo(float32(v0x), float32(v0y))
	path.CubicTo(float32(v1x), float32(v1y), float32(v2x), float32(v2y), float32(v3x), float32(v3y))

	strokeOp := &vector.StrokeOptions{Width: 1.5 * float32(op.GeoM.Element(0, 0))}
	vertices, indices := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)

	// **THE FIX**: Manually set the color for each vertex. This is the correct modern API.
	// The options struct for DrawTriangles does not have a color field.
	r, g, b, a := clr.RGBA()
	cr := float32(r) / 0xffff
	cg := float32(g) / 0xffff
	cb := float32(b) / 0xffff
	ca := float32(a) / 0xffff
	for i := range vertices {
		vertices[i].ColorR = cr
		vertices[i].ColorG = cg
		vertices[i].ColorB = cb
		vertices[i].ColorA = ca
	}

	screen.DrawTriangles(vertices, indices, getWhitePixel(), &ebiten.DrawTrianglesOptions{})
}
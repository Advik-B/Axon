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

// --- AESTHETICS & COLOR PALETTE (Inspired by Unreal Engine) ---
var (
	colorBg      = color.RGBA{R: 24, G: 25, B: 26, A: 255}
	colorGrid    = color.RGBA{R: 40, G: 42, B: 44, A: 255}
	colorGridSub = color.RGBA{R: 32, G: 34, B: 36, A: 255}

	colorNodeBody            = color.RGBA{R: 35, G: 38, B: 41, A: 230}
	colorNodeShadow          = color.RGBA{R: 0, G: 0, B: 0, A: 100}
	colorNodeBorder          = color.RGBA{R: 10, G: 10, B: 10, A: 255}
	colorText                = color.White
	colorTextDim             = color.Gray{Y: 180}
	colorTextImpl            = color.RGBA{R: 156, G: 163, B: 175, A: 255}
	nodeHeaderHeight float32 = 30.0
	nodeCornerRadius float32 = 8.0
	nodeShadowOffset float32 = 5.0

	colorExec      = color.White
	colorPortLabel = color.Gray{Y: 200}
	dataTypeColors = map[string]color.Color{
		"int":     color.RGBA{R: 0, G: 184, B: 212, A: 255},
		"string":  color.RGBA{R: 217, G: 70, B: 239, A: 255},
		"bool":    color.RGBA{R: 220, G: 38, B: 38, A: 255},
		"[]byte":  color.RGBA{R: 132, G: 204, B: 22, A: 255},
		"error":   color.RGBA{R: 245, G: 158, B: 11, A: 255},
		"float":   color.RGBA{R: 52, G: 211, B: 153, A: 255},
		"default": color.RGBA{R: 139, G: 92, B: 246, A: 255},
	}
	nodeColors = map[axon.NodeType]color.Color{
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

	// **NEW**: Map for node type titles
	nodeTypeTitles = map[axon.NodeType]string{
		axon.NodeType_FUNCTION:   "FUNCTION CALL",
		axon.NodeType_CONSTANT:   "CONSTANT",
		axon.NodeType_OPERATOR:   "OPERATOR",
		axon.NodeType_STRUCT_DEF: "STRUCT DEFINITION",
		axon.NodeType_FUNC_DEF:   "FUNCTION DEFINITION",
	}
)

var (
	whitePixelOnce sync.Once
	whitePixel     *ebiten.Image
)

func getWhitePixel() *ebiten.Image {
	whitePixelOnce.Do(func() {
		whitePixel = ebiten.NewImage(1, 1)
		whitePixel.Fill(color.White)
	})
	return whitePixel
}

// drawNode renders a single, beautifully styled node.
func drawNode(screen *ebiten.Image, node *LayoutNode, face, smallFace text.Face, op *ebiten.DrawImageOptions) {
	tx, ty := op.GeoM.Apply(float64(node.Rect.Min.X), float64(node.Rect.Min.Y))
	x, y := float32(tx), float32(ty)
	zoom := float32(op.GeoM.Element(0, 0))
	w, h := float32(node.Rect.Dx())*zoom, float32(node.Rect.Dy())*zoom
	radius := nodeCornerRadius * zoom

	drawFilledRoundRect(screen, x+nodeShadowOffset, y+nodeShadowOffset, w, h, radius, colorNodeShadow)
	drawFilledRoundRect(screen, x, y, w, h, radius, colorNodeBody)

	nodeColor, ok := nodeColors[node.Type]
	if !ok {
		nodeColor = defaultNodeColor
	}
	drawFilledRoundRect(screen, x, y, w, nodeHeaderHeight*zoom, radius, nodeColor)
	vector.DrawFilledRect(screen, x, y+(nodeHeaderHeight*zoom-radius), w, radius, nodeColor, false)

	strokeRoundRect(screen, x, y, w, h, radius, 1, colorNodeBorder)

	// --- Text Drawing ---
	titleOp := &text.DrawOptions{}
	titleOp.GeoM.Translate(float64(x+10), float64(y+22))
	titleOp.ColorScale.ScaleWithColor(colorText)
	text.Draw(screen, node.Label, face, titleOp)

	// **NEW**: Draw the small node type title on the header
	if typeTitle, ok := nodeTypeTitles[node.Type]; ok {
		typeOp := &text.DrawOptions{}
		typeAdvance, _ := text.Measure(typeTitle, smallFace, 0)
		typeOp.GeoM.Translate(float64(x+w)-typeAdvance-10, float64(y+22))
		typeOp.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 100}) // Semi-transparent white
		text.Draw(screen, typeTitle, smallFace, typeOp)
	}

	if node.Type == axon.NodeType_FUNCTION && node.ImplReference != "" {
		implOp := &text.DrawOptions{}
		implOp.GeoM.Translate(float64(x+10), float64(y+nodeHeaderHeight*zoom+20))
		implOp.ColorScale.ScaleWithColor(colorTextImpl)
		text.Draw(screen, node.ImplReference, smallFace, implOp)
	}

	drawPorts(screen, node, smallFace, op)
}

func drawPorts(screen *ebiten.Image, node *LayoutNode, face text.Face, op *ebiten.DrawImageOptions) {
	if p, ok := node.InputPorts["exec_in"]; ok {
		drawExecPin(screen, p, false, colorExec, op)
	}
	if p, ok := node.OutputPorts["exec_out"]; ok {
		drawExecPin(screen, p, true, colorExec, op)
	}

	for name, p := range node.InputPorts {
		if name == "exec_in" {
			continue
		}
		var portType string
		for _, portDef := range node.Inputs {
			if portDef.Name == name {
				portType = portDef.TypeName
				break
			}
		}
		drawDataPin(screen, p, portType, name, false, face, op)
	}
	for name, p := range node.OutputPorts {
		if name == "exec_out" {
			continue
		}
		var portType string
		for _, portDef := range node.Outputs {
			if portDef.Name == name {
				portType = portDef.TypeName
				break
			}
		}
		drawDataPin(screen, p, portType, name, true, face, op)
	}
}

func drawExecPin(screen *ebiten.Image, p image.Point, isOutput bool, clr color.Color, op *ebiten.DrawImageOptions) {
	zoom := float32(op.GeoM.Element(0, 0))
	size := 8 * zoom
	var path vector.Path
	tx, ty := op.GeoM.Apply(float64(p.X), float64(p.Y))
	if isOutput {
		path.MoveTo(float32(tx), float32(ty)-size/2)
		path.LineTo(float32(tx)+size, float32(ty))
		path.LineTo(float32(tx), float32(ty)+size/2)
		path.LineTo(float32(tx), float32(ty)-size/2)
	} else {
		path.MoveTo(float32(tx), float32(ty)-size/2)
		path.LineTo(float32(tx)-size, float32(ty))
		path.LineTo(float32(tx), float32(ty)+size/2)
		path.LineTo(float32(tx), float32(ty)-size/2)
	}
	vertices, indices := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{Width: 2 * zoom})
	colorVerts(vertices, clr)
	screen.DrawTriangles(vertices, indices, getWhitePixel(), &ebiten.DrawTrianglesOptions{})
}

func drawDataPin(screen *ebiten.Image, p image.Point, typeName, label string, isOutput bool, face text.Face, op *ebiten.DrawImageOptions) {
	zoom := float32(op.GeoM.Element(0, 0))
	clr, ok := dataTypeColors[typeName]
	if !ok {
		clr = dataTypeColors["default"]
	}
	tx, ty := op.GeoM.Apply(float64(p.X), float64(p.Y))
	vector.DrawFilledCircle(screen, float32(tx), float32(ty), (portRadius+1)*zoom, color.Black, false)
	vector.DrawFilledCircle(screen, float32(tx), float32(ty), portRadius*zoom, clr, false)

	labelOp := &text.DrawOptions{}
	labelOp.ColorScale.ScaleWithColor(colorPortLabel)

	// **THE FIX**: Measure the text's advance (width) and a capital letter's height for centering.
	advance, _ := text.Measure(label, face, 0)
	_, capHeight := text.Measure("A", face, 0) // Use height of a capital letter for vertical alignment

	// Calculate the correct Y offset to center the text baseline with the circle's center.
	yOffset := float64(ty) + capHeight/2

	if isOutput {
		labelOp.GeoM.Translate(float64(tx)-advance-float64(portRadius*zoom+5), yOffset)
	} else {
		labelOp.GeoM.Translate(float64(tx)+float64(portRadius*zoom+5), yOffset)
	}
	text.Draw(screen, label, face, labelOp)
}

func drawBezierCurve(screen *ebiten.Image, p0, p3 image.Point, clr color.Color, op *ebiten.DrawImageOptions) {
	dx, dy := p3.X-p0.X, p3.Y-p0.Y
	var p1, p2 image.Point
	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		p1, p2 = image.Pt(p0.X+dx/2, p0.Y), image.Pt(p3.X-dx/2, p3.Y)
	} else {
		p1, p2 = image.Pt(p0.X, p0.Y+dy/2), image.Pt(p3.X, p3.Y-dy/2)
	}

	var path vector.Path
	v0x, v0y := op.GeoM.Apply(float64(p0.X), float64(p0.Y))
	v1x, v1y := op.GeoM.Apply(float64(p1.X), float64(p1.Y))
	v2x, v2y := op.GeoM.Apply(float64(p2.X), float64(p2.Y))
	v3x, v3y := op.GeoM.Apply(float64(p3.X), float64(p3.Y))
	path.MoveTo(float32(v0x), float32(v0y))
	path.CubicTo(float32(v1x), float32(v1y), float32(v2x), float32(v2y), float32(v3x), float32(v3y))

	casingOp := &vector.StrokeOptions{Width: 5 * float32(op.GeoM.Element(0, 0))}
	casingVerts, casingIndices := path.AppendVerticesAndIndicesForStroke(nil, nil, casingOp)
	colorVerts(casingVerts, color.Black)
	screen.DrawTriangles(casingVerts, casingIndices, getWhitePixel(), &ebiten.DrawTrianglesOptions{})

	strokeOp := &vector.StrokeOptions{Width: 2.5 * float32(op.GeoM.Element(0, 0))}
	verts, indices := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)
	colorVerts(verts, clr)
	screen.DrawTriangles(verts, indices, getWhitePixel(), &ebiten.DrawTrianglesOptions{})
}

func colorVerts(v []ebiten.Vertex, clr color.Color) {
	r, g, b, a := clr.RGBA()
	cr, cg, cb, ca := float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff
	for i := range v {
		v[i].ColorR, v[i].ColorG, v[i].ColorB, v[i].ColorA = cr, cg, cb, ca
	}
}

func createRoundedRectPath(x, y, w, h, r float32) *vector.Path {
	var path vector.Path
	path.MoveTo(x+r, y)
	path.LineTo(x+w-r, y)
	path.QuadTo(x+w, y, x+w, y+r)
	path.LineTo(x+w, y+h-r)
	path.QuadTo(x+w, y+h, x+w-r, y+h)
	path.LineTo(x+r, y+h)
	path.QuadTo(x, y+h, x, y+h-r)
	path.LineTo(x, y+r)
	path.QuadTo(x, y, x+r, y)
	return &path
}

func drawFilledRoundRect(screen *ebiten.Image, x, y, w, h, r float32, clr color.Color) {
	path := createRoundedRectPath(x, y, w, h, r)
	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)
	colorVerts(vertices, clr)
	screen.DrawTriangles(vertices, indices, getWhitePixel(), &ebiten.DrawTrianglesOptions{})
}

func strokeRoundRect(screen *ebiten.Image, x, y, w, h, r, strokeWidth float32, clr color.Color) {
	path := createRoundedRectPath(x, y, w, h, r)
	strokeOp := &vector.StrokeOptions{Width: strokeWidth}
	vertices, indices := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)
	colorVerts(vertices, clr)
	screen.DrawTriangles(vertices, indices, getWhitePixel(), &ebiten.DrawTrianglesOptions{})
}

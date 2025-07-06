package previewer

import (
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	colorBg         = color.RGBA{R: 24, G: 25, B: 26, A: 255}
	colorNodeBg     = color.RGBA{R: 45, G: 48, B: 51, A: 255}
	colorNodeBorder = color.RGBA{R: 80, G: 85, B: 90, A: 255}
	colorText       = color.White
	colorTextDim    = color.Gray{Y: 150}
	colorExecEdge   = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	colorDataEdge   = color.RGBA{R: 76, G: 151, B: 255, A: 255}
	colorPort       = color.RGBA{R: 150, G: 155, B: 160, A: 255}
)

var (
	whitePixelOnce sync.Once
	whitePixel *ebiten.Image
)


func getWhitePixel() *ebiten.Image {
    whitePixelOnce.Do(func() {
        whitePixel = ebiten.NewImage(1, 1)
        whitePixel.Fill(color.White)
    })
    return whitePixel
}

// drawNode renders a single node, its text, and its ports to the screen.
func drawNode(screen *ebiten.Image, node *LayoutNode, face text.Face, op *ebiten.DrawImageOptions) {
	tx, ty := op.GeoM.Apply(float64(node.Rect.Min.X), float64(node.Rect.Min.Y))
	x, y := float32(tx), float32(ty)
	zoom := float32(op.GeoM.Element(0, 0))
	w, h := float32(node.Rect.Dx())*zoom, float32(node.Rect.Dy())*zoom

	vector.DrawFilledRect(screen, x, y, w, h, colorNodeBg, false)
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

	for _, p := range node.InputPorts {
		tx, ty := op.GeoM.Apply(float64(p.X), float64(p.Y))
		vector.DrawFilledCircle(screen, float32(tx), float32(ty), portRadius*zoom, colorPort, false)
	}
	for _, p := range node.OutputPorts {
		tx, ty := op.GeoM.Apply(float64(p.X), float64(p.Y))
		vector.DrawFilledCircle(screen, float32(tx), float32(ty), portRadius*zoom, colorPort, false)
	}
}

func drawBezierCurve(screen *ebiten.Image, p0, p3 image.Point, clr color.Color, op *ebiten.DrawImageOptions) {
    p1 := image.Pt(p0.X+hSpacing/2, p0.Y)
    p2 := image.Pt(p3.X-hSpacing/2, p3.Y)

    var path vector.Path
    v0x, v0y := op.GeoM.Apply(float64(p0.X), float64(p0.Y))
    v1x, v1y := op.GeoM.Apply(float64(p1.X), float64(p1.Y))
    v2x, v2y := op.GeoM.Apply(float64(p2.X), float64(p2.Y))
    v3x, v3y := op.GeoM.Apply(float64(p3.X), float64(p3.Y))

    path.MoveTo(float32(v0x), float32(v0y))
    path.CubicTo(float32(v1x), float32(v1y), float32(v2x), float32(v2y), float32(v3x), float32(v3y))

    strokeOp := &vector.StrokeOptions{
        Width: 1.5 * float32(op.GeoM.Element(0, 0)), // Adjust based on zoom
    }

    vertices, indices := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)

    // Convert color.Color to linear RGBA
    r, g, b, a := clr.RGBA()
    for i := range vertices {
        vertices[i].ColorR = float32(r) / 0xffff
        vertices[i].ColorG = float32(g) / 0xffff
        vertices[i].ColorB = float32(b) / 0xffff
        vertices[i].ColorA = float32(a) / 0xffff
    }

    white := getWhitePixel()
    screen.DrawTriangles(vertices, indices, white, &ebiten.DrawTrianglesOptions{})
}

package renderer

import (
	"embed"
	"fmt"
	"image"
	"image/color"

	"github.com/Advik-B/Axon/pkg/axon"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

//go:embed Go-Mono.ttf
var embeddedFont embed.FS

// --- Constants for Styling ---
const (
	Padding      = 40.0
	NodeHeader   = 30.0
	PortRadius   = 6.0
	PortSpacing  = 25.0
	CurveOffset  = 80.0
	DefaultNodeW = 200.0
	DefaultNodeH = 100.0
)

var (
	ColorBg         = color.RGBA{R: 45, G: 52, B: 54, A: 255}
	ColorNodeBg     = color.RGBA{R: 99, G: 110, B: 114, A: 255}
	ColorNodeStroke = color.RGBA{R: 178, G: 190, B: 195, A: 255}
	ColorFont       = color.White
	ColorDataEdge   = color.RGBA{R: 223, G: 230, B: 233, A: 255}
	ColorExecEdge   = color.RGBA{R: 116, G: 185, B: 255, A: 255}
)

type nodeLayout struct {
	*axon.Node
	x, y, w, h float64
}

// GenerateImage renders an Axon graph to an in-memory image.
func GenerateImage(graph *axon.Graph) (image.Image, error) {
	// Load font from embedded assets
	fontBytes, err := embeddedFont.ReadFile("Go-Mono.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded font: %w", err)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	layouts := make(map[string]*nodeLayout)
	var maxX, maxY float64

	// First pass: determine canvas size and store layouts
	for _, node := range graph.Nodes {
		w, h := DefaultNodeW, DefaultNodeH
		if node.VisualInfo != nil {
			if node.VisualInfo.Width > 0 {
				w = float64(node.VisualInfo.Width)
			}
			if node.VisualInfo.Height > 0 {
				h = float64(node.VisualInfo.Height)
			}
		}

		layout := &nodeLayout{
			Node: node,
			x:    float64(node.VisualInfo.GetX()),
			y:    float64(node.VisualInfo.GetY()),
			w:    w,
			h:    h,
		}
		layouts[node.Id] = layout

		if layout.x+layout.w > maxX {
			maxX = layout.x + layout.w
		}
		if layout.y+layout.h > maxY {
			maxY = layout.y + layout.h
		}
	}

	// Create drawing context
	dc := gg.NewContext(int(maxX+Padding), int(maxY+Padding))
	dc.SetColor(ColorBg)
	dc.Clear()
	dc.SetColor(ColorFont)

	// Load font face
	face := truetype.NewFace(font, &truetype.Options{Size: 14})
	dc.SetFontFace(face)

	// Second pass: draw everything
	drawEdges(dc, graph, layouts)
	drawNodes(dc, layouts)

	return dc.Image(), nil
}

func drawEdges(dc *gg.Context, graph *axon.Graph, layouts map[string]*nodeLayout) {
	// Draw Data Edges
	dc.SetColor(ColorDataEdge)
	dc.SetLineWidth(2)
	for _, edge := range graph.DataEdges {
		fromLayout, toLayout := layouts[edge.FromNodeId], layouts[edge.ToNodeId]
		if fromLayout == nil || toLayout == nil {
			continue
		}
		fromY := fromLayout.y + NodeHeader + PortSpacing/2 + float64(getPortIndex(fromLayout.Outputs, edge.FromPort))*PortSpacing
		toY := toLayout.y + NodeHeader + PortSpacing/2 + float64(getPortIndex(toLayout.Inputs, edge.ToPort))*PortSpacing

		fromX, toX := fromLayout.x+fromLayout.w, toLayout.x

		// FIX: Use MoveTo and CubicTo, not the non-existent DrawCubic.
		dc.MoveTo(fromX, fromY)
		dc.CubicTo(fromX+CurveOffset, fromY, toX-CurveOffset, toY, toX, toY)
		dc.Stroke()
	}

	// Draw Exec Edges
	dc.SetColor(ColorExecEdge)
	dc.SetLineWidth(3)
	for _, edge := range graph.ExecEdges {
		fromLayout, toLayout := layouts[edge.FromNodeId], layouts[edge.ToNodeId]
		if fromLayout == nil || toLayout == nil {
			continue
		}

		// FIX: Correctly calculate and use all coordinates to draw a curve
		// from the bottom-center of the source node to the top-center of the target.
		fromX := fromLayout.x + fromLayout.w/2
		fromY := fromLayout.y + fromLayout.h
		toX := toLayout.x + toLayout.w/2
		toY := toLayout.y

		// Use a quadratic curve for a simple arc.
		controlX := (fromX + toX) / 2
		controlY := fromY + CurveOffset/2 // A gentler curve than data edges

		dc.MoveTo(fromX, fromY)
		dc.QuadraticTo(controlX, controlY, toX, toY)
		dc.Stroke()
	}
}

// ... (drawNodes, drawPorts, getPortIndex, and nodeTypeColors are unchanged) ...
func drawNodes(dc *gg.Context, layouts map[string]*nodeLayout) {
	for _, layout := range layouts {
		// Node Body
		dc.SetColor(ColorNodeBg)
		dc.DrawRoundedRectangle(layout.x, layout.y, layout.w, layout.h, 10)
		dc.Fill()

		// Node Header
		headerColor, ok := nodeTypeColors[layout.Type]
		if !ok {
			headerColor = color.Gray{Y: 100}
		}
		dc.SetColor(headerColor)
		dc.DrawRectangle(layout.x, layout.y, layout.w, NodeHeader)
		dc.Fill()

		// Stroke
		dc.SetColor(ColorNodeStroke)
		dc.SetLineWidth(1)
		dc.DrawRoundedRectangle(layout.x, layout.y, layout.w, layout.h, 10)
		dc.Stroke()

		// Node Label
		dc.SetColor(ColorFont)
		dc.DrawStringAnchored(layout.Label, layout.x+layout.w/2, layout.y+NodeHeader/2, 0.5, 0.5)

		// Draw Ports
		drawPorts(dc, layout.Inputs, layout, true)
		drawPorts(dc, layout.Outputs, layout, false)
	}
}

func drawPorts(dc *gg.Context, ports []*axon.Port, layout *nodeLayout, isInput bool) {
	var x float64
	var anchor float64
	if isInput {
		x = layout.x
		anchor = 0
	} else {
		x = layout.x + layout.w
		anchor = 1
	}

	for i, port := range ports {
		y := layout.y + NodeHeader + PortSpacing/2 + float64(i)*PortSpacing
		dc.SetColor(ColorNodeStroke)
		dc.DrawCircle(x, y, PortRadius)
		dc.Fill()

		textX := x
		if isInput {
			textX += PortRadius + 5
		} else {
			textX -= PortRadius + 5
		}
		dc.SetColor(ColorFont)
		dc.DrawStringAnchored(port.Name, textX, y, anchor, 0.5)
	}
}

func getPortIndex(ports []*axon.Port, name string) int {
	for i, p := range ports {
		if p.Name == name {
			return i
		}
	}
	return 0
}

// Map node types to colors for the header
var nodeTypeColors = map[axon.NodeType]color.Color{
	axon.NodeType_START:    color.RGBA{R: 46, G: 204, B: 113, A: 255},
	axon.NodeType_END:      color.RGBA{R: 231, G: 76, B: 60, A: 255},
	axon.NodeType_FUNCTION: color.RGBA{R: 52, G: 152, B: 219, A: 255},
	axon.NodeType_CONSTANT: color.RGBA{R: 155, G: 89, B: 182, A: 255},
	axon.NodeType_OPERATOR: color.RGBA{R: 243, G: 156, B: 18, A: 255},
	axon.NodeType_IGNORE:   color.RGBA{R: 127, G: 140, B: 141, A: 255},
}
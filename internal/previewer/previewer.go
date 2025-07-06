package previewer

import (
	"log"

	"github.com/Advik-B/Axon/pkg/axon"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

// Previewer is the Ebitengine Game implementation.
type Previewer struct {
	graph       *axon.Graph
	layoutNodes map[string]*LayoutNode
	fontFace    text.Face

	// Camera controls
	camX, camY  float64
	camZoom     float64
	isDragging  bool
	dragStartX, dragStartY int
}

// NewPreviewer creates and initializes a new previewer application.
func NewPreviewer(graph *axon.Graph) (*Previewer, error) {
	// Initialize the font face for the new text/v2 package. This call is simpler and does not return an error.
	face := text.NewGoXFace(basicfont.Face7x13)

	p := &Previewer{
		graph:       graph,
		layoutNodes: CalculateLayout(graph),
		fontFace:    face,
		camZoom:     0.8, // Start slightly zoomed out
	}

	// Center the camera on the start node if possible
	if startNode, ok := p.layoutNodes["start"]; ok {
		p.camX = float64(startNode.Rect.Min.X + startNode.Rect.Dx()/2)
		p.camY = float64(startNode.Rect.Min.Y + startNode.Rect.Dy()/2)
	}

	return p, nil
}

// Update handles the game logic (input processing).
func (p *Previewer) Update() error {
	// Zooming with mouse wheel
	_, wy := ebiten.Wheel()
	if wy != 0 {
		mx, my := ebiten.CursorPosition()
		newZoom := p.camZoom * (1 + wy*0.1)
		if newZoom > 0.2 && newZoom < 4.0 { // Clamp zoom levels
			// Move camera center to mouse position before zoom
			p.camX += float64(mx) / p.camZoom
			p.camY += float64(my) / p.camZoom
			// Apply zoom
			p.camZoom = newZoom
			// Move camera back
			p.camX -= float64(mx) / p.camZoom
			p.camY -= float64(my) / p.camZoom
		}
	}

	// Panning with mouse drag
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !p.isDragging {
			p.isDragging = true
			p.dragStartX, p.dragStartY = ebiten.CursorPosition()
		}
		endX, endY := ebiten.CursorPosition()
		dx := float64(endX - p.dragStartX) / p.camZoom
		dy := float64(endY - p.dragStartY) / p.camZoom
		p.camX -= dx
		p.camY -= dy
		p.dragStartX, p.dragStartY = endX, endY
	} else {
		p.isDragging = false
	}

	return nil
}

// Draw renders the graph to the screen.
func (p *Previewer) Draw(screen *ebiten.Image) {
	screen.Fill(colorBg)

	// Base camera transform
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	op.GeoM.Translate(-p.camX, -p.camY)
	op.GeoM.Scale(p.camZoom, p.camZoom)
	op.GeoM.Translate(float64(sw)/2, float64(sh)/2)

	// Draw Edges first (they are behind nodes)
	for _, edge := range p.graph.ExecEdges {
		fromNode, ok1 := p.layoutNodes[edge.FromNodeId]
		toNode, ok2 := p.layoutNodes[edge.ToNodeId]
		if ok1 && ok2 {
			p0 := fromNode.OutputPorts["exec_out"]
			p3 := toNode.InputPorts["exec_in"]
			drawBezierCurve(screen, p0, p3, colorExecEdge, op)
		}
	}
	for _, edge := range p.graph.DataEdges {
		fromNode, ok1 := p.layoutNodes[edge.FromNodeId]
		toNode, ok2 := p.layoutNodes[edge.ToNodeId]
		if ok1 && ok2 {
			p0 := fromNode.OutputPorts[edge.FromPort]
			p3 := toNode.InputPorts[edge.ToPort]
			drawBezierCurve(screen, p0, p3, colorDataEdge, op)
		}
	}

	// Draw Nodes on top
	for _, node := range p.layoutNodes {
		drawNode(screen, node, p.fontFace, op)
	}
}

// Layout defines the logical screen size.
func (p *Previewer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (p *Previewer) run() {
	if err := ebiten.RunGame(p); err != nil {
		log.Fatal(err)
	}
}
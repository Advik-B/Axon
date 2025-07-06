package previewer

import (
	"image"
	"log"
	"math"

	"github.com/Advik-B/Axon/pkg/axon"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// Previewer is the Ebitengine Game implementation.
type Previewer struct {
	graph              *axon.Graph
	physicsNodes       map[string]*PhysicsNode
	titleFace          text.Face
	smallFace          text.Face
	camX, camY         float64
	camZoom            float64
	isDraggingNode     bool
	isPanning          bool
	dragStartX, dragStartY int
	draggedNode        *PhysicsNode
	lastWidth, lastHeight int
	currentOrientation LayoutOrientation
}

func NewPreviewer(graph *axon.Graph) (*Previewer, error) {
	titleFace := text.NewGoXFace(basicfont.Face7x13)
	smallFace := text.NewGoXFace(basicfont.Face7x13)

	p := &Previewer{
		graph:        graph,
		physicsNodes: initializePhysicsNodes(graph),
		titleFace:    titleFace,
		smallFace:    smallFace,
		camZoom:      0.7,
	}
	if startNode, ok := p.physicsNodes["start"]; ok {
		p.camX = startNode.Position.X + float64(startNode.Rect.Dx()/2)
		p.camY = startNode.Position.Y + float64(startNode.Rect.Dy()/2)
	}
	return p, nil
}

func (p *Previewer) Update() error {
	simulatePhysics(p.physicsNodes, p.graph.DataEdges, p.graph.ExecEdges, p.draggedNode, p.currentOrientation)
	p.handleZoom()
	p.handleDragAndPan()
	return nil
}

func (p *Previewer) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.lastWidth || outsideHeight != p.lastHeight {
		p.lastWidth, p.lastHeight = outsideWidth, outsideHeight
		newOrientation := Horizontal
		if outsideHeight > outsideWidth {
			newOrientation = Vertical
		}
		if newOrientation != p.currentOrientation {
			p.currentOrientation = newOrientation
			UpdateLayoutTargets(p.physicsNodes, p.graph, p.currentOrientation)
		}
	}
	return outsideWidth, outsideHeight
}

func (p *Previewer) Draw(screen *ebiten.Image) {
	p.drawBackgroundGrid(screen)
	op := &ebiten.DrawImageOptions{}
	sw, sh := ebiten.WindowSize()
	op.GeoM.Translate(-p.camX, -p.camY)
	op.GeoM.Scale(p.camZoom, p.camZoom)
	op.GeoM.Translate(float64(sw)/2, float64(sh)/2)

	for _, edge := range p.graph.ExecEdges {
		fromNode, ok1 := p.physicsNodes[edge.FromNodeId]
		toNode, ok2 := p.physicsNodes[edge.ToNodeId]
		if ok1 && ok2 {
			p0 := fromNode.OutputPorts["exec_out"]
			p3 := toNode.InputPorts["exec_in"]
			drawBezierCurve(screen, p0, p3, colorExec, op)
		}
	}
	for _, edge := range p.graph.DataEdges {
		fromNode, ok1 := p.physicsNodes[edge.FromNodeId]
		toNode, ok2 := p.physicsNodes[edge.ToNodeId]
		if ok1 && ok2 {
			p0 := fromNode.OutputPorts[edge.FromPort]
			p3 := toNode.InputPorts[edge.ToPort]
			portType := ""
			for _, portDef := range fromNode.Outputs {
				if portDef.Name == edge.FromPort { portType = portDef.TypeName; break }
			}
			clr, ok := dataTypeColors[portType]
			if !ok { clr = dataTypeColors["default"] }
			drawBezierCurve(screen, p0, p3, clr, op)
		}
	}

	for _, node := range p.physicsNodes {
		drawNode(screen, node.LayoutNode, p.titleFace, p.smallFace, op)
	}
}

func (p *Previewer) drawBackgroundGrid(screen *ebiten.Image) {
	screen.Fill(colorBg)
	sw, sh := screen.Size()
	topLeftX, topLeftY := p.worldCoords(0, 0)
	bottomRightX, bottomRightY := p.worldCoords(sw, sh)
	gridSize := 100.0
	subGridSize := 20.0

	startX := math.Floor(topLeftX/subGridSize) * subGridSize
	for x := startX; x < bottomRightX; x += subGridSize {
		sx1, sy1 := p.screenCoords(x, topLeftY)
		sx2, sy2 := p.screenCoords(x, bottomRightY)
		vector.StrokeLine(screen, float32(sx1), float32(sy1), float32(sx2), float32(sy2), 0.5, colorGridSub, false)
	}
	startY := math.Floor(topLeftY/subGridSize) * subGridSize
	for y := startY; y < bottomRightY; y += subGridSize {
		sx1, sy1 := p.screenCoords(topLeftX, y)
		sx2, sy2 := p.screenCoords(bottomRightX, y)
		vector.StrokeLine(screen, float32(sx1), float32(sy1), float32(sx2), float32(sy2), 0.5, colorGridSub, false)
	}

	startX = math.Floor(topLeftX/gridSize) * gridSize
	for x := startX; x < bottomRightX; x += gridSize {
		sx1, sy1 := p.screenCoords(x, topLeftY)
		sx2, sy2 := p.screenCoords(x, bottomRightY)
		vector.StrokeLine(screen, float32(sx1), float32(sy1), float32(sx2), float32(sy2), 1, colorGrid, false)
	}
	startY = math.Floor(topLeftY/gridSize) * gridSize
	for y := startY; y < bottomRightY; y += gridSize {
		sx1, sy1 := p.screenCoords(topLeftX, y)
		sx2, sy2 := p.screenCoords(bottomRightX, y)
		vector.StrokeLine(screen, float32(sx1), float32(sy1), float32(sx2), float32(sy2), 1, colorGrid, false)
	}
}

func (p *Previewer) worldCoords(screenX, screenY int) (float64, float64) {
	sw, sh := ebiten.WindowSize()
	wx := (float64(screenX)-float64(sw)/2)/p.camZoom + p.camX
	wy := (float64(screenY)-float64(sh)/2)/p.camZoom + p.camY
	return wx, wy
}

func (p *Previewer) screenCoords(worldX, worldY float64) (float32, float32) {
	sw, sh := ebiten.WindowSize()
	sx := (worldX-p.camX)*p.camZoom + float64(sw)/2
	sy := (worldY-p.camY)*p.camZoom + float64(sh)/2
	return float32(sx), float32(sy)
}

func (p *Previewer) handleDragAndPan() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if !p.isDraggingNode && !p.isPanning {
			wx, wy := p.worldCoords(mx, my)
			for _, n := range p.physicsNodes {
				if image.Pt(int(wx), int(wy)).In(n.Rect) {
					p.isDraggingNode = true
					p.draggedNode = n
					break
				}
			}
			if !p.isDraggingNode {
				p.isPanning = true
			}
			p.dragStartX, p.dragStartY = mx, my
		}
		if p.isDraggingNode && p.draggedNode != nil {
			wx, wy := p.worldCoords(mx, my)
			p.draggedNode.Position.X = wx - float64(p.draggedNode.Rect.Dx()/2)
			p.draggedNode.Position.Y = wy - float64(p.draggedNode.Rect.Dy()/2)
			p.draggedNode.updateRect(p.currentOrientation)
		} else if p.isPanning {
			endX, endY := mx, my
			dx := float64(endX - p.dragStartX) / p.camZoom
			dy := float64(endY - p.dragStartY) / p.camZoom
			p.camX -= dx
			p.camY -= dy
			p.dragStartX, p.dragStartY = endX, endY
		}
	} else {
		// **THE FIX**: When the mouse button is released, update the target position.
		if p.isDraggingNode && p.draggedNode != nil {
			// Make the dragged position the new "home" for the node.
			p.draggedNode.TargetPosition = p.draggedNode.Position
		}
		p.isDraggingNode = false
		p.isPanning = false
		p.draggedNode = nil
	}
}

func (p *Previewer) handleZoom() {
	_, wy := ebiten.Wheel()
	if wy != 0 {
		mx, my := ebiten.CursorPosition()
		newZoom := p.camZoom * (1 + wy*0.1)
		if newZoom > 0.1 && newZoom < 5.0 {
			p.camX += float64(mx) / p.camZoom
			p.camY += float64(my) / p.camZoom
			p.camZoom = newZoom
			p.camX -= float64(mx) / p.camZoom
			p.camY -= float64(my) / p.camZoom
		}
	}
}

func initializePhysicsNodes(graph *axon.Graph) map[string]*PhysicsNode {
	physicsNodes := make(map[string]*PhysicsNode)
	execAdj, nodeMap := buildAdjacency(graph)
	layers := calculateLayers(graph, execAdj, nodeMap)
	orientation := Horizontal
	for l, layerNodes := range layers {
		layerHeight := len(layerNodes)*(nodeHeight+vSpacing) - vSpacing
		startY := -layerHeight / 2
		x := l * (nodeWidth + hSpacing)
		for i, node := range layerNodes {
			y := startY + i*(nodeHeight+vSpacing)
			pn := &PhysicsNode{
				LayoutNode: &LayoutNode{
					Node:        node,
					InputPorts:  make(map[string]image.Point),
					OutputPorts: make(map[string]image.Point),
				},
				Position:       Vec2{X: float64(x), Y: float64(y)},
				TargetPosition: Vec2{X: float64(x), Y: float64(y)},
			}
			pn.updateRect(orientation)
			physicsNodes[node.Id] = pn
		}
	}
	return physicsNodes
}

func (p *Previewer) run() {
	if err := ebiten.RunGame(p); err != nil {
		log.Fatal(err)
	}
}
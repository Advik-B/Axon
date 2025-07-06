package previewer

import (
	"image"
	"log"

	"github.com/Advik-B/Axon/pkg/axon"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

type Previewer struct {
	graph              *axon.Graph
	physicsNodes       map[string]*PhysicsNode
	fontFace           text.Face
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
	face := text.NewGoXFace(basicfont.Face7x13)

	p := &Previewer{
		graph:        graph,
		physicsNodes: initializePhysicsNodes(graph), // Directly initialize all nodes
		fontFace:     face,
		camZoom:      0.8,
	}

	if startNode, ok := p.physicsNodes["start"]; ok {
		p.camX = startNode.Position.X + float64(startNode.Rect.Dx()/2)
		p.camY = startNode.Position.Y + float64(startNode.Rect.Dy()/2)
	}

	return p, nil
}


// Update handles the game logic, including physics and input.
func (p *Previewer) Update() error {
	// **THE FIX**: Pass the current orientation to the physics simulation.
	simulatePhysics(p.physicsNodes, p.graph.DataEdges, p.graph.ExecEdges, p.draggedNode, p.currentOrientation)

	p.handleZoom()
	p.handleDragAndPan()
	return nil
}


func (p *Previewer) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.lastWidth || outsideHeight != p.lastHeight {
		p.lastWidth = outsideWidth
		p.lastHeight = outsideHeight

		var newOrientation LayoutOrientation
		if outsideHeight > outsideWidth {
			newOrientation = Vertical
		} else {
			newOrientation = Horizontal
		}

		if newOrientation != p.currentOrientation {
			p.currentOrientation = newOrientation
			UpdateLayoutTargets(p.physicsNodes, p.graph, p.currentOrientation)
		}
	}
	return outsideWidth, outsideHeight
}


func (p *Previewer) Draw(screen *ebiten.Image) {
	screen.Fill(colorBg)

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
			drawBezierCurve(screen, p0, p3, colorExecEdge, op)
		}
	}
	for _, edge := range p.graph.DataEdges {
		fromNode, ok1 := p.physicsNodes[edge.FromNodeId]
		toNode, ok2 := p.physicsNodes[edge.ToNodeId]
		if ok1 && ok2 {
			p0 := fromNode.OutputPorts[edge.FromPort]
			p3 := toNode.InputPorts[edge.ToPort]
			drawBezierCurve(screen, p0, p3, colorDataEdge, op)
		}
	}

	for _, node := range p.physicsNodes {
		drawNode(screen, node.LayoutNode, p.fontFace, op)
	}
}


// ... (worldCoords and handleZoom are unchanged) ...
func (p *Previewer) worldCoords(screenX, screenY int) (float64, float64) {
	sw, sh := ebiten.WindowSize()
	wx := (float64(screenX) - float64(sw)/2)/p.camZoom + p.camX
	wy := (float64(screenY) - float64(sh)/2)/p.camZoom + p.camY
	return wx, wy
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

// handleDragAndPan also needs to pass the orientation when updating the rect.
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
			// **THE FIX**: Pass the orientation here as well.
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
		p.isDraggingNode = false
		p.isPanning = false
		p.draggedNode = nil
	}
}


func initializePhysicsNodes(graph *axon.Graph) map[string]*PhysicsNode {
	physicsNodes := make(map[string]*PhysicsNode)
	execAdj, nodeMap := buildAdjacency(graph)
	layers := calculateLayers(graph, execAdj, nodeMap)

	orientation := Horizontal

	for l, layerNodes := range layers {
		if orientation == Horizontal {
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
		} else { // Vertical
			layerWidth := len(layerNodes)*(nodeWidth+hSpacing) - hSpacing
			startX := -layerWidth / 2
			y := l * (nodeHeight + vSpacing)
			for i, node := range layerNodes {
				x := startX + i*(nodeWidth+hSpacing)
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
	}
	return physicsNodes
}

func (p *Previewer) run() {
	if err := ebiten.RunGame(p); err != nil {
		log.Fatal(err)
	}
}
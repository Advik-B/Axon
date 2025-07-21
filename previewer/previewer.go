package previewer

import (
	"bytes"
	"fmt"
	"github.com/Advik-B/Axon/transpiler"
	"image"
	"image/color"
	"log"
	"math"
	"strings"

	"github.com/Advik-B/Axon/pkg/axon"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/gomonobold"
)

const (
	// Panel dimensions
	codePanelWidthLandscape = 450
	codePanelHeightPortrait = 300
	// Code rendering properties
	codeLineHeight = 18
	codePadding    = 10
)

// Previewer is the Ebitengine Game implementation.
type Previewer struct {
	graph        *axon.Graph
	physicsNodes map[string]*PhysicsNode
	titleFace    text.Face
	smallFace    text.Face
	codeFace     text.Face

	// Camera and interaction state
	camX, camY, camZoom    float64
	isDraggingNode         bool
	isPanning              bool
	dragStartX, dragStartY int
	draggedNode            *PhysicsNode
	lastWidth, lastHeight  int
	currentOrientation     LayoutOrientation

	// Code panel state
	showCodePanel     bool
	codePanelImage    *ebiten.Image
	codeScrollY       float64
	codeContentHeight float64
	transpiledCode    string
}

func NewPreviewer(graph *axon.Graph) (*Previewer, error) {
	fontBytes := gomonobold.TTF

	// Create the font face source from TTF
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to load font source: %w", err)
	}

	// Construct text/v2 font faces
	titleFace := &text.GoTextFace{
		Source: source,
		Size:   16,
	}
	smallFace := &text.GoTextFace{
		Source: source,
		Size:   13,
	}
	codeFace := &text.GoTextFace{
		Source: source,
		Size:   12,
	}

	p := &Previewer{
		graph:        graph,
		physicsNodes: initializePhysicsNodes(graph),
		titleFace:    titleFace,
		smallFace:    smallFace,
		codeFace:     codeFace,
		camZoom:      0.7,
	}

	if startNode, ok := p.physicsNodes["start"]; ok {
		p.camX = startNode.Position.X + float64(startNode.Rect.Dx()/2)
		p.camY = startNode.Position.Y + float64(startNode.Rect.Dy()/2)
	}

	if err := p.updateCodePanel(); err != nil {
		log.Printf("Could not generate initial code preview: %v", err)
	}

	return p, nil
}

func (p *Previewer) Update() error {
	simulatePhysics(p.physicsNodes, p.graph.DataEdges, p.graph.ExecEdges, p.draggedNode, p.currentOrientation)
	p.handleZoom()
	p.handleInput()
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
				if portDef.Name == edge.FromPort {
					portType = portDef.TypeName
					break
				}
			}
			clr, ok := dataTypeColors[portType]
			if !ok {
				clr = dataTypeColors["default"]
			}
			drawBezierCurve(screen, p0, p3, clr, op)
		}
	}

	for _, node := range p.physicsNodes {
		drawNode(screen, node.LayoutNode, p.titleFace, p.smallFace, op)
	}

	if p.showCodePanel {
		p.drawCodePanel(screen)
	}
}

// chromaToRGBA converts a chroma.Colour to a standard color.RGBA
func chromaToRGBA(c chroma.Colour) color.RGBA {
	return color.RGBA{R: c.Red(), G: c.Green(), B: c.Blue(), A: 255}
}

func (p *Previewer) updateCodePanel() error {
	code, err := transpiler.Transpile(p.graph)
	if err != nil {
		p.transpiledCode = fmt.Sprintf("// Transpilation Error:\n// %v", err)
	} else {
		p.transpiledCode = code
	}

	lexer := lexers.Get("go")
	iterator, err := lexer.Tokenise(nil, p.transpiledCode)
	if err != nil {
		return err
	}
	style := styles.Get("monokai")
	bgColor := chromaToRGBA(style.Get(chroma.Background).Background)

	lines := strings.Split(p.transpiledCode, "\n")
	p.codeContentHeight = float64(len(lines)) * codeLineHeight
	imgWidth := 2000
	p.codePanelImage = ebiten.NewImage(imgWidth, int(p.codeContentHeight)+codePadding*2)
	p.codePanelImage.Fill(bgColor)

	face := p.codeFace
	drawOpts := &text.DrawOptions{}
	x := float64(codePadding)
	y := float64(codePadding) + codeLineHeight

	for token := iterator(); token != chroma.EOF; token = iterator() {
		entry := style.Get(token.Type)
		if entry.Colour.IsSet() {
			drawOpts.ColorScale.Reset()
			drawOpts.ColorScale.ScaleWithColor(chromaToRGBA(entry.Colour))
		} else {
			drawOpts.ColorScale.Reset() // Default to white
		}

		tokenLines := strings.Split(token.Value, "\n")
		for i, line := range tokenLines {
			if i > 0 {
				x = codePadding
				y += codeLineHeight
			}
			drawOpts.GeoM.Reset()
			drawOpts.GeoM.Translate(x, y)
			advance, _ := text.Measure(line, face, 0)
			text.Draw(p.codePanelImage, line, face, drawOpts)
			x += advance
		}
	}
	return nil
}

func (p *Previewer) drawCodePanel(screen *ebiten.Image) {
	sw, sh := screen.Size()
	var panelRect image.Rectangle
	if p.currentOrientation == Horizontal {
		panelRect = image.Rect(sw-codePanelWidthLandscape, 0, sw, sh)
	} else {
		panelRect = image.Rect(0, sh-codePanelHeightPortrait, sw, sh)
	}

	vector.DrawFilledRect(screen, float32(panelRect.Min.X), float32(panelRect.Min.Y), float32(panelRect.Dx()), float32(panelRect.Dy()), color.RGBA{20, 21, 22, 240}, false)

	if p.codePanelImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(panelRect.Min.X), float64(panelRect.Min.Y)-p.codeScrollY)
		screen.DrawImage(p.codePanelImage, op)
	}

	vector.StrokeRect(screen, float32(panelRect.Min.X-1), float32(panelRect.Min.Y), float32(panelRect.Dx()+1), float32(panelRect.Dy()), 1, color.Black, false)
}

func (p *Previewer) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.showCodePanel = !p.showCodePanel
		if p.showCodePanel {
			if err := p.updateCodePanel(); err != nil {
				log.Printf("Error updating code panel: %v", err)
			}
		}
	}

	mx, my := ebiten.CursorPosition()
	isCursorOverPanel := false
	var panelRect image.Rectangle
	if p.showCodePanel {
		if p.currentOrientation == Horizontal {
			panelRect = image.Rect(p.lastWidth-codePanelWidthLandscape, 0, p.lastWidth, p.lastHeight)
		} else {
			panelRect = image.Rect(0, p.lastHeight-codePanelHeightPortrait, p.lastWidth, p.lastHeight)
		}
		if image.Pt(mx, my).In(panelRect) {
			isCursorOverPanel = true
		}
	}

	if isCursorOverPanel {
		_, wheelY := ebiten.Wheel()
		p.codeScrollY -= wheelY * codeLineHeight

		var panelHeight float64 = float64(panelRect.Dy())
		maxScroll := p.codeContentHeight - panelHeight + codePadding*2
		if maxScroll < 0 {
			maxScroll = 0
		}
		if p.codeScrollY > maxScroll {
			p.codeScrollY = maxScroll
		}
		if p.codeScrollY < 0 {
			p.codeScrollY = 0
		}
	}

	_, wheelY := ebiten.Wheel()
	if isCursorOverPanel && wheelY != 0 {
		return
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if isCursorOverPanel {
			return
		}
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
			dx := float64(endX-p.dragStartX) / p.camZoom
			dy := float64(endY-p.dragStartY) / p.camZoom
			p.camX -= dx
			p.camY -= dy
			p.dragStartX, p.dragStartY = endX, endY
		}
	} else {
		if p.isDraggingNode && p.draggedNode != nil {
			p.draggedNode.TargetPosition = p.draggedNode.Position
		}
		p.isDraggingNode = false
		p.isPanning = false
		p.draggedNode = nil
	}
}

func (p *Previewer) drawBackgroundGrid(screen *ebiten.Image) {
	screen.Fill(colorBg)
	sw, sh := screen.Size()
	topLeftX, topLeftY := p.worldCoords(0, 0)
	bottomRightX, bottomRightY := p.worldCoords(sw, sh)
	gridSize := 100.0
	subGridSize := 20.0

	for x := math.Floor(topLeftX/subGridSize) * subGridSize; x < bottomRightX; x += subGridSize {
		sx1, sy1 := p.screenCoords(x, topLeftY)
		sx2, sy2 := p.screenCoords(x, bottomRightY)
		vector.StrokeLine(screen, float32(sx1), float32(sy1), float32(sx2), float32(sy2), 0.5, colorGridSub, false)
	}
	for y := math.Floor(topLeftY/subGridSize) * subGridSize; y < bottomRightY; y += subGridSize {
		sx1, sy1 := p.screenCoords(topLeftX, y)
		sx2, sy2 := p.screenCoords(bottomRightX, y)
		vector.StrokeLine(screen, float32(sx1), float32(sy1), float32(sx2), float32(sy2), 0.5, colorGridSub, false)
	}

	for x := math.Floor(topLeftX/gridSize) * gridSize; x < bottomRightX; x += gridSize {
		sx1, sy1 := p.screenCoords(x, topLeftY)
		sx2, sy2 := p.screenCoords(x, bottomRightY)
		vector.StrokeLine(screen, float32(sx1), float32(sy1), float32(sx2), float32(sy2), 1, colorGrid, false)
	}
	for y := math.Floor(topLeftY/gridSize) * gridSize; y < bottomRightY; y += gridSize {
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

func (p *Previewer) handleZoom() {
	mx, my := ebiten.CursorPosition()
	isCursorOverPanel := false
	if p.showCodePanel {
		var panelRect image.Rectangle
		if p.currentOrientation == Horizontal {
			panelRect = image.Rect(p.lastWidth-codePanelWidthLandscape, 0, p.lastWidth, p.lastHeight)
		} else {
			panelRect = image.Rect(0, p.lastHeight-codePanelHeightPortrait, p.lastWidth, p.lastHeight)
		}
		if image.Pt(mx, my).In(panelRect) {
			isCursorOverPanel = true
		}
	}
	if isCursorOverPanel {
		return
	}

	_, wy := ebiten.Wheel()
	if wy != 0 {
		newZoom := p.camZoom * math.Pow(1.1, wy)
		if newZoom > 0.1 && newZoom < 5.0 {
			oldWx, oldWy := p.worldCoords(mx, my)
			p.camZoom = newZoom
			newWx, newWy := p.worldCoords(mx, my)
			p.camX -= newWx - oldWx
			p.camY -= newWy - oldWy
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

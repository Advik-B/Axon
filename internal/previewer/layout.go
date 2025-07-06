package previewer

import (
	"image"
	"math"

	"github.com/Advik-B/Axon/pkg/axon"
)

// Layout constants
const (
	nodeHeight  = 90
	nodeWidth         = 200
	minNodeHeight     = 90
	hSpacing          = 100
	vSpacing          = 70
	portRadius        = 5
	portRowHeight     = 22 // Vertical space allocated for each port row
	nodePaddingTop    = 15
	nodePaddingBottom = 15
)

// LayoutOrientation defines the direction of the graph flow.
type LayoutOrientation int

const (
	Horizontal LayoutOrientation = iota
	Vertical
)

// LayoutNode stores a node's visual info, which is updated by the physics simulation.
type LayoutNode struct {
	*axon.Node
	Rect        image.Rectangle
	InputPorts  map[string]image.Point
	OutputPorts map[string]image.Point
}

// updateRect recalculates a node's visual rectangle and port positions based on its physics position and orientation.
func (n *PhysicsNode) updateRect(orientation LayoutOrientation) {
	x, y := int(math.Round(n.Position.X)), int(math.Round(n.Position.Y))

	// Dynamically calculate node height based on the number of ports
	numInputRows := len(n.Inputs)
	numOutputRows := len(n.Outputs)
	if n.Type != axon.NodeType_END { numOutputRows++ } // Account for exec_out
	if n.Type != axon.NodeType_START { numInputRows++ } // Account for exec_in
	
	bodyRowCount := math.Max(float64(numInputRows), float64(numOutputRows))
	dynamicHeight := int(nodeHeaderHeight) + nodePaddingTop + nodePaddingBottom + (int(bodyRowCount) * portRowHeight)
	finalHeight := int(math.Max(float64(dynamicHeight), float64(minNodeHeight)))

	n.Rect = image.Rect(x, y, x+nodeWidth, y+finalHeight)

	// Dynamically update port positions
	if orientation == Horizontal {
		// Exec ports are at the top of the body
		n.InputPorts["exec_in"] = image.Pt(x, y+int(nodeHeaderHeight)+nodePaddingTop)
		n.OutputPorts["exec_out"] = image.Pt(x+nodeWidth, y+int(nodeHeaderHeight)+nodePaddingTop)
		// Data ports follow
		for i, p := range n.Inputs {
			portY := y + int(nodeHeaderHeight) + nodePaddingTop + (i+1)*portRowHeight
			n.InputPorts[p.Name] = image.Pt(x, portY)
		}
		for i, p := range n.Outputs {
			portY := y + int(nodeHeaderHeight) + nodePaddingTop + (i+1)*portRowHeight
			n.OutputPorts[p.Name] = image.Pt(x+nodeWidth, portY)
		}
	} else { // Vertical
		n.InputPorts["exec_in"] = image.Pt(x+nodeWidth/2, y)
		n.OutputPorts["exec_out"] = image.Pt(x+nodeWidth/2, y+finalHeight)
		numInputs := len(n.Inputs)
		startX_in := x + nodeWidth/2 - (numInputs-1)*portRowHeight/2
		for i, p := range n.Inputs {
			n.InputPorts[p.Name] = image.Pt(startX_in+i*portRowHeight, y)
		}
		numOutputs := len(n.Outputs)
		startX_out := x + nodeWidth/2 - (numOutputs-1)*portRowHeight/2
		for i, p := range n.Outputs {
			n.OutputPorts[p.Name] = image.Pt(startX_out+i*portRowHeight, y+finalHeight)
		}
	}
}

// UpdateLayoutTargets calculates the ideal target positions for all nodes based on the given orientation.
func UpdateLayoutTargets(nodes map[string]*PhysicsNode, graph *axon.Graph, orientation LayoutOrientation) {
	execAdj, nodeMap := buildAdjacency(graph)
	layers := calculateLayers(graph, execAdj, nodeMap)

	for l, layerNodes := range layers {
		if orientation == Horizontal {
			layerHeight := len(layerNodes)*(minNodeHeight+vSpacing) - vSpacing
			startY := -layerHeight / 2
			x := l * (nodeWidth + hSpacing)
			for i, node := range layerNodes {
				y := startY + i*(minNodeHeight+vSpacing)
				if pn, ok := nodes[node.Id]; ok {
					pn.TargetPosition = Vec2{X: float64(x), Y: float64(y)}
				}
			}
		} else { // Vertical
			layerWidth := len(layerNodes)*(nodeWidth+hSpacing) - hSpacing
			startX := -layerWidth / 2
			y := l * (minNodeHeight + vSpacing)
			for i, node := range layerNodes {
				x := startX + i*(nodeWidth+hSpacing)
				if pn, ok := nodes[node.Id]; ok {
					pn.TargetPosition = Vec2{X: float64(x), Y: float64(y)}
				}
			}
		}
	}
}

// --- Helper functions for layout calculation ---

func buildAdjacency(graph *axon.Graph) (map[string][]string, map[string]*axon.Node) {
	execAdj := make(map[string][]string)
	nodeMap := make(map[string]*axon.Node)
	for _, node := range graph.Nodes {
		nodeMap[node.Id] = node
		execAdj[node.Id] = []string{}
	}
	for _, edge := range graph.ExecEdges {
		execAdj[edge.FromNodeId] = append(execAdj[edge.FromNodeId], edge.ToNodeId)
	}
	return execAdj, nodeMap
}

func calculateLayers(graph *axon.Graph, execAdj map[string][]string, nodeMap map[string]*axon.Node) map[int][]*axon.Node {
	layers := make(map[int][]*axon.Node)
	visited := make(map[string]bool)
	queue := []*axon.Node{}

	for _, node := range graph.Nodes {
		if node.Type == axon.NodeType_START || node.Type == axon.NodeType_FUNC_DEF {
			queue = append(queue, node)
			visited[node.Id] = true
		}
	}

	layerIndex := 0
	for len(queue) > 0 {
		layerSize := len(queue)
		for i := 0; i < layerSize; i++ {
			node := queue[0]
			queue = queue[1:]
			layers[layerIndex] = append(layers[layerIndex], node)
			for _, neighborID := range execAdj[node.Id] {
				if !visited[neighborID] {
					visited[neighborID] = true
					queue = append(queue, nodeMap[neighborID])
				}
			}
		}
		layerIndex++
	}

	globalsLayer := -1
	for _, node := range graph.Nodes {
		if !visited[node.Id] {
			layers[globalsLayer] = append(layers[globalsLayer], node)
		}
	}
	return layers
}

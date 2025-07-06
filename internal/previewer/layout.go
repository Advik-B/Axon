package previewer

import (
	"image"
	"math"

	"github.com/Advik-B/Axon/pkg/axon"
)

// Layout constants
const (
	nodeWidth   = 200
	nodeHeight  = 90
	hSpacing    = 100
	vSpacing    = 70
	portRadius  = 5
	portSpacing = 22
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
	n.Rect = image.Rect(x, y, x+nodeWidth, y+nodeHeight)

	// Dynamically update port positions based on orientation
	if orientation == Horizontal { // Left-to-Right layout
		n.InputPorts["exec_in"] = image.Pt(x, y+20)
		n.OutputPorts["exec_out"] = image.Pt(x+nodeWidth, y+20)
		numInputs := len(n.Node.Inputs)
		startY_in := y + nodeHeight/2 - (numInputs-1)*portSpacing/2
		for i, p := range n.Node.Inputs {
			n.InputPorts[p.Name] = image.Pt(x, startY_in+i*portSpacing)
		}
		numOutputs := len(n.Node.Outputs)
		startY_out := y + nodeHeight/2 - (numOutputs-1)*portSpacing/2
		for i, p := range n.Node.Outputs {
			n.OutputPorts[p.Name] = image.Pt(x+nodeWidth, startY_out+i*portSpacing)
		}
	} else { // Top-to-Bottom layout
		n.InputPorts["exec_in"] = image.Pt(x+25, y)
		n.OutputPorts["exec_out"] = image.Pt(x+25, y+nodeHeight)
		numInputs := len(n.Node.Inputs)
		startX_in := x + nodeWidth/2 - (numInputs-1)*portSpacing/2
		for i, p := range n.Node.Inputs {
			n.InputPorts[p.Name] = image.Pt(startX_in+i*portSpacing, y)
		}
		numOutputs := len(n.Node.Outputs)
		startX_out := x + nodeWidth/2 - (numOutputs-1)*portSpacing/2
		for i, p := range n.Node.Outputs {
			n.OutputPorts[p.Name] = image.Pt(startX_out+i*portSpacing, y+nodeHeight)
		}
	}
}

// UpdateLayoutTargets calculates the ideal target positions for all nodes based on the given orientation.
func UpdateLayoutTargets(nodes map[string]*PhysicsNode, graph *axon.Graph, orientation LayoutOrientation) {
	execAdj, nodeMap := buildAdjacency(graph)
	layers := calculateLayers(graph, execAdj, nodeMap)

	for l, layerNodes := range layers {
		if orientation == Horizontal {
			layerHeight := len(layerNodes)*(nodeHeight+vSpacing) - vSpacing
			startY := -layerHeight / 2
			x := l * (nodeWidth + hSpacing)
			for i, node := range layerNodes {
				y := startY + i*(nodeHeight+vSpacing)
				if pn, ok := nodes[node.Id]; ok {
					pn.TargetPosition = Vec2{X: float64(x), Y: float64(y)}
				}
			}
		} else { // Vertical
			layerWidth := len(layerNodes)*(nodeWidth+hSpacing) - hSpacing
			startX := -layerWidth / 2
			y := l * (nodeHeight + vSpacing)
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

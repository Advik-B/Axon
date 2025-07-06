package previewer

import (
	"image"
	"math"

	"github.com/Advik-B/Axon/pkg/axon"
)

const (
	nodeWidth   = 180
	nodeHeight  = 80
	hSpacing    = 80
	vSpacing    = 60
	portRadius  = 5
	portSpacing = 20
)

// LayoutNode stores a node's visual info, which is updated by the physics simulation.
type LayoutNode struct {
	*axon.Node
	Rect        image.Rectangle
	InputPorts  map[string]image.Point
	OutputPorts map[string]image.Point
}

// updateRect recalculates a node's visual rectangle and port positions based on its physics position.
func (n *PhysicsNode) updateRect() {
	x, y := int(math.Round(n.Position.X)), int(math.Round(n.Position.Y))
	n.Rect = image.Rect(x, y, x+nodeWidth, y+nodeHeight)

	// Update port positions relative to the new rect
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
}

// CalculateLayout performs the initial hierarchical layout and returns a map of physics-enabled nodes.
func CalculateLayout(graph *axon.Graph) map[string]*PhysicsNode {
	physicsNodes := make(map[string]*PhysicsNode)
	execAdj := make(map[string][]string)
	nodeMap := make(map[string]*axon.Node)
	for _, node := range graph.Nodes {
		nodeMap[node.Id] = node
		execAdj[node.Id] = []string{}
	}
	for _, edge := range graph.ExecEdges {
		execAdj[edge.FromNodeId] = append(execAdj[edge.FromNodeId], edge.ToNodeId)
	}

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

	for l, layerNodes := range layers {
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
				Position: Vec2{X: float64(x), Y: float64(y)},
			}
			pn.updateRect()
			physicsNodes[node.Id] = pn
		}
	}
	return physicsNodes
}
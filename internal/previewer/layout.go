package previewer

import (
	"image"

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

// LayoutNode stores a node and its calculated position and dimensions.
type LayoutNode struct {
	*axon.Node
	Rect        image.Rectangle
	InputPorts  map[string]image.Point
	OutputPorts map[string]image.Point
}

// CalculateLayout performs a hierarchical layout of the graph nodes.
func CalculateLayout(graph *axon.Graph) map[string]*LayoutNode {
	layoutNodes := make(map[string]*LayoutNode)
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
			rect := image.Rect(x, y, x+nodeWidth, y+nodeHeight)
			ln := &LayoutNode{
				Node:        node,
				Rect:        rect,
				InputPorts:  make(map[string]image.Point),
				OutputPorts: make(map[string]image.Point),
			}

			// Add conventional ports for execution flow, which aren't in the proto
			ln.InputPorts["exec_in"] = image.Pt(x, y+20)
			ln.OutputPorts["exec_out"] = image.Pt(x+nodeWidth, y+20)

			// Calculate data port positions, offset from the execution ports
			for j, p := range node.Inputs {
				ln.InputPorts[p.Name] = image.Pt(x, y+nodeHeight/2-len(node.Inputs)*portSpacing/2+(j*portSpacing)+portSpacing/2+10)
			}
			for j, p := range node.Outputs {
				ln.OutputPorts[p.Name] = image.Pt(x+nodeWidth, y+nodeHeight/2-len(node.Outputs)*portSpacing/2+(j*portSpacing)+portSpacing/2+10)
			}
			layoutNodes[node.Id] = ln
		}
	}

	return layoutNodes
}
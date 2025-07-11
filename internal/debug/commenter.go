package debug

import (
	"fmt"
	"strings"

	"github.com/Advik-B/Axon/pkg/axon"
)

// GenerateDebugGraph translates a standard axon.Graph into a DebugGraph with generated comments.
func GenerateDebugGraph(graph *axon.Graph) *DebugGraph {
	nodeLabelMap := make(map[string]string)
	for _, node := range graph.Nodes {
		nodeLabelMap[node.Id] = node.Label
	}

	debugNodes := make([]*DebugNode, len(graph.Nodes))
	for i, node := range graph.Nodes {
		debugNodes[i] = &DebugNode{
			HeadComment: generateNodeComment(node, graph, nodeLabelMap),
			Node:        node,
		}
	}

	return &DebugGraph{
		HeadComment: fmt.Sprintf(" Axon Debug Graph | Name: %s | ID: %s ", graph.Name, graph.Id),
		ID:          graph.Id,
		Name:        graph.Name,
		Nodes:       debugNodes,
		DataEdges:   graph.DataEdges,
		ExecEdges:   graph.ExecEdges,
	}
}

// generateNodeComment creates a descriptive multi-line comment for a single node.
func generateNodeComment(node *axon.Node, graph *axon.Graph, nodeLabels map[string]string) string {
	var sb strings.Builder

	// **THE FIX**: Create a more meaningful detail string.
	var details string
	switch node.Type {
	case axon.NodeType_CONSTANT:
		details = fmt.Sprintf("value: %s", node.Config["value"])
	case axon.NodeType_FUNCTION:
		details = fmt.Sprintf("calls: %s", node.ImplReference)
	case axon.NodeType_OPERATOR:
		details = fmt.Sprintf("op: '%s'", node.Config["op"])
	default:
		details = node.Type.String() // Fallback to the type name
	}

	sb.WriteString(fmt.Sprintf(" Node: %s (%s)", node.Label, details))

	// Describe input data connections
	for _, inputPort := range node.Inputs {
		for _, edge := range graph.DataEdges {
			if edge.ToNodeId == node.Id && edge.ToPort == inputPort.Name {
				sourceLabel := nodeLabels[edge.FromNodeId]
				sb.WriteString(fmt.Sprintf("\n - Input '%s' receives data from '%s.%s'", inputPort.Name, sourceLabel, edge.FromPort))
				break
			}
		}
	}

	// Describe execution flow
	for _, edge := range graph.ExecEdges {
		if edge.FromNodeId == node.Id {
			targetLabel := nodeLabels[edge.ToNodeId]
			sb.WriteString(fmt.Sprintf("\n - After this, execution flows to: '%s'", targetLabel))
		}
	}

	return sb.String()
}
package transpiler

import (
	"fmt"

	"github.com/Advik-B/Axon/pkg/axon"
)

// sortNodesByExec performs a topological sort on the graph using execution edges.
// This determines the precise order of operations.
func sortNodesByExec(graph *axon.Graph) ([]*axon.Node, error) {
	nodeMap := make(map[string]*axon.Node)
	for _, node := range graph.Nodes {
		nodeMap[node.Id] = node
	}

	adjList := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize every node
	for _, node := range graph.Nodes {
		inDegree[node.Id] = 0
		adjList[node.Id] = []string{}
	}

	// Build adjacency list and in-degree map from exec edges
	for _, edge := range graph.ExecEdges {
		adjList[edge.FromNodeId] = append(adjList[edge.FromNodeId], edge.ToNodeId)
		inDegree[edge.ToNodeId]++
	}

	// Find all nodes with an in-degree of 0. These are the starting points.
	// A valid graph should have one START node with an in-degree of 0.
	var queue []string
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	if len(queue) == 0 && len(graph.Nodes) > 0 {
		return nil, fmt.Errorf("execution graph has a cycle or no entry point (a START node with no incoming exec edges)")
	}

	var sortedNodes []*axon.Node
	for len(queue) > 0 {
		nodeID := queue[0]
		queue = queue[1:]

		node, exists := nodeMap[nodeID]
		if !exists {
			return nil, fmt.Errorf("node with id '%s' from exec edge not found in graph", nodeID)
		}
		sortedNodes = append(sortedNodes, node)

		for _, neighborID := range adjList[nodeID] {
			inDegree[neighborID]--
			if inDegree[neighborID] == 0 {
				queue = append(queue, neighborID)
			}
		}
	}

	if len(sortedNodes) != len(graph.Nodes) {
		return nil, fmt.Errorf("cycle detected in execution graph, not all nodes could be reached")
	}

	return sortedNodes, nil
}

// findSourceEdge finds the edge that connects to a given node's input port.
func findSourceEdge(allEdges []*axon.DataEdge, toNodeID, toPortName string) (*axon.DataEdge, error) {
	for _, edge := range allEdges {
		if edge.ToNodeId == toNodeID && edge.ToPort == toPortName {
			return edge, nil
		}
	}
	return nil, fmt.Errorf("no data edge found connecting to %s.%s", toNodeID, toPortName)
}

// findSourceVar finds the Go variable name connected to a specific input port of a node.
func findSourceVar(state *transpilationState, toNodeID, toPortName string) (string, error) {
	edge, err := findSourceEdge(state.graph.DataEdges, toNodeID, toPortName)
	if err != nil {
		return "", err
	}

	// The source key is a unique identifier for an output port.
	sourceKey := fmt.Sprintf("%s.%s", edge.FromNodeId, edge.FromPort)
	if varName, ok := state.outputVarMap[sourceKey]; ok {
		return varName, nil
	}

	// If the variable is not in the map, it means it's a constant that hasn't been transpiled yet.
	// We must generate its code on-demand.
	fromNode := state.nodeMap[edge.FromNodeId]
	if fromNode.Type == axon.NodeType_CONSTANT {
		code, err := generateConstant(state, fromNode)
		if err != nil {
			return "", err
		}
		// This is a bit of a hack: we should ideally process constants first.
		// For now, we prepend its code (which won't work). A better solution requires
		// a multi-pass transpiler or sorting by data dependency first.
		// For now, we just rely on the map being populated correctly.
		_ = code // We just needed to populate the map.
		if varName, ok := state.outputVarMap[sourceKey]; ok {
			return varName, nil
		}
	}

	return "", fmt.Errorf("no source variable found for %s.%s (source key %s)", toNodeID, toPortName, sourceKey)
}

// isOutputUsed checks if a specific output port is connected to any other node.
func isOutputUsed(graph *axon.Graph, nodeID, portName string) bool {
	for _, edge := range graph.DataEdges {
		if edge.FromNodeId == nodeID && edge.FromPort == portName {
			return true
		}
	}
	return false
}
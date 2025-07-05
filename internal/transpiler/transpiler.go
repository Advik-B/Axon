package transpiler

import (
	"fmt"
	"strings"

	"github.com/Advik-B/Axon/pkg/axon" // Use your actual go.mod path here
)

// Transpile converts an Axon graph into a Go source file.
func Transpile(graph *axon.Graph) (string, error) {
	var sb strings.Builder

	// Write package header and main function wrapper
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n") // Add a default import
	sb.WriteString("func main() {\n")

	// --- Transpilation Logic ---
	// Build a map for easy node lookup by ID
	nodeMap := make(map[string]*axon.Node)
	for _, node := range graph.Nodes {
		nodeMap[node.Id] = node
	}

	// Build dependency graph for topological sort
	inDegree := make(map[string]int)
	adjList := make(map[string][]string)
	for _, node := range graph.Nodes {
		inDegree[node.Id] = 0
		adjList[node.Id] = []string{}
	}

	// For data edges, the 'from' node is a prerequisite for the 'to' node
	for _, edge := range graph.DataEdges {
		adjList[edge.FromNodeId] = append(adjList[edge.FromNodeId], edge.ToNodeId)
		inDegree[edge.ToNodeId]++
	}

	// Simple topological sort using Kahn's algorithm
	queue := []string{}
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	var sortedNodes []*axon.Node
	for len(queue) > 0 {
		nodeID := queue[0]
		queue = queue[1:]
		sortedNodes = append(sortedNodes, nodeMap[nodeID])

		for _, neighborID := range adjList[nodeID] {
			inDegree[neighborID]--
			if inDegree[neighborID] == 0 {
				queue = append(queue, neighborID)
			}
		}
	}

	// Variable mapping: maps a node's output port to the Go variable name it produces.
	// For this simple transpiler, the variable name is the node's label.
	outputVarMap := make(map[string]string) // key: from_node_id.from_port, val: variable_name

	// Generate code for each node in sorted order
	for _, node := range sortedNodes {
		code, err := transpileNode(node, graph.DataEdges, outputVarMap)
		if err != nil {
			return "", fmt.Errorf("error transpiling node %s: %w", node.Id, err)
		}
		sb.WriteString(code)
	}

	sb.WriteString("}\n")
	return sb.String(), nil
}

// transpileNode generates the Go code for a single node.
func transpileNode(node *axon.Node, allEdges []*axon.DataEdge, outputVarMap map[string]string) (string, error) {
	switch node.Type {
	case axon.NodeType_CONSTANT:
		// Example: x := 5
		val := node.Config["value"]
		// Store the variable this node produces for other nodes to use
		outputVarMap[fmt.Sprintf("%s.out", node.Id)] = node.Label
		return fmt.Sprintf("\t%s := %s\n", node.Label, val), nil

	case axon.NodeType_OPERATOR:
		// Example: z := x + y
		op := node.Config["op"]
		inputA, errA := findSourceVar(node.Id, "a", allEdges, outputVarMap)
		inputB, errB := findSourceVar(node.Id, "b", allEdges, outputVarMap)

		if errA != nil || errB != nil {
			return "", fmt.Errorf("could not resolve inputs for operator node %s", node.Id)
		}
		// Store the variable this node produces
		outputVarMap[fmt.Sprintf("%s.out", node.Id)] = node.Label
		return fmt.Sprintf("\t%s := %s %s %s\n", node.Label, inputA, op, inputB), nil

	case axon.NodeType_FUNCTION:
		// A simple example for a built-in print function
		if node.Label == "Print" {
			arg, err := findSourceVar(node.Id, "in", allEdges, outputVarMap)
			if err != nil {
				return "", fmt.Errorf("could not find input for Print node: %w", err)
			}
			return fmt.Sprintf("\tfmt.Println(%s)\n", arg), nil
		}

	}
	// Default case for unhandled nodes
	return fmt.Sprintf("\t// Node %s (type %s) not implemented\n", node.Id, node.Type), nil
}

// findSourceVar finds the variable name connected to a specific input port of a node.
func findSourceVar(nodeID, portName string, allEdges []*axon.DataEdge, outputVarMap map[string]string) (string, error) {
	for _, edge := range allEdges {
		if edge.ToNodeId == nodeID && edge.ToPort == portName {
			// Found the edge connecting to our input port.
			// Now find the variable name from the source node's output port.
			sourceKey := fmt.Sprintf("%s.%s", edge.FromNodeId, edge.FromPort)
			if varName, ok := outputVarMap[sourceKey]; ok {
				return varName, nil
			}
		}
	}
	return "", fmt.Errorf("no source variable found for %s.%s", nodeID, portName)
}
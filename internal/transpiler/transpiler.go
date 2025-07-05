package transpiler

import (
	"fmt"
	"strings"

	"github.com/Advik-B/Axon/pkg/axon"
)

// Transpile converts an Axon graph into a Go source file string.
func Transpile(graph *axon.Graph) (string, error) {
	state := &transpilationState{
		graph:        graph,
		nodeMap:      make(map[string]*axon.Node),
		outputVarMap: make(map[string]string),
		importManager: &importManager{
			imports: make(map[string]struct{}),
		},
	}

	// 1. Index all nodes by their ID for quick lookup
	for _, node := range graph.Nodes {
		state.nodeMap[node.Id] = node
	}

	// 2. Validate the graph structure and get the execution order.
	// This checks for START/END nodes, cycles, dangling paths, and unreachable code.
	sortedNodes, err := validateAndSortExecGraph(graph, state.nodeMap)
	if err != nil {
		return "", fmt.Errorf("graph validation failed: %w", err)
	}

	// 3. Generate the Go code for the body of the main function
	var bodyBuilder strings.Builder
	for _, node := range sortedNodes {
		code, err := generateNodeCode(state, node)
		if err != nil {
			return "", fmt.Errorf("error generating code for node %s (%s): %w", node.Label, node.Id, err)
		}
		bodyBuilder.WriteString(code)
	}

	// 4. Assemble the final file content
	var finalCode strings.Builder
	finalCode.WriteString("package main\n\n")

	// Add imports collected during code generation
	finalCode.WriteString(state.importManager.generateImportBlock())
	finalCode.WriteString("\n")

	finalCode.WriteString("func main() {\n")
	finalCode.WriteString(bodyBuilder.String())
	finalCode.WriteString("}\n")

	return finalCode.String(), nil
}
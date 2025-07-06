package transpiler

import (
	"fmt"
	"strings"

	"github.com/Advik-B/Axon/pkg/axon"
)

// Transpile converts an Axon graph into a complete Go source file.
func Transpile(graph *axon.Graph) (string, error) {
	state, err := newState(graph)
	if err != nil {
		return "", fmt.Errorf("failed to initialize transpiler state: %w", err)
	}

	// 1. Identify all execution graphs (main function + global functions) and globals.
	entryPoints, globals, err := findExecutionScopes(graph)
	if err != nil {
		return "", fmt.Errorf("graph validation failed: %w", err)
	}

	var finalCode strings.Builder

	// 2. Write the package and import block
	finalCode.WriteString("package main\n\n")
	if len(graph.Imports) > 0 {
		finalCode.WriteString("import (\n")
		for _, imp := range graph.Imports {
			finalCode.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
		}
		finalCode.WriteString(")\n\n")
	}

	// 3. Transpile all global definitions (structs, constants, functions)
	globalCode, err := generateGlobals(state, globals, entryPoints)
	if err != nil {
		return "", err
	}
	finalCode.WriteString(globalCode)

	// 4. Find the main execution flow (the one starting with a START node)
	mainFlows, ok := entryPoints[axon.NodeType_START]
	if !ok || len(mainFlows) == 0 {
		// This case is for libraries that only define functions but have no main.
		fmt.Println("INFO: No START node found. Generating a library file without a main function.")
		return finalCode.String(), nil
	}

	// 5. Transpile the main function body
	finalCode.WriteString("func main() {\n")
	bodyCode, err := generateFunctionBody(state, mainFlows[0].nodes)
	if err != nil {
		return "", err
	}
	finalCode.WriteString(bodyCode)
	finalCode.WriteString("}\n")

	return finalCode.String(), nil
}

// generateGlobals transpiles all top-level definitions.
func generateGlobals(state *transpilationState, globals []*axon.Node, funcs map[axon.NodeType][]*flow) (string, error) {
	var sb strings.Builder

	// Generate Structs first
	for _, node := range globals {
		if node.Type == axon.NodeType_STRUCT_DEF {
			code, err := generateStructDef(state, node)
			if err != nil {
				return "", err
			}
			sb.WriteString(code)
		}
	}

	// Generate global constants
	for _, node := range globals {
		if node.Type == axon.NodeType_CONSTANT {
			// Register global constants so they are available everywhere.
			state.outputVarMap[fmt.Sprintf("%s.%s", node.Id, node.Outputs[0].Name)] = node.Label
			code, err := generateConstant(state, node, true) // true for global
			if err != nil {
				return "", err
			}
			sb.WriteString(code)
		}
	}

	// Generate global functions/methods
	if funcDefs, ok := funcs[axon.NodeType_FUNC_DEF]; ok {
		for _, flow := range funcDefs {
			code, err := generateFuncDef(state, flow.entryNode, flow.nodes)
			if err != nil {
				return "", err
			}
			sb.WriteString(code)
		}
	}
	return sb.String(), nil
}

// generateFunctionBody transpiles the nodes within a function's scope.
func generateFunctionBody(state *transpilationState, nodes []*axon.Node) (string, error) {
	var sb strings.Builder
	// Reverse the nodes for correct execution order from the reversed DFS
	reversedNodes := make([]*axon.Node, len(nodes))
	for i := range nodes {
		reversedNodes[i] = nodes[len(nodes)-1-i]
	}

	for _, node := range reversedNodes {
		// Transpile comments attached to the node
		sb.WriteString(generateCommentBlock(state, node))

		code, err := generateNodeCode(state, node)
		if err != nil {
			return "", fmt.Errorf("error generating code for node %s (%s): %w", node.Label, node.Id, err)
		}
		sb.WriteString(code)
	}
	return sb.String(), nil
}
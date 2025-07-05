package transpiler

import (
	"fmt"
	"strings"

	"github.com/Advik-B/Axon/pkg/axon"
)

// generateNodeCode dispatches to the correct code generation function based on node type.
func generateNodeCode(state *transpilationState, node *axon.Node) (string, error) {
	switch node.Type {
	case axon.NodeType_START:
		return "// Execution starts\n", nil // Start node is a marker, no code needed
	case axon.NodeType_CONSTANT:
		return generateConstant(state, node)
	case axon.NodeType_OPERATOR:
		return generateOperator(state, node)
	case axon.NodeType_FUNCTION:
		return generateFunction(state, node)
	default:
		return fmt.Sprintf("\t// Node type '%s' not implemented for node '%s'\n", node.Type, node.Id), nil
	}
}

// generateConstant generates code for a CONSTANT node. Example: `myVar := "hello"`
func generateConstant(state *transpilationState, node *axon.Node) (string, error) {
	val, ok := node.Config["value"]
	if !ok {
		return "", fmt.Errorf("constant node %s has no 'value' in config", node.Id)
	}
	if len(node.Outputs) == 0 {
		return "", fmt.Errorf("constant node %s must have at least one output port", node.Id)
	}

	// The variable name is the node's label.
	varName := node.Label
	// The output port is typically named "out".
	outputPortName := node.Outputs[0].Name

	// Register the output variable so other nodes can find it.
	state.outputVarMap[fmt.Sprintf("%s.%s", node.Id, outputPortName)] = varName

	return fmt.Sprintf("\t%s := %s\n", varName, val), nil
}

// generateOperator generates code for an OPERATOR node. Example: `c := a + b`
func generateOperator(state *transpilationState, node *axon.Node) (string, error) {
	op, ok := node.Config["op"]
	if !ok {
		return "", fmt.Errorf("operator node %s has no 'op' in config", node.Id)
	}

	// Find the variable names for the inputs "a" and "b".
	inputA, errA := findSourceVar(state, node.Id, "a")
	inputB, errB := findSourceVar(state, node.Id, "b")
	if errA != nil || errB != nil {
		return "", fmt.Errorf("could not resolve inputs for operator node %s", node.Id)
	}

	varName := node.Label
	outputPortName := node.Outputs[0].Name
	state.outputVarMap[fmt.Sprintf("%s.%s", node.Id, outputPortName)] = varName

	return fmt.Sprintf("\t%s := %s %s %s\n", varName, inputA, op, inputB), nil
}

// generateFunction generates code for a FUNCTION node. Example: `data, err := os.ReadFile(path)`
func generateFunction(state *transpilationState, node *axon.Node) (string, error) {
	if node.ImplReference == "" {
		return "", fmt.Errorf("function node %s is missing 'impl_reference'", node.Id)
	}

	// Parse "pkg.FuncName" and add "pkg" to imports
	funcName, err := state.importManager.addImportFromReference(node.ImplReference)
	if err != nil {
		return "", err
	}

	// Resolve all input arguments
	var args []string
	for _, inputPort := range node.Inputs {
		arg, err := findSourceVar(state, node.Id, inputPort.Name)
		if err != nil {
			return "", fmt.Errorf("could not resolve input '%s' for function node %s: %w", inputPort.Name, node.Id, err)
		}
		// Special handling for functions expecting strings from byte arrays
		if strings.HasSuffix(funcName, ".ToUpper") || strings.HasSuffix(funcName, ".ToLower") {
			sourceNode, sourcePort, _ := findSourceEdge(state.graph.DataEdges, node.Id, inputPort.Name)
			if sourceNode != nil && len(sourceNode.Outputs) > 0 && sourceNode.Outputs[0].Type == axon.DataType_BYTE_ARRAY {
				arg = fmt.Sprintf("string(%s)", arg)
			}
		}
		args = append(args, arg)
	}
	argString := strings.Join(args, ", ")

	// Prepare output variables
	var outputVars []string
	for i, outputPort := range node.Outputs {
		// Use node label with index for multiple outputs, e.g., `fileContents0`, `fileContents1`
		varName := node.Label
		if len(node.Outputs) > 1 {
			varName = fmt.Sprintf("%s%d", node.Label, i)
		}
		// Don't create a var for the error if it's not used.
		if outputPort.Type == axon.DataType_ERROR && !isOutputUsed(state.graph, node.Id, outputPort.Name) {
			varName = "_"
		}
		outputVars = append(outputVars, varName)
		// Register the output variable for other nodes to use
		state.outputVarMap[fmt.Sprintf("%s.%s", node.Id, outputPort.Name)] = varName
	}
	outputString := strings.Join(outputVars, ", ")

	// Assemble the final line of code
	if len(outputVars) > 0 {
		return fmt.Sprintf("\t%s := %s(%s)\n", outputString, funcName, argString), nil
	} else {
		return fmt.Sprintf("\t%s(%s)\n", funcName, argString), nil
	}
}
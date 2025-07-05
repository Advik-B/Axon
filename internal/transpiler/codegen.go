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
		return "\t// Execution starts\n", nil
	case axon.NodeType_END:
		return "\t// Execution ends\n", nil
	case axon.NodeType_IGNORE:
		return "", nil // IGNORE node has no code output; its existence is checked by its source
	case axon.NodeType_CONSTANT:
		return generateConstant(state, node)
	case axon.NodeType_OPERATOR:
		return generateOperator(state, node)
	case axon.NodeType_FUNCTION:
		return generateFunction(state, node)
	case axon.NodeType_RETURN:
		return generateReturn(state, node)
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
		return "", fmt.Errorf("constant node %s must have an output port", node.Id)
	}

	varName := node.Label
	outputPort := node.Outputs[0]
	state.outputVarMap[fmt.Sprintf("%s.%s", node.Id, outputPort.Name)] = varName

	// Use type information if available, otherwise let Go infer.
	if outputPort.TypeName != "" {
		state.importManager.addImportFromType(outputPort.TypeName)
		return fmt.Sprintf("\tvar %s %s = %s\n", varName, outputPort.TypeName, val), nil
	}
	return fmt.Sprintf("\t%s := %s\n", varName, val), nil
}

// generateOperator generates code for an OPERATOR node. Example: `c := a + b`
func generateOperator(state *transpilationState, node *axon.Node) (string, error) {
	op, ok := node.Config["op"]
	if !ok {
		return "", fmt.Errorf("operator node %s has no 'op' in config", node.Id)
	}
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

// generateFunction enforces explicit error handling.
func generateFunction(state *transpilationState, node *axon.Node) (string, error) {
	if node.ImplReference == "" {
		return "", fmt.Errorf("function node %s is missing 'impl_reference'", node.Id)
	}
	state.importManager.addImportFromType(node.ImplReference)

	// Resolve all input arguments
	var args []string
	for _, inputPort := range node.Inputs {
		arg, err := findSourceVar(state, node.Id, inputPort.Name)
		if err != nil {
			return "", fmt.Errorf("could not resolve input '%s' for func %s: %w", inputPort.Name, node.Id, err)
		}
		args = append(args, arg)
	}
	argString := strings.Join(args, ", ")

	// Prepare output variables and enforce explicit ignore
	var outputVars []string
	for i, outputPort := range node.Outputs {
		varName := fmt.Sprintf("%s%d", node.Label, i)
		if len(node.Outputs) == 1 {
			varName = node.Label
		}

		if !isOutputUsed(state.graph, node.Id, outputPort.Name) {
			return "", fmt.Errorf("output '%s' of function node '%s' is not used or ignored. Connect it to another node or an IGNORE node", outputPort.Name, node.Label)
		}

		// Check if the output is connected to an IGNORE node
		isIgnored := false
		for _, edge := range state.graph.DataEdges {
			if edge.FromNodeId == node.Id && edge.FromPort == outputPort.Name {
				if targetNode, ok := state.nodeMap[edge.ToNodeId]; ok && targetNode.Type == axon.NodeType_IGNORE {
					isIgnored = true
					break
				}
			}
		}

		if isIgnored {
			outputVars = append(outputVars, "_")
		} else {
			outputVars = append(outputVars, varName)
		}
		state.outputVarMap[fmt.Sprintf("%s.%s", node.Id, outputPort.Name)] = varName
	}

	outputString := strings.Join(outputVars, ", ")
	funcName := strings.Replace(node.ImplReference, "/", ".", -1) // Simplify for codegen

	if len(outputVars) > 0 {
		return fmt.Sprintf("\t%s := %s(%s)\n", outputString, funcName, argString), nil
	}
	return fmt.Sprintf("\t%s(%s)\n", funcName, argString), nil
}

// generateReturn generates a return statement.
func generateReturn(state *transpilationState, node *axon.Node) (string, error) {
	if len(node.Inputs) == 0 {
		return "\treturn\n", nil
	}
	var returnVars []string
	for _, inputPort := range node.Inputs {
		varName, err := findSourceVar(state, node.Id, inputPort.Name)
		if err != nil {
			return "", fmt.Errorf("could not find source for RETURN node input '%s': %w", inputPort.Name, err)
		}
		returnVars = append(returnVars, varName)
	}
	return fmt.Sprintf("\treturn %s\n", strings.Join(returnVars, ", ")), nil
}
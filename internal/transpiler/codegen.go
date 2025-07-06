package transpiler

import (
	"fmt"
	"strings"

	"github.com/Advik-B/Axon/pkg/axon"
)

// generateNodeCode dispatches to the correct code generation function based on node type.
func generateNodeCode(state *transpilationState, node *axon.Node) (string, error) {
	switch node.Type {
	// Flow control and definition nodes are markers; their logic is handled by the graph validation and scoping.
	case axon.NodeType_START, axon.NodeType_END, axon.NodeType_FUNC_DEF, axon.NodeType_STRUCT_DEF:
		return "", nil
	case axon.NodeType_IGNORE:
		return "", nil // IGNORE has no code output
	case axon.NodeType_CONSTANT:
		return generateConstant(state, node, false) // false for local scope
	case axon.NodeType_OPERATOR:
		return generateOperator(state, node)
	case axon.NodeType_FUNCTION:
		return generateFunctionCall(state, node)
	case axon.NodeType_RETURN:
		return generateReturn(state, node)
	default:
		return fmt.Sprintf("\t// Node type '%s' not implemented for node '%s'\n", node.Type, node.Id), nil
	}
}

// generateStructDef generates a `type ... struct` block.
func generateStructDef(state *transpilationState, node *axon.Node) (string, error) {
	var sb strings.Builder
	sb.WriteString(generateCommentBlock(state, node))
	sb.WriteString(fmt.Sprintf("type %s struct {\n", node.Label))
	for _, fieldPort := range node.Inputs {
		sb.WriteString(fmt.Sprintf("\t%s %s\n", fieldPort.Name, fieldPort.TypeName))
	}
	sb.WriteString("}\n\n")
	return sb.String(), nil
}

// generateFuncDef generates a complete function or method definition.
func generateFuncDef(state *transpilationState, entryNode *axon.Node, bodyNodes []*axon.Node) (string, error) {
	var sb strings.Builder
	sb.WriteString(generateCommentBlock(state, entryNode))

	var receiver, params, returnTypes []string
	// Check for a receiver (convention: an input port named 'receiver')
	for _, port := range entryNode.Inputs {
		if port.Name == "receiver" {
			// A short, conventional receiver name like 'u' for 'User'
			receiverName := strings.ToLower(string(port.TypeName[strings.LastIndex(port.TypeName, "*")+1]))
			receiver = append(receiver, fmt.Sprintf("(%s %s)", receiverName, port.TypeName))
			state.outputVarMap[fmt.Sprintf("%s.receiver", entryNode.Id)] = receiverName
			break
		}
	}

	// Function parameters are defined as output ports on the FUNC_DEF node
	for _, port := range entryNode.Outputs {
		params = append(params, fmt.Sprintf("%s %s", port.Name, port.TypeName))
		// Register the parameter as a known variable for the function body
		state.outputVarMap[fmt.Sprintf("%s.%s", entryNode.Id, port.Name)] = port.Name
	}

	// Find the RETURN node to determine return types
	for _, node := range bodyNodes {
		if node.Type == axon.NodeType_RETURN {
			for _, port := range node.Inputs {
				returnTypes = append(returnTypes, port.TypeName)
			}
			break // Assume one return node per function for simplicity
		}
	}

	returnStr := ""
	if len(returnTypes) > 0 {
		if len(returnTypes) > 1 {
			returnStr = fmt.Sprintf("(%s)", strings.Join(returnTypes, ", "))
		} else {
			returnStr = returnTypes[0]
		}
	}

	sb.WriteString(fmt.Sprintf("func %s %s(%s) %s {\n", strings.Join(receiver, ""), entryNode.Label, strings.Join(params, ", "), returnStr))
	bodyCode, err := generateFunctionBody(state, bodyNodes)
	if err != nil {
		return "", err
	}
	sb.WriteString(bodyCode)
	sb.WriteString("}\n\n")
	return sb.String(), nil
}

// generateConstant generates code for a CONSTANT node.
func generateConstant(state *transpilationState, node *axon.Node, isGlobal bool) (string, error) {
	val, ok := node.Config["value"]
	if !ok {
		return "", fmt.Errorf("constant node %s has no 'value' in config", node.Id)
	}
	varName := node.Label
	state.outputVarMap[fmt.Sprintf("%s.%s", node.Id, node.Outputs[0].Name)] = varName

	if isGlobal {
		return fmt.Sprintf("const %s = %s\n\n", varName, val), nil
	}
	return fmt.Sprintf("\t%s := %s\n", varName, val), nil
}

// generateFunctionCall generates a call to a function or method.
func generateFunctionCall(state *transpilationState, node *axon.Node) (string, error) {
	if node.ImplReference == "" {
		return "", fmt.Errorf("function call node %s is missing 'impl_reference'", node.Id)
	}

	var args []string
	for _, inputPort := range node.Inputs {
		arg, err := findSourceVar(state, node.Id, inputPort.Name)
		if err != nil {
			return "", fmt.Errorf("could not resolve input '%s' for func call %s: %w", inputPort.Name, node.Id, err)
		}
		args = append(args, arg)
	}
	argString := strings.Join(args, ", ")

	var outputVars []string
	for i, outputPort := range node.Outputs {
		if !isOutputUsed(state.graph, node.Id, outputPort.Name) {
			return "", fmt.Errorf("output '%s' of function call '%s' is not used or explicitly ignored", outputPort.Name, node.Label)
		}
		varName := fmt.Sprintf("%s_out%d", node.Label, i)
		if len(node.Outputs) == 1 {
			varName = node.Label
		}

		isIgnored := isPortConnectedToIgnore(state, node.Id, outputPort.Name)
		if isIgnored {
			outputVars = append(outputVars, "_")
		} else {
			outputVars = append(outputVars, varName)
		}
		state.outputVarMap[fmt.Sprintf("%s.%s", node.Id, outputPort.Name)] = varName
	}

	outputString := strings.Join(outputVars, ", ")
	if len(outputVars) > 0 {
		return fmt.Sprintf("\t%s := %s(%s)\n", outputString, node.ImplReference, argString), nil
	}
	return fmt.Sprintf("\t%s(%s)\n", node.ImplReference, argString), nil
}

// generateOperator generates code for a binary operation or struct instantiation.
func generateOperator(state *transpilationState, node *axon.Node) (string, error) {
	op, ok := node.Config["op"]
	if !ok {
		return "", fmt.Errorf("operator node %s has no 'op' in config", node.Id)
	}

	// Handle special struct instantiation operator
	if strings.HasPrefix(op, "&") || strings.HasPrefix(op, "") && 'A' <= op[0] && op[0] <= 'Z' {
		var fields []string
		for _, port := range node.Inputs {
			arg, err := findSourceVar(state, node.Id, port.Name)
			if err != nil {
				return "", err
			}
			fields = append(fields, arg)
		}
		varName := node.Label
		state.outputVarMap[fmt.Sprintf("%s.%s", node.Id, node.Outputs[0].Name)] = varName
		return fmt.Sprintf("\t%s := %s{%s}\n", varName, op, strings.Join(fields, ", ")), nil
	}

	// Handle standard binary operators
	inputA, errA := findSourceVar(state, node.Id, "a")
	inputB, errB := findSourceVar(state, node.Id, "b")
	if errA != nil || errB != nil {
		return "", fmt.Errorf("could not resolve inputs for operator node %s", node.Id)
	}
	varName := node.Label
	state.outputVarMap[fmt.Sprintf("%s.%s", node.Id, node.Outputs[0].Name)] = varName
	return fmt.Sprintf("\t%s := %s %s %s\n", varName, inputA, op, inputB), nil
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

// generateCommentBlock finds comments for a node and formats them.
func generateCommentBlock(state *transpilationState, node *axon.Node) string {
	if len(node.CommentIds) == 0 {
		return ""
	}
	var sb strings.Builder
	prefix := "\t" // Assume local scope
	if node.Type == axon.NodeType_FUNC_DEF || node.Type == axon.NodeType_STRUCT_DEF {
		prefix = "" // Global scope
	}

	for _, commentID := range node.CommentIds {
		if comment, ok := state.commentMap[commentID]; ok {
			for _, line := range strings.Split(strings.TrimSpace(comment.Content), "\n") {
				sb.WriteString(fmt.Sprintf("%s// %s\n", prefix, line))
			}
		}
	}
	return sb.String()
}
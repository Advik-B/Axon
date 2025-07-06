package transpiler

import (
	"fmt"
	"github.com/Advik-B/Axon/pkg/axon"
)

// flow represents a single, self-contained execution graph (like main or a function).
type flow struct {
	entryNode *axon.Node
	nodes     []*axon.Node
}

// findExecutionScopes identifies all separate execution flows and global definitions.
func findExecutionScopes(graph *axon.Graph) (map[axon.NodeType][]*flow, []*axon.Node, error) {
	nodeMap := make(map[string]*axon.Node)
	adjList := make(map[string][]string)
	for _, node := range graph.Nodes {
		nodeMap[node.Id] = node
		adjList[node.Id] = []string{}
	}
	for _, edge := range graph.ExecEdges {
		adjList[edge.FromNodeId] = append(adjList[edge.FromNodeId], edge.ToNodeId)
	}

	entryPoints := make(map[axon.NodeType][]*flow)
	visited := make(map[string]bool)

	for _, node := range graph.Nodes {
		if node.Type == axon.NodeType_START || node.Type == axon.NodeType_FUNC_DEF {
			if visited[node.Id] {
				continue
			}
			var pathNodes []*axon.Node
			var pathVisited = make(map[string]bool)
			err := dfs(node.Id, adjList, nodeMap, pathVisited, &pathNodes)
			if err != nil {
				return nil, nil, fmt.Errorf("validation failed for flow starting at '%s': %w", node.Label, err)
			}
			if err := validateFlowTermination(node, pathNodes, nodeMap, adjList); err != nil {
				return nil, nil, err
			}
			for _, n := range pathNodes {
				visited[n.Id] = true
			}
			entryPoints[node.Type] = append(entryPoints[node.Type], &flow{
				entryNode: node,
				nodes:     pathNodes,
			})
		}
	}

	// Identify globals: nodes not visited during any flow traversal
	var globals []*axon.Node
	for _, node := range graph.Nodes {
		if !visited[node.Id] {
			// **FIX**: An IGNORE node is a valid data sink and is allowed to be "unreachable" in the execution flow.
			if node.Type == axon.NodeType_IGNORE {
				// Mark it as visited so it doesn't get flagged as an invalid global.
				visited[node.Id] = true
				continue
			}
			if node.Type != axon.NodeType_CONSTANT && node.Type != axon.NodeType_STRUCT_DEF {
				return nil, nil, fmt.Errorf("unreachable node '%s' is not a valid global type (CONSTANT or STRUCT_DEF)", node.Label)
			}
			globals = append(globals, node)
		}
	}

	return entryPoints, globals, nil
}

// dfs traverses a flow from an entry point and returns the nodes in topological order.
func dfs(nodeID string, adjList map[string][]string, nodeMap map[string]*axon.Node, visited map[string]bool, pathNodes *[]*axon.Node) error {
	if visited[nodeID] {
		return nil
	}
	visited[nodeID] = true

	for _, neighborID := range adjList[nodeID] {
		if err := dfs(neighborID, adjList, nodeMap, visited, pathNodes); err != nil {
			return err
		}
	}
	*pathNodes = append(*pathNodes, nodeMap[nodeID])
	return nil
}

// validateFlowTermination ensures that a given flow terminates correctly.
func validateFlowTermination(entryNode *axon.Node, pathNodes []*axon.Node, nodeMap map[string]*axon.Node, adjList map[string][]string) error {
	terminatorType := axon.NodeType_END
	if entryNode.Type == axon.NodeType_FUNC_DEF {
		terminatorType = axon.NodeType_RETURN
	}

	hasTerminator := false
	for _, node := range pathNodes {
		if len(adjList[node.Id]) == 0 {
			if node.Type != terminatorType {
				return fmt.Errorf("flow starting at '%s' has a dangling path at node '%s'. It must end with a %s node", entryNode.Label, node.Label, terminatorType)
			}
			hasTerminator = true
		}
	}

	if !hasTerminator && len(pathNodes) > 1 {
		return fmt.Errorf("flow starting at '%s' has a cycle or does not have a valid %s terminator node", entryNode.Label, terminatorType)
	}
	return nil
}

// findSourceVar finds the Go variable name connected to a specific input port of a node.
func findSourceVar(state *transpilationState, toNodeID, toPortName string) (string, error) {
	for _, edge := range state.graph.DataEdges {
		if edge.ToNodeId == toNodeID && edge.ToPort == toPortName {
			sourceKey := fmt.Sprintf("%s.%s", edge.FromNodeId, edge.FromPort)
			if varName, ok := state.outputVarMap[sourceKey]; ok {
				return varName, nil
			}
			return "", fmt.Errorf("unresolved source variable for %s.%s (source key %s). This may indicate a flaw in the execution order or a missing global", toNodeID, toPortName, sourceKey)
		}
	}
	return "", fmt.Errorf("no data edge found connecting to %s.%s", toNodeID, toPortName)
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

// isPortConnectedToIgnore checks if a specific output port is connected to an IGNORE node.
func isPortConnectedToIgnore(state *transpilationState, nodeID, portName string) bool {
	for _, edge := range state.graph.DataEdges {
		if edge.FromNodeId == nodeID && edge.FromPort == portName {
			if targetNode, ok := state.nodeMap[edge.ToNodeId]; ok && targetNode.Type == axon.NodeType_IGNORE {
				return true
			}
		}
	}
	return false
}
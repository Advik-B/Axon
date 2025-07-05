package transpiler

import (
	"fmt"
	"strings"

	"github.com/Advik-B/Axon/pkg/axon"
)

// validateAndSortExecGraph ensures the graph is well-formed and returns a valid execution order.
// It checks for:
// 1. At least one START node.
// 2. All execution paths from a START node must eventually reach an END node.
// 3. All nodes are reachable from a START node (no dead code).
func validateAndSortExecGraph(graph *axon.Graph, nodeMap map[string]*axon.Node) ([]*axon.Node, error) {
	adjList := make(map[string][]string)
	revAdjList := make(map[string][]string) // For reverse traversal from END nodes
	var startNodes []*axon.Node

	for _, node := range graph.Nodes {
		adjList[node.Id] = []string{}
		revAdjList[node.Id] = []string{}
		if node.Type == axon.NodeType_START {
			startNodes = append(startNodes, node)
		}
	}
	if len(startNodes) == 0 && len(graph.Nodes) > 0 {
		return nil, fmt.Errorf("no START node found in graph")
	}

	for _, edge := range graph.ExecEdges {
		adjList[edge.FromNodeId] = append(adjList[edge.FromNodeId], edge.ToNodeId)
		revAdjList[edge.ToNodeId] = append(revAdjList[edge.ToNodeId], edge.FromNodeId)
	}

	// DFS to find all reachable nodes and check for dangling paths
	var sortedNodes []*axon.Node
	visited := make(map[string]bool)
	path := make(map[string]bool) // For cycle detection

	var dfs func(nodeID string) error
	dfs = func(nodeID string) error {
		if path[nodeID] {
			return fmt.Errorf("cycle detected in execution graph involving node %s", nodeID)
		}
		if visited[nodeID] {
			return nil
		}

		path[nodeID] = true
		visited[nodeID] = true

		node := nodeMap[nodeID]
		if node.Type == axon.NodeType_END {
			delete(path, nodeID)
			sortedNodes = append(sortedNodes, node)
			return nil // Valid termination
		}

		successors := adjList[nodeID]
		if len(successors) == 0 {
			return fmt.Errorf("dangling execution path at node %s ('%s'). All paths must terminate at an END node", node.Id, node.Label)
		}

		for _, neighborID := range successors {
			if err := dfs(neighborID); err != nil {
				return err
			}
		}

		delete(path, nodeID)
		// Prepend node to get topological sort order
		sortedNodes = append([]*axon.Node{node}, sortedNodes...)
		return nil
	}

	for _, startNode := range startNodes {
		if err := dfs(startNode.Id); err != nil {
			return nil, err
		}
	}

	// Check for unreachable nodes (nodes that exist but weren't visited)
	if len(visited) != len(graph.Nodes) {
		var unreachable []string
		for _, node := range graph.Nodes {
			if !visited[node.Id] {
				unreachable = append(unreachable, fmt.Sprintf("'%s' (%s)", node.Label, node.Id))
			}
		}
		return nil, fmt.Errorf("unreachable nodes detected (dead code): %s", strings.Join(unreachable, ", "))
	}

	return sortedNodes, nil
}

// findSourceVar finds the Go variable name connected to a specific input port of a node.
func findSourceVar(state *transpilationState, toNodeID, toPortName string) (string, error) {
	for _, edge := range state.graph.DataEdges {
		if edge.ToNodeId == toNodeID && edge.ToPort == toPortName {
			sourceKey := fmt.Sprintf("%s.%s", edge.FromNodeId, edge.FromPort)
			if varName, ok := state.outputVarMap[sourceKey]; ok {
				return varName, nil
			}
			return "", fmt.Errorf("unresolved source variable for %s.%s (source key %s). This may indicate a flaw in the execution order", toNodeID, toPortName, sourceKey)
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
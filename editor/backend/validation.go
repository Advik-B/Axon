package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Advik-B/Axon/parser"
	"github.com/Advik-B/Axon/pkg/axon"
)

// ValidationResult represents the result of graph validation
type ValidationResult struct {
	IsValid bool                `json:"isValid"`
	Errors  []ValidationError   `json:"errors"`
	Info    ValidationInfo      `json:"info"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Type        string `json:"type"`
	Message     string `json:"message"`
	NodeID      string `json:"nodeId,omitempty"`
	EdgeID      string `json:"edgeId,omitempty"`
	Severity    string `json:"severity"` // error, warning, info
}

// ValidationInfo provides additional information about the graph
type ValidationInfo struct {
	NodeCount       int      `json:"nodeCount"`
	DataEdgeCount   int      `json:"dataEdgeCount"`
	ExecEdgeCount   int      `json:"execEdgeCount"`
	StartNodes      []string `json:"startNodes"`
	EndNodes        []string `json:"endNodes"`
	IsolatedNodes   []string `json:"isolatedNodes"`
	CyclicEdges     []string `json:"cyclicEdges"`
	UnusedImports   []string `json:"unusedImports"`
}

// ValidateGraph performs comprehensive validation of an Axon graph
func ValidateGraph(graphData map[string]interface{}) (*ValidationResult, error) {
	// Convert to Axon graph format
	data, err := json.Marshal(graphData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal graph data: %w", err)
	}
	
	graph, err := parser.LoadGraphFromBytes(data)
	if err != nil {
		return &ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{
					Type:     "parse_error",
					Message:  fmt.Sprintf("Failed to parse graph: %v", err),
					Severity: "error",
				},
			},
		}, nil
	}

	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
		Info: ValidationInfo{
			NodeCount:     len(graph.Nodes),
			DataEdgeCount: len(graph.DataEdges),
			ExecEdgeCount: len(graph.ExecEdges),
			StartNodes:    []string{},
			EndNodes:      []string{},
			IsolatedNodes: []string{},
			CyclicEdges:   []string{},
			UnusedImports: []string{},
		},
	}

	// Validate nodes
	nodeMap := make(map[string]*axon.Node)
	for _, node := range graph.Nodes {
		nodeMap[node.Id] = node
		
		// Check for required fields
		if node.Id == "" {
			result.addError("validation_error", "Node missing ID", node.Id, "error")
		}

		// Track start and end nodes
		if node.Type == axon.NodeType_START {
			result.Info.StartNodes = append(result.Info.StartNodes, node.Id)
		}
		if node.Type == axon.NodeType_END {
			result.Info.EndNodes = append(result.Info.EndNodes, node.Id)
		}

		// Validate node-specific requirements
		switch node.Type {
		case axon.NodeType_CONSTANT:
			if node.Config == nil || node.Config["value"] == "" {
				result.addError("validation_error", "CONSTANT node missing value config", node.Id, "error")
			}
		case axon.NodeType_OPERATOR:
			if node.Config == nil || node.Config["op"] == "" {
				result.addError("validation_error", "OPERATOR node missing op config", node.Id, "error")
			}
		case axon.NodeType_FUNCTION:
			if node.ImplReference == "" {
				result.addError("validation_error", "FUNCTION node missing impl_reference", node.Id, "error")
			}
		}
	}

	// Validate execution flow
	if len(result.Info.StartNodes) == 0 {
		result.addError("flow_error", "Graph must have at least one START node", "", "error")
	}
	
	if len(result.Info.StartNodes) > 1 {
		result.addError("flow_warning", "Graph has multiple START nodes", "", "warning")
	}

	if len(result.Info.EndNodes) == 0 {
		result.addError("flow_warning", "Graph should have at least one END node", "", "warning")
	}

	// Validate edges
	for _, edge := range graph.DataEdges {
		if !validateEdge(edge.FromNodeId, edge.ToNodeId, nodeMap, result) {
			continue
		}
		
		// Check port existence
		fromNode := nodeMap[edge.FromNodeId]
		toNode := nodeMap[edge.ToNodeId]
		
		if !hasOutputPort(fromNode, edge.FromPort) {
			result.addError("edge_error", 
				fmt.Sprintf("Source node '%s' does not have output port '%s'", edge.FromNodeId, edge.FromPort), 
				"", "error")
		}
		
		if !hasInputPort(toNode, edge.ToPort) {
			result.addError("edge_error", 
				fmt.Sprintf("Target node '%s' does not have input port '%s'", edge.ToNodeId, edge.ToPort), 
				"", "error")
		}
	}

	for _, edge := range graph.ExecEdges {
		validateEdge(edge.FromNodeId, edge.ToNodeId, nodeMap, result)
	}

	// Check for isolated nodes (nodes with no connections)
	connectedNodes := make(map[string]bool)
	for _, edge := range graph.DataEdges {
		connectedNodes[edge.FromNodeId] = true
		connectedNodes[edge.ToNodeId] = true
	}
	for _, edge := range graph.ExecEdges {
		connectedNodes[edge.FromNodeId] = true
		connectedNodes[edge.ToNodeId] = true
	}
	
	for _, node := range graph.Nodes {
		if !connectedNodes[node.Id] && node.Type != axon.NodeType_START && node.Type != axon.NodeType_END {
			result.Info.IsolatedNodes = append(result.Info.IsolatedNodes, node.Id)
			result.addError("flow_warning", 
				fmt.Sprintf("Node '%s' is isolated (no connections)", node.Id), 
				node.Id, "warning")
		}
	}

	// Check for unused imports
	usedPackages := make(map[string]bool)
	for _, node := range graph.Nodes {
		if node.Type == axon.NodeType_FUNCTION && node.ImplReference != "" {
			if strings.Contains(node.ImplReference, ".") {
				pkg := strings.Split(node.ImplReference, ".")[0]
				usedPackages[pkg] = true
			}
		}
	}
	
	for _, imp := range graph.Imports {
		if !usedPackages[imp] {
			result.Info.UnusedImports = append(result.Info.UnusedImports, imp)
			result.addError("optimization_info", 
				fmt.Sprintf("Import '%s' is unused", imp), 
				"", "info")
		}
	}

	// Graph is invalid if there are any errors
	for _, err := range result.Errors {
		if err.Severity == "error" {
			result.IsValid = false
			break
		}
	}

	return result, nil
}

func (r *ValidationResult) addError(errType, message, nodeID, severity string) {
	r.Errors = append(r.Errors, ValidationError{
		Type:     errType,
		Message:  message,
		NodeID:   nodeID,
		Severity: severity,
	})
}

func validateEdge(fromID, toID string, nodeMap map[string]*axon.Node, result *ValidationResult) bool {
	if fromID == "" {
		result.addError("edge_error", "Edge missing from_node_id", "", "error")
		return false
	}
	
	if toID == "" {
		result.addError("edge_error", "Edge missing to_node_id", "", "error")
		return false
	}
	
	if _, exists := nodeMap[fromID]; !exists {
		result.addError("edge_error", 
			fmt.Sprintf("Edge references non-existent source node '%s'", fromID), 
			"", "error")
		return false
	}
	
	if _, exists := nodeMap[toID]; !exists {
		result.addError("edge_error", 
			fmt.Sprintf("Edge references non-existent target node '%s'", toID), 
			"", "error")
		return false
	}
	
	return true
}

func hasOutputPort(node *axon.Node, portName string) bool {
	if node == nil {
		return false
	}
	for _, output := range node.Outputs {
		if output.Name == portName {
			return true
		}
	}
	return false
}

func hasInputPort(node *axon.Node, portName string) bool {
	if node == nil {
		return false
	}
	for _, input := range node.Inputs {
		if input.Name == portName {
			return true
		}
	}
	return false
}
package transpiler

import (
	"github.com/Advik-B/Axon/pkg/axon"
)

// transpilationState holds the complete context during a single transpilation run.
type transpilationState struct {
	graph      *axon.Graph
	nodeMap    map[string]*axon.Node
	commentMap map[string]*axon.Comment
	// Maps "nodeID.portName" -> "goVariableName"
	outputVarMap map[string]string
}

func newState(graph *axon.Graph) (*transpilationState, error) {
	state := &transpilationState{
		graph:        graph,
		nodeMap:      make(map[string]*axon.Node),
		commentMap:   make(map[string]*axon.Comment),
		outputVarMap: make(map[string]string),
	}

	for _, node := range graph.Nodes {
		state.nodeMap[node.Id] = node
	}
	for _, comment := range graph.Comments {
		state.commentMap[comment.Id] = comment
	}

	return state, nil
}
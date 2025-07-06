package debug

import "github.com/Advik-B/Axon/pkg/axon"

// DebugNode represents a node with an attached comment for YAML output.
type DebugNode struct {
	HeadComment string `yaml:"-"` // This field is used by yaml.v3 to add a comment block before the node.
	axon.Node   `yaml:",inline"`
}

// DebugGraph is a wrapper around the core axon.Graph for YAML serialization with comments.
type DebugGraph struct {
	HeadComment string       `yaml:"-"`
	ID          string       `yaml:"id"`
	Name        string       `yaml:"name"`
	Nodes       []*DebugNode `yaml:"nodes"`
	DataEdges   []*DataEdge  `yaml:"data_edges"`
	ExecEdges   []*ExecEdge  `yaml:"exec_edges"`
}

// DataEdge and ExecEdge are included for completeness if we ever want to comment them.
type DataEdge = axon.DataEdge
type ExecEdge = axon.ExecEdge
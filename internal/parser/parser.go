package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Advik-B/Axon/pkg/axon"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

// Graph is an alias to the generated struct for easier use in other packages.
type Graph = axon.Graph

// LoadGraphFromFile reads an .ax (JSON), .axb (binary), or .axd (YAML) file
// and parses it into an Axon Graph struct.
func LoadGraphFromFile(filePath string) (*Graph, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var graph Graph
	fileExt := filepath.Ext(filePath)

	switch fileExt {
	case ".ax":
		unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
		if err := unmarshaler.Unmarshal(bytes, &graph); err != nil {
			return nil, fmt.Errorf("failed to parse .ax (JSON) file: %w", err)
		}
		return &graph, nil

	case ".axb":
		if err := proto.Unmarshal(bytes, &graph); err != nil {
			return nil, fmt.Errorf("failed to parse .axb (binary) file: %w", err)
		}
		return &graph, nil

	case ".axd":
		// yaml.v3 can unmarshal into a struct even if it doesn't have yaml tags.
		if err := yaml.Unmarshal(bytes, &graph); err != nil {
			return nil, fmt.Errorf("failed to parse .axd (YAML) file: %w", err)
		}
		return &graph, nil

	default:
		return nil, fmt.Errorf("unsupported file extension '%s': must be .ax, .axb, or .axd", fileExt)
	}
}
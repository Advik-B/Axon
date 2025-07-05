package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Advik-B/Axon/pkg/axon"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Graph is an alias to the generated struct for easier use in other packages.
type Graph = axon.Graph

// LoadGraphFromFile reads an .ax (JSON) or .axb (binary) file and parses
// it into an Axon Graph struct, auto-detecting the format from the extension.
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
		// Handle JSON format
		unmarshaler := protojson.UnmarshalOptions{
			DiscardUnknown: true,
		}
		if err := unmarshaler.Unmarshal(bytes, &graph); err != nil {
			return nil, fmt.Errorf("failed to parse .ax (JSON) file: %w", err)
		}
		return &graph, nil

	case ".axb":
		// Handle binary Protobuf format
		if err := proto.Unmarshal(bytes, &graph); err != nil {
			return nil, fmt.Errorf("failed to parse .axb (binary) file: %w", err)
		}
		return &graph, nil

	default:
		return nil, fmt.Errorf("unsupported file extension '%s': must be .ax or .axb", fileExt)
	}
}
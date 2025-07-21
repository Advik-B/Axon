package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Advik-B/Axon/pkg/axon"
	"github.com/ulikunitz/xz"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

// Graph is an alias to the generated struct for easier use in other packages.
type Graph = axon.Graph

// LoadGraphFromFile reads a graph file, auto-detecting the format (.ax, .axb, .axd, .axc).
func LoadGraphFromFile(filePath string) (*Graph, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var graph Graph
	fileExt := filepath.Ext(filePath)
	var bytes []byte

	switch fileExt {
	case ".ax", ".axd":
		// For text-based formats, read the raw bytes.
		bytes, err = io.ReadAll(file)
		if err != nil {
			return nil, err
		}
	case ".axb":
		// For binary format, read the raw bytes.
		bytes, err = io.ReadAll(file)
		if err != nil {
			return nil, err
		}
	case ".axc":
		// For compressed format, decompress first.
		xzReader, err := xz.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("failed to create xz reader for .axc file: %w", err)
		}
		bytes, err = io.ReadAll(xzReader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress .axc file: %w", err)
		}
		// The decompressed bytes are in .axb format, fall through to the proto.Unmarshal logic below.
		fileExt = ".axb" // Treat decompressed data as binary
	default:
		return nil, fmt.Errorf("unsupported file extension '%s': must be .ax, .axb, .axd, or .axc", fileExt)
	}

	// Now, unmarshal the bytes based on the (potentially updated) format.
	switch fileExt {
	case ".ax":
		unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
		if err := unmarshaler.Unmarshal(bytes, &graph); err != nil {
			return nil, fmt.Errorf("failed to parse .ax (JSON) file: %w", err)
		}
	case ".axb":
		if err := proto.Unmarshal(bytes, &graph); err != nil {
			return nil, fmt.Errorf("failed to parse .axb (binary) file: %w", err)
		}
	case ".axd":
		if err := yaml.Unmarshal(bytes, &graph); err != nil {
			return nil, fmt.Errorf("failed to parse .axd (YAML) file: %w", err)
		}
	}

	return &graph, nil
}
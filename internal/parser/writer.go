package parser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Advik-B/Axon/internal/debug"
	"github.com/ulikunitz/xz"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

// SaveGraphToFile saves a given Graph struct to a file, automatically
// selecting the correct format based on the output file's extension.
func SaveGraphToFile(graph *Graph, filePath string) error {
	fileExt := filepath.Ext(filePath)
	var outputBytes []byte
	var err error

	switch fileExt {
	case ".ax":
		jsonMarshaler := protojson.MarshalOptions{Indent: "  ", UseProtoNames: true}
		outputBytes, err = jsonMarshaler.Marshal(graph)
		if err != nil {
			return fmt.Errorf("failed to marshal to .ax (JSON): %w", err)
		}

	case ".axb":
		outputBytes, err = proto.Marshal(graph)
		if err != nil {
			return fmt.Errorf("failed to marshal to .axb (binary): %w", err)
		}

	case ".axd":
		debugGraph := debug.GenerateDebugGraph(graph)
		outputBytes, err = yaml.Marshal(debugGraph)
		if err != nil {
			return fmt.Errorf("failed to marshal to .axd (YAML): %w", err)
		}

	case ".axc":
		// First, marshal to the intermediate binary (.axb) format.
		binaryData, err := proto.Marshal(graph)
		if err != nil {
			return fmt.Errorf("failed to marshal to intermediate binary for .axc: %w", err)
		}

		// Now, compress the binary data using XZ.
		var compressedBuf bytes.Buffer
		xzWriter, err := xz.NewWriter(&compressedBuf)
		if err != nil {
			return fmt.Errorf("failed to create xz writer for .axc: %w", err)
		}

		if _, err := xzWriter.Write(binaryData); err != nil {
			return fmt.Errorf("failed to write compressed data: %w", err)
		}
		// It's crucial to close the writer to flush the stream.
		if err := xzWriter.Close(); err != nil {
			return fmt.Errorf("failed to finalize xz stream: %w", err)
		}

		outputBytes = compressedBuf.Bytes()

	default:
		return fmt.Errorf("unsupported output file extension '%s': must be .ax, .axb, .axd, or .axc", fileExt)
	}

	return os.WriteFile(filePath, outputBytes, 0644)
}
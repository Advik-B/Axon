package parser

import (
	"encoding/json"
	"io"
	"os"

	"github.com/Advik-B/Axon/pkg/axon"
	"google.golang.org/protobuf/encoding/protojson"
)

// LoadGraphFromFile reads an .ax file and parses it into an Axon Graph struct.
// It uses protojson to unmarshal directly into the protobuf-generated structs.
func LoadGraphFromFile(filePath string) (*axon.Graph, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var graph axon.Graph
	// Use a custom unmarshaler to handle enums as strings
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	if err := unmarshaler.Unmarshal(bytes, &graph); err != nil {
		// Fallback for simple JSON in case protojson fails, for robustness.
		if errJson := json.Unmarshal(bytes, &graph); errJson != nil {
			return nil, err // Return the original protojson error
		}
	}

	return &graph, nil
}
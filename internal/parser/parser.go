package parser

import (
	"encoding/json"
	"os"
	"io"
	"github.com/Advik-B/Axon/pkg/axon" // Use your actual go.mod path here
)

// LoadGraphFromFile reads an .ax file and parses it into an Axon Graph struct.
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
	if err := json.Unmarshal(bytes, &graph); err != nil {
		return nil, err
	}

	return &graph, nil
}
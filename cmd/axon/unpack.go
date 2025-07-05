package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"github.com/Advik-B/Axon/internal/parser"
)

// unpackCmd represents the unpack command
var unpackCmd = &cobra.Command{
	Use:   "unpack [path/to/graph.axb]",
	Short: "Unpacks a binary .axb file into a JSON-based .ax file.",
	Long: `Reads a binary .axb file and decodes it into the human-readable .ax (JSON)
format. This is useful for debugging, manual editing, or version control.`,
	Args: cobra.ExactArgs(1),
	Run:  runUnpack,
}

// Graph is an alias to the generated struct for use within this command.
type Graph = parser.Graph

func runUnpack(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// 1. Validate file exists and has the correct extension
	if !strings.HasSuffix(filePath, ".axb") {
		fmt.Println("❌ Error: Input file must be a .axb (binary) file.")
		os.Exit(1)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("❌ Error: Input file not found at '%s'\n", filePath)
		os.Exit(1)
	}
	fmt.Printf("📂 Unpacking file: %s\n", filePath)

	// 2. Read the binary .axb file
	binaryData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error reading binary file: %v\n", err)
		os.Exit(1)
	}

	// 3. Unmarshal the binary data into a graph struct
	var graph Graph
	if err := proto.Unmarshal(binaryData, &graph); err != nil {
		fmt.Printf("❌ Error unmarshaling binary data into graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Successfully parsed binary graph.")

	// 4. Marshal the graph struct into pretty-printed JSON
	jsonMarshaler := protojson.MarshalOptions{
		Indent: "  ",
		UseProtoNames: true,
	}
	jsonData, err := jsonMarshaler.Marshal(&graph)
	if err != nil {
		fmt.Printf("❌ Error marshaling graph to JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Serialized graph to JSON format.")

	// 5. Write the output to a .ax file
	outputPath := strings.TrimSuffix(filePath, ".axb") + ".ax"
	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing to output file %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Successfully unpacked graph to %s\n", outputPath)
}

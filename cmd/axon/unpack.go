package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Advik-B/Axon/internal/parser"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// unpackCmd represents the unpack command
var unpackCmd = &cobra.Command{
	Use:   "unpack [path/to/graph.axb | path/to/graph.axd]",
	Short: "Unpacks a binary .axb or debug .axd file into a JSON .ax file.",
	Long: `Reads a binary .axb or a debug .axd file and decodes it into the
human-readable .ax (JSON) format. This is useful for inspection or conversion.`,
	Args: cobra.ExactArgs(1),
	Run:  runUnpack,
}

func runUnpack(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// Validate file exists and has a supported extension
	if !strings.HasSuffix(filePath, ".axb") && !strings.HasSuffix(filePath, ".axd") {
		fmt.Println("❌ Error: Input file must be a .axb (binary) or .axd (debug YAML) file.")
		os.Exit(1)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("❌ Error: Input file not found at '%s'\n", filePath)
		os.Exit(1)
	}
	fmt.Printf("📂 Unpacking file: %s\n", filePath)

	// Load the graph. The parser now understands .axb and .axd.
	graph, err := parser.LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error parsing graph file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Successfully parsed source graph.")

	// Marshal the graph struct into pretty-printed JSON
	jsonMarshaler := protojson.MarshalOptions{
		Indent:        "  ",
		UseProtoNames: true,
	}
	jsonData, err := jsonMarshaler.Marshal(graph)
	if err != nil {
		fmt.Printf("❌ Error marshaling graph to JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Serialized graph to JSON format.")

	// Determine output path
	outputPath := strings.TrimSuffix(filePath, ".axb")
	outputPath = strings.TrimSuffix(outputPath, ".axd") + ".ax"

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing to output file %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Successfully unpacked graph to %s\n", outputPath)
}
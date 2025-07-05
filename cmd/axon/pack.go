package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Advik-B/Axon/internal/parser"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   "pack [path/to/graph.ax]",
	Short: "Packs a JSON-based .ax file into a binary .axb file.",
	Long: `Reads a human-readable .ax (JSON) file and serializes it into the highly
efficient binary Protobuf format, saving it as an .axb file. This is ideal
for faster loading and smaller storage.`,
	Args: cobra.ExactArgs(1),
	Run:  runPack,
}

func runPack(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// 1. Validate file exists and has the correct extension
	if !strings.HasSuffix(filePath, ".ax") {
		fmt.Println("❌ Error: Input file must be a .ax (JSON) file.")
		os.Exit(1)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("❌ Error: Input file not found at '%s'\n", filePath)
		os.Exit(1)
	}
	fmt.Printf("📦 Packing file: %s\n", filePath)

	// 2. Parse the source .ax file (which is JSON)
	graph, err := parser.LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error parsing graph file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	fmt.Println("   - Successfully parsed JSON graph.")

	// 3. Marshal the graph struct into binary Protobuf format
	binaryData, err := proto.Marshal(graph)
	if err != nil {
		fmt.Printf("❌ Error marshaling graph to binary format: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Serialized graph to binary format.")

	// 4. Write the output to a .axb file
	outputPath := strings.TrimSuffix(filePath, ".ax") + ".axb"
	err = os.WriteFile(outputPath, binaryData, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing to output file %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Successfully packed graph to %s\n", outputPath)
}
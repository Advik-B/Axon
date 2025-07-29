package main

import (
	"fmt"
	parser2 "github.com/Advik-B/Axon/parser"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// unpackCmd represents the unpack command
var unpackCmd = &cobra.Command{
	Use:   "unpack [path/to/graph.axc]",
	Short: "Unpacks a compressed .axc file into a binary .axb file.",
	Long: `Reads a compressed .axc file and decompresses it into the uncompressed
binary .axb format.`,
	Args: cobra.ExactArgs(1),
	Run:  runUnpack,
}

func runUnpack(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// 1. Validate input file type.
	if !strings.HasSuffix(filePath, ".axc") {
		fmt.Println("❌ Error: Input file for unpacking must be a .axc file.")
		os.Exit(1)
	}

	// 2. Load the graph from the compressed format.
	fmt.Printf("📂 Decompressing source graph: %s\n", filePath)
	graph, err := parser2.LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error reading compressed file: %v\n", err)
		os.Exit(1)
	}

	// 3. Determine the output path.
	baseName := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	outputPath := baseName + ".axb"
	fmt.Printf("   -> Saving to binary: %s\n", outputPath)

	// 4. Save the graph to the binary .axb format.
	err = parser2.SaveGraphToFile(graph, outputPath)
	if err != nil {
		fmt.Printf("❌ Error writing binary file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Successfully unpacked graph to %s\n", outputPath)
}

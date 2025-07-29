package main

import (
	"fmt"
	parser2 "github.com/Advik-B/Axon/parser"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   "pack [path/to/graph.ax | .axd | .axb]",
	Short: "Packs any Axon graph file into a compressed .axc file.",
	Long: `Reads any valid Axon graph format (.ax, .axd, .axb) and compresses it
into the highly efficient .axc format using XZ compression.`,
	Args: cobra.ExactArgs(1),
	Run:  runPack,
}

func runPack(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// 1. Load the graph from any supported format. The parser handles the complexity.
	fmt.Printf("📦 Reading source graph: %s\n", filePath)
	graph, err := parser2.LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error parsing graph file: %v\n", err)
		os.Exit(1)
	}

	// 2. Determine the output path.
	baseName := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	outputPath := baseName + ".axc"
	fmt.Printf("   -> Compressing to: %s\n", outputPath)

	// 3. Save the graph to the .axc format. The writer handles the marshal-then-compress logic.
	err = parser2.SaveGraphToFile(graph, outputPath)
	if err != nil {
		fmt.Printf("❌ Error writing compressed file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Successfully packed graph to %s\n", outputPath)
}

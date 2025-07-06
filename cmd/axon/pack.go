package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Advik-B/Axon/internal/debug"
	"github.com/Advik-B/Axon/internal/parser"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

func init() {
	// Add the --debug flag to the pack command
	packCmd.Flags().BoolP("debug", "d", false, "Pack to a human-readable YAML debug file (.axd)")
}

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   "pack [path/to/graph.ax]",
	Short: "Packs a JSON .ax file into a binary .axb or debug .axd file.",
	Long: `Reads a human-readable .ax (JSON) file and serializes it.
- Default: Packs to a highly efficient binary .axb file.
- With --debug flag: Packs to a commented, human-readable YAML .axd file.`,
	Args: cobra.ExactArgs(1),
	Run:  runPack,
}

func runPack(cmd *cobra.Command, args []string) {
	filePath := args[0]
	useDebugFormat, _ := cmd.Flags().GetBool("debug")

	// Validate file exists and has the correct extension
	if !strings.HasSuffix(filePath, ".ax") {
		fmt.Println("❌ Error: Input file for packing must be a .ax (JSON) file.")
		os.Exit(1)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("❌ Error: Input file not found at '%s'\n", filePath)
		os.Exit(1)
	}

	// Parse the source .ax file
	graph, err := parser.LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error parsing graph file %s: %v\n", filePath, err)
		os.Exit(1)
	}

	if useDebugFormat {
		packToDebug(graph, filePath)
	} else {
		packToBinary(graph, filePath)
	}
}

func packToBinary(graph *parser.Graph, sourcePath string) {
	fmt.Printf("📦 Packing to binary format: %s\n", sourcePath)

	binaryData, err := proto.Marshal(graph)
	if err != nil {
		fmt.Printf("❌ Error marshaling graph to binary format: %v\n", err)
		os.Exit(1)
	}

	outputPath := strings.TrimSuffix(sourcePath, ".ax") + ".axb"
	err = os.WriteFile(outputPath, binaryData, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing to output file %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Successfully packed graph to %s\n", outputPath)
}

func packToDebug(graph *parser.Graph, sourcePath string) {
	fmt.Printf("📝 Packing to debug (YAML) format: %s\n", sourcePath)

	// Generate the special debug struct with comments
	debugGraph := debug.GenerateDebugGraph(graph)

	// Marshal to YAML using yaml.v3, which will render the comments
	yamlData, err := yaml.Marshal(debugGraph)
	if err != nil {
		fmt.Printf("❌ Error marshaling graph to YAML format: %v\n", err)
		os.Exit(1)
	}

	outputPath := strings.TrimSuffix(sourcePath, ".ax") + ".axd"
	err = os.WriteFile(outputPath, yamlData, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing to output file %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Successfully packed debug graph to %s\n", outputPath)
}
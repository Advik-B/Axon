package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Advik-B/Axon/internal/parser"
	"github.com/spf13/cobra"
)

func init() {
	convertCmd.Flags().StringP("input", "i", "", "Input graph file (.ax, .axb, .axd, .axc)")
	convertCmd.Flags().StringP("output", "o", "", "Output graph file (.ax, .axb, .axd, .axc)")
}

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [input] [output]",
	Short: "Converts between Axon graph formats (.ax, .axb, .axd, .axc).",
	Long: `A flexible utility to convert Axon graphs between the human-readable JSON (.ax),
the efficient binary (.axb), the commented debug YAML (.axd), and the
compressed binary (.axc) formats.`,
	Run: runConvert,
}

func runConvert(cmd *cobra.Command, args []string) {
	var inputPath, outputPath string
	var err error

	inputFlag, _ := cmd.Flags().GetString("input")
	outputFlag, _ := cmd.Flags().GetString("output")
	hasFlags := inputFlag != "" || outputFlag != ""
	hasPositionalArgs := len(args) > 0

	if hasFlags && hasPositionalArgs {
		fmt.Println("❌ Error: Cannot use positional arguments and flags (-i, -o) at the same time.")
		os.Exit(1)
	}

	if hasFlags {
		if inputFlag == "" || outputFlag == "" {
			fmt.Println("❌ Error: Both --input (-i) and --output (-o) flags are required when using flags.")
			os.Exit(1)
		}
		inputPath = inputFlag
		outputPath = outputFlag
	} else if len(args) == 2 {
		inputPath = args[0]
		outputPath = args[1]
	} else {
		fmt.Println("❌ Error: Invalid arguments. Use 'axon convert <input> <output>' or 'axon convert -i <input> -o <output>'.")
		os.Exit(1)
	}

	if inputPath == outputPath {
		fmt.Println("❌ Error: Input and output file paths cannot be the same.")
		os.Exit(1)
	}

	// Updated validation logic
	validExts := map[string]bool{".ax": true, ".axb": true, ".axd": true, ".axc": true}
	inputExt := filepath.Ext(inputPath)
	outputExt := filepath.Ext(outputPath)

	if !validExts[inputExt] {
		fmt.Printf("❌ Error: Invalid input file type '%s'. Must be .ax, .axb, .axd, or .axc.\n", inputExt)
		os.Exit(1)
	}
	if outputExt == ".go" {
		fmt.Println("❌ Error: Cannot convert to .go. Please use the 'axon transpile' command instead.")
		os.Exit(1)
	}
	if !validExts[outputExt] {
		fmt.Printf("❌ Error: Invalid output file type '%s'. Must be .ax, .axb, .axd, or .axc.\n", outputExt)
		os.Exit(1)
	}
	if _, err = os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Printf("❌ Error: Input file not found at '%s'\n", inputPath)
		os.Exit(1)
	}

	fmt.Printf("🔄 Converting %s -> %s...\n", inputPath, outputPath)

	graph, err := parser.LoadGraphFromFile(inputPath)
	if err != nil {
		fmt.Printf("❌ Error reading input file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Successfully parsed input graph.")

	err = parser.SaveGraphToFile(graph, outputPath)
	if err != nil {
		fmt.Printf("❌ Error writing output file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Successfully wrote output graph.")

	fmt.Println("\n✅ Conversion complete.")
}
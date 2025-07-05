package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Advik-B/Axon/internal/parser"
	"github.com/Advik-B/Axon/internal/transpiler"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [path/to/graph.ax]",
	Short: "Transpiles an Axon graph (.ax) file to Go code.",
	Long: `Build reads an .ax file, validates its structure, and transpiles it into
a runnable Go program located in the 'out' directory.

It checks for valid execution paths, explicit error handling, and type consistency
before generating the final output.`,
	Args: cobra.ExactArgs(1), // Requires exactly one argument: the file path.
	Run:  runBuild,
}

// runBuild contains the sequential logic for the build process.
func runBuild(cmd *cobra.Command, args []string) {
	filePath := args[0]
	startTime := time.Now()

	fmt.Println("🚀 Starting Axon build process...")

	// 1. Validate file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("❌ Error: Input file not found at '%s'\n", filePath)
		os.Exit(1)
	}
	fmt.Printf("   - Found graph file: %s\n", filePath)

	// 2. Parse the .ax file
	fmt.Println("   - Parsing graph...")
	graph, err := parser.LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error parsing graph file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	fmt.Printf("   - Successfully parsed graph: %s\n", graph.Name)

	// 3. Transpile the graph to Go code
	fmt.Println("   - Transpiling to Go...")
	goCode, err := transpiler.Transpile(graph)
	if err != nil {
		fmt.Printf("❌ Error transpiling graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Transpilation successful.")

	// 4. Write the output to a file
	fmt.Println("   - Writing output file...")
	outputDir := "out"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("❌ Error creating output directory %s: %v\n", outputDir, err)
		os.Exit(1)
	}

	outputFile := filepath.Join(outputDir, "main.go")
	err = os.WriteFile(outputFile, []byte(goCode), 0644)
	if err != nil {
		fmt.Printf("❌ Error writing to output file %s: %v\n", outputFile, err)
		os.Exit(1)
	}
	fmt.Printf("   - Go code written to %s\n", outputFile)

	duration := time.Since(startTime)
	fmt.Printf("\n✅ Build Succeeded in %.2fs!\n", duration.Seconds())
	fmt.Printf("   Run the output with: go run %s\n", outputFile)
}
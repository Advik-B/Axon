package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Advik-B/Axon/internal/parser"
	"github.com/Advik-B/Axon/internal/renderer"
	"github.com/fogleman/gg"
	"github.com/spf13/cobra"
)

func init() {
	// Add --output flag to the preview command
	previewCmd.Flags().StringP("output", "o", "", "Output PNG file path (default: <graph_name>.png)")
}

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview [path/to/graph.ax|axb|axd]",
	Short: "Generates a PNG image preview of an Axon graph.",
	Long: `Reads any Axon graph file format and uses its visual metadata to render
a high-quality PNG image. This allows for a quick visual inspection of the
graph's structure without needing a GUI.`,
	Args: cobra.ExactArgs(1),
	Run:  runPreview,
}

func runPreview(cmd *cobra.Command, args []string) {
	filePath := args[0]
	fmt.Printf("🖼️  Generating preview for: %s\n", filePath)

	// 1. Load the graph from file (parser handles all formats)
	graph, err := parser.LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error parsing graph file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Graph loaded successfully.")

	// 2. Generate the image in memory
	img, err := renderer.GenerateImage(graph)
	if err != nil {
		fmt.Printf("❌ Error rendering image: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Image rendered successfully.")

	// 3. Determine output path
	outputPath, _ := cmd.Flags().GetString("output")
	if outputPath == "" {
		baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
		outputPath = baseName + ".png"
	}
	fmt.Printf("   - Saving image to: %s\n", outputPath)

	// 4. Save the image to a PNG file
	if err := gg.SavePNG(outputPath, img); err != nil {
		fmt.Printf("❌ Error saving PNG file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Preview saved successfully to %s\n", outputPath)
}
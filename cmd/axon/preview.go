package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Advik-B/Axon/internal/parser"
	"github.com/Advik-B/Axon/internal/previewer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview [path/to/graph.ax | .axb | .axd | .axc]",
	Short: "Generates a visual preview of an Axon graph in a window.",
	Long: `Reads any valid Axon graph file and generates a visual representation of it.

This command uses an automatic layout algorithm to arrange the nodes for maximum
clarity, ignoring any stored visual information. The viewer is interactive,
supporting panning (mouse drag) and zooming (mouse wheel).`,
	Args: cobra.ExactArgs(1),
	Run:  runPreview,
}

func runPreview(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// 1. Load the graph from any supported format.
	fmt.Printf("🔎 Loading graph for preview: %s\n", filePath)
	graph, err := parser.LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error parsing graph file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Graph loaded successfully. Launching preview window...")

	// 2. Initialize the Ebitengine previewer with the graph data.
	previewApp, err := previewer.NewPreviewer(graph)
	if err != nil {
		log.Fatalf("❌ Failed to initialize previewer: %v", err)
	}

	// 3. Configure and run the Ebitengine window.
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle(fmt.Sprintf("Axon Preview - %s", graph.Name))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(previewApp); err != nil {
		log.Fatalf("❌ Ebitengine exited with an error: %v", err)
	}

	fmt.Println("👋 Preview window closed.")
}
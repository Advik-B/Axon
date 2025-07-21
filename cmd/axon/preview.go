package main

import (
	"fmt"
	"github.com/Advik-B/Axon/parser"
	"github.com/Advik-B/Axon/previewer"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview [path/to/graph.ax | .axb | .axd | .axc]",
	Short: "Generates a dynamic, physics-based preview of an Axon graph.",
	Long: `Reads any valid Axon graph file and generates a dynamic visual representation.

This viewer uses a spring-mass physics simulation to create a 'fluid' or 'blob-like'
feel. Nodes are connected by springs and will react to being moved.

Controls:
  - Drag Node:  Click and drag a node to move it.
  - Pan View:   Click and drag the background.
  - Zoom View:  Use the mouse wheel.`,
	Args: cobra.ExactArgs(1),
	Run:  runPreview,
}

func runPreview(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// 1. Load the graph from any supported format.
	fmt.Printf("🔎 Loading graph for physics preview: %s\n", filePath)
	graph, err := parser.LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error parsing graph file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   - Graph loaded. Launching interactive physics preview...")

	// 2. Initialize the Ebitengine previewer with the graph data.
	previewApp, err := previewer.NewPreviewer(graph)
	if err != nil {
		log.Fatalf("❌ Failed to initialize previewer: %v", err)
	}

	// 3. Configure and run the Ebitengine window.
	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowTitle(fmt.Sprintf("Axon Physics Preview - %s", graph.Name))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(previewApp); err != nil {
		log.Fatalf("❌ Ebitengine exited with an error: %v", err)
	}

	fmt.Println("👋 Preview window closed.")
}

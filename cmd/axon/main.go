package main

import (
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "build" {
		fmt.Println("Usage: axon build <path/to/graph.ax>")
		os.Exit(1)
	}

	filePath := os.Args[2]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Error: File not found at '%s'\n", filePath)
		os.Exit(1)
	}

	model := initialModel(filePath)
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
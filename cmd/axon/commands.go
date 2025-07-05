package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Advik-B/Axon/internal/parser"
	"github.com/Advik-B/Axon/internal/transpiler"
	"github.com/Advik-B/Axon/pkg/axon"
	tea "github.com/charmbracelet/bubbletea"
)

// --- Custom Messages ---
// These messages are sent from our commands back to the model's Update function.

type parseResultMsg struct {
	graph    *axon.Graph
	duration time.Duration
	err      error
}

type transpileResultMsg struct {
	code     string
	duration time.Duration
	err      error
}

type writeResultMsg struct {
	duration time.Duration
	err      error
}

// --- Commands ---
// These functions perform the work and return a message.

func runParseStage(filePath string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		graph, err := parser.LoadGraphFromFile(filePath)
		return parseResultMsg{
			graph:    graph,
			duration: time.Since(start),
			err:      err,
		}
	}
}

func runTranspileStage(graph *axon.Graph) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		code, err := transpiler.Transpile(graph)
		return transpileResultMsg{
			code:     code,
			duration: time.Since(start),
			err:      err,
		}
	}
}

func runWriteStage(goCode string, outputPath string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return writeResultMsg{duration: time.Since(start), err: err}
		}

		err := os.WriteFile(outputPath, []byte(goCode), 0644)
		return writeResultMsg{
			duration: time.Since(start),
			err:      err,
		}
	}
}
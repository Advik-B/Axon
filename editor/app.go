package main

import (
    "context"
    "github.com/Advik-B/Axon/parser"
    "github.com/Advik-B/Axon/pkg/axon"
    "github.com/Advik-B/Axon/transpiler"
    "github.com/wailsapp/wails/v2/pkg/runtime" // For dialogs
    "google.golang.org/protobuf/proto"
)

// App struct (no changes)
type App struct {
    ctx   context.Context
    graph *axon.Graph
}

// LoadGraph prompts for a file, parses it, and returns the SERIALIZED graph
func (a *App) LoadGraph() ([]byte, error) {
    // 1. Open a file dialog
    filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
        Title: "Open Axon Graph",
        Filters: []runtime.FileFilter{
            {DisplayName: "Axon Files (*.ax, *.axb, ...)", Pattern: "*.ax;*.axb;*.axd;*.axc"},
        },
    })
    if err != nil {
        return nil, err
    }
    if filePath == "" { // User cancelled
        return nil, nil
    }

    // 2. Parse the file into the Go struct
    graph, err := parser.LoadGraphFromFile(filePath)
    if err != nil {
        return nil, err
    }
    a.graph = graph // Update backend state

    // 3. Marshal the struct into binary format and return the bytes
    return proto.Marshal(a.graph)
}

// SaveGraph receives SERIALIZED graph data, unmarshals it, and saves it.
func (a *App) SaveGraph(graphBytes []byte) (string, error) {
    // 1. Open a save file dialog
    filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
        Title:           "Save Axon Graph",
        DefaultFilename: "graph.ax",
    })
    if err != nil {
        return "", err
    }
    if filePath == "" { // User cancelled
        return "", nil
    }

    // 2. Unmarshal the bytes from the frontend into the Go struct
    var graphToSave axon.Graph
    if err := proto.Unmarshal(graphBytes, &graphToSave); err != nil {
        return "", err
    }

    // 3. Save the Go struct to a file (using your existing parser)
    if err := parser.SaveGraphToFile(&graphToSave, filePath); err != nil {
        return "", err
    }
    a.graph = &graphToSave // Update backend state

    return "Graph saved successfully!", nil
}

// TranspileCurrentGraph receives serialized data and returns code
func (a *App) TranspileCurrentGraph(graphBytes []byte) (string, error) {
    var graphToTranspile axon.Graph
    if err := proto.Unmarshal(graphBytes, &graphToTranspile); err != nil {
        return "", err
    }
    return transpiler.Transpile(&graphToTranspile)
}

func NewApp() *App {
    return &App{}
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Initialize with a default empty graph
    a.graph = &axon.Graph{
        Id:   "new-graph",
        Name: "New Graph",
    }
}

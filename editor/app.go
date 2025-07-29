// editor/app.go
package main

import (
    "context"
    "github.com/Advik-B/Axon/parser"
    "github.com/Advik-B/Axon/pkg/axon"
    "github.com/Advik-B/Axon/transpiler"
)

// App struct
type App struct {
    ctx   context.Context
    graph *axon.Graph // Holds the state of the currently loaded graph
}

// NewApp creates a new App application struct
func NewApp() *App {
    return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Initialize with a default empty graph
    a.graph = &axon.Graph{
        Id:   "new-graph",
        Name: "New Graph",
    }
}

// LoadGraph prompts the user to open a file and returns the parsed graph.
// NOTE: Wails has runtime dialogs for this.
func (a *App) LoadGraph(filePath string) (*axon.Graph, error) {
    graph, err := parser.LoadGraphFromFile(filePath)
    if err != nil {
        return nil, err
    }
    a.graph = graph // Update backend state
    return a.graph, nil
}

// SaveGraph saves the current graph state to a file.
func (a *App) SaveGraph(filePath string, graphData *axon.Graph) error {
    a.graph = graphData // Update backend state from frontend
    return parser.SaveGraphToFile(a.graph, filePath)
}

// GetInitialGraph returns the default graph on startup.
func (a *App) GetInitialGraph() *axon.Graph {
    return a.graph
}

// TranspileCurrentGraph transpiles the graph held in state.
func (a *App) TranspileCurrentGraph(graphData *axon.Graph) (string, error) {
    a.graph = graphData // Ensure backend state is up-to-date
    return transpiler.Transpile(a.graph)
}

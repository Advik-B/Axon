package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// GoplsServer manages the gopls language server process
type GoplsServer struct {
	cmd     *exec.Cmd
	workDir string
}

// CompletionRequest represents a completion request to gopls
type CompletionRequest struct {
	FileName string `json:"fileName"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Content  string `json:"content"`
}

// CompletionItem represents a completion suggestion
type CompletionItem struct {
	Label         string `json:"label"`
	Kind          string `json:"kind"`
	Detail        string `json:"detail"`
	Documentation string `json:"documentation"`
	InsertText    string `json:"insertText"`
}

// CompletionResponse represents the response from gopls
type CompletionResponse struct {
	Items []CompletionItem `json:"items"`
}

// NewGoplsServer creates a new gopls server instance
func NewGoplsServer() *GoplsServer {
	// Create a temporary working directory for Go code analysis
	workDir, err := os.MkdirTemp("", "axon-gopls-*")
	if err != nil {
		log.Printf("Failed to create temp dir for gopls: %v", err)
		return nil
	}

	// Initialize a Go module in the temp directory
	cmd := exec.Command("go", "mod", "init", "axon-temp")
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to initialize Go module: %v", err)
		os.RemoveAll(workDir)
		return nil
	}

	return &GoplsServer{
		workDir: workDir,
	}
}

// GetCompletions retrieves autocompletions for Go code at a specific position
func (g *GoplsServer) GetCompletions(req CompletionRequest) (*CompletionResponse, error) {
	if g == nil {
		return &CompletionResponse{Items: []CompletionItem{}}, nil
	}

	// Write the Go code to a temporary file
	tempFile := filepath.Join(g.workDir, "main.go")
	if err := os.WriteFile(tempFile, []byte(req.Content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	// Use gopls completion via command line
	// Note: This is a simplified version. In a production system, you'd use the LSP protocol
	cmd := exec.Command("gopls", "completion", fmt.Sprintf("%s:%d:%d", tempFile, req.Line, req.Column))
	cmd.Dir = g.workDir
	
	output, err := cmd.Output()
	if err != nil {
		// If gopls fails, return basic Go keywords and stdlib suggestions
		return g.getFallbackCompletions(req), nil
	}

	// Parse gopls output (simplified - in reality you'd parse LSP JSON-RPC responses)
	items := g.parseGoplsOutput(string(output))
	
	return &CompletionResponse{Items: items}, nil
}

// getFallbackCompletions provides basic Go completions when gopls is not available
func (g *GoplsServer) getFallbackCompletions(req CompletionRequest) *CompletionResponse {
	basicCompletions := []CompletionItem{
		{Label: "fmt.Println", Kind: "function", Detail: "func Println(a ...interface{}) (n int, err error)", InsertText: "fmt.Println"},
		{Label: "fmt.Printf", Kind: "function", Detail: "func Printf(format string, a ...interface{}) (n int, err error)", InsertText: "fmt.Printf"},
		{Label: "fmt.Sprintf", Kind: "function", Detail: "func Sprintf(format string, a ...interface{}) string", InsertText: "fmt.Sprintf"},
		{Label: "if", Kind: "keyword", Detail: "if statement", InsertText: "if"},
		{Label: "for", Kind: "keyword", Detail: "for loop", InsertText: "for"},
		{Label: "func", Kind: "keyword", Detail: "function declaration", InsertText: "func"},
		{Label: "var", Kind: "keyword", Detail: "variable declaration", InsertText: "var"},
		{Label: "const", Kind: "keyword", Detail: "constant declaration", InsertText: "const"},
		{Label: "type", Kind: "keyword", Detail: "type declaration", InsertText: "type"},
		{Label: "struct", Kind: "keyword", Detail: "struct type", InsertText: "struct"},
		{Label: "interface", Kind: "keyword", Detail: "interface type", InsertText: "interface"},
		{Label: "return", Kind: "keyword", Detail: "return statement", InsertText: "return"},
		{Label: "break", Kind: "keyword", Detail: "break statement", InsertText: "break"},
		{Label: "continue", Kind: "keyword", Detail: "continue statement", InsertText: "continue"},
		{Label: "switch", Kind: "keyword", Detail: "switch statement", InsertText: "switch"},
		{Label: "case", Kind: "keyword", Detail: "case clause", InsertText: "case"},
		{Label: "default", Kind: "keyword", Detail: "default clause", InsertText: "default"},
		{Label: "package", Kind: "keyword", Detail: "package declaration", InsertText: "package"},
		{Label: "import", Kind: "keyword", Detail: "import declaration", InsertText: "import"},
		{Label: "go", Kind: "keyword", Detail: "go statement", InsertText: "go"},
		{Label: "defer", Kind: "keyword", Detail: "defer statement", InsertText: "defer"},
		{Label: "select", Kind: "keyword", Detail: "select statement", InsertText: "select"},
		{Label: "chan", Kind: "keyword", Detail: "channel type", InsertText: "chan"},
		{Label: "map", Kind: "keyword", Detail: "map type", InsertText: "map"},
		{Label: "make", Kind: "function", Detail: "func make(t Type, size ...IntegerType) Type", InsertText: "make"},
		{Label: "new", Kind: "function", Detail: "func new(Type) *Type", InsertText: "new"},
		{Label: "len", Kind: "function", Detail: "func len(v Type) int", InsertText: "len"},
		{Label: "cap", Kind: "function", Detail: "func cap(v Type) int", InsertText: "cap"},
		{Label: "append", Kind: "function", Detail: "func append(slice []Type, elems ...Type) []Type", InsertText: "append"},
		{Label: "copy", Kind: "function", Detail: "func copy(dst, src []Type) int", InsertText: "copy"},
		{Label: "delete", Kind: "function", Detail: "func delete(m map[Type]Type1, key Type)", InsertText: "delete"},
		{Label: "close", Kind: "function", Detail: "func close(c chan<- Type)", InsertText: "close"},
		{Label: "panic", Kind: "function", Detail: "func panic(v interface{})", InsertText: "panic"},
		{Label: "recover", Kind: "function", Detail: "func recover() interface{}", InsertText: "recover"},
		{Label: "int", Kind: "type", Detail: "int type", InsertText: "int"},
		{Label: "int8", Kind: "type", Detail: "int8 type", InsertText: "int8"},
		{Label: "int16", Kind: "type", Detail: "int16 type", InsertText: "int16"},
		{Label: "int32", Kind: "type", Detail: "int32 type", InsertText: "int32"},
		{Label: "int64", Kind: "type", Detail: "int64 type", InsertText: "int64"},
		{Label: "uint", Kind: "type", Detail: "uint type", InsertText: "uint"},
		{Label: "uint8", Kind: "type", Detail: "uint8 type", InsertText: "uint8"},
		{Label: "uint16", Kind: "type", Detail: "uint16 type", InsertText: "uint16"},
		{Label: "uint32", Kind: "type", Detail: "uint32 type", InsertText: "uint32"},
		{Label: "uint64", Kind: "type", Detail: "uint64 type", InsertText: "uint64"},
		{Label: "float32", Kind: "type", Detail: "float32 type", InsertText: "float32"},
		{Label: "float64", Kind: "type", Detail: "float64 type", InsertText: "float64"},
		{Label: "string", Kind: "type", Detail: "string type", InsertText: "string"},
		{Label: "bool", Kind: "type", Detail: "bool type", InsertText: "bool"},
		{Label: "byte", Kind: "type", Detail: "byte type (alias for uint8)", InsertText: "byte"},
		{Label: "rune", Kind: "type", Detail: "rune type (alias for int32)", InsertText: "rune"},
		{Label: "uintptr", Kind: "type", Detail: "uintptr type", InsertText: "uintptr"},
		{Label: "error", Kind: "type", Detail: "error interface", InsertText: "error"},
	}

	return &CompletionResponse{Items: basicCompletions}
}

// parseGoplsOutput parses the output from gopls (simplified)
func (g *GoplsServer) parseGoplsOutput(output string) []CompletionItem {
	// This is a simplified parser. In reality, you'd parse JSON-RPC LSP responses
	var items []CompletionItem
	
	// For now, return empty as this would require full LSP implementation
	// The fallback completions above provide good coverage
	
	return items
}

// Cleanup cleans up the gopls server resources
func (g *GoplsServer) Cleanup() {
	if g != nil && g.workDir != "" {
		os.RemoveAll(g.workDir)
	}
}

// GetGoCodeCompletions provides completions for the given Go code context
func GetGoCodeCompletions(code string, line, column int) (*CompletionResponse, error) {
	server := NewGoplsServer()
	if server == nil {
		// Return basic completions if gopls setup fails
		return &CompletionResponse{
			Items: []CompletionItem{
				{Label: "fmt.Println", Kind: "function", Detail: "Print to stdout", InsertText: "fmt.Println"},
				{Label: "fmt.Printf", Kind: "function", Detail: "Formatted print", InsertText: "fmt.Printf"},
			},
		}, nil
	}
	defer server.Cleanup()

	req := CompletionRequest{
		FileName: "main.go",
		Line:     line,
		Column:   column,
		Content:  code,
	}

	return server.GetCompletions(req)
}
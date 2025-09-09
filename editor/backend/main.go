package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Advik-B/Axon/parser"
	"github.com/Advik-B/Axon/transpiler"
)

type Server struct {
	port string
}

type GraphResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) setupRoutes() {
	// Serve static files from frontend dist
	fs := http.FileServer(http.Dir("../frontend/dist/"))
	http.Handle("/", fs)

	// API endpoints
	http.HandleFunc("/api/graphs", s.handleGraphs)
	http.HandleFunc("/api/graphs/validate", s.handleValidateGraph)
	http.HandleFunc("/api/graphs/transpile", s.handleTranspileGraph)
	http.HandleFunc("/api/graphs/load", s.handleLoadGraph)
	http.HandleFunc("/api/graphs/save", s.handleSaveGraph)
	http.HandleFunc("/api/completions", s.handleCompletions)
	http.HandleFunc("/api/health", s.handleHealth)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := GraphResponse{
		Success: true,
		Message: "Axon editor backend is running",
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGraphs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "GET":
		// List available graph files
		s.listGraphs(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listGraphs(w http.ResponseWriter, r *http.Request) {
	graphsDir := "../../examples"
	files, err := filepath.Glob(filepath.Join(graphsDir, "*.ax"))
	if err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Error reading graphs directory: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	var graphs []string
	for _, file := range files {
		graphs = append(graphs, filepath.Base(file))
	}

	response := GraphResponse{
		Success: true,
		Message: "Graphs listed successfully",
		Data:    graphs,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleLoadGraph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := r.URL.Query().Get("file")
	if filename == "" {
		response := GraphResponse{
			Success: false,
			Message: "Missing 'file' parameter",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Load from examples directory for now
	filePath := filepath.Join("../../examples", filename)
	graph, err := parser.LoadGraphFromFile(filePath)
	if err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Error loading graph: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := GraphResponse{
		Success: true,
		Message: "Graph loaded successfully",
		Data:    graph,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleSaveGraph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var graph map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&graph); err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing graph data: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	filename := r.URL.Query().Get("file")
	if filename == "" {
		response := GraphResponse{
			Success: false,
			Message: "Missing 'file' parameter",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Save to examples directory for now
	filePath := filepath.Join("../../examples", filename)
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Error marshaling graph: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Error saving graph: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := GraphResponse{
		Success: true,
		Message: "Graph saved successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleValidateGraph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var graph map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&graph); err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing graph data: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Use enhanced validation
	validationResult, err := ValidateGraph(graph)
	if err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Graph validation failed: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := GraphResponse{
		Success: validationResult.IsValid,
		Message: "Graph validation completed",
		Data:    validationResult,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleTranspileGraph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var graph map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&graph); err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing graph data: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to JSON and back to parse as Axon graph
	data, _ := json.Marshal(graph)
	axonGraph, err := parser.LoadGraphFromBytes(data)
	if err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Graph parsing failed: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Transpile the graph
	goCode, err := transpiler.Transpile(axonGraph)
	if err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Transpilation failed: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := GraphResponse{
		Success: true,
		Message: "Transpilation successful",
		Data: map[string]interface{}{
			"go_code": goCode,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleCompletions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code   string `json:"code"`
		Line   int    `json:"line"`
		Column int    `json:"column"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing completion request: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get completions from gopls
	completions, err := GetGoCodeCompletions(req.Code, req.Line, req.Column)
	if err != nil {
		response := GraphResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get completions: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := GraphResponse{
		Success: true,
		Message: "Completions retrieved successfully",
		Data:    completions,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) Start() {
	s.setupRoutes()
	
	fmt.Printf("🚀 Axon Editor Server starting on http://localhost:%s\n", s.port)
	fmt.Printf("📊 API endpoints available at http://localhost:%s/api/*\n", s.port)
	
	if err := http.ListenAndServe(":"+s.port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	
	server := NewServer(port)
	server.Start()
}
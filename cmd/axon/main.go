package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Advik-B/Axon/internal/parser"
	"github.com/Advik-B/Axon/internal/transpiler"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: axon build <path/to/graph.ax>")
		os.Exit(1)
	}

	command := os.Args[1]
	filePath := os.Args[2]

	if command != "build" {
		log.Fatalf("Unknown command: %s. Expected 'build'.", command)
	}

	// 1. Parse the .ax file
	graph, err := parser.LoadGraphFromFile(filePath)
	if err != nil {
		log.Fatalf("Error parsing graph file %s: %v", filePath, err)
	}

	fmt.Printf("Successfully parsed graph: %s\n", graph.Name)

	// 2. Transpile the graph to Go code
	goCode, err := transpiler.Transpile(graph)
	if err != nil {
		log.Fatalf("Error transpiling graph: %v", err)
	}

	fmt.Println("Transpilation successful. Writing to output file...")

	// 3. Write the output to a file
	outputPath := "out/"
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, 0755)
	}

	outputFile := outputPath + "main.go"
	err = ioutil.WriteFile(outputFile, []byte(goCode), 0644)
	if err != nil {
		log.Fatalf("Error writing to output file %s: %v", outputFile, err)
	}

	fmt.Printf("Go code written to %s\n", outputFile)
}
package transpiler

import (
	"fmt"
	"strings"
	"github.com/Advik-B/Axon/pkg/axon"
)

// transpilationState holds the complete context during a single transpilation run.
type transpilationState struct {
	graph         *axon.Graph
	nodeMap       map[string]*axon.Node
	outputVarMap  map[string]string // Maps "nodeID.portName" -> "goVariableName"
	importManager *importManager
}

// importManager tracks and generates the import block for the Go file.
type importManager struct {
	imports map[string]struct{} // Using a map as a set for unique imports.
}

// parsePackagePath extracts a valid Go package path from a type string.
// Examples:
// "string" -> ""
// "*http.Request" -> "net/http"
// "[]*os.File" -> "os"
// "map[string]io.Reader" -> "io"
func parsePackagePath(typeStr string) string {
	// Clean up common prefixes
	typeStr = strings.TrimLeft(typeStr, "[]*& ")
	// Find the last dot, which separates the package from the type
	lastDot := strings.LastIndex(typeStr, ".")
	if lastDot == -1 {
		return "" // It's a built-in type (string, int, etc.)
	}
	// The potential package path is everything before the last dot
	potentialPath := typeStr[:lastDot]
	// If there's a space, it's likely a complex type like `map[string]io.Reader`
	lastSpace := strings.LastIndex(potentialPath, " ")
	if lastSpace != -1 {
		potentialPath = potentialPath[lastSpace+1:]
	}
	return potentialPath
}

// addImportFromType parses a Go type string and adds its package to the imports.
func (im *importManager) addImportFromType(typeStr string) {
	pkgPath := parsePackagePath(typeStr)
	if pkgPath != "" {
		im.imports[pkgPath] = struct{}{}
	}
}

// generateImportBlock creates the `import (...)` block for the final Go file.
func (im *importManager) generateImportBlock() string {
	if len(im.imports) == 0 {
		return ""
	}
	if len(im.imports) == 1 {
		for pkg := range im.imports {
			return fmt.Sprintf("import \"%s\"\n", pkg)
		}
	}
	var sb strings.Builder
	sb.WriteString("import (\n")
	for pkg := range im.imports {
		sb.WriteString(fmt.Sprintf("\t\"%s\"\n", pkg))
	}
	sb.WriteString(")\n")
	return sb.String()
}
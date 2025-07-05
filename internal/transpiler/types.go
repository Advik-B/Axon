package transpiler

import (
	"fmt"
	"path"
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

// addImportFromReference parses a Go function reference (e.g., "fmt.Println")
// and adds the required package ("fmt") to the list of imports.
// It returns the full function reference.
func (im *importManager) addImportFromReference(ref string) (string, error) {
	parts := strings.Split(ref, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid implementation reference '%s'; expected format 'package.Function'", ref)
	}
	pkgPath := findStdLibPackagePath(parts[0])
	if pkgPath == "" {
		return "", fmt.Errorf("package '%s' not found in standard library map", parts[0])
	}
	im.imports[pkgPath] = struct{}{}
	return ref, nil
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

// A simple map to find full package paths for common stdlib packages.
// This can be expanded as needed.
var stdLibPackageMap = map[string]string{
	"fmt":     "fmt",
	"os":      "os",
	"strings": "strings",
	"io":      "io",
	"http":    "net/http",
	"json":    "encoding/json",
	"math":    "math",
	"rand":    "math/rand",
	"time":    "time",
}

func findStdLibPackagePath(pkgName string) string {
	if p, ok := stdLibPackageMap[pkgName]; ok {
		return p
	}
	// Fallback for nested packages like `http.MethodGet`
	return path.Dir(pkgName)
}
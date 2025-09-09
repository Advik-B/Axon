# 🧠 Axon Visual Editor

A full-fledged visual editor for the Axon visual programming language, built with **Svelte** and **Go**.

![Axon Visual Editor](https://github.com/user-attachments/assets/d86f5393-9f66-4b8e-81a9-3b28930c8f3e)

## ✨ Features

### 🎨 Visual Graph Editor
- **Drag-and-drop interface** for creating and editing Axon graphs
- **Real-time node visualization** with color-coded node types
- **Interactive edges** for data flow and execution flow
- **Pan and zoom** support for large graphs
- **Grid-based layout** for precise positioning

### 🔧 Smart Development Tools
- **Go Language Server Integration** - Powered by `gopls` for intelligent autocompletion
- **Real-time Graph Validation** - Instant feedback on graph correctness
- **Live Code Transpilation** - See your visual graphs as Go code in real-time
- **Comprehensive Error Reporting** - Detailed validation with error types and suggestions

### 📊 Advanced Graph Analysis
- **Flow Analysis** - Validates execution paths from START to END nodes
- **Port Validation** - Ensures data connections are type-safe
- **Dependency Detection** - Identifies unused imports and isolated nodes
- **Performance Hints** - Optimization suggestions for better code generation

### 💾 Project Management
- **Multiple File Formats** - Load and save `.ax`, `.axb`, `.axd`, and `.axc` files
- **Example Graphs** - Pre-built examples to get started quickly
- **Version Control Friendly** - Human-readable JSON format for easy diff/merge

## 🚀 Quick Start

### Prerequisites

- **Go 1.24+** - [Download Here](https://golang.org/doc/install)
- **Node.js 18+** - [Download Here](https://nodejs.org/)
- **gopls** (optional) - For enhanced autocompletion

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Advik-B/Axon.git
   cd Axon/editor
   ```

2. **Run the editor:**
   ```bash
   ./start.sh
   ```

3. **Open your browser:**
   Navigate to [http://localhost:8080](http://localhost:8080)

## 🏗️ Architecture

### Backend (Go)
- **Web Server** - Serves the frontend and API endpoints
- **Graph Parser** - Handles multiple Axon file formats
- **Transpiler Integration** - Direct integration with Axon's Go transpiler
- **gopls Integration** - Language server protocol for Go autocompletion
- **Validation Engine** - Comprehensive graph validation and analysis

### Frontend (Svelte)
- **Component-based Architecture** - Modular, reusable UI components
- **D3.js Integration** - Powerful graph visualization and interaction
- **Real-time Updates** - Reactive stores for instant UI updates
- **TypeScript** - Type-safe development for better reliability

## 📁 Project Structure

```
editor/
├── backend/                 # Go backend server
│   ├── main.go             # Main server and API endpoints
│   ├── gopls.go            # gopls language server integration
│   └── validation.go       # Graph validation engine
├── frontend/               # Svelte frontend application
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/ # Svelte components
│   │   │   ├── stores/     # State management
│   │   │   └── types/      # TypeScript definitions
│   │   └── routes/         # SvelteKit routes
│   └── dist/               # Built frontend (generated)
├── start.sh                # Startup script
└── README.md               # This file
```

## 🔌 API Endpoints

### Graph Management
- `GET /api/graphs` - List available graph files
- `GET /api/graphs/load?file=name.ax` - Load a specific graph
- `POST /api/graphs/save?file=name.ax` - Save graph to file
- `POST /api/graphs/validate` - Validate graph structure
- `POST /api/graphs/transpile` - Transpile graph to Go code

### Development Tools
- `POST /api/completions` - Get Go language autocompletions
- `GET /api/health` - Server health check

## 🎯 Node Types

| Type | Description | Color | Inputs | Outputs |
|------|-------------|-------|---------|---------|
| **START** | Entry point | 🟢 Green | None | Execution |
| **END** | Exit point | 🔴 Red | Execution | None |
| **CONSTANT** | Literal values | 🔵 Blue | None | Data |
| **VARIABLE** | Stored values | 🟠 Orange | Data | Data |
| **OPERATOR** | Math/logic ops | 🟣 Purple | Data | Data |
| **FUNCTION** | Go function calls | 🔷 Cyan | Data | Data |
| **IF** | Conditional branch | 🟡 Yellow | Data | Execution |
| **LOOP** | Iteration | 🟢 Green | Data | Execution |

## 🔗 Edge Types

- **Execution Edges** (Yellow, Dashed) - Control program flow
- **Data Edges** (Blue, Solid) - Pass data between nodes

## 🛠️ Development

### Backend Development
```bash
cd backend
go run *.go
```

### Frontend Development
```bash
cd frontend
npm install
npm run dev
```

### Building for Production
```bash
cd frontend
npm run build
```

## 🎓 Usage Examples

### Creating a Simple Addition Program

1. **Add Constants** - Drag two CONSTANT nodes
2. **Add Operator** - Add an OPERATOR node with "+" operation
3. **Connect Data** - Connect constant outputs to operator inputs
4. **Add Function** - Add fmt.Println to display result
5. **Connect Execution** - Wire START → OPERATOR → FUNCTION → END
6. **Transpile** - Click "Transpile" to see generated Go code

### Loading Example Graphs

1. Click "Open Graph" in the toolbar
2. Select from available examples (add.ax, stdlib_example.ax)
3. Explore the loaded graph structure
4. Try modifying and transpiling

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Run tests: `go test ./...` (backend) and `npm test` (frontend)
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the main repository [LICENSE](../LICENSE.txt) file for details.

## 🙏 Acknowledgments

- **Axon Language** - Built on top of the excellent Axon visual programming language
- **Svelte** - Reactive frontend framework
- **D3.js** - Powerful data visualization library
- **gopls** - Go language server for intelligent code completion
- **Ebitengine** - Used in the original Axon previewer

---

🚀 **Ready to start visual programming with Go?** Launch the editor and create your first graph!
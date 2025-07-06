# Axon: Visual Programming, For Real.

<p align="center">
  <img src="https://raw.githubusercontent.com/Advik-B/Axon/main/assets/axon-banner.png" alt="Axon Banner" />
</p>

<p align="center">
  <strong>Go from a visual idea to idiomatic, runnable Go code.</strong>
  <br />
  Axon is a visual, node-based programming language that bridges the gap between intuitive, flow-based logic and high-performance, real-world applications.
  <br /><<strong>Build complex, scalable Go applications with the simplicity of a flowchart.</strong>
  <br />
  Axon is a visual, node-based programming language that transpiles directly to clean, idiomatic Go.
  <br />
  <br />
  <a href="#-features"><strong>Features</strong></a> ·
  <a href="#-how-it-works"><strong>How It Works</strong></a> ·
  <a href="#-getting-started"><strong>Getting Started</strong></a> ·
  <a href="#-the-axon-cli"><strong>CLI Commands</strong></a>
</p>

---

<p align="center">
  <em>An interactive, physics-based preview of an Axon graph. Drag nodes and watch the connections flow!</em>
  <br>
  <img src="https://i.imgur.com/your-physics-preview-animation.gif" alt="Axon Physics Preview Animation" width="700"/>
</p>

---

## 🧠 Why Axon?

Tired of wrestling with boilerplate and syntax? Axon empowers you to focus on logic and flow, not just code. It's perfect for:

-   **Rapidbr />
  <a href="#-core-features">Features</a> •
  <a href="#-interactive-preview">Live Preview</a> •
  <a href="#-the-axon-workflow">Workflow</a> •
  <a href="#-getting-started">Getting Started</a> •
  <a href="#-cli-commands">CLI Commands</a>
</p>

---

## What is Axon?

Tired of boilerplate? Teaching a new programmer? Want to *see* your logic?

**Axon** is your answer. Unlike toy block-based languages, Axon is designed for serious development. It provides a structured, visual environment to design complex systems, which it then transpiles directly into clean, human-readable, and performant Go code.

It’s the clarity of flowcharts with the power of a compiled language.

<p align="center">
  <em>We strongly recommend replacing this image with a GIF of the `axon preview` command in action!</em>
  <img src="https://raw.githubusercontent.com/Advik-B/Axon/main/assets/preview.png" alt="Axon Physics-Based Previewer" width="700"/>
</p>

## ✨ Core Features

*   🧠 **Intuitive Node-Based Logic**: Build programs by connecting nodes. Control execution flow and data flow explicitly and visually.
*   🚀 **Transpiles to Idiomatic Go**: Generate clean, efficient, and readable Go code that you can compile and run anywhere. No black boxes.
*   🔍 **Interactive Physics Preview**: Launch a dynamic, physics-based visualization of your graph. Drag nodes, see connections spring into place, and get a feel for your program's structure.
*   💾 **Git-Friendly Format**: Graph files (`.ax`) are stored as human-readable JSON, making version control, diffing, and collaboration a breeze.
*   📦 **Multiple File Formats**: Choose the best format for your needs:
    *   `.ax`: Human-readable JSON for editing and version control.
    *   `.axb`: Raw binary for fast loading.
    *   `.axd`: Commented YAML for debugging.
    *   `.axc`: Compressed binary for distribution.
*   🧩 **Extensible by Design**: Easily define custom nodes that map to your own Go functions or external libraries.

---

## 🛠️ The Axon Workflow

Axon is built around a simple, powerful command-line interface.

### 1. Design: Create a `.ax` file

Define your logic in a simple JSON format. This example creates two constants, adds them, and prints the result.

**`add.ax`**
```json
{
  "id": "basic-addition-v2",
  "name": "Add Numbers with Execution Prototyping**: Build and test complex logic in minutes.
-   **Visual Learners**: Understand your application's architecture at a glance.
-   **Teaching Go**: Introduce core programming concepts without the steep learning curve of textual code.
-   **Complex Systems**: Design intricate data pipelines and concurrent workflows intuitively.
-   **Team Collaboration**: Share and version-control your logic in a clean, human-readable format.

## 🚀 Features

-   **Visual by Design**: Construct programs by connecting nodes. What you see is what you get.
-   **Transpiles to Idiomatic Go**: Generates clean, human-readable, and performant Go code. No black boxes.
-   **Interactive Physics Preview**: Launch a dynamic, interactive visualization of your graph with the `preview` command.
-   **Git-Friendly Format**: Uses a human-readable `.ax` (JSON) format, making version control and collaboration seamless.
-   **Full Go Integration**: Use any function from the Go standard library or your own custom packages.
-   **Multiple File Formats**: Choose between human-readable JSON (`.ax`), commented YAML (`.axd`), high-performance binary (`.axb`), or compressed archives (`.axc`).
-   **Powerful CLI**: A comprehensive command-line interface to build, convert, pack, and preview your graphs.

---

## 🔧 How It Works

Design your logic visually in an `.ax` file, which represents nodes and their connections (edges). Axon's transpiler then walks this graph to generate a `main.go` file.

#### 1. The Visual Graph (`add.ax`)

```json
{
  "id": "basic-addition-v2",
  "name": "Add Numbers with Execution Flow",
  "imports": ["fmt"],
  "nodes": [
    { "id": "start", "type": "START" },
    { "id": "const1", "type": "CONSTANT", "label": "x", "config": { "value": "5" }, ... },
    { "id": "const2", "type": "CONSTANT", "label": "y", "config": { "value": "3" }, ... },
    { "id": "sum", "type": "OPERATOR", "label": "z", "config": { "op": "+" }, ... },
    { "id": "printer", "type": "FUNCTION", "impl_reference": "fmt.Println", ... },
    { "id": "end", "type": "END" }
  ],
  "data_edges": [
    { "from_node_id": "const1", "to_node_id": "sum", ... },
    { "from_node_id": "const2", "to_node_id": "sum", ... },
    { "from_node_id": "sum", "to_node_id": "printer", ... }
  ],
  "exec_edges": [
    { "from_node_id": "start", "to_node_id": "sum" },
    { "from_node_id": "sum", "to_node_id": "printer" },
    { "from_node_id": "printer", "to_node_id": "end" }
 Flow",
  "imports": ["fmt"],
  "nodes": [
    { "id": "start", "type": "START" },
    { "id": "end", "type": "END" },
    { "id": "const1", "type": "CONSTANT", "label": "x", "outputs": [{ "type_name": "int" }], "config": { "value": "5" }},
    { "id": "const2", "type": "CONSTANT", "label": "y", "outputs": [{ "type_name": "int" }], "config": { "value": "3" }},
    { "id": "sum", "type": "OPERATOR", "label": "z", "inputs": [{ "name": "a", "type_name": "int" }, { "name": "b", "type_name": "int" }], "outputs": [{ "type_name": "int" }], "config": { "op": "+" }},
    { "id": "printer", "type": "FUNCTION", "label": "PrintResult", "impl_reference": "fmt.Println", "inputs": [{ "name": "a", "type_name": "int" }] }
  ],
  "data_edges": [
    { "from_node_id": "const1", "from_port": "out", "to_node_id": "sum", "to_port": "a" },
    { "from_node_id": "const2", "from_port": "out", "to_node_id": "sum", "to_port": "b" },
    { "from_node_id": "sum", "from_port": "out", "to_node_id": "printer", "to_port": "a" }
  ],
  "exec_edges": [
    { "from_node_id": "start", "to_node_id": "sum" },
    { "from_node_id": "sum", "to_node_id": "printer" },
    { "from_node_id": "printer", "to_node_id": "end" }
  ]
}
```

### 2. Visualize: Preview your graph

Before compiling, see your graph come to life!

```bash
axon preview examples/add.ax
```

This launches the interactive, physics-based previewer where you can pan, zoom, and drag nodes.

### 3. Build: Transpile to Go

When you're ready, transpile your visual graph into a real Go program.

```bash
axon build examples/add.ax
```

This generates the following clean Go code in `out/main.go`:

```go
package main

import (
	"fmt"
)

func main() {
	z := 5 + 3
	fmt.Println(z)
}
```

### 4. Run: Execute your code

Run your new Go program just like any other.

```bash
go run out/main.go
# Output: 8
```

---

## 🚀 Getting Started

### Prerequisites

-   [Go](https://golang.org/doc/install) (version 1.24 or later)
-   [Protocol Buffers Compiler](https://grpc.io/docs/protoc-installation/) (`protoc`)

### Installation

1.  **Clone the repository:**
    ```bash  ]
}
```

#### 2. The Transpiler (`axon build`)

The `build` command processes the graph, follows the execution and data edges, and generates Go code.

#### 3. The Go Output (`out/main.go`)

```go
package main

import (
	"fmt"
)

func main() {
	x := 5
	y := 3
	z := x + y
	fmt.Println(z)
}
```

---

## 🛠️ Getting Started

### 1. Prerequisites

You need `Go` and `protoc` installed on your system.

### 2. Clone and Build

```bash
# Clone the repository
git clone https://github.com/Advik-B/Axon.git
cd Axon

# Install dependencies and tools
make install-deps

# Build the Axon CLI
make build
```

### 3. Build Your First Graph

Transpile one of the example graphs into Go code.

```bash
./axon build examples/add.ax
```

This will generate `out/main.go`.

### 4. Run the Output

```bash
go run out/main.go
# Output: 8
```

### 5. Preview it!

Launch the interactive physics-based preview.

```bash
./axon preview examples/add.ax
```

---

## ✨ The Axon CLI

Axon comes with a powerful set of tools to manage your projects.

| Command                                   | Description                                                                                               |
| ----------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `axon build [file]`                       | **Transpiles** any Axon graph (`.ax`, `.axb`, `.axd`, `.axc`) into a runnable `out/main.go` file.            |
| `axon preview [file]`                     | **Launches** a beautiful, interactive, physics-based visualization of your graph.                           |
| `axon pack [file]`                        | **Compresses** any graph format into a highly efficient `.axc` binary archive using XZ compression.       |
| `axon unpack [file.axc]`                  | **Decompresses** an `.axc` archive back into the standard `.axb` binary format.                             |
| `axon convert <in-file> <out-file>`       | **Converts** between all Axon formats (`.ax`, `.axd`, `.axb`, `.axc`).                                      |

---

## 🗺️ Roadmap

-   [ ] **Visual Editor**: A full-fledged GUI for creating and editing `.ax` graphs.
-   [ ] **Live Reload**: Automatically update the preview and transpiled code on file changes.
-   [ ] **Expanded Standard Library**: More built-in nodes for common operations.
-   [ ] **Community Node Repository**: A place to share and discover custom nodes.
-   [ ] **Plugin Architecture**: Extend Axon's core functionality with custom plugins.

---

## 📄 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

## 🙏 Credits

Axon is inspired by the elegance of [Unreal Engine's Blueprints](https://docs.unrealengine.com/) and the ambitious vision of [Eve](http://witheve.com/), but is designed from the ground up to be a practical and powerful tool for the Go ecosystem.
    git clone https://github.com/Advik-B/Axon.git
    cd Axon
    ```

2.  **Install dependencies:**
    The `Makefile` helps you install the necessary Go tools.
    ```bash
    make install-deps
    ```

3.  **Build the `axon` binary:**
    ```bash
    make build
    ```

4.  **Verify the build:**
    You should now have an `axon` (or `axon.exe`) executable in the root directory.
    ```bash
    ./axon --help
    ```

---

## 🕹️ CLI Commands

Axon comes with a suite of powerful commands:

| Command                               | Description                                                                                               |
| ------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `axon build [file]`                   | Transpiles an Axon graph (`.ax`, `.axb`, `.axd`, `.axc`) into a runnable `out/main.go` file.                  |
| `axon preview [file]`                 | Launches the interactive, physics-based graph visualizer. **Try this first!**                             |
| `axon pack [file]`                    | Compresses any Axon graph format into the efficient `.axc` binary format for distribution.                |
| `axon unpack [file.axc]`              | Decompresses a `.axc` file into the standard `.axb` binary format.                                        |
| `axon convert [in-file] [out-file]`   | Converts between any Axon format (`.ax`, `.axb`, `.axd`, `.axc`).                                           |

## 📄 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

---

## Credits & Inspiration

Axon stands on the shoulders of giants. It is heavily inspired by the visual scripting systems in [Unreal Engine](https://docs.unrealengine.com/en-US/Engine/Blueprints/index.html) and the ambitious ideas of [Eve](http://witheve.com/), with a pragmatic focus on generating production-quality Go code.
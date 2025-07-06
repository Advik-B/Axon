# ✨ Axon: Visual Programming, For Real.

<p align="center">
  <img src="https://github.com/user-attachments/assets/b8a5a4c1-d011-4e03-9a37-ea273e48b9c2" alt="Axon Banner" style="width: 40%; height: auto;" />
	<br/>
  <sub>Yes, this logo is AI Generated, boo me, Im too poor to hire an artist</sub>
</p>

<p align="center">
  <strong>Go from a visual idea to idiomatic, runnable Go code.</strong>
  <br />
  Axon is a visual, node-based programming language that bridges the gap between intuitive, flow-based logic and high-performance, real-world applications.
  <br />
  <br />
  <a href="#-core-features"><strong>Features</strong></a> ·
  <a href="#-interactive-preview"><strong>Live Preview</strong></a> ·
  <a href="#%EF%B8%8F-the-axon-workflow"><strong>Workflow</strong></a> ·
  <a href="#-getting-started"><strong>Getting Started</strong></a> ·
  <a href="#%EF%B8%8F-cli-commands"><strong>CLI Commands</strong></a>
</p>

---

## 🚀 What is Axon?

Tired of boilerplate? Teaching a new programmer? Want to *see* your logic?

**Axon** is your answer. Unlike toy block-based languages, Axon is designed for serious development. It provides a structured, visual environment to design complex systems, which it then transpiles directly into clean, human-readable, and performant Go code.

It’s the clarity of flowcharts with the power of a compiled language.

## 📺 Interactive Preview

The star of the show. Axon's previewer isn't just a static diagram. It's a living, breathing visualization of your graph's structure, complete with a spring-mass physics simulation and a live-transpiled code view.


<p align="center">
  <img src="https://github.com/user-attachments/assets/f451e867-3cd4-490b-97b0-f61e2f8ea852" alt="Axon Physics-Based Previewer" width="1000"/>
  <br/>
  <em>Pan, zoom, and drag nodes. The graph fluidly rearranges itself while showing the generated Go code in a side panel.</em>
  <br/>
  <sub>Yes the code for this is something I am not proud of. Yes. this is super jank. Yes it looks ugly</sub>
</p>

## ✨ Core Features

*   🧠 **Intuitive Node-Based Logic**: Build programs by connecting nodes. Control execution flow (`ExecEdge`) and data flow (`DataEdge`) explicitly and visually.
*   🚀 **Transpiles to Idiomatic Go**: Generate clean, efficient, and readable Go code that you can compile and run anywhere. No black boxes.
*   🔍 **Interactive Physics Preview**: Launch a dynamic, physics-based visualization of your graph. Drag nodes, see connections spring into place, and get a feel for your program's structure.
*   💾 **Git-Friendly Format**: Graph files (`.ax`) are stored as human-readable JSON, making version control, diffing, and collaboration a breeze.
*   📦 **Multiple File Formats**: Choose the best format for your needs:
    *   `.ax`: Human-readable JSON for editing and version control.
    *   `.axb`: Raw Protobuf binary for fast loading.
    *   `.axd`: Commented YAML for debugging.
    *   `.axc`: Compressed binary for distribution.
*   🧩 **Extensible by Design**: Easily define custom nodes that map to your own Go functions or any function from the Go standard library.

---

## 🛠️ The Axon Workflow

Axon is built around a simple, powerful command-line interface.

### 1. Design: Create a `.ax` file

*⚠️ Small note, in the (not too distant) future, I will create a proper node graph editor using `Tauri`*

*For now, I am sticking to writing the json files manually* 

*The node graph preview is just some jank (spagetti) code I threw together*

Define your logic in a simple JSON format. This example creates two constants, adds them, and prints the result.

**`add.ax`**
```json
{
  "id": "basic-addition-v2",
  "name": "Add Numbers with Execution Flow",
  "imports": ["fmt"],
  "nodes": [
    { "id": "start", "type": "START" },
    { "id": "end", "type": "END" },
    { "id": "const1", "type": "CONSTANT", "label": "x", "outputs": [{ "type_name": "int" }], "config": { "value": "5" }},
    { "id": "const2", "type": "CONSTANT", "label": "y", "outputs": [{ "type_name": "int" }], "config": { "value": "3" }},
    { "id": "sum", "type": "OPERATOR", "label": "z", "inputs": [{ "name": "a", "type_name": "int" }, { "name": "b", "type_name": "int" }], "outputs": [{ "name": "out", "type_name": "int" }], "config": { "op": "+" }},
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
	x := 5
	y := 3
	z := x + y
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
-   Make (for windows, refer to the [chocolatey package](https://community.chocolatey.org/packages/make))

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/Advik-B/Axon.git
    cd Axon
    ```

2.  **Install dependencies and tools:**
    The `Makefile` simplifies installing the Go tools needed for development.
    ```bash
    make install-deps
    ``` 

3.  **Build the `axon` binary:**
    ```bash
    make
    ```

4.  **Verify the installation:**
    You should now have an `axon` (or `axon.exe` on Windows) executable.
    ```bash
    ./axon --help
    ```

## 🕹️ CLI Commands

Axon comes with a powerful and flexible set of tools:

| Command                               | Description                                                                                               |
| ------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `axon build [file]`                   | **Transpiles** any Axon graph (`.ax`, `.axb`, `.axd`, `.axc`) into a runnable `out/main.go` file.            |
| `axon preview [file]`                 | **Launches** a beautiful, interactive, physics-based visualization of your graph.                           |
| `axon pack [file]`                    | **Compresses** any graph format into a highly efficient `.axc` binary archive using XZ compression.       |
| `axon unpack [file.axc]`              | **Decompresses** an `.axc` archive back into the standard `.axb` binary format.                             |
| `axon convert <in-file> <out-file>`   | **Converts** between all Axon formats (`.ax`, `.axd`, `.axb`, `.axc`).                                      |

---

## 🗺️ Roadmap

-   [x] **Core Transpiler**: A fully working transpiler.
-   [x] **Preview**: A readonly visualiser for the graph 
-   [ ] **Visual Editor**: A full-fledged GUI for creating and editing `.ax` graphs from scratch.
-   [ ] **Live Reload**: Automatically update the previewer when graph files change.
-   [ ] **Plugin Architecture**: A formal way to extend Axon's core transpiler and previewer functionality.

---

## 📄 License

This project is licensed under the **MIT License**.

## 🫂Credits & Inspiration

Axon is heavily inspired by the elegance and power of [Unreal Engine's Blueprints](https://docs.unrealengine.com/).

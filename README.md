# Axon — Visual Node-Based Programming Language

**Axon** is a visual, node-based programming language designed to make programming intuitive, accessible, and powerful. Unlike traditional block-based languages like Scratch, Axon supports complex, scalable systems and **transpiles directly to idiomatic Go code**.

Axon is ideal for:

- 🧠 Teaching computational thinking
- 🔧 Building real-world applications without writing text code
- 🧱 Creating modular logic through flow-based programming

---

## 🚀 Features

- ⚙️ **Node-based architecture**: Build logic by connecting visual nodes.
- 🔄 **Data and Execution Flow**: Explicit control and data wiring via `DataEdge` and `ExecEdge`.
- 📦 **Transpiles to Go**: Fully functional Go code output — compile and run with `go run`.
- 📝 **Git-friendly `.ax` format**: Human-readable format for versioning, collaboration, and CI.
- 🧩 **Custom Nodes**: Extend Axon with your own Go functions or external libraries.

---

## 📁 File Format: `.ax`

Axon uses a JSON format that mirrors the `axon.proto` schema (defined using Protocol Buffers). Here's an example node graph:

```json
{
  "id": "basic-addition",
  "name": "Add Numbers",
  "nodes": [
    {
      "id": "const1",
      "type": "CONSTANT",
      "label": "x",
      "outputs": [
        { "name": "out", "type": "INTEGER", "is_optional": false }
      ],
      "config": { "value": "5" }
    },
    {
      "id": "const2",
      "type": "CONSTANT",
      "label": "y",
      "outputs": [
        { "name": "out", "type": "INTEGER", "is_optional": false }
      ],
      "config": { "value": "3" }
    },
    {
      "id": "sum",
      "type": "OPERATOR",
      "label": "z",
      "inputs": [
        { "name": "a", "type": "INTEGER", "is_optional": false },
        { "name": "b", "type": "INTEGER", "is_optional": false }
      ],
      "outputs": [
        { "name": "out", "type": "INTEGER", "is_optional": false }
      ],
      "config": { "op": "+" }
    }
  ],
  "data_edges": [
    { "from_node_id": "const1", "from_port": "out", "to_node_id": "sum", "to_port": "a" },
    { "from_node_id": "const2", "from_port": "out", "to_node_id": "sum", "to_port": "b" }
  ],
  "exec_edges": []
}
```

---

## 🛠️ Getting Started

### 1. Clone and Build

```bash
git clone https://github.com/Advik-B/Axion.git
cd axon-lang
go build -o axon ./cmd/axon
```

### 2. Transpile a Graph

```bash
./axon build examples/add.ax
```

Outputs:

```go
package main

func main() {
    x := 5
    y := 3
    z := x + y
}
```

### 3. Run the Output

```bash
go run out/main.go
```




<!-- ## 🤝 Git Best Practices

- ✅ Always commit `.axs` (text)
- ❌ Avoid committing `.pb` or `.go` outputs
- Add `.gitignore`:

```gitignore
*.pb
*.pb.go
out/
*.exe
``` -->

---

## ✨ Roadmap



---

## 📄 License

MIT License — free to use, extend, and remix.

---

## Credits

[Unreal Blueprints](https://docs.unrealengine.com/), and [Eve](http://witheve.com/), but designed to be practical, expressive, and truly powerful.

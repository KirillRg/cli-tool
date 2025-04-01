# CLI Tool for Insomnia Collections → k6 AST

🚀 A command-line tool that parses **Insomnia (v5)** collections in YAML format and transforms them into Go-based **Abstract Syntax Tree (AST)** structures, preparing them for future code generation (e.g., k6 load testing scripts).

---

## 💡 Purpose

This CLI is part of a Master's thesis project (*HSE, 2025*) aimed at automating the transformation of API collections into load testing scripts. It addresses limitations of existing tools:

- ❌ `postman-to-k6` is incompatible with Insomnia and not accessible due to licensing restrictions in Russia.
- ❌ `openapi-to-k6` lacks logic handling, only defines requests.
- ✅ Our tool is **vendor-neutral**, extensible, and CLI-friendly.
- ✅ It uses its own **platform-independent AST**

---
## 📁 Project Structure
```graphql
cli-tool/
├── cmd/
│   ├── root.go           # CLI root using cobra
│   └── generate.go       # generate command logic
├── internal/
│   └── ast/
│       ├── ast.go        # AST structures
│       └── generator.go  # AST generation methods
│   └── parser/
│       ├── insomnia.go   # Data models for Insomnia collection
│       └── parser.go     # YAML parsing logic
├── Insomnia_Test_Collection.yaml  # Example input collection
├── go.mod / go.sum       # Go dependencies
├── LICENSE
├── main.go               # CLI entrypoint
└── README.md             # Project info
```
---
## ⚙️ Installation & Usage
### Prerequisites
- Go 1.21+
- Cobra CLI installed:
  ```bash
  go install github.com/spf13/cobra-cli@latest
  ```
### Run the tool
  ```bash
  go run main.go generate --input Insomnia_Test_Collection.yaml
  ```
### Example Output
  ```css
 Generating scripts from file: Insomnia_Test_Collection.yaml
 Parsed collection: &{Type:..., Name:..., Collection:[...]}
  ```
---
## 👤 Author
🎓 Rogoza Kirill Andreevich
Master's Program: System and Software Engineering (SPI)
HSE University, 2025
GitHub: @KirillRg


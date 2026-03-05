# CLI Tool for Insomnia Collections → k6 AST

🚀 A command-line tool that parses **Insomnia (v5)** collections in YAML format and transforms them into Go-based **Abstract Syntax Tree (AST)** structures, preparing them for future code creation (e.g., k6 load testing scripts).

---

## 💡 Purpose

This CLI is part of a Master's thesis project (*HSE, 2025*) aimed at automating the transformation of API collections into load testing scripts. It addresses limitations of existing tools:

- ❌ `postman-to-k6` is incompatible with Insomnia and not accessible due to licensing restrictions in Russia.
- ❌ `openapi-to-k6` lacks logic handling, only defines requests.
- ✅ Our tool is **vendor-neutral**, extensible, and CLI-friendly.
- ✅ It uses its own **platform-independent ESTREE standarted AST**

---
## 📁 Project Structure
```graphql
cli-tool/
├── cmd/
│   ├── root.go           # CLI root using cobra
│   └── build.go       # command logic
├── internal/
│   └── translator/
│       ├── translator.go   # k6 JS translatyion from AST
│   └── ast/
│       ├── ast.go        # AST structures
│       └── generator.go  # AST generation methods
│   └── parser/
│       ├── insomnia.go   # Data models for Insomnia collection
│       └── parser.go     # YAML parsing logic
├── test/
│   ├── becnchmark_test.go   # Benchmarks for E2E AST generation
│   └── result.txt     # Test results
├── Insomnia_Test_Collection_With_Environment.yaml  # Example input collection with Insomnia Environment variables
├── Insomnia_Test_Collection.yaml  # Basic example input collection
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
  go run main.go translate --input Insomnia_Test_Collection.yaml
  ```
### Example Output
  ```css
  //TO DO

  ```
---
## 🔬 Benchmarking

This CLI was evaluated using Go microbenchmarks to ensure it meets non-functional performance requirements under both average and extreme load conditions.

### Benchmarked Scenarios

| Scenario                      | Description                                         |
|------------------------------|-----------------------------------------------------|
| `Benchmark_10RequestASTGeneration` | Simulates average Insomnia collection (~10 requests) |
| `Benchmark_Constant_10000Requests` | Simulates stress-load with 10,000 sequential requests  |

### How to Run

Benchmarks are located in the [`test/`](./test) directory. To run them and save results:

```bash
go test ./test -bench=Benchmark_ -benchtime=50x -benchmem > test/result.txt
```
---
## 👤 Author
🎓 Rogoza Kirill Andreevich
Master's Program: System and Software Engineering (SPI)
HSE University, 2025
GitHub: @KirillRg


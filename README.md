# CLI Tool for Insomnia Collections → k6 AST

A CLI tool that translates **Insomnia (v5)** collections (YAML) into ready-to-run **k6** load testing scripts. It uses **Abstract Syntax Tree (AST)** as an intermediate representation and also provides an **[Insomnia plugin](https://github.com/KirillRg/insomnia-plugin-k6-translator)** as a UI wrapper over the CLI core.

---

## Purpose

This CLI is part of a Master's thesis project (*HSE, 2025 - 2026*) aimed at automating the transformation of API collections into load testing scripts. The key idea is to use **ESTree-compliant AST** as an intermediate representation — this separates parsing from code generation and makes the transformation pipeline more maintainable and verifiable.

The main goal is to provide a new integration bridge between **Insomnia** (as a practical API client) and **k6** (as a modern load testing tool)

---
## Project Structure
```graphql
cli-tool/
├── cmd/
│   ├── root.go        # CLI root using cobra
│   └── build.go       # command/flags logic
├── internal/
│   └── ast/
│       ├── ast.go        # AST structures with ESTREE links
│       ├── generator.go  # AST generation methods for collection
│       ├── json_export.go  # AST to JSON conversion
│       └── options.go  # AST generation methods for profile

│   └── parser/
│       ├── insomnia.go   # Data models for Insomnia collection
│       └── parser.go     # YAML parsing logic
│   └── profile/
│       ├── builder.go   # k6 load profile construction
│       └── model.go     # load profile structure definitions
│   └── translator/
│       ├── node_translators.go   # set of translation rules for nodes
│       └── translator.go   # k6 JS translation from AST
├── test/
│   ├── becnchmark_test.go   # Benchmarks for E2E k6 translation with different stages
│   └── func_result.txt     # Functional test results
│   └── nonfunc_result.txt     # Becnchmark test results
├── Insomnia_Test_Collection_With_Environment.yaml  # Example input collection with Insomnia Environment variables
├── Insomnia_Test_Collection.yaml  # Basic example input collection
├── go.mod / go.sum       # Go dependencies
├── LICENSE
├── main.go               # CLI entrypoint
└── README.md             # Project info
```
---
## Installation & Usage
### Prerequisites
- Go 1.21+
- Cobra CLI installed:
  ```bash
  go install github.com/spf13/cobra-cli@latest
  ```
### Run the tool with output inside project
  ```bash
  go run main.go translate --input Insomnia_Test_Collection.yaml --vus 100 --duration 100s
  ```
### Example Output
  ```css
  Translated k6 scrypt: //output_file_path
  ```
---
## Benchmarking

This CLI was evaluated using Go microbenchmarks to ensure it meets non-functional performance requirements under both average and extreme load conditions.

### Benchmarked Scenarios

| Scenario                      | Description                                         |
|------------------------------|-----------------------------------------------------|
| `Benchmark_10RequestASTGeneration` | Simulates stress-load for scenario 'parse -> AST' on 10 requests collection |
| `Benchmark_Constant_10000Requests` | Simulates stress-load for scenario 'parse -> AST' on 10,000 requests collection |
| `Benchmark_10RequestFullPipeline` | Simulates stress-load for scenario 'parse -> AST -> script' on 10 requests collection  |
| `Benchmark_Constant_10000RequestsFullPipeline` | Simulates stress-load for scenario 'parse -> AST -> script' on 10,000 requests collection |
| `Benchmark_10RequestFullPipelineWithWrite` | Simulates stress-load for scenario 'parse -> AST -> script -> file' on 10 requests collection |
| `Benchmark_Constant_10000RequestsFullPipelineWithWrite` | Simulates stress-load for scenario 'parse -> AST -> script -> file' on 10,000 requests collection |

### How to Run

Benchmarks are located in the [`test/`](./test) directory. To run them and save results:

```bash
go test ./test -bench=Benchmark_ -benchtime=50x -benchmem > test/nonfunc_result.txt
```
---
## Functional Testing

The project also includes a set of functional automated checks covering:

- unit tests for `profile`, `parser`, `ast`, and `translator`
- validation tests for CLI input flags
- one E2E test that verifies the generated k6 script is accepted by `k6 inspect` *k6 download required

### How to Run

To execute all functional checks and save the results into a text file:

```bash
go test ./test/autotests/... -v > test/func_result.txt
```
---
## 👤 Author
🎓 Rogoza Kirill Andreevich
Master's Program: System and Software Engineering (SPI)
HSE University, 2025
GitHub: @KirillRg


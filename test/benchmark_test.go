package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/KirillRg/cli-tool/internal/ast"
	"github.com/KirillRg/cli-tool/internal/parser"
	"github.com/KirillRg/cli-tool/internal/profile"
	"github.com/KirillRg/cli-tool/internal/translator"
)

func Benchmark_10RequestASTGeneration(b *testing.B) {
	yamlData := generateFakeCollectionYAML(10)
	tmpFile := writeTempFile(b, yamlData)

	collection, err := parser.ParseInsomniaCollection(tmpFile)
	if err != nil {
		b.Fatal(err)
	}
	fmt.Printf("Loaded 10-req collection: %d requests\n", len(collection.Collection))

	lp, err := buildBenchmarkProfile(tmpFile)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection, err := parser.ParseInsomniaCollection(tmpFile)
		if err != nil {
			b.Fatal(err)
		}
		_ = ast.GenerateAST(collection, lp)
	}
}

func Benchmark_Constant_10000Requests(b *testing.B) {
	yamlData := generateFakeCollectionYAML(10000)
	tmpFile := writeTempFile(b, yamlData)

	collection, err := parser.ParseInsomniaCollection(tmpFile)
	if err != nil {
		b.Fatal(err)
	}
	fmt.Printf("Loaded 10K-req collection: %d requests\n", len(collection.Collection))

	lp, err := buildBenchmarkProfile(tmpFile)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection, err := parser.ParseInsomniaCollection(tmpFile)
		if err != nil {
			b.Fatal(err)
		}
		_ = ast.GenerateAST(collection, lp)
	}
}

func Benchmark_10RequestFullPipeline(b *testing.B) {
	yamlData := generateFakeCollectionYAML(10)
	tmpFile := writeTempFile(b, yamlData)

	collection, err := parser.ParseInsomniaCollection(tmpFile)
	if err != nil {
		b.Fatal(err)
	}
	fmt.Printf("Loaded 10-req collection for full pipeline: %d requests\n", len(collection.Collection))

	lp, err := buildBenchmarkProfile(tmpFile)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection, err := parser.ParseInsomniaCollection(tmpFile)
		if err != nil {
			b.Fatal(err)
		}

		tree := ast.GenerateAST(collection, lp)

		_, err = translator.TranslateProgram(tree)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Constant_10000RequestsFullPipeline(b *testing.B) {
	yamlData := generateFakeCollectionYAML(10000)
	tmpFile := writeTempFile(b, yamlData)

	collection, err := parser.ParseInsomniaCollection(tmpFile)
	if err != nil {
		b.Fatal(err)
	}
	fmt.Printf("Loaded 10K-req collection for full pipeline: %d requests\n", len(collection.Collection))

	lp, err := buildBenchmarkProfile(tmpFile)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection, err := parser.ParseInsomniaCollection(tmpFile)
		if err != nil {
			b.Fatal(err)
		}

		tree := ast.GenerateAST(collection, lp)

		_, err = translator.TranslateProgram(tree)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_10RequestFullPipelineWithWrite(b *testing.B) {
	yamlData := generateFakeCollectionYAML(10)
	tmpFile := writeTempFile(b, yamlData)
	outFile := "bench_output_10.js"

	collection, err := parser.ParseInsomniaCollection(tmpFile)
	if err != nil {
		b.Fatal(err)
	}
	fmt.Printf("Loaded 10-req collection for full pipeline + write: %d requests\n", len(collection.Collection))

	lp, err := buildBenchmarkProfile(tmpFile)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection, err := parser.ParseInsomniaCollection(tmpFile)
		if err != nil {
			b.Fatal(err)
		}

		tree := ast.GenerateAST(collection, lp)

		code, err := translator.TranslateProgram(tree)
		if err != nil {
			b.Fatal(err)
		}

		err = os.WriteFile(outFile, []byte(code), 0644)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Constant_10000RequestsFullPipelineWithWrite(b *testing.B) {
	yamlData := generateFakeCollectionYAML(10000)
	tmpFile := writeTempFile(b, yamlData)
	outFile := "bench_output_10000.js"

	collection, err := parser.ParseInsomniaCollection(tmpFile)
	if err != nil {
		b.Fatal(err)
	}
	fmt.Printf("Loaded 10K-req collection for full pipeline + write: %d requests\n", len(collection.Collection))

	lp, err := buildBenchmarkProfile(tmpFile)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection, err := parser.ParseInsomniaCollection(tmpFile)
		if err != nil {
			b.Fatal(err)
		}

		tree := ast.GenerateAST(collection, lp)

		code, err := translator.TranslateProgram(tree)
		if err != nil {
			b.Fatal(err)
		}

		err = os.WriteFile(outFile, []byte(code), 0644)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func buildBenchmarkProfile(inputPath string) (*profile.LoadProfile, error) {
	return profile.Build(profile.BuildInput{
		InputPath:     inputPath,
		OutputPath:    "test",
		VUs:           1,
		Duration:      "10s",
		Iterations:    0,
		StagesRaw:     nil,
		ThresholdsRaw: nil,
	})
}

// Генерация фейковой коллекции с N запросами
func generateFakeCollectionYAML(count int) string {
	var b strings.Builder
	b.WriteString("type: collection.insomnia.rest/5.0\nname: bench\ncollection:\n")

	for i := 0; i < count; i++ {
		b.WriteString(fmt.Sprintf(`  - url: http://fake.local/request%d
    name: req_%d
    method: GET
    parameters:
      - name: p%d
        value: v%d
        disabled: false
    headers:
      - name: H%d
        value: val%d
        disabled: false
    settings:
      renderRequestBody: true
      encodeUrl: true
`, i, i, i, i, i, i))
	}

	b.WriteString("environments:\n  data: {}\n")
	return b.String()
}

// Запись YAML во временный файл
func writeTempFile(tb testing.TB, content string) string {
	f, err := os.CreateTemp("", "bench-*.yaml")
	if err != nil {
		tb.Fatal(err)
	}
	_, err = f.WriteString(content)
	if err != nil {
		tb.Fatal(err)
	}
	f.Close()
	return f.Name()
}

package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/KirillRg/cli-tool/internal/ast"
	"github.com/KirillRg/cli-tool/internal/parser"
)

// Бенчмарк: сколько 10-запросных коллекций обрабатывается за секунду
func Benchmark_10RequestASTGeneration(b *testing.B) {
	yamlData := generateFakeCollectionYAML(10)
	tmpFile := writeTempFile(b, yamlData)

	collection, err := parser.ParseInsomniaCollection(tmpFile)
	if err != nil {
		b.Fatal(err)
	}
	fmt.Printf("Loaded 10-req collection: %d requests\n", len(collection.Collection))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection, err := parser.ParseInsomniaCollection(tmpFile)
		if err != nil {
			b.Fatal(err)
		}
		_ = ast.GenerateAST(collection)
	}
}

// Бенчмарк: постоянная нагрузка на 10 000 запросов
func Benchmark_Constant_10000Requests(b *testing.B) {
	yamlData := generateFakeCollectionYAML(10000)
	tmpFile := writeTempFile(b, yamlData)

	collection, err := parser.ParseInsomniaCollection(tmpFile)
	if err != nil {
		b.Fatal(err)
	}
	fmt.Printf("Loaded 10K-req collection: %d requests\n", len(collection.Collection))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection, err := parser.ParseInsomniaCollection(tmpFile)
		if err != nil {
			b.Fatal(err)
		}
		_ = ast.GenerateAST(collection)
	}
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

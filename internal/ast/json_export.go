package ast

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func WriteProgramJSON(program Program, outputPath string) error {
	jsonBytes, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal AST to JSON: %w", err)
	}

	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(outputPath, jsonBytes, 0o644); err != nil {
		return fmt.Errorf("failed to write AST JSON file: %w", err)
	}

	return nil
}

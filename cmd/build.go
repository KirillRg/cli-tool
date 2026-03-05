package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KirillRg/cli-tool/internal/ast"
	"github.com/KirillRg/cli-tool/internal/parser"
	"github.com/KirillRg/cli-tool/internal/translator"
	"github.com/spf13/cobra"
)

var inputFilePath string

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translates load scripts from Insomnia collections",
	Long: `Translates k6 load scripts from provided Insomnia collection file.
			Currently supports basic JSON and YAML Insomnia collections.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Translating scripts from file:", inputFilePath)

		// 1) Parse collection file
		collection, err := parser.ParseInsomniaCollection(inputFilePath)
		if err != nil {
			fmt.Println("Failed to parse file:", err)
			return
		}
		fmt.Printf("\n___ ___ ___\n")
		fmt.Printf("Parsed collection: %+v\n", collection)

		// 2) Build ESTree AST
		tree := ast.GenerateAST(collection)
		fmt.Printf("\n___ ___ ___\n")
		fmt.Printf("Created AST:\n%+v\n", tree)

		// 3) Translate JS from AST
		script, err := translator.TranslateProgram(tree)
		if err != nil {
			fmt.Println("Failed to translate JS:", err)
			return
		}

		// 4) Write output to a fixed file on every run
		outDir := "result"
		outFile := filepath.Join(outDir, "k6_script.js")
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			fmt.Println("Failed to create output directory:", err)
			return
		}
		if err := os.WriteFile(outFile, []byte(script), 0o644); err != nil {
			fmt.Println("Failed to write script:", err)
			return
		}

		fmt.Printf("\n___ ___ ___\n")
		fmt.Printf("Translated k6 script: %s\n", outFile)
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "Input Insomnia collection file (required)")
	translateCmd.MarkFlagRequired("input")
}

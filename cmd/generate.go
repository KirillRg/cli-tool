package cmd

import (
	"fmt"

	"github.com/KirillRg/cli-tool/internal/ast"
	"github.com/KirillRg/cli-tool/internal/parser"
	"github.com/spf13/cobra"
)

var inputFilePath string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates load scripts from Insomnia collections",
	Long: `Generates k6 load scripts from provided Insomnia collection file.
			Currently supports basic JSON and YAML Insomnia collections.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating scripts from file:", inputFilePath)

		// 1. Парсим файл коллекции
		collection, err := parser.ParseInsomniaCollection(inputFilePath)
		if err != nil {
			fmt.Println("Failed to parse file:", err)
			return
		}
		fmt.Printf("\n___ ___ ___\n")
		fmt.Printf("Parsed collection: %+v\n", collection)

		// 2. Генерируем AST
		tree := ast.GenerateAST(collection)
		fmt.Printf("\n___ ___ ___\n")
		fmt.Printf("Generated AST:\n%+v\n", tree)

		// 3. TODO: Здесь будет Writer, чтобы получить JS-код
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "Input Insomnia collection file (required)")
	generateCmd.MarkFlagRequired("input")
}

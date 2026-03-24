package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/KirillRg/cli-tool/internal/ast"
	"github.com/KirillRg/cli-tool/internal/parser"
	"github.com/KirillRg/cli-tool/internal/profile"
	"github.com/KirillRg/cli-tool/internal/translator"
	"github.com/spf13/cobra"
)

// Флаги
var (
	inputFilePath  string
	outputFilePath string

	vusFlag        int
	durationFlag   string
	iterationsFlag int

	stageFlags     []string
	thresholdFlags []string
)

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translates load scripts from Insomnia collections",
	Long: `Translates k6 load scripts from provided Insomnia collection file.
Currently supports basic JSON and YAML Insomnia collections.`,

	SilenceUsage: true,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateTranslateFlags()
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := profile.Build(profile.BuildInput{
			InputPath:     inputFilePath,
			OutputPath:    outputFilePath,
			VUs:           vusFlag,
			Duration:      durationFlag,
			Iterations:    iterationsFlag,
			StagesRaw:     stageFlags,
			ThresholdsRaw: thresholdFlags,
		})
		if err != nil {
			return fmt.Errorf("failed to build load profile: %w", err)
		}

		collection, err := parser.ParseInsomniaCollection(inputFilePath)
		if err != nil {
			return fmt.Errorf("failed to parse input file: %w", err)
		}

		tree := ast.GenerateAST(collection, profile)

		script, err := translator.TranslateProgram(tree)
		if err != nil {
			return fmt.Errorf("failed to translate JS: %w", err)
		}

		outDir := filepath.Dir(outputFilePath)
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		if err := os.WriteFile(outputFilePath, []byte(script), 0o644); err != nil {
			return fmt.Errorf("failed to write script: %w", err)
		}

		fmt.Printf("Translated k6 script: %s\n", outputFilePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)

	// Базовые input/output флаги
	translateCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "Input Insomnia collection file (required)")
	translateCmd.Flags().StringVarP(&outputFilePath, "output", "o", "result/k6_script.js", "Output k6 script path")

	// Флаги профиля нагрузки
	translateCmd.Flags().IntVar(&vusFlag, "vus", 0, "Number of virtual users")
	translateCmd.Flags().StringVar(&durationFlag, "duration", "", "Test duration, for example 30s or 1m")
	translateCmd.Flags().IntVar(&iterationsFlag, "iterations", 0, "Total number of iterations")
	translateCmd.Flags().StringSliceVar(&stageFlags, "stage", nil, "Ramp stage in format <duration>:<target>, repeatable")
	translateCmd.Flags().StringSliceVar(&thresholdFlags, "threshold", nil, "Threshold rule in format <metric>=<condition>, repeatable")

	if err := translateCmd.MarkFlagRequired("input"); err != nil {
		panic(err)
	}
}

// Валидация флагов
func validateTranslateFlags() error {

	if vusFlag < 0 {
		return fmt.Errorf("--vus cannot be negative")
	}

	if iterationsFlag < 0 {
		return fmt.Errorf("--iterations cannot be negative")
	}

	modeCount := 0

	if durationFlag != "" {
		modeCount++
	}
	if iterationsFlag > 0 {
		modeCount++
	}
	if len(stageFlags) > 0 {
		modeCount++
	}

	if modeCount == 0 {
		return fmt.Errorf("load profile is not specified: use either --vus with --duration, --vus with --iterations, or one or more --stage flags")
	}

	if modeCount > 1 {
		return fmt.Errorf("conflicting load profile flags: use only one mode from (--vus + --duration), (--vus + --iterations), or (--stage ...)")
	}

	if durationFlag != "" {
		if vusFlag <= 0 {
			return fmt.Errorf("--vus must be greater than 0 when --duration is used")
		}
	}

	if iterationsFlag > 0 {
		if vusFlag <= 0 {
			return fmt.Errorf("--vus must be greater than 0 when --iterations is used")
		}
	}

	if len(stageFlags) > 0 {

		if vusFlag > 0 {
			return fmt.Errorf("--vus cannot be used together with --stage in the current implementation")
		}

		if err := validateStageFlags(stageFlags); err != nil {
			return err
		}
	}

	if err := validateThresholdFlags(thresholdFlags); err != nil {
		return err
	}

	return nil
}

func validateStageFlags(stages []string) error {
	for _, rawStage := range stages {
		rawStage = strings.TrimSpace(rawStage)
		if rawStage == "" {
			return fmt.Errorf("empty --stage value is not allowed")
		}

		parts := strings.SplitN(rawStage, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --stage value %q: expected format <duration>:<target>", rawStage)
		}

		duration := strings.TrimSpace(parts[0])
		targetRaw := strings.TrimSpace(parts[1])

		if duration == "" {
			return fmt.Errorf("invalid --stage value %q: duration cannot be empty", rawStage)
		}

		target, err := strconv.Atoi(targetRaw)
		if err != nil {
			return fmt.Errorf("invalid --stage value %q: target must be an integer", rawStage)
		}

		if target < 0 {
			return fmt.Errorf("invalid --stage value %q: target cannot be negative", rawStage)
		}
	}

	return nil
}

func validateThresholdFlags(thresholds []string) error {
	for _, rawThreshold := range thresholds {
		rawThreshold = strings.TrimSpace(rawThreshold)
		if rawThreshold == "" {
			return fmt.Errorf("empty --threshold value is not allowed")
		}

		parts := strings.SplitN(rawThreshold, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --threshold value %q: expected format <metric>=<condition>", rawThreshold)
		}

		metric := strings.TrimSpace(parts[0])
		condition := strings.TrimSpace(parts[1])

		if metric == "" {
			return fmt.Errorf("invalid --threshold value %q: metric cannot be empty", rawThreshold)
		}

		if condition == "" {
			return fmt.Errorf("invalid --threshold value %q: condition cannot be empty", rawThreshold)
		}
	}

	return nil
}

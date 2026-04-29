package integrations

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	astpkg "github.com/KirillRg/cli-tool/internal/ast"
	"github.com/KirillRg/cli-tool/internal/parser"
	"github.com/KirillRg/cli-tool/internal/profile"
	"github.com/KirillRg/cli-tool/internal/translator"
)

func TestGeneratedScript_K6Inspect(t *testing.T) {
	if _, err := exec.LookPath("k6"); err != nil {
		t.Skip("k6 is not installed in PATH")
	}

	collectionPath := filepath.Join("..", "..", "..", "Insomnia_Test_Collection_With_Environment.yaml")

	collection, err := parser.ParseInsomniaCollection(collectionPath)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	loadProfile, err := profile.Build(profile.BuildInput{
		InputPath:  "input.yaml",
		OutputPath: "output.js",
		VUs:        10,
		Duration:   "30s",
	})
	if err != nil {
		t.Fatalf("unexpected profile error: %v", err)
	}

	generatedProgram := astpkg.GenerateAST(collection, loadProfile)

	generatedScript, err := translator.TranslateProgram(generatedProgram)
	if err != nil {
		t.Fatalf("unexpected translation error: %v", err)
	}

	scriptPath := filepath.Join(t.TempDir(), "generated-script.js")
	if err := os.WriteFile(scriptPath, []byte(generatedScript), 0o600); err != nil {
		t.Fatalf("failed to write generated script: %v", err)
	}

	inspectCommand := exec.Command("k6", "inspect", scriptPath)
	inspectOutput, err := inspectCommand.CombinedOutput()
	if err != nil {
		t.Fatalf("k6 inspect failed: %v\n%s", err, string(inspectOutput))
	}
}

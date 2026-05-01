package validation

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuildCommand_FlagValidation(t *testing.T) {
	projectRoot := filepath.Join("..", "..", "..")
	collectionPath := "Insomnia_Test_Collection_With_Environment.yaml"

	testCases := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name: "valid duration profile",
			args: []string{
				"run", ".",
				"translate",
				"--input", collectionPath,
				"--output", "result.js",
				"--vus", "10",
				"--duration", "30s",
			},
		},
		{
			name: "valid iterations profile",
			args: []string{
				"run", ".",
				"translate",
				"--input", collectionPath,
				"--output", "result.js",
				"--vus", "5",
				"--iterations", "100",
			},
		},
		{
			name: "valid stages profile",
			args: []string{
				"run", ".",
				"translate",
				"--input", collectionPath,
				"--output", "result.js",
				"--stage", "30s:5",
				"--stage", "1m:20",
			},
		},
		{
			name: "invalid mixed profile modes",
			args: []string{
				"run", ".",
				"translate",
				"--input", collectionPath,
				"--output", "result.js",
				"--vus", "10",
				"--duration", "30s",
				"--iterations", "100",
				"--stage", "30s:5",
			},
			expectError: true,
		},
		{
			name: "invalid threshold format",
			args: []string{
				"run", ".",
				"translate",
				"--input", collectionPath,
				"--output", "result.js",
				"--vus", "10",
				"--duration", "30s",
				"--threshold", "http_req_duration:p(95)<500",
			},
			expectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			command := exec.Command("go", testCase.args...)
			command.Dir = projectRoot

			output, err := command.CombinedOutput()
			if testCase.expectError && err == nil {
				t.Fatalf("expected error, got nil\noutput:\n%s", string(output))
			}

			if !testCase.expectError && err != nil {
				t.Fatalf("unexpected error: %v\noutput:\n%s", err, string(output))
			}
		})
	}
}

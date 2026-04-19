package unit

import (
	"path/filepath"
	"testing"

	"github.com/KirillRg/cli-tool/internal/parser"
)

func TestSubstituteEnvVariables(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		env      map[string]string
		expected string
	}{
		{
			name:  "single variable",
			input: "http://localtest/{{ _.env_param_one }}",
			env: map[string]string{
				"env_param_one": "valueOfEnvParamONE",
			},
			expected: "http://localtest/valueOfEnvParamONE",
		},
		{
			name:  "multiple variables",
			input: "{{ _.env_param_one }} {{ _.env_param_two }}",
			env: map[string]string{
				"env_param_one": "valueOfEnvParamONE",
				"env_param_two": "valueOfEnvParamTWO",
			},
			expected: "valueOfEnvParamONE valueOfEnvParamTWO",
		},
		{
			name:  "unknown variable stays unchanged",
			input: "{{ _.unknown_value }}",
			env: map[string]string{
				"env_param_one": "valueOfEnvParamONE",
			},
			expected: "{{ _.unknown_value }}",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := parser.SubstituteEnvVariables(testCase.input, testCase.env)

			if result != testCase.expected {
				t.Fatalf("expected %q, got %q", testCase.expected, result)
			}
		})
	}
}

func TestParseInsomniaCollection_WithEnvironmentVariables(t *testing.T) {
	collectionPath := filepath.Join("..", "..", "..", "Insomnia_Test_Collection_With_Environment.yaml") // точки чтобы пройти на директории повыше
	collection, err := parser.ParseInsomniaCollection(collectionPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(collection.Collection) != 4 {
		t.Fatalf("expected 4 requests, got %d", len(collection.Collection))
	}

	if collection.Collection[0].Parameters[0].Value != "valueOfEnvParamTHREE" {
		t.Fatalf("expected substituted parameter value, got %q", collection.Collection[0].Parameters[0].Value)
	}

	if collection.Collection[1].Body.Text != "{\n\t\"key_1\":\"value\",\n\t\"key_2\":\"valueOfEnvParamTWO\"\n}" {
		t.Fatalf("expected substituted body, got %q", collection.Collection[1].Body.Text)
	}

	if collection.Collection[2].URL != "http://localtest/deleting/valueOfEnvParamONE" {
		t.Fatalf("expected substituted delete url, got %q", collection.Collection[2].URL)
	}
}

package unit

import (
	"path/filepath"
	"testing"

	astpkg "github.com/KirillRg/cli-tool/internal/ast"
	"github.com/KirillRg/cli-tool/internal/parser"
	"github.com/KirillRg/cli-tool/internal/profile"
)

func TestGenerateAST_ImportExportBlocks(t *testing.T) {
	collectionPath := filepath.Join("..", "..", "..", "Insomnia_Test_Collection_With_Environment.yaml")

	collection, err := parser.ParseInsomniaCollection(collectionPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loadProfile, err := profile.Build(profile.BuildInput{
		InputPath:  "input.yaml",
		OutputPath: "output.js",
		VUs:        10,
		Duration:   "30s",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	program := astpkg.GenerateAST(collection, loadProfile)

	if len(program.Body) != 3 {
		t.Fatalf("expected 3 top-level nodes, got %d", len(program.Body))
	}

	if _, ok := program.Body[0].(*astpkg.ImportDeclaration); !ok {
		t.Fatalf("expected ImportDeclaration, got %T", program.Body[0])
	}

	exportOptions, ok := program.Body[1].(*astpkg.ExportNamedDeclaration)
	if !ok {
		t.Fatalf("expected ExportNamedDeclaration, got %T", program.Body[1])
	}

	optionsDeclaration, ok := exportOptions.Declaration.(*astpkg.VariableDeclaration)
	if !ok {
		t.Fatalf("expected VariableDeclaration, got %T", exportOptions.Declaration)
	}

	optionsIdentifier, ok := optionsDeclaration.Declarations[0].ID.(*astpkg.Identifier)
	if !ok || optionsIdentifier.Name != "options" {
		t.Fatalf("expected exported variable 'options', got %T", optionsDeclaration.Declarations[0].ID)
	}

	if _, ok := program.Body[2].(*astpkg.ExportDefaultDeclaration); !ok {
		t.Fatalf("expected ExportDefaultDeclaration, got %T", program.Body[2])
	}
}

func TestGenerateAST_DefaultFunctionRequests(t *testing.T) {
	collectionPath := filepath.Join("..", "..", "..", "Insomnia_Test_Collection_With_Environment.yaml")

	collection, err := parser.ParseInsomniaCollection(collectionPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loadProfile, err := profile.Build(profile.BuildInput{
		InputPath:  "input.yaml",
		OutputPath: "output.js",
		VUs:        10,
		Duration:   "30s",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	program := astpkg.GenerateAST(collection, loadProfile)

	exportDefault, ok := program.Body[2].(*astpkg.ExportDefaultDeclaration)
	if !ok {
		t.Fatalf("expected ExportDefaultDeclaration, got %T", program.Body[2])
	}

	defaultFunction, ok := exportDefault.Declaration.(*astpkg.FunctionExpression)
	if !ok {
		t.Fatalf("expected FunctionExpression, got %T", exportDefault.Declaration)
	}

	if len(defaultFunction.Body.Body) != len(collection.Collection) {
		t.Fatalf("expected %d request statements, got %d", len(collection.Collection), len(defaultFunction.Body.Body))
	}

	firstStatement, ok := defaultFunction.Body.Body[0].(*astpkg.ExpressionStatement)
	if !ok {
		t.Fatalf("expected first statement to be ExpressionStatement, got %T", defaultFunction.Body.Body[0])
	}

	firstCall, ok := firstStatement.Expression.(*astpkg.CallExpression)
	if !ok {
		t.Fatalf("expected first expression to be CallExpression, got %T", firstStatement.Expression)
	}

	callee, ok := firstCall.Callee.(*astpkg.MemberExpression)
	if !ok {
		t.Fatalf("expected call callee to be MemberExpression, got %T", firstCall.Callee)
	}

	property, ok := callee.Property.(*astpkg.Identifier)
	if !ok || property.Name != "request" {
		t.Fatalf("expected callee property to be request, got %T", callee.Property)
	}
}

func TestBuildOptionsAST_ProfileModes(t *testing.T) {
	testCases := []struct {
		name              string
		loadProfile       *profile.LoadProfile
		expectedFirstKeys []string
	}{
		{
			name: "constant vus profile",
			loadProfile: &profile.LoadProfile{
				Mode:     profile.ModeConstantVUs,
				VUs:      10,
				Duration: "30s",
			},
			expectedFirstKeys: []string{"vus", "duration"},
		},
		{
			name: "shared iterations profile",
			loadProfile: &profile.LoadProfile{
				Mode:       profile.ModeSharedIterations,
				VUs:        5,
				Iterations: 100,
			},
			expectedFirstKeys: []string{"vus", "iterations"},
		},
		{
			name: "stages profile",
			loadProfile: &profile.LoadProfile{
				Mode: profile.ModeStages,
				Stages: []profile.StageConfig{
					{Duration: "30s", Target: 5},
					{Duration: "1m", Target: 20},
				},
			},
			expectedFirstKeys: []string{"stages"},
		},
		{
			name: "constant vus with thresholds",
			loadProfile: &profile.LoadProfile{
				Mode:     profile.ModeConstantVUs,
				VUs:      10,
				Duration: "30s",
				Thresholds: map[string][]string{
					"http_req_duration": {"p(95)<500"},
				},
			},
			expectedFirstKeys: []string{"vus", "duration", "thresholds"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			optionsAST := astpkg.BuildOptionsAST(testCase.loadProfile)

			if len(optionsAST.Properties) != len(testCase.expectedFirstKeys) {
				t.Fatalf("expected %d properties, got %d", len(testCase.expectedFirstKeys), len(optionsAST.Properties))
			}

			for propertyIndex, expectedKey := range testCase.expectedFirstKeys {
				identifierKey, isIdentifier := optionsAST.Properties[propertyIndex].Key.(*astpkg.Identifier)
				if !isIdentifier {
					t.Fatalf("expected identifier key at index %d, got %T", propertyIndex, optionsAST.Properties[propertyIndex].Key)
				}

				if identifierKey.Name != expectedKey {
					t.Fatalf("expected key %q at index %d, got %q", expectedKey, propertyIndex, identifierKey.Name)
				}
			}
		})
	}
}

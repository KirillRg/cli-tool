package unit

import (
	"strings"
	"testing"

	"github.com/KirillRg/cli-tool/internal/ast"
	"github.com/KirillRg/cli-tool/internal/translator"
)

func TestTranslateProgram_MainK6Blocks(t *testing.T) {
	program := ast.Program{
		Type:       "Program",
		SourceType: "module",
		Body: []ast.Node{
			&ast.ImportDeclaration{
				Type: "ImportDeclaration",
				Specifiers: []ast.ImportSpecifier{
					&ast.ImportDefaultSpecifier{
						Type: "ImportDefaultSpecifier",
						Local: &ast.Identifier{
							Type: "Identifier",
							Name: "http",
						},
					},
				},
				Source: &ast.Literal{
					Type:  "Literal",
					Value: "k6/http",
				},
			},
			&ast.ExportNamedDeclaration{
				Type: "ExportNamedDeclaration",
				Declaration: &ast.VariableDeclaration{
					Type: "VariableDeclaration",
					Kind: "const",
					Declarations: []*ast.VariableDeclarator{
						{
							Type: "VariableDeclarator",
							ID: &ast.Identifier{
								Type: "Identifier",
								Name: "options",
							},
							Init: &ast.ObjectExpression{
								Type: "ObjectExpression",
								Properties: []*ast.Property{
									{
										Type: "Property",
										Key: &ast.Identifier{
											Type: "Identifier",
											Name: "vus",
										},
										Value: &ast.Literal{
											Type:  "Literal",
											Value: 10,
										},
									},
									{
										Type: "Property",
										Key: &ast.Identifier{
											Type: "Identifier",
											Name: "duration",
										},
										Value: &ast.Literal{
											Type:  "Literal",
											Value: "30s",
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.ExportDefaultDeclaration{
				Type: "ExportDefaultDeclaration",
				Declaration: &ast.FunctionExpression{
					Type: "FunctionExpression",
					Body: &ast.BlockStatement{
						Type: "BlockStatement",
						Body: []ast.Statement{
							&ast.ExpressionStatement{
								Type: "ExpressionStatement",
								Expression: &ast.CallExpression{
									Type: "CallExpression",
									Callee: &ast.MemberExpression{
										Type: "MemberExpression",
										Object: &ast.Identifier{
											Type: "Identifier",
											Name: "http",
										},
										Property: &ast.Identifier{
											Type: "Identifier",
											Name: "request",
										},
									},
									Arguments: []ast.Expression{
										&ast.Literal{Type: "Literal", Value: "GET"},
										&ast.Literal{Type: "Literal", Value: "http://localtest/getting"},
										&ast.Literal{Type: "Literal", Value: nil},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := translator.TranslateProgram(program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFragments := []string{
		`import http from "k6/http";`,
		`export const options`,
		`vus: 10`,
		`duration: "30s"`,
		`export default function()`,
		`http.request("GET", "http://localtest/getting", null);`,
	}

	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(result, expectedFragment) {
			t.Fatalf("expected fragment %q in result:\n%s", expectedFragment, result)
		}
	}
}

func TestTranslateProgram_StagesOptions(t *testing.T) {
	program := ast.Program{
		Type:       "Program",
		SourceType: "module",
		Body: []ast.Node{
			&ast.ExportNamedDeclaration{
				Type: "ExportNamedDeclaration",
				Declaration: &ast.VariableDeclaration{
					Type: "VariableDeclaration",
					Kind: "const",
					Declarations: []*ast.VariableDeclarator{
						{
							Type: "VariableDeclarator",
							ID: &ast.Identifier{
								Type: "Identifier",
								Name: "options",
							},
							Init: &ast.ObjectExpression{
								Type: "ObjectExpression",
								Properties: []*ast.Property{
									{
										Type: "Property",
										Key: &ast.Identifier{
											Type: "Identifier",
											Name: "stages",
										},
										Value: &ast.ArrayExpression{
											Type: "ArrayExpression",
											Elements: []ast.Expression{
												&ast.ObjectExpression{
													Type: "ObjectExpression",
													Properties: []*ast.Property{
														{
															Type: "Property",
															Key: &ast.Identifier{
																Type: "Identifier",
																Name: "duration",
															},
															Value: &ast.Literal{
																Type:  "Literal",
																Value: "30s",
															},
														},
														{
															Type: "Property",
															Key: &ast.Identifier{
																Type: "Identifier",
																Name: "target",
															},
															Value: &ast.Literal{
																Type:  "Literal",
																Value: 5,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := translator.TranslateProgram(program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFragments := []string{
		`stages`,
		`duration: "30s"`,
		`target: 5`,
	}

	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(result, expectedFragment) {
			t.Fatalf("expected fragment %q in result:\n%s", expectedFragment, result)
		}
	}
}

package ast

import (
	"github.com/KirillRg/cli-tool/internal/parser"
)

func GenerateAST(collection *parser.InsomniaCollection) *Program {
	var statements []StatementNode

	for _, req := range collection.Collection {
		method := req.Method
		url := req.URL

		call := CallExpression{
			Callee: Identifier{Name: getHttpFunctionName(method)},
			Arguments: []ExpressionNode{
				Literal{Value: url},
			},
		}

		stmt := ExpressionStatement{Expression: call}
		statements = append(statements, stmt)
	}

	return &Program{
		Body: statements,
	}
}

func getHttpFunctionName(method string) string {
	switch method {
	case "GET":
		return "http.get"
	case "POST":
		return "http.post"
	case "PUT":
		return "http.put"
	case "DELETE":
		return "http.del"
	default:
		return "http.unknown"
	}
}

package ast

import (
	"github.com/KirillRg/cli-tool/internal/parser"
)

// Генерация AST для коллекции
func GenerateAST(collection *parser.InsomniaCollection) Program {
	var body []Node
	for _, req := range collection.Collection {
		body = append(body, GenerateRequestAST(req))
	}
	return Program{Body: body}
}

// Генерация универсального запроса в AST
func GenerateRequestAST(req parser.RequestItem) Node {
	return ExpressionStatement{
		Expression: CallExpression{
			Callee: MemberExpression{
				Object:   Identifier{Name: "http"},
				Property: Identifier{Name: req.Method},
			},
			Arguments: []Node{
				Literal{Value: req.URL},
				ObjectExpression{
					Properties: []Property{
						{Key: Identifier{Name: "headers"}, Value: GenerateHeaders(req.Headers)},
						{Key: Identifier{Name: "params"}, Value: GenerateParameters(req.Parameters)},
						{Key: Identifier{Name: "body"}, Value: GenerateBody(req.Body)},
					},
				},
			},
		},
	}
}

// Генерация параметров запроса
func GenerateParameters(params []parser.RequestParam) Node {
	properties := []Property{}
	for _, param := range params {
		if !param.Disabled {
			properties = append(properties, Property{
				Key:   Identifier{Name: param.Name},
				Value: Literal{Value: param.Value},
			})
		}
	}
	return ObjectExpression{Properties: properties}
}

// Генерация заголовков запроса
func GenerateHeaders(headers []parser.RequestHeader) Node {
	properties := []Property{}
	for _, header := range headers {
		if !header.Disabled {
			properties = append(properties, Property{
				Key:   Identifier{Name: header.Name},
				Value: Literal{Value: header.Value},
			})
		}
	}
	return ObjectExpression{Properties: properties}
}

// Генерация тела запроса
func GenerateBody(body parser.RequestBody) Node {
	return ObjectExpression{
		Properties: []Property{
			{Key: Identifier{Name: "mimeType"}, Value: Literal{Value: body.MimeType}},
			{Key: Identifier{Name: "text"}, Value: Literal{Value: body.Text}},
		},
	}
}

package ast

import "github.com/KirillRg/cli-tool/internal/parser"

// Генерация AST для коллекции
func GenerateAST(collection *parser.InsomniaCollection) Program {
	var body []Statement
	for _, req := range collection.Collection {
		body = append(body, GenerateRequestAST(req))
	}
	return Program{Body: body}
}

func GenerateRequestAST(req parser.RequestItem) Statement {
	args := []Expression{Literal{Value: req.URL}}

	var properties []Property

	if len(req.Headers) > 0 {
		headers := GenerateHeaders(req.Headers)
		if len(headers.Properties) > 0 {
			properties = append(properties, Property{
				Key:   Identifier{Name: "headers"},
				Value: headers,
			})
		}
	}

	if len(req.Parameters) > 0 {
		params := GenerateParameters(req.Parameters)
		if len(params.Properties) > 0 {
			properties = append(properties, Property{
				Key:   Identifier{Name: "params"},
				Value: params,
			})
		}
	}

	if req.Body.MimeType != "" || req.Body.Text != "" {
		properties = append(properties, Property{
			Key:   Identifier{Name: "body"},
			Value: GenerateBody(req.Body),
		})
	}

	if len(properties) > 0 {
		args = append(args, ObjectExpression{Properties: properties})
	}

	return ExpressionStatement{
		Expression: CallExpression{
			Callee: MemberExpression{
				Object:   Identifier{Name: "http"},
				Property: Identifier{Name: req.Method},
			},
			Arguments: args,
		},
	}
}

func GenerateHeaders(headers []parser.RequestHeader) ObjectExpression {
	var properties []Property
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

func GenerateParameters(params []parser.RequestParam) ObjectExpression {
	var properties []Property
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

func GenerateBody(body parser.RequestBody) ObjectExpression {
	var properties []Property
	if body.MimeType != "" {
		properties = append(properties, Property{
			Key:   Identifier{Name: "mimeType"},
			Value: Literal{Value: body.MimeType},
		})
	}
	if body.Text != "" {
		properties = append(properties, Property{
			Key:   Identifier{Name: "text"},
			Value: Literal{Value: body.Text},
		})
	}
	return ObjectExpression{Properties: properties}
}

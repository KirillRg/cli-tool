package ast

import "github.com/KirillRg/cli-tool/internal/parser"

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
	// Универсальная генерация для всех типов запросов
	args := []Node{Literal{Value: req.URL}}

	properties := []Property{}

	// Генерация заголовков, если они есть и активны
	if len(req.Headers) > 0 {
		headers := GenerateHeaders(req.Headers)
		if len(headers.Properties) > 0 {
			properties = append(properties, Property{
				Key:   Identifier{Name: "headers"},
				Value: headers,
			})
		}
	}

	// Генерация параметров, если они есть и активны
	if len(req.Parameters) > 0 {
		params := GenerateParameters(req.Parameters)
		if len(params.Properties) > 0 {
			properties = append(properties, Property{
				Key:   Identifier{Name: "params"},
				Value: params,
			})
		}
	}

	// Генерация тела запроса, если оно есть
	if req.Body.MimeType != "" || req.Body.Text != "" {
		properties = append(properties, Property{
			Key:   Identifier{Name: "body"},
			Value: GenerateBody(req.Body),
		})
	}

	// Добавляем свойства (headers, params, body) в аргументы
	if len(properties) > 0 {
		args = append(args, ObjectExpression{Properties: properties})
	}

	// Формируем окончательную структуру вызова
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

// Генерация заголовков запроса
func GenerateHeaders(headers []parser.RequestHeader) ObjectExpression {
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

// Генерация параметров запроса
func GenerateParameters(params []parser.RequestParam) ObjectExpression {
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

// Генерация тела запроса
func GenerateBody(body parser.RequestBody) ObjectExpression {
	properties := []Property{}
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

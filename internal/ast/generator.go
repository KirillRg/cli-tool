package ast

import (
	"net/url"
	"strings"
	"unicode"

	"github.com/KirillRg/cli-tool/internal/parser"
)

// Генерация AST для коллекции
func GenerateAST(collection *parser.InsomniaCollection) Program {
	var body []Statement
	for _, req := range collection.Collection {
		body = append(body, GenerateRequestAST(req))
	}
	return Program{Body: body}
}

func GenerateRequestAST(req parser.RequestItem) Statement {

	urlWithQuery := appendQueryString(req.URL, req.Parameters)

	args := []Expression{
		Literal{Value: req.Method},
		Literal{Value: urlWithQuery},
	}

	bodyArg := GenerateBody(req.Body)
	args = append(args, bodyArg)

	headersForRequest := req.Headers
	if req.Body.MimeType != "" && !hasEnabledHeader(headersForRequest, "Content-Type") {
		headersForRequest = append(headersForRequest, parser.RequestHeader{
			Name:     "Content-Type",
			Value:    req.Body.MimeType,
			Disabled: false,
		})
	}

	var paramsProperties []Property
	if len(headersForRequest) > 0 {
		headers := GenerateHeaders(headersForRequest)
		if len(headers.Properties) > 0 {
			paramsProperties = append(paramsProperties, Property{
				Key:   Identifier{Name: "headers"},
				Value: headers,
			})
		}
	}

	// Добавляем 4-й аргумент только если есть содержимое (например headers)
	if len(paramsProperties) > 0 {
		args = append(args, ObjectExpression{Properties: paramsProperties})
	}

	return ExpressionStatement{
		Expression: CallExpression{
			Callee: MemberExpression{
				Object:   Identifier{Name: "http"},
				Property: Identifier{Name: "request"},
			},
			Arguments: args,
		},
	}
}

func GenerateHeaders(headers []parser.RequestHeader) ObjectExpression {
	var properties []Property
	for _, header := range headers {
		if header.Disabled {
			continue
		}
		var key Expression
		if isValidJSIdentifier(header.Name) {
			key = Identifier{Name: header.Name}
		} else {
			key = Literal{Value: header.Name}
		}

		properties = append(properties, Property{
			Key:   key,
			Value: Literal{Value: header.Value},
		})
	}
	return ObjectExpression{Properties: properties}
}

func GenerateBody(body parser.RequestBody) Expression {
	if body.Text != "" {
		return Literal{Value: body.Text}
	}

	if body.MimeType != "" {
		return Literal{Value: ""}
	}

	return Identifier{Name: "null"}
}

// HELPERS
// Единый метод добавления query string через строковую обработку.
func appendQueryString(rawURL string, params []parser.RequestParam) string {

	type pair struct{ k, v string }
	var pairs []pair
	for _, p := range params {
		if p.Disabled || p.Name == "" {
			continue
		}
		pairs = append(pairs, pair{
			k: url.QueryEscape(p.Name),
			v: url.QueryEscape(p.Value),
		})
	}
	if len(pairs) == 0 {
		return rawURL
	}

	base := rawURL
	frag := ""
	if i := strings.IndexByte(rawURL, '#'); i >= 0 {
		base = rawURL[:i]
		frag = rawURL[i:]
	}

	sep := "?"
	if strings.Contains(base, "?") {
		if strings.HasSuffix(base, "?") || strings.HasSuffix(base, "&") {
			sep = ""
		} else {
			sep = "&"
		}
	}

	var b strings.Builder
	b.WriteString(base)
	b.WriteString(sep)

	for i, kv := range pairs {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(kv.k)
		b.WriteByte('=')
		b.WriteString(kv.v)
	}

	b.WriteString(frag)
	return b.String()
}

// Проверка на
func isValidJSIdentifier(s string) bool {
	if s == "" {
		return false
	}

	for i, r := range s {
		if i == 0 {
			if !(r == '_' || r == '$' || unicode.IsLetter(r)) {
				return false
			}
		} else {
			if !(r == '_' || r == '$' || unicode.IsLetter(r) || unicode.IsDigit(r)) {
				return false
			}
		}
	}
	return true
}

// Проверка включенного header (связь с body и Content-Type в нём)
func hasEnabledHeader(headers []parser.RequestHeader, name string) bool {
	for _, h := range headers {
		if h.Disabled {
			continue
		}
		if strings.EqualFold(h.Name, name) {
			return true
		}
	}
	return false
}

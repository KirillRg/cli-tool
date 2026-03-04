package ast

import (
	"net/url"
	"strings"
	"unicode"

	"github.com/KirillRg/cli-tool/internal/parser"
)

// Генерация AST для коллекции
func GenerateAST(collection *parser.InsomniaCollection) Program {
	// 1) Statements для default function body
	var requestStatements []Statement
	for _, req := range collection.Collection {
		requestStatements = append(requestStatements, GenerateRequestAST(req))
	}

	// 2) import http from "k6/http";
	importHttp := &ImportDeclaration{
		Type: "ImportDeclaration",
		Specifiers: []ImportSpecifier{
			&ImportDefaultSpecifier{
				Type:  "ImportDefaultSpecifier",
				Local: &Identifier{Type: "Identifier", Name: "http"},
			},
		},
		Source: &Literal{Type: "Literal", Value: "k6/http"},
	}

	// 3) export const options = { vus: 1, duration: "10s" };
	optionsObj := &ObjectExpression{
		Type: "ObjectExpression",
		Properties: []*Property{
			{
				Type:      "Property",
				Key:       &Identifier{Type: "Identifier", Name: "vus"},
				Value:     &Literal{Type: "Literal", Value: 1},
				Kind:      "init",
				Method:    false,
				Shorthand: false,
				Computed:  false,
			},
			{
				Type:      "Property",
				Key:       &Identifier{Type: "Identifier", Name: "duration"},
				Value:     &Literal{Type: "Literal", Value: "10s"},
				Kind:      "init",
				Method:    false,
				Shorthand: false,
				Computed:  false,
			},
		},
	}

	optionsDecl := &VariableDeclaration{
		Type: "VariableDeclaration",
		Kind: "const",
		Declarations: []*VariableDeclarator{
			{
				Type: "VariableDeclarator",
				ID:   &Identifier{Type: "Identifier", Name: "options"},
				Init: optionsObj,
			},
		},
	}

	exportOptions := &ExportNamedDeclaration{
		Type:        "ExportNamedDeclaration",
		Declaration: optionsDecl,
	}

	// 4) export default function () { ...requests... }
	defaultFunc := &FunctionExpression{
		Type:   "FunctionExpression",
		ID:     nil,
		Params: []Pattern{},
		Body: &BlockStatement{
			Type: "BlockStatement",
			Body: requestStatements,
		},
	}

	exportDefault := &ExportDefaultDeclaration{
		Type:        "ExportDefaultDeclaration",
		Declaration: defaultFunc,
	}

	// 5) Program { type:"Program", sourceType:"module", body:[import, export options, export default] }
	return Program{
		Type:       "Program",
		SourceType: "module",
		Body: []Node{
			importHttp,
			exportOptions,
			exportDefault,
		},
	}
}

// GenerateRequestAST строит ExpressionStatement с CallExpression(http.request(...))
func GenerateRequestAST(req parser.RequestItem) Statement {
	urlWithQuery := appendQueryString(req.URL, req.Parameters)

	args := []Expression{
		&Literal{Type: "Literal", Value: req.Method},
		&Literal{Type: "Literal", Value: urlWithQuery},
		GenerateBody(req.Body),
	}

	headersForRequest := req.Headers
	if req.Body.MimeType != "" && !hasEnabledHeader(headersForRequest, "Content-Type") {
		headersForRequest = append(headersForRequest, parser.RequestHeader{
			Name:     "Content-Type",
			Value:    req.Body.MimeType,
			Disabled: false,
		})
	}

	var paramsProps []*Property
	if len(headersForRequest) > 0 {
		headers := GenerateHeaders(headersForRequest)
		if len(headers.Properties) > 0 {
			paramsProps = append(paramsProps, &Property{
				Type:      "Property",
				Key:       &Identifier{Type: "Identifier", Name: "headers"},
				Value:     headers,
				Kind:      "init",
				Method:    false,
				Shorthand: false,
				Computed:  false,
			})
		}
	}

	if len(paramsProps) > 0 {
		args = append(args, &ObjectExpression{
			Type:       "ObjectExpression",
			Properties: paramsProps,
		})
	}

	call := &CallExpression{
		Type: "CallExpression",
		Callee: &MemberExpression{
			Type:     "MemberExpression",
			Object:   &Identifier{Type: "Identifier", Name: "http"},
			Property: &Identifier{Type: "Identifier", Name: "request"},
			Computed: false,
		},
		Arguments: args,
	}

	return &ExpressionStatement{
		Type:       "ExpressionStatement",
		Expression: call,
	}
}

func GenerateHeaders(headers []parser.RequestHeader) *ObjectExpression {
	var props []*Property
	for _, header := range headers {
		if header.Disabled {
			continue
		}

		var key Expression
		if isValidJSIdentifier(header.Name) {
			key = &Identifier{Type: "Identifier", Name: header.Name}
		} else {
			key = &Literal{Type: "Literal", Value: header.Name}
		}

		props = append(props, &Property{
			Type:      "Property",
			Key:       key,
			Value:     &Literal{Type: "Literal", Value: header.Value},
			Kind:      "init",
			Method:    false,
			Shorthand: false,
			Computed:  false,
		})
	}

	return &ObjectExpression{
		Type:       "ObjectExpression",
		Properties: props,
	}
}

func GenerateBody(body parser.RequestBody) Expression {
	if body.Text != "" {
		return &Literal{Type: "Literal", Value: body.Text}
	}
	if body.MimeType != "" {
		return &Literal{Type: "Literal", Value: ""}
	}
	return &Literal{Type: "Literal", Value: nil}
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

// Проверка на строковое соответствие JS
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

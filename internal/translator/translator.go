package translator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/KirillRg/cli-tool/internal/ast"
)

// TranslateProgram converts an ESTree Program into executable JavaScript code.
func TranslateProgram(program ast.Program) (string, error) {
	translationContext := newTranslationContext()

	if err := translationContext.translateNode(program); err != nil {
		return "", err
	}

	resultCode := translationContext.String()
	if !strings.HasSuffix(resultCode, "\n") {
		resultCode += "\n"
	}

	return resultCode, nil
}

type translationContext struct {
	builder     strings.Builder
	indentLevel int
	registry    translatorRegistry
}

func newTranslationContext() *translationContext {
	context := &translationContext{}
	context.registry = buildTranslatorRegistry()
	return context
}

func (context *translationContext) String() string {
	return context.builder.String()
}

func (context *translationContext) write(text string) {
	context.builder.WriteString(text)
}

func (context *translationContext) writeNewline() {
	context.builder.WriteByte('\n')
}

func (context *translationContext) writeIndent() {
	for indentIndex := 0; indentIndex < context.indentLevel; indentIndex++ {
		context.builder.WriteString("  ")
	}
}

func (context *translationContext) withIncreasedIndent(action func() error) error {
	context.indentLevel++
	err := action()
	context.indentLevel--
	return err
}

func (context *translationContext) translateNode(node ast.Node) error {
	if node == nil {
		return errors.New("translator: node is nil")
	}

	nodeType := reflect.TypeOf(node)
	nodeTranslator, hasTranslator := context.registry[nodeType]
	if !hasTranslator {
		return fmt.Errorf("translator: unsupported node type %T", node)
	}

	return nodeTranslator.Translate(context, node)
}

type nodeTranslator interface {
	Translate(context *translationContext, node ast.Node) error
}

type translatorRegistry map[reflect.Type]nodeTranslator

func buildTranslatorRegistry() translatorRegistry {
	return translatorRegistry{
		reflect.TypeOf(ast.Program{}):                   programTranslator{},
		reflect.TypeOf(&ast.ImportDeclaration{}):        importDeclarationTranslator{},
		reflect.TypeOf(&ast.ExportNamedDeclaration{}):   exportNamedDeclarationTranslator{},
		reflect.TypeOf(&ast.ExportDefaultDeclaration{}): exportDefaultDeclarationTranslator{},
		reflect.TypeOf(&ast.VariableDeclaration{}):      variableDeclarationTranslator{},
		reflect.TypeOf(&ast.VariableDeclarator{}):       variableDeclaratorTranslator{},
		reflect.TypeOf(&ast.FunctionExpression{}):       functionExpressionTranslator{},
		reflect.TypeOf(&ast.BlockStatement{}):           blockStatementTranslator{},
		reflect.TypeOf(&ast.ExpressionStatement{}):      expressionStatementTranslator{},
		reflect.TypeOf(&ast.CallExpression{}):           callExpressionTranslator{},
		reflect.TypeOf(&ast.MemberExpression{}):         memberExpressionTranslator{},
		reflect.TypeOf(&ast.ObjectExpression{}):         objectExpressionTranslator{},
		reflect.TypeOf(&ast.Property{}):                 propertyTranslator{},
		reflect.TypeOf(&ast.Identifier{}):               identifierTranslator{},
		reflect.TypeOf(&ast.Literal{}):                  literalTranslator{},
		reflect.TypeOf(&ast.ArrayExpression{}):          arrayExpressionTranslator{},
	}
}

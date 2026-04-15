package translator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/KirillRg/cli-tool/internal/ast"
)

type programTranslator struct{}

func (programTranslator) Translate(context *translationContext, node ast.Node) error {
	programNode, isProgram := node.(ast.Program)
	if !isProgram {
		return fmt.Errorf("translator: expected ast.Program, got %T", node)
	}

	for statementIndex, statementNode := range programNode.Body {
		if err := context.translateNode(statementNode); err != nil {
			return err
		}

		if statementIndex != len(programNode.Body)-1 {
			context.writeNewline()
		}
	}

	return nil
}

type importDeclarationTranslator struct{}

func (importDeclarationTranslator) Translate(context *translationContext, node ast.Node) error {
	importDeclaration, isImportDeclaration := node.(*ast.ImportDeclaration)
	if !isImportDeclaration {
		return fmt.Errorf("translator: expected *ast.ImportDeclaration, got %T", node)
	}

	if importDeclaration.Source == nil {
		return errors.New("translator: import declaration source is nil")
	}

	if len(importDeclaration.Specifiers) != 1 {
		return fmt.Errorf(
			"translator: only one import specifier is supported, got %d",
			len(importDeclaration.Specifiers),
		)
	}

	defaultSpecifier, isDefaultSpecifier := importDeclaration.Specifiers[0].(*ast.ImportDefaultSpecifier)
	if !isDefaultSpecifier || defaultSpecifier.Local == nil {
		return errors.New("translator: only ImportDefaultSpecifier with local identifier is supported")
	}

	context.write("import ")
	context.write(defaultSpecifier.Local.Name)
	context.write(" from ")

	if err := context.translateNode(importDeclaration.Source); err != nil {
		return err
	}

	context.write(";")
	context.writeNewline()
	return nil
}

type exportNamedDeclarationTranslator struct{}

func (exportNamedDeclarationTranslator) Translate(context *translationContext, node ast.Node) error {
	exportDeclaration, isExportDeclaration := node.(*ast.ExportNamedDeclaration)
	if !isExportDeclaration {
		return fmt.Errorf("translator: expected *ast.ExportNamedDeclaration, got %T", node)
	}

	if exportDeclaration.Declaration == nil {
		return errors.New("translator: export named declaration is missing nested declaration")
	}

	context.write("export ")
	return context.translateNode(exportDeclaration.Declaration)
}

type exportDefaultDeclarationTranslator struct{}

func (exportDefaultDeclarationTranslator) Translate(context *translationContext, node ast.Node) error {
	exportDeclaration, isExportDeclaration := node.(*ast.ExportDefaultDeclaration)
	if !isExportDeclaration {
		return fmt.Errorf("translator: expected *ast.ExportDefaultDeclaration, got %T", node)
	}

	if exportDeclaration.Declaration == nil {
		return errors.New("translator: export default declaration is missing nested declaration")
	}

	context.write("export default ")

	if err := context.translateNode(exportDeclaration.Declaration); err != nil {
		return err
	}

	context.writeNewline()
	return nil
}

type variableDeclarationTranslator struct{}

func (variableDeclarationTranslator) Translate(context *translationContext, node ast.Node) error {
	variableDeclaration, isVariableDeclaration := node.(*ast.VariableDeclaration)
	if !isVariableDeclaration {
		return fmt.Errorf("translator: expected *ast.VariableDeclaration, got %T", node)
	}

	context.write(variableDeclaration.Kind)
	context.write(" ")

	for declarationIndex, declaratorNode := range variableDeclaration.Declarations {
		if declaratorNode == nil {
			return errors.New("translator: variable declarator is nil")
		}

		if err := context.translateNode(declaratorNode); err != nil {
			return err
		}

		if declarationIndex != len(variableDeclaration.Declarations)-1 {
			context.write(", ")
		}
	}

	context.write(";")
	context.writeNewline()
	return nil
}

type variableDeclaratorTranslator struct{}

func (variableDeclaratorTranslator) Translate(context *translationContext, node ast.Node) error {
	variableDeclarator, isVariableDeclarator := node.(*ast.VariableDeclarator)
	if !isVariableDeclarator {
		return fmt.Errorf("translator: expected *ast.VariableDeclarator, got %T", node)
	}

	identifierNode, isIdentifier := variableDeclarator.ID.(*ast.Identifier)
	if !isIdentifier {
		return fmt.Errorf(
			"translator: variable declarator identifier must be *ast.Identifier, got %T",
			variableDeclarator.ID,
		)
	}

	context.write(identifierNode.Name)
	context.write(" = ")

	if variableDeclarator.Init == nil {
		context.write("null")
		return nil
	}

	return context.translateNode(variableDeclarator.Init)
}

type functionExpressionTranslator struct{}

func (functionExpressionTranslator) Translate(context *translationContext, node ast.Node) error {
	functionExpression, isFunctionExpression := node.(*ast.FunctionExpression)
	if !isFunctionExpression {
		return fmt.Errorf("translator: expected *ast.FunctionExpression, got %T", node)
	}

	if functionExpression.Body == nil {
		return errors.New("translator: function expression body is nil")
	}

	context.write("function")

	if functionExpression.ID != nil {
		context.write(" ")
		context.write(functionExpression.ID.Name)
	}

	context.write("(")

	for parameterIndex, parameterNode := range functionExpression.Params {
		identifierParameter, isIdentifierParameter := parameterNode.(*ast.Identifier)
		if !isIdentifierParameter {
			return fmt.Errorf(
				"translator: function parameter must be *ast.Identifier, got %T",
				parameterNode,
			)
		}

		if parameterIndex > 0 {
			context.write(", ")
		}

		context.write(identifierParameter.Name)
	}

	context.write(") ")
	return context.translateNode(functionExpression.Body)
}

type blockStatementTranslator struct{}

func (blockStatementTranslator) Translate(context *translationContext, node ast.Node) error {
	blockStatement, isBlockStatement := node.(*ast.BlockStatement)
	if !isBlockStatement {
		return fmt.Errorf("translator: expected *ast.BlockStatement, got %T", node)
	}

	context.write("{")
	context.writeNewline()

	return context.withIncreasedIndent(func() error {
		for _, statementNode := range blockStatement.Body {
			if statementNode == nil {
				return errors.New("translator: block statement contains nil child")
			}

			context.writeIndent()

			if err := context.translateNode(statementNode); err != nil {
				return err
			}
		}

		context.write("}")
		context.writeNewline()
		return nil
	})
}

type expressionStatementTranslator struct{}

func (expressionStatementTranslator) Translate(context *translationContext, node ast.Node) error {
	expressionStatement, isExpressionStatement := node.(*ast.ExpressionStatement)
	if !isExpressionStatement {
		return fmt.Errorf("translator: expected *ast.ExpressionStatement, got %T", node)
	}

	if expressionStatement.Expression == nil {
		return errors.New("translator: expression statement expression is nil")
	}

	if err := context.translateNode(expressionStatement.Expression); err != nil {
		return err
	}

	context.write(";")
	context.writeNewline()
	return nil
}

type callExpressionTranslator struct{}

func (callExpressionTranslator) Translate(context *translationContext, node ast.Node) error {
	callExpression, isCallExpression := node.(*ast.CallExpression)
	if !isCallExpression {
		return fmt.Errorf("translator: expected *ast.CallExpression, got %T", node)
	}

	if callExpression.Callee == nil {
		return errors.New("translator: call expression callee is nil")
	}

	if err := context.translateNode(callExpression.Callee); err != nil {
		return err
	}

	context.write("(")

	for argumentIndex, argumentNode := range callExpression.Arguments {
		if argumentIndex > 0 {
			context.write(", ")
		}

		if argumentNode == nil {
			context.write("null")
			continue
		}

		if err := context.translateNode(argumentNode); err != nil {
			return err
		}
	}

	context.write(")")
	return nil
}

type memberExpressionTranslator struct{}

func (memberExpressionTranslator) Translate(context *translationContext, node ast.Node) error {
	memberExpression, isMemberExpression := node.(*ast.MemberExpression)
	if !isMemberExpression {
		return fmt.Errorf("translator: expected *ast.MemberExpression, got %T", node)
	}

	if memberExpression.Object == nil {
		return errors.New("translator: member expression object is nil")
	}

	if memberExpression.Property == nil {
		return errors.New("translator: member expression property is nil")
	}

	if err := context.translateNode(memberExpression.Object); err != nil {
		return err
	}

	if memberExpression.Computed {
		context.write("[")
		if err := context.translateNode(memberExpression.Property); err != nil {
			return err
		}
		context.write("]")
		return nil
	}

	identifierProperty, isIdentifierProperty := memberExpression.Property.(*ast.Identifier)
	if isIdentifierProperty {
		context.write(".")
		context.write(identifierProperty.Name)
		return nil
	}

	context.write("[")
	if err := context.translateNode(memberExpression.Property); err != nil {
		return err
	}
	context.write("]")
	return nil
}

type objectExpressionTranslator struct{}

func (objectExpressionTranslator) Translate(context *translationContext, node ast.Node) error {
	objectExpression, isObjectExpression := node.(*ast.ObjectExpression)
	if !isObjectExpression {
		return fmt.Errorf("translator: expected *ast.ObjectExpression, got %T", node)
	}

	if len(objectExpression.Properties) == 0 {
		context.write("{}")
		return nil
	}

	context.write("{")
	context.writeNewline()

	return context.withIncreasedIndent(func() error {
		for propertyIndex, propertyNode := range objectExpression.Properties {
			if propertyNode == nil {
				return errors.New("translator: object expression property is nil")
			}

			context.writeIndent()

			if err := context.translateNode(propertyNode); err != nil {
				return err
			}

			if propertyIndex != len(objectExpression.Properties)-1 {
				context.write(",")
			}

			context.writeNewline()
		}

		context.writeIndent()
		context.write("}")
		return nil
	})
}

type propertyTranslator struct{}

func (propertyTranslator) Translate(context *translationContext, node ast.Node) error {
	propertyNode, isProperty := node.(*ast.Property)
	if !isProperty {
		return fmt.Errorf("translator: expected *ast.Property, got %T", node)
	}

	if propertyNode.Key == nil {
		return errors.New("translator: property key is nil")
	}

	switch typedKey := propertyNode.Key.(type) {
	case *ast.Identifier:
		context.write(typedKey.Name)
	case *ast.Literal:
		if err := context.translateNode(typedKey); err != nil {
			return err
		}
	default:
		return fmt.Errorf("translator: unsupported property key type %T", propertyNode.Key)
	}

	context.write(": ")

	if propertyNode.Value == nil {
		context.write("null")
		return nil
	}

	return context.translateNode(propertyNode.Value)
}

type identifierTranslator struct{}

func (identifierTranslator) Translate(context *translationContext, node ast.Node) error {
	identifierNode, isIdentifier := node.(*ast.Identifier)
	if !isIdentifier {
		return fmt.Errorf("translator: expected *ast.Identifier, got %T", node)
	}

	context.write(identifierNode.Name)
	return nil
}

type literalTranslator struct{}

func (literalTranslator) Translate(context *translationContext, node ast.Node) error {
	literalNode, isLiteral := node.(*ast.Literal)
	if !isLiteral {
		return fmt.Errorf("translator: expected *ast.Literal, got %T", node)
	}

	if literalNode.Value == nil {
		context.write("null")
		return nil
	}

	switch literalValue := literalNode.Value.(type) {
	case string:
		context.write(strconv.Quote(literalValue))
	case bool:
		if literalValue {
			context.write("true")
		} else {
			context.write("false")
		}
	case float64, float32, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		context.write(fmt.Sprintf("%v", literalValue))
	default:
		context.write(fmt.Sprintf("%v", literalValue))
	}

	return nil
}

type arrayExpressionTranslator struct{}

func (arrayExpressionTranslator) Translate(context *translationContext, node ast.Node) error {
	arrayExpression, isArrayExpression := node.(*ast.ArrayExpression)
	if !isArrayExpression {
		return fmt.Errorf("translator: expected *ast.ArrayExpression, got %T", node)
	}

	if len(arrayExpression.Elements) == 0 {
		context.write("[]")
		return nil
	}

	context.write("[")
	context.writeNewline()

	return context.withIncreasedIndent(func() error {
		for elementIndex, elementNode := range arrayExpression.Elements {
			context.writeIndent()

			if elementNode == nil {
				context.write("null")
			} else {
				if err := context.translateNode(elementNode); err != nil {
					return err
				}
			}

			if elementIndex != len(arrayExpression.Elements)-1 {
				context.write(",")
			}

			context.writeNewline()
		}

		context.writeIndent()
		context.write("]")
		return nil
	})
}

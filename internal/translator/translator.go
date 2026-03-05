package translator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/KirillRg/cli-tool/internal/ast"
)

// TranslateProgram converts an ESTree Program into executable JavaScript code.

func TranslateProgram(p ast.Program) (string, error) {
	t := &translator{}
	if err := t.translateNode(p); err != nil {
		return "", err
	}

	out := t.b.String()
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	return out, nil
}

type translator struct {
	b      strings.Builder
	indent int
}

func (t *translator) translateNode(n ast.Node) error {
	switch v := n.(type) {
	case ast.Program:
		return t.translateProgram(v)
	case *ast.ImportDeclaration:
		return t.translateImportDeclaration(v)
	case *ast.ExportNamedDeclaration:
		return t.translateExportNamedDeclaration(v)
	case *ast.ExportDefaultDeclaration:
		return t.translateExportDefaultDeclaration(v)
	case *ast.VariableDeclaration:
		return t.translateVariableDeclaration(v)
	case *ast.VariableDeclarator:
		return t.translateVariableDeclarator(v)
	case *ast.FunctionExpression:
		return t.translateFunctionExpression(v)
	case *ast.BlockStatement:
		return t.translateBlockStatement(v)
	case *ast.ExpressionStatement:
		return t.translateExpressionStatement(v)
	case *ast.CallExpression:
		return t.translateCallExpression(v)
	case *ast.MemberExpression:
		return t.translateMemberExpression(v)
	case *ast.ObjectExpression:
		return t.translateObjectExpression(v)
	case *ast.Property:
		return t.translateProperty(v)
	case *ast.Identifier:
		t.write(v.Name)
		return nil
	case *ast.Literal:
		return t.translateLiteral(v)
	default:
		return fmt.Errorf("translator: unsupported node type %T", n)
	}
}

func (t *translator) translateProgram(p ast.Program) error {
	for i, n := range p.Body {
		if err := t.translateNode(n); err != nil {
			return err
		}
		if i != len(p.Body)-1 {
			t.write("\n")
		}
	}
	return nil
}

func (t *translator) translateImportDeclaration(n *ast.ImportDeclaration) error {
	if n.Source == nil {
		return errors.New("translator: ImportDeclaration.source is nil")
	}
	if len(n.Specifiers) != 1 {
		return fmt.Errorf("translator: only 1 import specifier supported, got %d", len(n.Specifiers))
	}
	sp, ok := n.Specifiers[0].(*ast.ImportDefaultSpecifier)
	if !ok || sp.Local == nil {
		return errors.New("translator: only ImportDefaultSpecifier supported")
	}

	t.write("import ")
	t.write(sp.Local.Name)
	t.write(" from ")
	if err := t.translateLiteral(n.Source); err != nil {
		return err
	}
	t.write(";\n")
	return nil
}

func (t *translator) translateExportNamedDeclaration(n *ast.ExportNamedDeclaration) error {
	t.write("export ")
	if n.Declaration == nil {
		return errors.New("translator: ExportNamedDeclaration.declaration is nil")
	}
	return t.translateNode(n.Declaration.(ast.Node))
}

func (t *translator) translateExportDefaultDeclaration(n *ast.ExportDefaultDeclaration) error {
	t.write("export default ")
	if n.Declaration == nil {
		return errors.New("translator: ExportDefaultDeclaration.declaration is nil")
	}
	if err := t.translateNode(n.Declaration.(ast.Node)); err != nil {
		return err
	}
	t.write("\n")
	return nil
}

func (t *translator) translateVariableDeclaration(n *ast.VariableDeclaration) error {
	t.write(n.Kind)
	t.write(" ")
	for i, d := range n.Declarations {
		if d == nil {
			return errors.New("translator: VariableDeclarator is nil")
		}
		if err := t.translateVariableDeclarator(d); err != nil {
			return err
		}
		if i != len(n.Declarations)-1 {
			t.write(", ")
		}
	}
	t.write(";\n")
	return nil
}

func (t *translator) translateVariableDeclarator(n *ast.VariableDeclarator) error {
	id, ok := n.ID.(*ast.Identifier)
	if !ok {
		return fmt.Errorf("translator: VariableDeclarator.id must be Identifier, got %T", n.ID)
	}
	t.write(id.Name)
	t.write(" = ")
	if n.Init == nil {
		t.write("null")
		return nil
	}
	return t.translateNode(n.Init.(ast.Node))
}

func (t *translator) translateFunctionExpression(n *ast.FunctionExpression) error {
	t.write("function")
	if n.ID != nil {
		t.write(" ")
		t.write(n.ID.Name)
	}
	t.write(" (")
	for i, p := range n.Params {
		id, ok := p.(*ast.Identifier)
		if !ok {
			return fmt.Errorf("translator: function param must be Identifier, got %T", p)
		}
		if i > 0 {
			t.write(", ")
		}
		t.write(id.Name)
	}
	t.write(") ")
	if n.Body == nil {
		return errors.New("translator: FunctionExpression.body is nil")
	}
	return t.translateBlockStatement(n.Body)
}

func (t *translator) translateBlockStatement(n *ast.BlockStatement) error {
	t.write("{\n")
	t.indent++
	for _, st := range n.Body {
		if st == nil {
			return errors.New("translator: statement is nil")
		}
		t.writeIndent()
		if err := t.translateNode(st.(ast.Node)); err != nil {
			return err
		}
	}
	t.indent--
	t.write("}\n")
	return nil
}

func (t *translator) translateExpressionStatement(n *ast.ExpressionStatement) error {
	if n.Expression == nil {
		return errors.New("translator: ExpressionStatement.expression is nil")
	}
	if err := t.translateNode(n.Expression.(ast.Node)); err != nil {
		return err
	}
	t.write(";\n")
	return nil
}

func (t *translator) translateCallExpression(n *ast.CallExpression) error {
	if n.Callee == nil {
		return errors.New("translator: CallExpression.callee is nil")
	}
	if err := t.translateNode(n.Callee.(ast.Node)); err != nil {
		return err
	}
	t.write("(")
	for i, a := range n.Arguments {
		if i > 0 {
			t.write(", ")
		}
		if a == nil {
			t.write("null")
			continue
		}
		if err := t.translateNode(a.(ast.Node)); err != nil {
			return err
		}
	}
	t.write(")")
	return nil
}

func (t *translator) translateMemberExpression(n *ast.MemberExpression) error {
	if n.Object == nil || n.Property == nil {
		return errors.New("translator: MemberExpression.object/property is nil")
	}
	if err := t.translateNode(n.Object.(ast.Node)); err != nil {
		return err
	}
	if n.Computed {
		t.write("[")
		if err := t.translateNode(n.Property.(ast.Node)); err != nil {
			return err
		}
		t.write("]")
		return nil
	}
	if id, ok := n.Property.(*ast.Identifier); ok {
		t.write(".")
		t.write(id.Name)
		return nil
	}
	t.write("[")
	if err := t.translateNode(n.Property.(ast.Node)); err != nil {
		return err
	}
	t.write("]")
	return nil
}

func (t *translator) translateObjectExpression(n *ast.ObjectExpression) error {
	if len(n.Properties) == 0 {
		t.write("{}")
		return nil
	}

	t.write("{\n")
	t.indent++
	for i, p := range n.Properties {
		if p == nil {
			return errors.New("translator: Property is nil")
		}
		t.writeIndent()
		if err := t.translateProperty(p); err != nil {
			return err
		}
		if i != len(n.Properties)-1 {
			t.write(",")
		}
		t.write("\n")
	}
	t.indent--
	t.writeIndent()
	t.write("}")
	return nil
}

func (t *translator) translateProperty(n *ast.Property) error {
	if n.Key == nil {
		return errors.New("translator: Property.key is nil")
	}
	switch k := n.Key.(type) {
	case *ast.Identifier:
		t.write(k.Name)
	case *ast.Literal:
		if err := t.translateLiteral(k); err != nil {
			return err
		}
	default:
		return fmt.Errorf("translator: unsupported Property.key type %T", n.Key)
	}
	t.write(": ")
	if n.Value == nil {
		t.write("null")
		return nil
	}
	return t.translateNode(n.Value.(ast.Node))
}

func (t *translator) translateLiteral(n *ast.Literal) error {
	if n.Value == nil {
		t.write("null")
		return nil
	}
	switch v := n.Value.(type) {
	case string:
		t.write(strconv.Quote(v))
		return nil
	case bool:
		if v {
			t.write("true")
		} else {
			t.write("false")
		}
		return nil
	case float64, float32, int, int32, int64, uint, uint32, uint64:
		t.write(fmt.Sprintf("%v", v))
		return nil
	default:
		t.write(fmt.Sprintf("%v", v))
		return nil
	}
}

func (t *translator) write(s string) {
	t.b.WriteString(s)
}

func (t *translator) writeIndent() {
	for i := 0; i < t.indent; i++ {
		t.b.WriteString("  ")
	}
}

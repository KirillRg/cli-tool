package ast

// Узлы по стандарту ESTree

type Node interface {
	Node()
}

// Program - корневой узел
// https://github.com/estree/estree/blob/master/es5.md#program

type Program struct {
	Body []Statement
}

func (p Program) Node() {}

// Statement - интерфейс для всех statements

type Statement interface {
	Node
	Statement()
}

// ExpressionStatement - statement содержащий Expression
// https://github.com/estree/estree/blob/master/es5.md#expressionstatement

type ExpressionStatement struct {
	Expression Expression
}

func (e ExpressionStatement) Node()      {}
func (e ExpressionStatement) Statement() {}

// Expression - интерфейс для выражений

type Expression interface {
	Node
	Expression()
}

// CallExpression - вызов метода или функции
// https://github.com/estree/estree/blob/master/es5.md#callexpression

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
}

func (c CallExpression) Node()       {}
func (c CallExpression) Expression() {}

// MemberExpression - доступ к свойству объекта
// https://github.com/estree/estree/blob/master/es5.md#memberexpression

type MemberExpression struct {
	Object   Expression
	Property Identifier
}

func (m MemberExpression) Node()       {}
func (m MemberExpression) Expression() {}

// Identifier - идентификатор
// https://github.com/estree/estree/blob/master/es5.md#identifier

type Identifier struct {
	Name string
}

func (i Identifier) Node()       {}
func (i Identifier) Expression() {}

// ObjectExpression - объект (key-value)
// https://github.com/estree/estree/blob/master/es5.md#objectexpression

type ObjectExpression struct {
	Properties []Property
}

func (o ObjectExpression) Node()       {}
func (o ObjectExpression) Expression() {}

// Property - пара ключ-значение объекта
// https://github.com/estree/estree/blob/master/es5.md#property

type Property struct {
	Key   Identifier
	Value Expression
}

func (p Property) Node() {}

// Literal - литерал (значения строк, чисел и т.д.)
// https://github.com/estree/estree/blob/master/es5.md#literal

type Literal struct {
	Value string
}

func (l Literal) Node()       {}
func (l Literal) Expression() {}

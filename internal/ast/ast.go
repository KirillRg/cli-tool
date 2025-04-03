package ast

// Структуры AST по стандарту ESTree

type Program struct {
	Body []Node
}

type ExpressionStatement struct {
	Expression CallExpression
}

type CallExpression struct {
	Callee    MemberExpression
	Arguments []Node
}

type MemberExpression struct {
	Object   Identifier
	Property Identifier
}

type Identifier struct {
	Name string
}

type ObjectExpression struct {
	Properties []Property
}

type Property struct {
	Key   Identifier
	Value Node
}

type Literal struct {
	Value string
}

type Node interface{}

package ast

// Kind определяет тип узла
type Kind string

const (
	ProgramKind             Kind = "Program"
	ExpressionStatementKind Kind = "ExpressionStatement"
	CallExpressionKind      Kind = "CallExpression"
	LiteralKind             Kind = "Literal"
	IdentifierKind          Kind = "Identifier"
)

type Node interface {
	GetKind() Kind
}

// ======== Интерфейсы ========

type StatementNode interface {
	Node
	statementNode()
}
type ExpressionNode interface {
	Node
	expressionNode()
}

// ======== Узлы ========

type Program struct {
	Body []StatementNode
}

func (p Program) GetKind() Kind { return ProgramKind }

// Линкуем Expression вместе с Statement`ом
type ExpressionStatement struct {
	Expression ExpressionNode
}

func (ExpressionStatement) GetKind() Kind  { return ExpressionStatementKind }
func (ExpressionStatement) statementNode() {}

// Как правило вызов функции
type CallExpression struct {
	Callee    ExpressionNode
	Arguments []ExpressionNode
}

func (CallExpression) GetKind() Kind   { return CallExpressionKind }
func (CallExpression) expressionNode() {}

//Имя переменной
type Identifier struct {
	Name string
}

func (Identifier) GetKind() Kind   { return IdentifierKind }
func (Identifier) expressionNode() {}

//Ее значение
type Literal struct {
	Value string
}

func (Literal) GetKind() Kind   { return LiteralKind }
func (Literal) expressionNode() {}

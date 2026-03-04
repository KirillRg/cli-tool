package ast

// Spec (Node objects): https://github.com/estree/estree/blob/master/es5.md#node-objects
type Node interface {
	node()
}

// --------------------
// Categories (по смыслу ESTree)
// --------------------

// Statement: https://github.com/estree/estree/blob/master/es5.md#statements
type Statement interface {
	Node
	statement()
}

// Expression: https://github.com/estree/estree/blob/master/es5.md#expressions
type Expression interface {
	Node
	expression()
}

// Pattern: https://github.com/estree/estree/blob/master/es5.md#patterns
type Pattern interface {
	Node
	pattern()
}

// Declaration: https://github.com/estree/estree/blob/master/es5.md#declarations
type Declaration interface {
	Statement
	declaration()
}

// ImportOrExportDeclaration: https://github.com/estree/estree/blob/master/es2015.md#modules
type ImportOrExportDeclaration interface {
	Node
	moduleDecl()
}

// --------------------
// Program (root)
// --------------------

// Spec (Program ES5): https://github.com/estree/estree/blob/master/es5.md#program
// Spec (Program extension sourceType/body): https://github.com/estree/estree/blob/master/es2015.md#program
type Program struct {
	Type       string `json:"type"`
	SourceType string `json:"sourceType,omitempty"`
	Body       []Node `json:"body"`
}

func (p Program) node() {}

// --------------------
// Literals / Identifiers
// --------------------

// Spec: https://github.com/estree/estree/blob/master/es5.md#identifier
type Identifier struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func (*Identifier) node()       {}
func (*Identifier) expression() {}
func (*Identifier) pattern()    {}

// Spec: https://github.com/estree/estree/blob/master/es5.md#literal
type Literal struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (*Literal) node()       {}
func (*Literal) expression() {}

// --------------------
// Statements
// --------------------

// Spec: https://github.com/estree/estree/blob/master/es5.md#expressionstatement
type ExpressionStatement struct {
	Type       string     `json:"type"`
	Expression Expression `json:"expression"`
}

func (*ExpressionStatement) node()      {}
func (*ExpressionStatement) statement() {}

// Spec: https://github.com/estree/estree/blob/master/es5.md#blockstatement
type BlockStatement struct {
	Type string      `json:"type"`
	Body []Statement `json:"body"`
}

func (*BlockStatement) node()      {}
func (*BlockStatement) statement() {}

// --------------------
// Expressions
// --------------------

// Spec: https://github.com/estree/estree/blob/master/es5.md#callexpression
type CallExpression struct {
	Type      string       `json:"type"`
	Callee    Expression   `json:"callee"`
	Arguments []Expression `json:"arguments"`
}

func (*CallExpression) node()       {}
func (*CallExpression) expression() {}

// Spec: https://github.com/estree/estree/blob/master/es5.md#memberexpression
type MemberExpression struct {
	Type     string     `json:"type"`
	Object   Expression `json:"object"`
	Property Expression `json:"property"`
	Computed bool       `json:"computed"`
}

func (*MemberExpression) node()       {}
func (*MemberExpression) expression() {}

// Spec: https://github.com/estree/estree/blob/master/es5.md#objectexpression
type ObjectExpression struct {
	Type       string      `json:"type"`
	Properties []*Property `json:"properties"`
}

func (*ObjectExpression) node()       {}
func (*ObjectExpression) expression() {}

// Spec (Property ES5): https://github.com/estree/estree/blob/master/es5.md#property
// Spec (Property extensions method/shorthand/computed): https://github.com/estree/estree/blob/master/es2015.md#property
type Property struct {
	Type      string     `json:"type"`
	Key       Expression `json:"key"`
	Value     Expression `json:"value"`
	Kind      string     `json:"kind"`
	Method    bool       `json:"method"`
	Shorthand bool       `json:"shorthand"`
	Computed  bool       `json:"computed"`
}

func (*Property) node() {}

// --------------------
// Functions
// --------------------

// Spec: https://github.com/estree/estree/blob/master/es5.md#functionexpression
type FunctionExpression struct {
	Type   string          `json:"type"`
	ID     *Identifier     `json:"id"`
	Params []Pattern       `json:"params"`
	Body   *BlockStatement `json:"body"`
}

func (*FunctionExpression) node()       {}
func (*FunctionExpression) expression() {}

// --------------------
// Declarations
// --------------------

// Spec (VariableDeclaration ES5): https://github.com/estree/estree/blob/master/es5.md#variabledeclaration
type VariableDeclaration struct {
	Type         string                `json:"type"`
	Declarations []*VariableDeclarator `json:"declarations"`
	Kind         string                `json:"kind"`
}

func (*VariableDeclaration) node()        {}
func (*VariableDeclaration) statement()   {}
func (*VariableDeclaration) declaration() {}

// Spec: https://github.com/estree/estree/blob/master/es5.md#variabledeclarator
type VariableDeclarator struct {
	Type string     `json:"type"`
	ID   Pattern    `json:"id"`
	Init Expression `json:"init"`
}

func (*VariableDeclarator) node() {}

// --------------------
// Modules: import / export
// --------------------

// Spec: https://github.com/estree/estree/blob/master/es2015.md#importdeclaration
type ImportDeclaration struct {
	Type       string            `json:"type"` // "ImportDeclaration"
	Specifiers []ImportSpecifier `json:"specifiers"`
	Source     *Literal          `json:"source"`
}

func (*ImportDeclaration) node()       {}
func (*ImportDeclaration) moduleDecl() {}

// Specifiers interfaces: https://github.com/estree/estree/blob/master/es2015.md#imports
type ImportSpecifier interface {
	Node
	importSpec()
}

// Spec: https://github.com/estree/estree/blob/master/es2015.md#importdefaultspecifier
type ImportDefaultSpecifier struct {
	Type  string      `json:"type"`
	Local *Identifier `json:"local"`
}

func (*ImportDefaultSpecifier) node()       {}
func (*ImportDefaultSpecifier) importSpec() {}

// Spec: https://github.com/estree/estree/blob/master/es2015.md#exportnameddeclaration
type ExportNamedDeclaration struct {
	Type        string      `json:"type"`
	Declaration Declaration `json:"declaration"`
}

func (*ExportNamedDeclaration) node()       {}
func (*ExportNamedDeclaration) moduleDecl() {}

// Spec: https://github.com/estree/estree/blob/master/es2015.md#exportdefaultdeclaration
type ExportDefaultDeclaration struct {
	Type        string `json:"type"`
	Declaration Node   `json:"declaration"`
}

func (*ExportDefaultDeclaration) node()       {}
func (*ExportDefaultDeclaration) moduleDecl() {}

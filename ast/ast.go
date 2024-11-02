package ast

type Node interface {
	node()
}

type Statement interface {
	Node
	statementNode()
}
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) node() {}

type LetStatement struct {
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) node()          {}
func (ls *LetStatement) statementNode() {}

type Identifier struct {
	Value string
}

func (i *Identifier) node()           {}
func (i *Identifier) expressionNode() {}

type ReturnStatement struct {
	Value Expression
}

func (rs *ReturnStatement) node()          {}
func (rs *ReturnStatement) statementNode() {}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) node()          {}
func (es *ExpressionStatement) statementNode() {}

type IntegerLiteral struct {
	Value int
}

func (il *IntegerLiteral) node()           {}
func (il *IntegerLiteral) expressionNode() {}

type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) node()           {}
func (pe *PrefixExpression) expressionNode() {}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) node()           {}
func (ie *InfixExpression) expressionNode() {}

type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) node()           {}
func (bl *BooleanLiteral) expressionNode() {}

type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) node()           {}
func (ie *IfExpression) expressionNode() {}

type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) node()          {}
func (bs *BlockStatement) statementNode() {}

type FunctionLiteral struct {
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) node()           {}
func (fl *FunctionLiteral) expressionNode() {}

type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) node()           {}
func (ce *CallExpression) expressionNode() {}

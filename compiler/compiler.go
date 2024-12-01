package compiler

import (
	"log"

	"github.com/vincentlabelle/monkey/ast"
	"github.com/vincentlabelle/monkey/code"
	"github.com/vincentlabelle/monkey/object"
	"github.com/vincentlabelle/monkey/symbol"
)

type Compiler struct {
	code        *Bytecode
	symbolTable *symbol.SymbolTable
}

func New() *Compiler {
	return &Compiler{
		code: &Bytecode{
			Instructions: code.Instructions{},
			Constants:    []object.Object{},
		},
		symbolTable: symbol.NewTable(),
	}
}

func (c *Compiler) Compile(program *ast.Program) *Bytecode {
	c.compileStatements(program.Statements)
	return c.code
}

func (c *Compiler) compileStatements(statements []ast.Statement) int {
	var pos int
	for _, statement := range statements {
		switch s := statement.(type) {
		case *ast.ExpressionStatement:
			pos = c.compileExpressionStatement(s)
		case *ast.LetStatement:
			pos = c.compileLetStatement(s)
		default:
			message := "cannot compile; encountered unexpected statement type"
			log.Fatal(message)
		}
	}
	return pos
}

func (c *Compiler) compileExpressionStatement(
	statement *ast.ExpressionStatement,
) int {
	c.compileExpression(statement.Expression)
	return c.emit(code.OpPop)
}

func (c *Compiler) compileExpression(expression ast.Expression) {
	switch e := expression.(type) {
	case *ast.IntegerLiteral:
		c.compileIntegerLiteral(e)
	case *ast.BooleanLiteral:
		c.compileBooleanLiteral(e)
	case *ast.InfixExpression:
		c.compileInfixExpression(e)
	case *ast.PrefixExpression:
		c.compilePrefixExpression(e)
	case *ast.IfExpression:
		c.compileIfExpression(e)
	case *ast.Identifier:
		c.compileIdentifier(e)
	default:
		message := "cannot compile; encountered unexpected expression type"
		log.Fatal(message)
	}
}

func (c *Compiler) compileIntegerLiteral(expression *ast.IntegerLiteral) {
	obj := object.NativeToInteger(expression.Value)
	pos := c.addConstant(obj)
	c.emit(code.OpConstant, pos)
}

func (c *Compiler) addConstant(obj object.Object) int {
	pos := len(c.code.Constants)
	c.code.Constants = append(c.code.Constants, obj)
	return pos
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	instruction := code.Make(op, operands...)
	pos := c.addInstruction(instruction)
	return pos
}

func (c *Compiler) addInstruction(instruction []byte) int {
	pos := len(c.code.Instructions)
	c.code.Instructions = append(c.code.Instructions, instruction...)
	return pos
}

func (c *Compiler) compileBooleanLiteral(expression *ast.BooleanLiteral) {
	if expression.Value {
		c.emit(code.OpTrue)
	} else {
		c.emit(code.OpFalse)
	}
}

func (c *Compiler) compileInfixExpression(expression *ast.InfixExpression) {
	c.compileExpression(expression.Left)
	c.compileExpression(expression.Right)
	c.compileInfixOperator(expression.Operator)
}

func (c *Compiler) compileInfixOperator(operator string) {
	op, ok := code.InfixOperator[operator]
	if !ok {
		message := "cannot compile; encountered unexpected infix operator"
		log.Fatal(message)
	}
	c.emit(op)
}

func (c *Compiler) compilePrefixExpression(expression *ast.PrefixExpression) {
	c.compileExpression(expression.Right)
	c.compilePrefixOperator(expression.Operator)
}

func (c *Compiler) compilePrefixOperator(operator string) {
	op, ok := code.PrefixOperator[operator]
	if !ok {
		message := "cannot compile; encountered unexpected prefix operator"
		log.Fatal(message)
	}
	c.emit(op)
}

func (c *Compiler) compileIfExpression(expression *ast.IfExpression) {
	jumpIfPos := c.compileIfCondition(expression.Condition)
	c.compileBlockStatement(expression.Consequence)
	jumpPos := c.compileIfAlternative(expression, jumpIfPos)
	c.changeJumpOperand(jumpPos)
}

func (c *Compiler) compileIfCondition(expression ast.Expression) int {
	c.compileExpression(expression)
	return c.emit(code.OpJumpIf, 9999) // 9999 to replace
}

func (c *Compiler) compileBlockStatement(statement *ast.BlockStatement) {
	pos := c.compileStatements(statement.Statements)
	if c.isOpPop(pos) {
		c.truncateInstructions(pos)
	}
}

func (c *Compiler) isOpPop(pos int) bool {
	return c.atOpcode(pos) == code.OpPop
}

func (c *Compiler) atOpcode(pos int) code.Opcode {
	return code.Opcode(c.code.Instructions[pos])
}

func (c *Compiler) truncateInstructions(pos int) {
	c.code.Instructions = c.code.Instructions[:pos]
}

func (c *Compiler) compileIfAlternative(
	expression *ast.IfExpression,
	jumpIfPos int,
) int {
	jumpPos := c.emit(code.OpJump, 9999) // 9999 to replace
	c.changeJumpOperand(jumpIfPos)
	c.innerCompileIfAlternative(expression.Alternative)
	return jumpPos
}

func (c *Compiler) innerCompileIfAlternative(statement *ast.BlockStatement) {
	if statement != nil {
		c.compileBlockStatement(statement)
	} else {
		c.emit(code.OpNull)
	}
}

func (c *Compiler) changeJumpOperand(pos int) {
	op := c.atOpcode(pos)
	operand := len(c.code.Instructions)
	instruction := code.Make(op, operand)
	c.replaceInstruction(pos, instruction)
}

func (c *Compiler) replaceInstruction(pos int, instruction []byte) {
	copy(c.code.Instructions[pos:], instruction)
}

func (c *Compiler) compileIdentifier(expression *ast.Identifier) {
	sym := c.resolveSymbol(expression)
	c.emit(code.OpGetGlobal, sym.Index)
}

func (c *Compiler) resolveSymbol(expression *ast.Identifier) symbol.Symbol {
	return c.symbolTable.Resolve(expression.Value)
}

func (c *Compiler) compileLetStatement(statement *ast.LetStatement) int {
	c.compileExpression(statement.Value)
	index := c.defineSymbol(statement.Name)
	return c.emit(code.OpSetGlobal, index)
}

func (c *Compiler) defineSymbol(expression *ast.Identifier) int {
	sym := c.symbolTable.Define(expression.Value)
	return sym.Index
}

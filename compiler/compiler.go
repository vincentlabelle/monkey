package compiler

import (
	"log"

	"github.com/vincentlabelle/monkey/ast"
	"github.com/vincentlabelle/monkey/code"
	"github.com/vincentlabelle/monkey/object"
)

type Compiler struct {
	code *Bytecode
}

func New() *Compiler {
	return &Compiler{
		code: &Bytecode{
			Instructions: code.Instructions{},
			Constants:    []object.Object{},
		},
	}
}

func (c *Compiler) Compile(program *ast.Program) {
	for _, statement := range program.Statements {
		switch s := statement.(type) {
		case *ast.ExpressionStatement:
			c.compileExpressionStatement(s)
		default:
			message := "cannot compile; encountered unexpected statement type"
			log.Fatal(message)
		}
	}
}

func (c *Compiler) compileExpressionStatement(
	statement *ast.ExpressionStatement,
) {
	c.compileExpression(statement.Expression)
	c.emit(code.OpPop)
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
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	pos := len(c.code.Instructions)
	c.code.Instructions = append(c.code.Instructions, ins...)
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
		message := "cannot compile; encountered unexpected operator"
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
		message := "cannot compile; encountered unexpected operator"
		log.Fatal(message)
	}
	c.emit(op)
}

func (c *Compiler) Bytecode() *Bytecode {
	return c.code
}

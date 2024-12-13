package compiler

import (
	"log"
	"maps"

	"github.com/vincentlabelle/monkey/ast"
	"github.com/vincentlabelle/monkey/code"
	"github.com/vincentlabelle/monkey/object"
	"github.com/vincentlabelle/monkey/symbol"
)

type Compiler struct {
	scopes      []code.Instructions
	scopeIndex  int
	constants   []object.Object
	symbolTable *symbol.SymbolTable
}

func New() *Compiler {
	return &Compiler{
		scopes:      []code.Instructions{{}},
		constants:   []object.Object{},
		symbolTable: symbol.NewTable(),
	}
}

func (c *Compiler) Compile(program *ast.Program) *Bytecode {
	c.compileStatements(program.Statements)
	instructions := c.currentInstructions()
	return &Bytecode{Instructions: instructions, Constants: c.constants}
}

func (c *Compiler) enterScope() {
	c.innerEnterScope()
	c.enterSymbolTable()
}

func (c *Compiler) innerEnterScope() {
	c.scopes = append(c.scopes, code.Instructions{})
	c.scopeIndex++
}

func (c *Compiler) enterSymbolTable() {
	c.symbolTable = symbol.NewInnerTable(c.symbolTable)
}

func (c *Compiler) compileStatements(statements []ast.Statement) int {
	var pos int
	for _, statement := range statements {
		switch s := statement.(type) {
		case *ast.ExpressionStatement:
			pos = c.compileExpressionStatement(s)
		case *ast.LetStatement:
			pos = c.compileLetStatement(s)
		case *ast.ReturnStatement:
			pos = c.compileReturnStatement(s)
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
	case *ast.StringLiteral:
		c.compileStringLiteral(e)
	case *ast.InfixExpression:
		c.compileInfixExpression(e)
	case *ast.PrefixExpression:
		c.compilePrefixExpression(e)
	case *ast.IfExpression:
		c.compileIfExpression(e)
	case *ast.Identifier:
		c.compileIdentifier(e)
	case *ast.ArrayLiteral:
		c.compileArrayLiteral(e)
	case *ast.HashLiteral:
		c.compileHashLiteral(e)
	case *ast.IndexExpression:
		c.compileIndexExpression(e)
	case *ast.FunctionLiteral:
		c.compileFunctionLiteral(e)
	case *ast.CallExpression:
		c.compileCallExpression(e)
	default:
		message := "cannot compile; encountered unexpected expression type"
		log.Fatal(message)
	}
}

func (c *Compiler) compileIntegerLiteral(expression *ast.IntegerLiteral) {
	obj := object.NativeToInteger(expression.Value)
	c.compileConstant(obj)
}

func (c *Compiler) compileConstant(obj object.Object) {
	pos := c.addConstant(obj)
	c.emit(code.OpConstant, pos)
}

func (c *Compiler) addConstant(obj object.Object) int {
	pos := len(c.constants)
	c.constants = append(c.constants, obj)
	return pos
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	instruction := code.Make(op, operands...)
	pos := c.addInstruction(instruction)
	return pos
}

func (c *Compiler) addInstruction(instruction []byte) int {
	instructions := c.currentInstructions()
	pos := len(instructions)
	instructions = append(instructions, instruction...)
	c.updateCurrentInstructions(instructions)
	return pos
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex]
}

func (c *Compiler) updateCurrentInstructions(instructions code.Instructions) {
	c.scopes[c.scopeIndex] = instructions
}

func (c *Compiler) compileBooleanLiteral(expression *ast.BooleanLiteral) {
	if expression.Value {
		c.emit(code.OpTrue)
	} else {
		c.emit(code.OpFalse)
	}
}

func (c *Compiler) compileStringLiteral(expression *ast.StringLiteral) {
	obj := object.NativeToString(expression.Value)
	c.compileConstant(obj)
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

func (c *Compiler) compileBlockStatement(statement *ast.BlockStatement) bool {
	pos := c.compileStatements(statement.Statements)
	if c.isOpPop(pos) {
		c.truncateInstructions(pos)
		return true
	}
	return false
}

func (c *Compiler) isOpPop(pos int) bool {
	return c.atOpcode(pos) == code.OpPop
}

func (c *Compiler) atOpcode(pos int) code.Opcode {
	instructions := c.currentInstructions()
	return code.Opcode(instructions[pos])
}

func (c *Compiler) truncateInstructions(pos int) {
	instructions := c.currentInstructions()
	c.updateCurrentInstructions(instructions[:pos])
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
	operand := len(c.currentInstructions())
	instruction := code.Make(op, operand)
	c.replaceInstruction(pos, instruction)
}

func (c *Compiler) replaceInstruction(pos int, instruction []byte) {
	instructions := c.currentInstructions()
	copy(instructions[pos:], instruction)
}

func (c *Compiler) compileIdentifier(expression *ast.Identifier) {
	sym := c.resolveSymbol(expression)
	c.loadSymbol(sym)
}

func (c *Compiler) resolveSymbol(expression *ast.Identifier) symbol.Symbol {
	sym, ok := c.symbolTable.Resolve(expression.Value)
	if !ok {
		message := "cannot compile; encountered symbol which cannot be resolved"
		log.Fatal(message)
	}
	return sym
}

func (c *Compiler) loadSymbol(sym symbol.Symbol) {
	switch sym.Scope {
	case symbol.BuiltinScope:
		c.emit(code.OpGetBuiltin, sym.Index)
	case symbol.GlobalScope:
		c.emit(code.OpGetGlobal, sym.Index)
	case symbol.FreeScope:
		c.emit(code.OpGetFree, sym.Index)
	case symbol.FunctionScope:
		c.emit(code.OpCurrentClosure)
	default:
		c.emit(code.OpGetLocal, sym.Index)
	}
}

func (c *Compiler) compileArrayLiteral(expression *ast.ArrayLiteral) {
	c.compileExpressions(expression.Elements)
	c.emit(code.OpArray, len(expression.Elements))
}

func (c *Compiler) compileExpressions(expressions []ast.Expression) {
	for _, expression := range expressions {
		c.compileExpression(expression)
	}
}

func (c *Compiler) compileHashLiteral(expression *ast.HashLiteral) {
	c.compileHashLiteralPairs(expression.Pairs)
	c.emit(code.OpHash, len(expression.Pairs))
}

func (c *Compiler) compileHashLiteralPairs(
	pairs map[ast.HashKey]ast.Expression,
) {
	keys := ast.SortHashKeys(maps.Keys(pairs))
	for _, key := range keys {
		c.compileExpression(key.Expression)
		c.compileExpression(pairs[key])
	}
}

func (c *Compiler) compileIndexExpression(expression *ast.IndexExpression) {
	c.compileExpression(expression.Left)
	c.compileExpression(expression.Index)
	c.emit(code.OpIndex)
}

func (c *Compiler) compileFunctionLiteral(expression *ast.FunctionLiteral) {
	instructions, count, free := c.innerCompileFunctionLiteral(expression)
	obj := &object.CompiledFunction{
		Instructions:  instructions,
		NumLocals:     count,
		NumParameters: len(expression.Parameters),
	}
	c.compileClosure(obj, free)
}

func (c *Compiler) innerCompileFunctionLiteral(
	expression *ast.FunctionLiteral,
) (code.Instructions, int, []symbol.Symbol) {
	c.enterScope()
	c.defineFunctionName(expression)
	c.compileFunctionParameters(expression.Parameters)
	c.compileFunctionBody(expression.Body)
	return c.leaveScope()
}

func (c *Compiler) defineFunctionName(expression *ast.FunctionLiteral) {
	if expression.Name != "" {
		c.symbolTable.DefineFunctionName(expression.Name)
	}
}

func (c *Compiler) compileFunctionParameters(expressions []*ast.Identifier) {
	for _, expression := range expressions {
		c.defineSymbol(expression)
	}
}

func (c *Compiler) compileFunctionBody(
	statement *ast.BlockStatement,
) {
	if len(statement.Statements) == 0 {
		c.emit(code.OpReturn)
	} else {
		c.compileNonEmptyFunctionBody(statement)
	}
}

func (c *Compiler) compileNonEmptyFunctionBody(statement *ast.BlockStatement) {
	truncated := c.compileBlockStatement(statement)
	if truncated {
		c.emit(code.OpReturnValue)
	}
}

func (c *Compiler) leaveScope() (code.Instructions, int, []symbol.Symbol) {
	instructions := c.innerLeaveScope()
	count, free := c.leaveSymbolTable()
	return instructions, count, free
}

func (c *Compiler) innerLeaveScope() code.Instructions {
	instructions := c.currentInstructions()
	c.scopes = c.scopes[:c.scopeIndex]
	c.scopeIndex--
	return instructions
}

func (c *Compiler) leaveSymbolTable() (int, []symbol.Symbol) {
	count := c.symbolTable.CountDefinitions()
	free := c.symbolTable.Free()
	c.symbolTable = c.symbolTable.Outer()
	return count, free
}

func (c *Compiler) compileClosure(obj object.Object, free []symbol.Symbol) {
	c.loadSymbols(free)
	pos := c.addConstant(obj)
	c.emit(code.OpClosure, pos, len(free))
}

func (c *Compiler) loadSymbols(symbols []symbol.Symbol) {
	for _, sym := range symbols {
		c.loadSymbol(sym)
	}
}

func (c *Compiler) compileCallExpression(expression *ast.CallExpression) {
	c.compileExpression(expression.Function)
	c.compileExpressions(expression.Arguments)
	c.emit(code.OpCall, len(expression.Arguments))
}

func (c *Compiler) compileLetStatement(statement *ast.LetStatement) int {
	sym := c.defineSymbol(statement.Name)
	c.compileExpression(statement.Value)
	op := c.getOpSet(sym)
	return c.emit(op, sym.Index)
}

func (c *Compiler) defineSymbol(expression *ast.Identifier) symbol.Symbol {
	return c.symbolTable.Define(expression.Value)
}

func (c *Compiler) getOpSet(sym symbol.Symbol) code.Opcode {
	if sym.Scope == symbol.GlobalScope {
		return code.OpSetGlobal
	}
	return code.OpSetLocal
}

func (c *Compiler) compileReturnStatement(statement *ast.ReturnStatement) int {
	c.compileExpression(statement.Value)
	return c.emit(code.OpReturnValue)
}

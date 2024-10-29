package parser

import (
	"log"
	"strconv"

	"github.com/vincentlabelle/monkey/ast"
	"github.com/vincentlabelle/monkey/lexer"
	"github.com/vincentlabelle/monkey/token"
)

type Parser struct {
	lex       *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{lex: lex}
	p.forward()
	p.forward()
	return p
}

func (p *Parser) forward() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	statements := []ast.Statement{}
	for !p.isCurToken(token.EOF) {
		statement := p.parseStatement()
		statements = append(statements, statement)
		p.forward()
	}
	return &ast.Program{Statements: statements}
}

func (p *Parser) isCurToken(type_ token.TokenType) bool {
	return p.curToken.Type == type_
}

func (p *Parser) parseStatement() ast.Statement {
	if p.isCurToken(token.LET) {
		return p.parseLetStatement()
	}
	if p.isCurToken(token.RETURN) {
		return p.parseReturnStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	p.forward()
	if !p.isCurToken(token.IDENT) {
		message := "cannot parse program; let must be followed by an identifier"
		log.Fatal(message)
	}
	name := p.parseIdentifier()
	p.forward()
	if !p.isCurToken(token.ASSIGN) {
		message := `cannot parse program; 
			identifier in let statement must be followed by assignement`
		log.Fatal(message)
	}
	p.forward()
	value := p.parseExpression(LOWEST)
	if p.isPeekToken(token.SEMICOLON) {
		p.forward()
	}
	return &ast.LetStatement{Name: name, Value: value}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Value: p.curToken.Literal}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	expression := p.parsePrefixExpression()
	for !p.isPeekToken(token.SEMICOLON) && precedence < p.peekPrecedence() {
		p.forward()
		expression = p.parseInfixExpression(expression)
	}
	return expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	var expression ast.Expression
	if p.isCurToken(token.IDENT) {
		expression = p.parseIdentifier()
	} else if p.isCurToken(token.INT) {
		expression = p.parseIntegerLiteral()
	} else if p.isCurToken(token.TRUE) || p.isCurToken(token.FALSE) {
		expression = p.parseBooleanLiteral()
	} else if p.isCurToken(token.BANG) || p.isCurToken(token.MINUS) {
		expression = p.parsePrefix()
	} else if p.isCurToken(token.LPAREN) {
		expression = p.parseGroupedExpression()
	} else if p.isCurToken(token.IF) {
		expression = p.parseIfExpression()
	} else if p.isCurToken(token.FUNCTION) {
		expression = p.parseFunctionLiteral()
	} else {
		message := "cannot parse program; cannot parse prefix expression for %v"
		log.Fatalf(message, p.curToken.Type)
	}
	return expression
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	value, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		message := "cannot parse program; unable to convert ASCII to integer"
		log.Fatal(message)
	}
	return &ast.IntegerLiteral{Value: value}
}

func (p *Parser) parseBooleanLiteral() *ast.BooleanLiteral {
	return &ast.BooleanLiteral{Value: p.isCurToken(token.TRUE)}
}

func (p *Parser) parsePrefix() *ast.PrefixExpression {
	operator := p.curToken.Literal
	p.forward()
	right := p.parseExpression(PREFIX)
	return &ast.PrefixExpression{Operator: operator, Right: right}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.forward()
	expression := p.parseExpression(LOWEST)
	p.forward()
	if !p.isCurToken(token.RPAREN) {
		message := "cannot parse program; missing ) to close grouped expression"
		log.Fatal(message)
	}
	return expression
}

func (p *Parser) parseIfExpression() *ast.IfExpression {
	condition, consequence := p.parseIf()
	alternative := p.parseElse()
	return &ast.IfExpression{
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
	}
}

func (p *Parser) parseIf() (ast.Expression, *ast.BlockStatement) {
	p.forward()
	if !p.isCurToken(token.LPAREN) {
		message := "cannot parse program; missing ( after if"
		log.Fatal(message)
	}
	condition := p.parseExpression(LOWEST)
	p.forward()
	if !p.isCurToken(token.LBRACE) {
		message := "cannot parse program; missing { after if"
		log.Fatal(message)
	}
	consequence := p.parseBlockStatement()
	return condition, consequence
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	p.forward()
	statements := []ast.Statement{}
	for !p.isCurToken(token.RBRACE) && !p.isCurToken(token.EOF) {
		statement := p.parseStatement()
		statements = append(statements, statement)
		p.forward()
	}
	return &ast.BlockStatement{Statements: statements}
}

func (p *Parser) parseElse() *ast.BlockStatement {
	var block *ast.BlockStatement
	if p.isPeekToken(token.ELSE) {
		p.forward()
		p.forward()
		if !p.isCurToken(token.LBRACE) {
			message := "cannot parse program; missing { after else"
			log.Fatal(message)
		}
		block = p.parseBlockStatement()
	}
	return block
}

func (p *Parser) parseFunctionLiteral() *ast.FunctionLiteral {
	p.forward()
	if !p.isCurToken(token.LPAREN) {
		message := "cannot parse program; missing ( after fn"
		log.Fatal(message)
	}
	parameters := p.parseFunctionParameters()
	p.forward()
	if !p.isCurToken(token.LBRACE) {
		message := "cannot parse program; missing { after fn"
		log.Fatal(message)
	}
	body := p.parseBlockStatement()
	return &ast.FunctionLiteral{Parameters: parameters, Body: body}
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	if p.isPeekToken(token.RPAREN) {
		p.forward()
		return []*ast.Identifier{}
	}
	return p.innerParseFunctionParameters()
}

func (p *Parser) isPeekToken(type_ token.TokenType) bool {
	return p.peekToken.Type == type_
}

func (p *Parser) innerParseFunctionParameters() []*ast.Identifier {
	p.forward()
	identifiers := []*ast.Identifier{p.parseFunctionParameter()}
	p.forward()
	for p.isCurToken(token.COMMA) {
		p.forward()
		identifier := p.parseFunctionParameter()
		identifiers = append(identifiers, identifier)
		p.forward()
	}
	if !p.isCurToken(token.RPAREN) {
		message := "cannot parse program; missing ) after fn"
		log.Fatal(message)
	}
	return identifiers
}

func (p *Parser) parseFunctionParameter() *ast.Identifier {
	if !p.isCurToken(token.IDENT) {
		message := `cannot parse program; 
			unexpected token in function parameters`
		log.Fatal(message)
	}
	return p.parseIdentifier()
}

func (p *Parser) peekPrecedence() int {
	return p.precedence(p.peekToken.Type)
}

func (p *Parser) precedence(type_ token.TokenType) int {
	if precedence, ok := precedences[type_]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	var expression ast.Expression
	if p.isCurToken(token.PLUS) ||
		p.isCurToken(token.MINUS) ||
		p.isCurToken(token.ASTERISK) ||
		p.isCurToken(token.SLASH) ||
		p.isCurToken(token.EQ) ||
		p.isCurToken(token.NE) ||
		p.isCurToken(token.LT) ||
		p.isCurToken(token.GT) {
		expression = p.parseInfix(left)
	} else if p.isCurToken(token.LPAREN) {
		return p.parseCallExpression(left)
	} else {
		expression = left
	}
	return expression
}

func (p *Parser) parseInfix(left ast.Expression) *ast.InfixExpression {
	operator := p.curToken.Literal
	precedence := p.curPrecedence()
	p.forward()
	right := p.parseExpression(precedence)
	return &ast.InfixExpression{Left: left, Operator: operator, Right: right}
}

func (p *Parser) curPrecedence() int {
	return p.precedence(p.curToken.Type)
}

func (p *Parser) parseCallExpression(
	left ast.Expression,
) *ast.CallExpression {
	arguments := p.parseCallArguments()
	return &ast.CallExpression{Function: left, Arguments: arguments}
}

func (p *Parser) parseCallArguments() []ast.Expression {
	if p.isPeekToken(token.RPAREN) {
		p.forward()
		return []ast.Expression{}
	}
	return p.innerParseCallArguments()
}

func (p *Parser) innerParseCallArguments() []ast.Expression {
	p.forward()
	arguments := []ast.Expression{p.parseExpression(LOWEST)}
	p.forward()
	for p.isCurToken(token.COMMA) {
		p.forward()
		argument := p.parseExpression(LOWEST)
		arguments = append(arguments, argument)
		p.forward()
	}
	if !p.isCurToken(token.RPAREN) {
		message := "cannot parse program; missing ) after call"
		log.Fatal(message)
	}
	return arguments
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	p.forward()
	value := p.parseExpression(LOWEST)
	if p.isPeekToken(token.SEMICOLON) {
		p.forward()
	}
	return &ast.ReturnStatement{Value: value}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expression := p.parseExpression(LOWEST)
	if p.isPeekToken(token.SEMICOLON) {
		p.forward()
	}
	return &ast.ExpressionStatement{Expression: expression}
}

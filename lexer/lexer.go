package lexer

import "github.com/vincentlabelle/monkey/token"

type Lexer struct {
	input    string
	position int
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhite()
	if l.isEOF() {
		return token.Token{Type: token.EOF, Literal: ""}
	}
	return l.nextToken()
}

func (l *Lexer) skipWhite() {
	for !l.isEOF() && l.isWhite() {
		l.forward(1)
	}
}

func (l *Lexer) isEOF() bool {
	return l.position >= len(l.input)
}

func (l *Lexer) isWhite() bool {
	char := l.getUnaryChar()
	return char == ' ' || char == '\n' || char == '\t' || char == '\r'
}

func (l *Lexer) getUnaryChar() byte {
	return l.input[l.position]
}

func (l *Lexer) forward(times int) {
	l.position += times
}

func (l *Lexer) nextToken() token.Token {
	var tok token.Token
	if tok_, ok := l.getBinary(); ok {
		tok = tok_
		l.forward(2)
	} else if tok_, ok := l.getUnary(); ok {
		tok = tok_
		l.forward(1)
	} else if tok_, ok := l.getInt(); ok {
		tok = tok_
	} else if tok_, ok := l.getLetter(); ok {
		tok = tok_
	} else if tok_, ok := l.getQuoted(); ok {
		tok = tok_
	} else {
		tok = token.Token{Type: token.ILLEGAL, Literal: l.getUnaryString()}
		l.forward(1)
	}
	return tok
}

func (l *Lexer) getBinary() (token.Token, bool) {
	str := l.getBinaryString()
	return l.getBinaryToken(str)
}

func (l *Lexer) getBinaryString() string {
	if l.isNextEOF() {
		return ""
	}
	return l.input[l.position : l.position+2]
}

func (l *Lexer) isNextEOF() bool {
	return l.position+1 >= len(l.input)
}

func (l *Lexer) getBinaryToken(str string) (token.Token, bool) {
	if type_, ok := token.Binary[str]; ok {
		return token.Token{Type: type_, Literal: str}, ok
	}
	return token.Token{}, false
}

func (l *Lexer) getUnary() (token.Token, bool) {
	char := l.getUnaryChar()
	return l.getUnaryToken(char)
}

func (l *Lexer) getUnaryToken(char byte) (token.Token, bool) {
	if type_, ok := token.Unary[char]; ok {
		return token.Token{Type: type_, Literal: string(char)}, ok
	}
	return token.Token{}, false
}

func (l *Lexer) getInt() (token.Token, bool) {
	if !l.isDigit() {
		return token.Token{}, false
	}
	return l.getIntToken(), true
}

func (l *Lexer) isDigit() bool {
	char := l.getUnaryChar()
	return '0' <= char && char <= '9'
}

func (l *Lexer) getIntToken() token.Token {
	literal := l.getIntLiteral()
	return token.Token{Type: token.INT, Literal: literal}
}

func (l *Lexer) getIntLiteral() string {
	literal := ""
	for !l.isEOF() && l.isDigit() {
		literal += l.getUnaryString()
		l.forward(1)
	}
	return literal
}

func (l *Lexer) getUnaryString() string {
	char := l.getUnaryChar()
	return string(char)
}

func (l *Lexer) getLetter() (token.Token, bool) {
	if !l.isLetter() {
		return token.Token{}, false
	}
	return l.getLetterToken(), true
}

func (l *Lexer) isLetter() bool {
	char := l.getUnaryChar()
	return 'a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z' ||
		char == '_'
}

func (l *Lexer) getLetterToken() token.Token {
	literal := l.getLetterLiteral()
	type_ := l.getLetterType(literal)
	return token.Token{Type: type_, Literal: literal}
}

func (l *Lexer) getLetterLiteral() string {
	literal := ""
	for !l.isEOF() && l.isLetter() {
		literal += l.getUnaryString()
		l.forward(1)
	}
	return literal
}

func (l *Lexer) getLetterType(literal string) token.TokenType {
	if type_, ok := token.Keywords[literal]; ok {
		return type_
	}
	return token.IDENT
}

func (l *Lexer) getQuoted() (token.Token, bool) {
	if !l.isQuote() {
		return token.Token{}, false
	}
	return l.getQuotedToken(), true
}

func (l *Lexer) isQuote() bool {
	return l.getUnaryChar() == '"'
}

func (l *Lexer) getQuotedToken() token.Token {
	literal := l.getQuotedLiteral()
	return token.Token{Type: token.STRING, Literal: literal}
}

func (l *Lexer) getQuotedLiteral() string {
	l.forward(1)
	literal := ""
	for !l.isEOF() && !l.isQuote() {
		literal += l.getUnaryString()
		l.forward(1)
	}
	l.forward(1)
	return literal
}

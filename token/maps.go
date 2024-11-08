package token

var Unary = map[byte]TokenType{
	'=': ASSIGN,
	'+': PLUS,
	'-': MINUS,
	'*': ASTERISK,
	'/': SLASH,
	'!': BANG,
	'<': LT,
	'>': GT,
	',': COMMA,
	';': SEMICOLON,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
}

var Binary = map[string]TokenType{
	"==": EQ,
	"!=": NE,
}

var Keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

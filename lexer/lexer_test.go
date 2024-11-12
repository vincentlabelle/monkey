package lexer

import (
	"testing"

	"github.com/vincentlabelle/monkey/token"
)

func Test(t *testing.T) {
	setup := []struct {
		input    string
		expected []token.Token
	}{
		{
			`=  +-*/!<>,;(){}==!=# `,
			[]token.Token{
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.ASTERISK, Literal: "*"},
				{Type: token.SLASH, Literal: "/"},
				{Type: token.BANG, Literal: "!"},
				{Type: token.LT, Literal: "<"},
				{Type: token.GT, Literal: ">"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EQ, Literal: "=="},
				{Type: token.NE, Literal: "!="},
				{Type: token.ILLEGAL, Literal: "#"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			`5+ 50 == 55`,
			[]token.Token{
				{Type: token.INT, Literal: "5"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.INT, Literal: "50"},
				{Type: token.EQ, Literal: "=="},
				{Type: token.INT, Literal: "55"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			`let five = 5;
			let ten = 10;

			let add = fn(x, y) {
				x + y;
			};

			let result = add(five, ten);
			!-/*5;
			5 < 10 > 5;

			if (5 < 10) {
				return true;
			} else {
				return false;
			}

			10 == 10;
			10 != 9;
			"foobar"
			"foo bar"
			[1, 2];
			{"foo": "bar"};`,
			[]token.Token{
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "five"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INT, Literal: "5"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "ten"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INT, Literal: "10"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "add"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.FUNCTION, Literal: "fn"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "result"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.IDENT, Literal: "add"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENT, Literal: "five"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "ten"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.BANG, Literal: "!"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.SLASH, Literal: "/"},
				{Type: token.ASTERISK, Literal: "*"},
				{Type: token.INT, Literal: "5"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.INT, Literal: "5"},
				{Type: token.LT, Literal: "<"},
				{Type: token.INT, Literal: "10"},
				{Type: token.GT, Literal: ">"},
				{Type: token.INT, Literal: "5"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.IF, Literal: "if"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.INT, Literal: "5"},
				{Type: token.LT, Literal: "<"},
				{Type: token.INT, Literal: "10"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.TRUE, Literal: "true"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.ELSE, Literal: "else"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.FALSE, Literal: "false"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.INT, Literal: "10"},
				{Type: token.EQ, Literal: "=="},
				{Type: token.INT, Literal: "10"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.INT, Literal: "10"},
				{Type: token.NE, Literal: "!="},
				{Type: token.INT, Literal: "9"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.STRING, Literal: "foobar"},
				{Type: token.STRING, Literal: "foo bar"},
				{Type: token.LBRACKET, Literal: "["},
				{Type: token.INT, Literal: "1"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.INT, Literal: "2"},
				{Type: token.RBRACKET, Literal: "]"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.STRING, Literal: "foo"},
				{Type: token.COLON, Literal: ":"},
				{Type: token.STRING, Literal: "bar"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for _, s := range setup {
		lex := New(s.input)
		for _, expected := range s.expected {
			actual := lex.NextToken()
			if actual != expected {
				t.Fatalf("expected=%v, actual=%v", expected, actual)
			}
		}
	}
}

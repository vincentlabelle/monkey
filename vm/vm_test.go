package vm

import (
	"testing"

	"github.com/vincentlabelle/monkey/ast"
	"github.com/vincentlabelle/monkey/compiler"
	"github.com/vincentlabelle/monkey/lexer"
	"github.com/vincentlabelle/monkey/object"
	"github.com/vincentlabelle/monkey/parser"
)

func Test(t *testing.T) {
	setup := []struct {
		input    string
		expected object.Object
	}{
		{`1;`, &object.Integer{Value: 1}},
		{`2;`, &object.Integer{Value: 2}},
		{`1 + 2;`, &object.Integer{Value: 3}},
		{`1 - 2;`, &object.Integer{Value: -1}},
		{`1 * 2;`, &object.Integer{Value: 2}},
		{`4 / 2;`, &object.Integer{Value: 2}},
		{`50 / 2 * 2 + 10 -5;`, &object.Integer{Value: 55}},
		{`5 + 5 + 5 + 5 - 10;`, &object.Integer{Value: 10}},
		{`2 * 2 * 2 * 2 * 2;`, &object.Integer{Value: 32}},
		{`5 * 2 + 10;`, &object.Integer{Value: 20}},
		{`5 + 2 * 10;`, &object.Integer{Value: 25}},
		{`5 * (2 + 10);`, &object.Integer{Value: 60}},
		{`true;`, object.TRUE},
		{`false;`, object.FALSE},
		{`1 < 2;`, object.TRUE},
		{`1 > 2;`, object.FALSE},
		{`1 < 1;`, object.FALSE},
		{`1 > 1;`, object.FALSE},
		{`1 == 1;`, object.TRUE},
		{`1 != 1;`, object.FALSE},
		{`1 == 2;`, object.FALSE},
		{`1 != 2;`, object.TRUE},
		{`true == true;`, object.TRUE},
		{`false == false;`, object.TRUE},
		{`true == false;`, object.FALSE},
		{`true != false;`, object.TRUE},
		{`(1 < 2) == true;`, object.TRUE},
		{`(1 < 2) == false;`, object.FALSE},
		{`(1 > 2) == true;`, object.FALSE},
		{`(1 > 2) == false;`, object.TRUE},
		{`-5;`, &object.Integer{Value: -5}},
		{`-10;`, &object.Integer{Value: -10}},
		{`-50 + 100 + -50;`, &object.Integer{Value: 0}},
		{`(5 + 10 * 2 + 15 / 3) * 2 + -10;`, &object.Integer{Value: 50}},
		{`!true;`, object.FALSE},
		{`!false;`, object.TRUE},
		{`!5;`, object.FALSE},
		{`!!true;`, object.TRUE},
		{`!!false;`, object.FALSE},
		{`!!5;`, object.TRUE},
	}

	for _, s := range setup {
		vm := new_(s.input)
		vm.Run()
		actual := vm.LastPopped()
		testObject(t, actual, s.expected)
	}
}

func new_(input string) *VM {
	code := compile(input)
	return New(code)
}

func compile(input string) *compiler.Bytecode {
	program := parse(input)
	c := compiler.New()
	c.Compile(program)
	return c.Bytecode()
}

func parse(input string) *ast.Program {
	lex := lexer.New(input)
	p := parser.New(lex)
	return p.ParseProgram()
}

func testObject(
	t *testing.T,
	actual object.Object,
	expected object.Object,
) {
	switch e := expected.(type) {
	case *object.Integer:
		a, ok := actual.(*object.Integer)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testIntegerObject(t, a, e)
	case *object.Boolean:
		a, ok := actual.(*object.Boolean)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testBooleanObject(t, a, e)
	default:
		t.Fatal("object type unknown")
	}
}

func testIntegerObject(
	t *testing.T,
	actual *object.Integer,
	expected *object.Integer,
) {
	if actual.Value != expected.Value {
		t.Fatalf(
			"integer value mismatch. got=%v, expected=%v",
			actual.Value,
			expected.Value,
		)
	}
}

func testBooleanObject(
	t *testing.T,
	actual *object.Boolean,
	expected *object.Boolean,
) {
	if actual.Value != expected.Value {
		t.Fatalf(
			"boolean value mismatch. got=%v, expected=%v",
			actual.Value,
			expected.Value,
		)
	}
}

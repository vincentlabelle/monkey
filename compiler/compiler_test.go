package compiler

import (
	"testing"

	"github.com/vincentlabelle/monkey/ast"
	"github.com/vincentlabelle/monkey/code"
	"github.com/vincentlabelle/monkey/lexer"
	"github.com/vincentlabelle/monkey/object"
	"github.com/vincentlabelle/monkey/parser"
)

func Test(t *testing.T) {
	setup := []struct {
		input             string
		expectedPieces    []code.Instructions
		expectedConstants []object.Object
	}{
		{
			`1 + 2;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
		},
		{
			`1 - 2;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
		},
		{
			`1 * 2;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
		},
		{
			`2 / 1;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 2},
				&object.Integer{Value: 1},
			},
		},
		{
			`1; 2;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
		},
		{
			`true;`,
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
			[]object.Object{},
		},
		{
			`false;`,
			[]code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
			[]object.Object{},
		},
		{
			`1 > 2;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
		},
		{
			`1 < 2;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpLowerThan),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
		},
		{
			`1 == 2;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
		},
		{
			`1 != 2;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
		},
		{
			`true == false;`,
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
			[]object.Object{},
		},
		{
			`true != false;`,
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
			[]object.Object{},
		},
		{
			`-1`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
			},
		},
		{
			`!true`,
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
			[]object.Object{},
		},
		{
			`if (true) { 10 }; 3333;`,
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpIf, 10),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJump, 11),
				code.Make(code.OpNull),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 10},
				&object.Integer{Value: 3333},
			},
		},
		{
			`if (true) { 10 } else { 20 }; 3333;`,
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpIf, 10),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJump, 13),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 10},
				&object.Integer{Value: 20},
				&object.Integer{Value: 3333},
			},
		},
	}

	for _, s := range setup {
		actual := compile(s.input)
		expected := combine(s.expectedPieces, s.expectedConstants)
		testBytecode(t, actual, expected)
	}
}

func compile(input string) *Bytecode {
	program := parse(input)
	c := New()
	return c.Compile(program)
}

func parse(input string) *ast.Program {
	lex := lexer.New(input)
	p := parser.New(lex)
	return p.ParseProgram()
}

func combine(
	pieces []code.Instructions,
	constants []object.Object,
) *Bytecode {
	instructions := code.Concatenate(pieces)
	return &Bytecode{Instructions: instructions, Constants: constants}
}

func testBytecode(t *testing.T, actual *Bytecode, expected *Bytecode) {
	testInstructions(t, actual.Instructions, expected.Instructions)
	testConstants(t, actual.Constants, expected.Constants)
}

func testInstructions(
	t *testing.T,
	actual code.Instructions,
	expected code.Instructions,
) {
	if len(actual) != len(expected) {
		t.Fatalf(
			"number of bytes mismatch. got=%q, expected=%q",
			actual,
			expected,
		)
	}
	for i, e := range expected {
		a := actual[i]
		if a != e {
			t.Fatalf("byte mismatch. got=%q, expected=%q", actual, expected)
		}
	}
}

func testConstants(
	t *testing.T,
	actual []object.Object,
	expected []object.Object,
) {
	if len(actual) != len(expected) {
		t.Fatalf(
			"number of constants mismatch. got=%v, expected=%v",
			len(actual),
			len(expected),
		)
	}
	for i, e := range expected {
		testObject(t, actual[i], e)
	}
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

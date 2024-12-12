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
		{
			`let one = 1; let two = 2;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
		},
		{
			`let one = 1; one;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
			},
		},
		{
			`let one = 1; let two = one; two;`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
			},
		},
		{
			`"monkey";`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.String{Value: "monkey"},
			},
		},
		{
			`"mon" + "key";`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.String{Value: "mon"},
				&object.String{Value: "key"},
			},
		},
		{
			`[];`,
			[]code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
			[]object.Object{},
		},
		{
			`[1, 2, 3];`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			},
		},
		{
			`[1 + 2, 3 - 4, 5 * 6];`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
				&object.Integer{Value: 5},
				&object.Integer{Value: 6},
			},
		},
		{
			`{};`,
			[]code.Instructions{
				code.Make(code.OpHash, 0),
				code.Make(code.OpPop),
			},
			[]object.Object{},
		},
		{
			`{1: 2, 3: 4, 5: 6};`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpHash, 3),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
				&object.Integer{Value: 5},
				&object.Integer{Value: 6},
			},
		},
		{
			`{1: 2 + 3, 4: 5 * 6};`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpHash, 2),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
				&object.Integer{Value: 5},
				&object.Integer{Value: 6},
			},
		},
		{
			`[1, 2, 3][1 + 1];`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 1},
				&object.Integer{Value: 1},
			},
		},
		{
			`{1: 2}[2 - 1];`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 2},
				&object.Integer{Value: 1},
			},
		},
		{
			`fn() { return 5 + 10; };`,
			[]code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 5},
				&object.Integer{Value: 10},
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpConstant, 0),
							code.Make(code.OpConstant, 1),
							code.Make(code.OpAdd),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     0,
					NumParameters: 0,
				},
			},
		},
		{
			`fn() { 5 + 10; };`,
			[]code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 5},
				&object.Integer{Value: 10},
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpConstant, 0),
							code.Make(code.OpConstant, 1),
							code.Make(code.OpAdd),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     0,
					NumParameters: 0,
				},
			},
		},
		{
			`fn() { 1; 2; };`,
			[]code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpConstant, 0),
							code.Make(code.OpPop),
							code.Make(code.OpConstant, 1),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     0,
					NumParameters: 0,
				},
			},
		},
		{
			`fn() {};`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.CompiledFunction{
					Instructions:  code.Make(code.OpReturn),
					NumLocals:     0,
					NumParameters: 0,
				},
			},
		},
		{
			`fn() { 24; }();`,
			[]code.Instructions{
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 24},
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpConstant, 0),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     0,
					NumParameters: 0,
				},
			},
		},
		{
			`let noArg = fn() { 24; }; noArg();`,
			[]code.Instructions{
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 24},
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpConstant, 0),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     0,
					NumParameters: 0,
				},
			},
		},
		{
			`let num = 55; fn() { num; };`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 55},
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpGetGlobal, 0),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     0,
					NumParameters: 0,
				},
			},
		},
		{
			`fn() { let num = 55; num; };`,
			[]code.Instructions{
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 55},
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpConstant, 0),
							code.Make(code.OpSetLocal, 0),
							code.Make(code.OpGetLocal, 0),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     1,
					NumParameters: 0,
				},
			},
		},
		{
			`fn() { let a = 55; let b = 77; a + b; };`,
			[]code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 55},
				&object.Integer{Value: 77},
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpConstant, 0),
							code.Make(code.OpSetLocal, 0),
							code.Make(code.OpConstant, 1),
							code.Make(code.OpSetLocal, 1),
							code.Make(code.OpGetLocal, 0),
							code.Make(code.OpGetLocal, 1),
							code.Make(code.OpAdd),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     2,
					NumParameters: 0,
				},
			},
		},
		{
			`let oneArg = fn(a) { a }; oneArg(24);`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpGetLocal, 0),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     1,
					NumParameters: 1,
				},
				&object.Integer{Value: 24},
			},
		},
		{
			`let manyArg = fn(a, b, c) { a; b; c; }; manyArg(24, 25, 26);`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpCall, 3),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpGetLocal, 0),
							code.Make(code.OpPop),
							code.Make(code.OpGetLocal, 1),
							code.Make(code.OpPop),
							code.Make(code.OpGetLocal, 2),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     3,
					NumParameters: 3,
				},
				&object.Integer{Value: 24},
				&object.Integer{Value: 25},
				&object.Integer{Value: 26},
			},
		},
		{
			`len([]); push([], 1);`,
			[]code.Instructions{
				code.Make(code.OpGetBuiltin, 0),
				code.Make(code.OpArray, 0),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
				code.Make(code.OpGetBuiltin, 5),
				code.Make(code.OpArray, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpCall, 2),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.Integer{Value: 1},
			},
		},
		{
			`fn() { len([]); };`,
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
			[]object.Object{
				&object.CompiledFunction{
					Instructions: code.Concatenate(
						[]code.Instructions{
							code.Make(code.OpGetBuiltin, 0),
							code.Make(code.OpArray, 0),
							code.Make(code.OpCall, 1),
							code.Make(code.OpReturnValue),
						},
					),
					NumLocals:     0,
					NumParameters: 0,
				},
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
	case *object.String:
		a, ok := actual.(*object.String)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testStringObject(t, a, e)
	case *object.CompiledFunction:
		a, ok := actual.(*object.CompiledFunction)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testCompiledFunctionObject(t, a, e)
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

func testStringObject(
	t *testing.T,
	actual *object.String,
	expected *object.String,
) {
	if actual.Value != expected.Value {
		t.Fatalf(
			"string value mismatch. got=%v, expected=%v",
			actual.Value,
			expected.Value,
		)
	}
}

func testCompiledFunctionObject(
	t *testing.T,
	actual *object.CompiledFunction,
	expected *object.CompiledFunction,
) {
	testInstructions(t, actual.Instructions, expected.Instructions)
	testNumLocals(t, actual.NumLocals, expected.NumLocals)
	testNumParameters(t, actual.NumParameters, expected.NumParameters)
}

func testNumLocals(t *testing.T, actual int, expected int) {
	if actual != expected {
		t.Fatalf(
			"number of locals mismatch. got=%v, expected=%v",
			actual,
			expected,
		)
	}
}

func testNumParameters(t *testing.T, actual int, expected int) {
	if actual != expected {
		t.Fatalf(
			"number of parameters mismatch. got=%v, expected=%v",
			actual,
			expected,
		)
	}
}

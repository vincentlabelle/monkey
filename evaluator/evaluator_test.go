package evaluator

import (
	"testing"

	"github.com/vincentlabelle/monkey/ast"
	"github.com/vincentlabelle/monkey/lexer"
	"github.com/vincentlabelle/monkey/object"
	"github.com/vincentlabelle/monkey/parser"
)

func Test(t *testing.T) {
	setup := []struct {
		input    string
		expected object.Object
	}{
		{`5;`, &object.Integer{Value: 5}},
		{`10;`, &object.Integer{Value: 10}},
		{`-5;`, &object.Integer{Value: -5}},
		{`-10;`, &object.Integer{Value: -10}},
		{`5 + 5 + 5 + 5 - 10;`, &object.Integer{Value: 10}},
		{`2 * 2 * 2 * 2 * 2;`, &object.Integer{Value: 32}},
		{`-50 + 100 + -50;`, &object.Integer{Value: 0}},
		{`5 * 2 + 10;`, &object.Integer{Value: 20}},
		{`5 + 2 * 10;`, &object.Integer{Value: 25}},
		{`20 + 2 * -10;`, &object.Integer{Value: 0}},
		{`50 / 2 * 2 + 10;`, &object.Integer{Value: 60}},
		{`2 * (5 + 10);`, &object.Integer{Value: 30}},
		{`3 * 3 * 3 + 10;`, &object.Integer{Value: 37}},
		{`3 * (3 * 3) + 10;`, &object.Integer{Value: 37}},
		{`(5 + 10 * 2 + 15 / 3) * 2 + -10;`, &object.Integer{Value: 50}},
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
		{`false != true;`, object.TRUE},
		{`(1 < 2) == true;`, object.TRUE},
		{`(1 < 2) == false;`, object.FALSE},
		{`(1 > 2) == true;`, object.FALSE},
		{`(1 > 2) == false;`, object.TRUE},
		{`!true;`, object.FALSE},
		{`!false;`, object.TRUE},
		{`!5;`, object.FALSE},
		{`!!true;`, object.TRUE},
		{`!!false;`, object.FALSE},
		{`!!5;`, object.TRUE},
		{`if (true) { 10 };`, &object.Integer{Value: 10}},
		{`if (false) { 10 };`, object.NULL},
		{`if (1) { 10 };`, &object.Integer{Value: 10}},
		{`if (1 < 2) { 10 };`, &object.Integer{Value: 10}},
		{`if (1 > 2) { 10 };`, object.NULL},
		{`if (1 > 2) { 10 } else { 20 };`, &object.Integer{Value: 20}},
		{`if (1 < 2) { 10 } else { 20 };`, &object.Integer{Value: 10}},
		{`return 10;`, &object.Integer{Value: 10}},
		{`return 10; 9;`, &object.Integer{Value: 10}},
		{`return 2 * 5; 9;`, &object.Integer{Value: 10}},
		{`9; return 2 * 5; 9;`, &object.Integer{Value: 10}},
		{`if (10 > 1) {return 10; }`, &object.Integer{Value: 10}},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}
			`,
			&object.Integer{Value: 10},
		},
		{
			`
			let f = fn(x) {
				return x;
				x + 10;
			};
			f(10);
			`,
			&object.Integer{Value: 10},
		},
		{
			`
			let f = fn(x) {
				let result = x + 10;
				return result;
				10;
			};
			f(10);
			`,
			&object.Integer{Value: 20},
		},
		{`let a = 5; a;`, &object.Integer{Value: 5}},
		{`let a = 5 * 5; a;`, &object.Integer{Value: 25}},
		{`let a = 5; let b = a; b;`, &object.Integer{Value: 5}},
		{
			`let a = 5; let b = a; let c = a + b + 5; c;`,
			&object.Integer{Value: 15},
		},
		{
			`let identity = fn(x) { x; }; identity(5);`,
			&object.Integer{Value: 5},
		},
		{
			`let identity = fn(x) { return x; }; identity(5);`,
			&object.Integer{Value: 5},
		},
		{
			`let double = fn(x) { x * 2; }; double(5);`,
			&object.Integer{Value: 10},
		},
		{
			`let add = fn(x, y) { x + y; }; add(5, 5);`,
			&object.Integer{Value: 10},
		},
		{
			`let add = fn(x, y) { return x + y; }; add(5, 5); add(10, 5)`,
			&object.Integer{Value: 15},
		},
		{
			`let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));`,
			&object.Integer{Value: 20},
		},
		{`fn(x) { x; }(5);`, &object.Integer{Value: 5}},
		{
			`
			let first = 10;
			let second = 10;
			let third = 10;

			let ourFunction = fn(first) {
				let second = 20;
				first + second + third;
			};

			ourFunction(20) + first + second;
			`,
			&object.Integer{Value: 70},
		},
		{
			`
			let newAdder = fn(x) {
				fn(y) { x + y };
			};
			let addTwo = newAdder(2);
			addTwo(2);
			`,
			&object.Integer{Value: 4},
		},
		{`"Hello World!";`, &object.String{Value: "Hello World!"}},
		{`"Hello" + " "  + "World!";`, &object.String{Value: "Hello World!"}},
		{`len("");`, &object.Integer{Value: 0}},
		{`len("four");`, &object.Integer{Value: 4}},
		{`len("hello world");`, &object.Integer{Value: 11}},
		{`len("hello world");`, &object.Integer{Value: 11}},
		{
			`[1, 2 * 2, 3 + 3];`,
			&object.Array{Elements: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 4},
				&object.Integer{Value: 6},
			}},
		},
		{`[1, 2, 3][0];`, &object.Integer{Value: 1}},
		{`[1, 2, 3][1];`, &object.Integer{Value: 2}},
		{`[1, 2, 3][2];`, &object.Integer{Value: 3}},
		{`let i = 0; [1][i];`, &object.Integer{Value: 1}},
		{`[1, 2, 3][1 + 1];`, &object.Integer{Value: 3}},
		{`let a = [1, 2, 3]; a[2];`, &object.Integer{Value: 3}},
		{`let a = [1, 2, 3]; a[0] + a[1] + a[2];`, &object.Integer{Value: 6}},
		{`let a = [1, 2, 3]; let i = a[0]; a[i]`, &object.Integer{Value: 2}},
		{`[1, 2, 3][3];`, object.NULL},
		{`[1, 2, 3][-1];`, object.NULL},
		{`len([]);`, &object.Integer{Value: 0}},
		{`len([1]);`, &object.Integer{Value: 1}},
		{`len([1, 2, 3]);`, &object.Integer{Value: 3}},
		{`first([]);`, object.NULL},
		{`first([1]);`, &object.Integer{Value: 1}},
		{`first([1, 2, 3]);`, &object.Integer{Value: 1}},
		{`last([]);`, object.NULL},
		{`last([1]);`, &object.Integer{Value: 1}},
		{`last([1, 2, 3]);`, &object.Integer{Value: 3}},
		{`rest([]);`, object.NULL},
		{`rest([1]);`, &object.Array{Elements: []object.Object{}}},
		{`rest([1, 2]);`, &object.Array{Elements: []object.Object{
			&object.Integer{Value: 2},
		}}},
		{`push([], 1);`, &object.Array{Elements: []object.Object{
			&object.Integer{Value: 1},
		}}},
		{`push([1], 2);`, &object.Array{Elements: []object.Object{
			&object.Integer{Value: 1},
			&object.Integer{Value: 2},
		}}},
		{`puts(1);`, object.NULL},
		{
			`
			let two = "two";
			{
				"one": 10 - 9,
				two: 1 + 1,
				"thr" + "ee": 6 / 2,
				4: 4,
				true: 5,
				false: 6
			};
			`, &object.Hash{Pairs: map[object.HashKey]object.HashPair{
				(&object.String{Value: "one"}).HashKey(): {
					Key:   &object.String{Value: "one"},
					Value: &object.Integer{Value: 1},
				},
				(&object.String{Value: "two"}).HashKey(): {
					Key:   &object.String{Value: "two"},
					Value: &object.Integer{Value: 2},
				},
				(&object.String{Value: "three"}).HashKey(): {
					Key:   &object.String{Value: "three"},
					Value: &object.Integer{Value: 3},
				},
				(&object.Integer{Value: 4}).HashKey(): {
					Key:   &object.Integer{Value: 4},
					Value: &object.Integer{Value: 4},
				},
				(object.TRUE).HashKey(): {
					Key:   object.TRUE,
					Value: &object.Integer{Value: 5},
				},
				(object.FALSE).HashKey(): {
					Key:   object.FALSE,
					Value: &object.Integer{Value: 6},
				},
			}},
		},
		{`{"foo": 5}["foo"];`, &object.Integer{Value: 5}},
		{`{"foo": 5}["bar"];`, object.NULL},
		{`let key = "foo"; {"foo": 5}[key];`, &object.Integer{Value: 5}},
		{`{}["foo"];`, object.NULL},
		{`{5: 5}[5];`, &object.Integer{Value: 5}},
		{`{true: 5}[true];`, &object.Integer{Value: 5}},
		{`{false: 5}[false];`, &object.Integer{Value: 5}},
	}
	for _, s := range setup {
		actual := eval(s.input)
		testObject(t, actual, s.expected)
	}
}

func eval(input string) object.Object {
	program := parse(input)
	env := object.NewEnvironment()
	return Eval(program, env)
}

func parse(input string) *ast.Program {
	lex := lexer.New(input)
	p := parser.New(lex)
	return p.ParseProgram()
}

func testObject(t *testing.T, actual object.Object, expected object.Object) {
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
		testInteger(t, a, e)
	case *object.Boolean:
		a, ok := actual.(*object.Boolean)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testBoolean(t, a, e)
	case *object.String:
		a, ok := actual.(*object.String)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testString(t, a, e)
	case *object.Array:
		a, ok := actual.(*object.Array)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testArray(t, a, e)
	case *object.Hash:
		a, ok := actual.(*object.Hash)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testHash(t, a, e)
	case *object.Null:
		_, ok := actual.(*object.Null)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
	default:
		t.Fatal("object type unknown")
	}
}

func testInteger(
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

func testBoolean(
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

func testString(
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

func testArray(
	t *testing.T,
	actual *object.Array,
	expected *object.Array,
) {
	testObjects(t, actual.Elements, expected.Elements)
}

func testObjects(
	t *testing.T,
	actual []object.Object,
	expected []object.Object,
) {
	if len(actual) != len(expected) {
		t.Fatalf(
			"number of objects mismatch. got=%v, expected=%v",
			len(actual),
			len(expected),
		)
	}
	for i := 0; i < len(actual); i++ {
		testObject(t, actual[i], expected[i])
	}
}

func testHash(
	t *testing.T,
	actual *object.Hash,
	expected *object.Hash,
) {
	for ek, ev := range expected.Pairs {
		av, ok := actual.Pairs[ek]
		if !ok {
			t.Fatalf("hash key mismatch. missing %v in actual", ek)
		}
		testHashPair(t, av, ev)
	}
}

func testHashPair(
	t *testing.T,
	actual object.HashPair,
	expected object.HashPair,
) {
	testObject(t, actual.Key, expected.Key)
	testObject(t, actual.Value, expected.Value)
}

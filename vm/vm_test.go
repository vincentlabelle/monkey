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
		{`if (true) { 10 };`, &object.Integer{Value: 10}},
		{`if (true) { 10 } else { 20 };`, &object.Integer{Value: 10}},
		{`if (false) { 10 } else { 20 };`, &object.Integer{Value: 20}},
		{`if (1) { 10 };`, &object.Integer{Value: 10}},
		{`if (1 < 2) { 10 };`, &object.Integer{Value: 10}},
		{`if (1 < 2) { 10 } else { 20 };`, &object.Integer{Value: 10}},
		{`if (1 > 2) { 10 } else { 20 };`, &object.Integer{Value: 20}},
		{`if (1 > 2) { 10 };`, object.NULL},
		{`if (false) { 10 };`, object.NULL},
		{`!(if (false) { 5 });`, object.TRUE},
		{
			`if (if (false) { 10 }) { 10 } else { 20 };`,
			&object.Integer{Value: 20},
		},
		{`let one = 1; one;`, &object.Integer{Value: 1}},
		{`let one = 1; let two = 2; one + two;`, &object.Integer{Value: 3}},
		{
			`let one = 1; let two = one + one; one + two;`,
			&object.Integer{Value: 3},
		},
		{`"monkey";`, &object.String{Value: "monkey"}},
		{`"mon" + "key";`, &object.String{Value: "monkey"}},
		{`"mon" + "key" + "banana";`, &object.String{Value: "monkeybanana"}},
		{
			`[];`,
			&object.Array{Elements: []object.Object{}},
		},
		{
			`[1, 2, 3];`,
			&object.Array{Elements: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			}},
		},
		{
			`[1 + 2, 3 * 4, 5 + 6];`,
			&object.Array{Elements: []object.Object{
				&object.Integer{Value: 3},
				&object.Integer{Value: 12},
				&object.Integer{Value: 11},
			}},
		},
		{
			`{};`,
			&object.Hash{},
		},
		{
			`{1: 2, 2: 3};`,
			&object.Hash{Pairs: map[object.HashKey]object.HashPair{
				(&object.Integer{Value: 1}).HashKey(): {
					Key:   &object.Integer{Value: 1},
					Value: &object.Integer{Value: 2},
				},
				(&object.Integer{Value: 2}).HashKey(): {
					Key:   &object.Integer{Value: 2},
					Value: &object.Integer{Value: 3},
				},
			}},
		},
		{
			`{1 + 1: 2 * 2, 3 + 3: 4 * 4};`,
			&object.Hash{Pairs: map[object.HashKey]object.HashPair{
				(&object.Integer{Value: 2}).HashKey(): {
					Key:   &object.Integer{Value: 2},
					Value: &object.Integer{Value: 4},
				},
				(&object.Integer{Value: 6}).HashKey(): {
					Key:   &object.Integer{Value: 6},
					Value: &object.Integer{Value: 16},
				},
			}},
		},
		{`[1, 2, 3][1];`, &object.Integer{Value: 2}},
		{`[1, 2, 3][0 + 2];`, &object.Integer{Value: 3}},
		{`[[1, 1, 1]][0][0];`, &object.Integer{Value: 1}},
		{`[][0];`, object.NULL},
		{`[1, 2, 3][99];`, object.NULL},
		{`[1][-1];`, object.NULL},
		{`{1: 1, 2: 2}[1];`, &object.Integer{Value: 1}},
		{`{1: 1, 2: 2}[2];`, &object.Integer{Value: 2}},
		{`{1: 1}[0];`, object.NULL},
		{`{}[0];`, object.NULL},
		{
			`let fivePlusTen = fn() { 5 + 10; }; fivePlusTen();`,
			&object.Integer{Value: 15},
		},
		{
			`let one = fn() { 1; }; let two = fn() { 2; }; one() + two();`,
			&object.Integer{Value: 3},
		},
		{
			`
			let a = fn() { 1; };
			let b = fn() { a() + 1; };
			let c = fn() { b() + 1; };
			c();
			`,
			&object.Integer{Value: 3},
		},
		{
			`let earlyExit = fn() { return 99; 100; }; earlyExit();`,
			&object.Integer{Value: 99},
		},
		{
			`let earlyExit = fn() { return 99; return 100; }; earlyExit();`,
			&object.Integer{Value: 99},
		},
		{`let noReturn = fn() {}; noReturn();`, object.NULL},
		{
			`
			let noReturn = fn() {};
			let noReturnTwo = fn() { noReturn(); };
			noReturn();
			noReturnTwo();
			`,
			object.NULL,
		},
		{
			`
			let returnsOne = fn() { 1; };
			let returnsOneReturner = fn() { returnsOne; };
			returnsOneReturner()();
			`,
			&object.Integer{Value: 1},
		},
		{
			`let one = fn() { let one = 1; one; }; one();`,
			&object.Integer{Value: 1},
		},
		{
			`
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			oneAndTwo();
			`,
			&object.Integer{Value: 3},
		},
		{
			`
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			let threeAndFour = fn() {
				let three = 3;
				let four = 4;
				three + four;
			};
			oneAndTwo() + threeAndFour();
			`,
			&object.Integer{Value: 10},
		},
		{
			`
			let firstFoobar = fn() { let foobar = 50; foobar; };
			let secondFoobar = fn() { let foobar = 100; foobar; };
			firstFoobar() + secondFoobar();
			`,
			&object.Integer{Value: 150},
		},
		{
			`
			let globalSeed = 50;
			let minusOne = fn() { let num = 1; globalSeed - num; };
			let minusTwo = fn() { let num = 2; globalSeed - num; };
			minusOne() + minusTwo();
			`,
			&object.Integer{Value: 97},
		},
		{
			`
			let returnsOneReturner = fn() {
				let returnsOne = fn() { 1; };
				returnsOne;
			};
			returnsOneReturner()();
			`,
			&object.Integer{Value: 1},
		},
		{
			`let identity = fn(a) { a; }; identity(4);`,
			&object.Integer{Value: 4},
		},
		{
			`let sum = fn(a, b) { a + b; }; sum(1, 2);`,
			&object.Integer{Value: 3},
		},
		{
			`let sum = fn(a, b) { a + b; }; sum(1, 2);`,
			&object.Integer{Value: 3},
		},
		{
			`let sum = fn(a, b) { let c = a + b; c; }; sum(1, 2);`,
			&object.Integer{Value: 3},
		},
		{
			`let sum = fn(a, b) { let c = a + b; c; }; sum(1, 2) + sum(3, 4);`,
			&object.Integer{Value: 10},
		},
		{
			`
			let sum = fn(a, b) { let c = a + b; c; };
			let outer = fn() { sum(1, 2) + sum(3, 4); };
			outer();
			`,
			&object.Integer{Value: 10},
		},
		{
			`
			let globalNum = 10;

			let sum = fn(a, b) {
				let c = a + b;
				c + globalNum;
			};

			let outer = fn() {
				sum(1, 2) + sum(3, 4) + globalNum;
			};

			outer() + globalNum;
			`,
			&object.Integer{Value: 50},
		},
		{`len("")`, &object.Integer{Value: 0}},
		{`len("four")`, &object.Integer{Value: 4}},
		{`len("hello world")`, &object.Integer{Value: 11}},
		{`len([1, 2, 3])`, &object.Integer{Value: 3}},
		{`len([])`, &object.Integer{Value: 0}},
		{`puts("hello", "world!")`, object.NULL},
		{`first([1, 2, 3])`, &object.Integer{Value: 1}},
		{`first([])`, object.NULL},
		{`last([1, 2, 3])`, &object.Integer{Value: 3}},
		{`last([])`, object.NULL},
		{
			`rest([1, 2, 3])`,
			&object.Array{
				Elements: []object.Object{
					&object.Integer{Value: 2},
					&object.Integer{Value: 3},
				}},
		},
		{`rest([])`, object.NULL},
		{
			`push([], 1)`,
			&object.Array{
				Elements: []object.Object{
					&object.Integer{Value: 1},
				}},
		},
		{
			`
			let newClosure = fn(a) { fn() { a; } };
			let closure = newClosure(99);
			closure();
			`,
			&object.Integer{Value: 99},
		},
		{
			`
			let newAdder = fn(a, b) {
				fn(c) { a + b + c; };
			};
			let adder = newAdder(1, 2);
			adder(8);
			`,
			&object.Integer{Value: 11},
		},
		{
			`
			let newAdder = fn(a, b) {
				let c = a + b;
				fn(d) { c + d; };
			};
			let adder = newAdder(1, 2);
			adder(8);
			`,
			&object.Integer{Value: 11},
		},
		{
			`
			let newAdderOuter = fn(a, b) {
				let c = a + b;
				fn(d) { 
					let e = c + d;
					fn(f) { e + f; };
				};
			};
			let newAdderInner = newAdderOuter(1, 2);
			let adder = newAdderInner(3);
			adder(8);
			`,
			&object.Integer{Value: 14},
		},
		{
			`
			let a = 1;
			let newAdderOuter = fn(b) {
				fn(c) { 
					fn(d) { a + b + c + d; };
				};
			};
			let newAdderInner = newAdderOuter(2);
			let adder = newAdderInner(3);
			adder(8);
			`,
			&object.Integer{Value: 14},
		},
		{
			`
			let newClosure = fn(a, b) {
				let one = fn() { a; };
				let two = fn() { b; };
				fn() { one() + two(); };
			};
			let closure = newClosure(9, 90);
			closure();
			`,
			&object.Integer{Value: 99},
		},
		{
			`
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			countDown(1);
			`,
			&object.Integer{Value: 0},
		},
		{
			`
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			let wrapper = fn() {
				countDown(1);
			};
			wrapper();
			`,
			&object.Integer{Value: 0},
		},
		{
			`
			let wrapper = fn() {
				let countDown = fn(x) {
					if (x == 0) {
						return 0;
					} else {
						countDown(x - 1);
					}
				};
				countDown(1);
			};
			wrapper();
			`,
			&object.Integer{Value: 0},
		},
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
	return c.Compile(program)
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
	case *object.Null:
		_, ok := actual.(*object.Null)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
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
	case *object.Array:
		a, ok := actual.(*object.Array)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testArrayObject(t, a, e)
	case *object.Hash:
		a, ok := actual.(*object.Hash)
		if !ok {
			t.Fatalf(
				"object type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testHashObject(t, a, e)
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

func testArrayObject(
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

func testHashObject(
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

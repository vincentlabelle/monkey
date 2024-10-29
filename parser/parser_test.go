package parser

import (
	"testing"

	"github.com/vincentlabelle/monkey/ast"
	"github.com/vincentlabelle/monkey/lexer"
)

func Test(t *testing.T) {
	setup := []struct {
		input    string
		expected *ast.Program
	}{
		{
			input: `25;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.IntegerLiteral{Value: 25},
					},
				},
			},
		},
		{
			input: `foobar;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.Identifier{Value: "foobar"},
					},
				},
			},
		},
		{
			input: `false;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.BooleanLiteral{Value: false},
					},
				},
			},
		},
		{
			input: `-25;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.PrefixExpression{
							Operator: "-",
							Right:    &ast.IntegerLiteral{Value: 25},
						},
					},
				},
			},
		},
		{
			input: `!true;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.PrefixExpression{
							Operator: "!",
							Right:    &ast.BooleanLiteral{Value: true},
						},
					},
				},
			},
		},
		{
			input: `!abc;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.PrefixExpression{
							Operator: "!",
							Right:    &ast.Identifier{Value: "abc"},
						},
					},
				},
			},
		},
		{
			input: `abc + 5;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.Identifier{Value: "abc"},
							Operator: "+",
							Right:    &ast.IntegerLiteral{Value: 5},
						},
					},
				},
			},
		},
		{
			input: `abc - 5;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.Identifier{Value: "abc"},
							Operator: "-",
							Right:    &ast.IntegerLiteral{Value: 5},
						},
					},
				},
			},
		},
		{
			input: `abc * 5;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.Identifier{Value: "abc"},
							Operator: "*",
							Right:    &ast.IntegerLiteral{Value: 5},
						},
					},
				},
			},
		},
		{
			input: `abc / 5;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.Identifier{Value: "abc"},
							Operator: "/",
							Right:    &ast.IntegerLiteral{Value: 5},
						},
					},
				},
			},
		},
		{
			input: `abc > 5;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.Identifier{Value: "abc"},
							Operator: ">",
							Right:    &ast.IntegerLiteral{Value: 5},
						},
					},
				},
			},
		},
		{
			input: `abc < 5;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.Identifier{Value: "abc"},
							Operator: "<",
							Right:    &ast.IntegerLiteral{Value: 5},
						},
					},
				},
			},
		},
		{
			input: `abc == 5;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.Identifier{Value: "abc"},
							Operator: "==",
							Right:    &ast.IntegerLiteral{Value: 5},
						},
					},
				},
			},
		},
		{
			input: `true != false;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.BooleanLiteral{Value: true},
							Operator: "!=",
							Right:    &ast.BooleanLiteral{Value: false},
						},
					},
				},
			},
		},
		{
			input: `return 5;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ReturnStatement{
						Value: &ast.IntegerLiteral{Value: 5},
					},
				},
			},
		},
		{
			input: `return true;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ReturnStatement{
						Value: &ast.BooleanLiteral{Value: true},
					},
				},
			},
		},
		{
			input: `return abc;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ReturnStatement{
						Value: &ast.Identifier{Value: "abc"},
					},
				},
			},
		},
		{
			input: `let abc = 25;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Name:  &ast.Identifier{Value: "abc"},
						Value: &ast.IntegerLiteral{Value: 25},
					},
				},
			},
		},
		{
			input: `let y = true;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Name:  &ast.Identifier{Value: "y"},
						Value: &ast.BooleanLiteral{Value: true},
					},
				},
			},
		},
		{
			input: `let foobar = y;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Name:  &ast.Identifier{Value: "foobar"},
						Value: &ast.Identifier{Value: "y"},
					},
				},
			},
		},
		{
			input: `-a * b;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.PrefixExpression{
								Operator: "-",
								Right:    &ast.Identifier{Value: "a"},
							},
							Operator: "*",
							Right:    &ast.Identifier{Value: "b"},
						},
					},
				},
			},
		},
		{
			input: `!-a;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.PrefixExpression{
							Operator: "!",
							Right: &ast.PrefixExpression{
								Operator: "-",
								Right:    &ast.Identifier{Value: "a"},
							},
						},
					},
				},
			},
		},
		{
			input: `a + b + c;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left:     &ast.Identifier{Value: "a"},
								Operator: "+",
								Right:    &ast.Identifier{Value: "b"},
							},
							Operator: "+",
							Right:    &ast.Identifier{Value: "c"},
						},
					},
				},
			},
		},
		{
			input: `a + b - c;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left:     &ast.Identifier{Value: "a"},
								Operator: "+",
								Right:    &ast.Identifier{Value: "b"},
							},
							Operator: "-",
							Right:    &ast.Identifier{Value: "c"},
						},
					},
				},
			},
		},
		{
			input: `a * b * c;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left:     &ast.Identifier{Value: "a"},
								Operator: "*",
								Right:    &ast.Identifier{Value: "b"},
							},
							Operator: "*",
							Right:    &ast.Identifier{Value: "c"},
						},
					},
				},
			},
		},
		{
			input: `a * b / c;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left:     &ast.Identifier{Value: "a"},
								Operator: "*",
								Right:    &ast.Identifier{Value: "b"},
							},
							Operator: "/",
							Right:    &ast.Identifier{Value: "c"},
						},
					},
				},
			},
		},
		{
			input: `a + b / c;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.Identifier{Value: "a"},
							Operator: "+",
							Right: &ast.InfixExpression{
								Left:     &ast.Identifier{Value: "b"},
								Operator: "/",
								Right:    &ast.Identifier{Value: "c"},
							},
						},
					},
				},
			},
		},
		{
			input: `a + b * c + d / e - f;`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left: &ast.InfixExpression{
									Left:     &ast.Identifier{Value: "a"},
									Operator: "+",
									Right: &ast.InfixExpression{
										Left:     &ast.Identifier{Value: "b"},
										Operator: "*",
										Right:    &ast.Identifier{Value: "c"},
									},
								},
								Operator: "+",
								Right: &ast.InfixExpression{
									Left:     &ast.Identifier{Value: "d"},
									Operator: "/",
									Right:    &ast.Identifier{Value: "e"},
								},
							},
							Operator: "-",
							Right:    &ast.Identifier{Value: "f"},
						},
					},
				},
			},
		},
		{
			input: `3 + 4; -5 * 5`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.IntegerLiteral{Value: 3},
							Operator: "+",
							Right:    &ast.IntegerLiteral{Value: 4},
						},
					},
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.PrefixExpression{
								Operator: "-",
								Right:    &ast.IntegerLiteral{Value: 5},
							},
							Operator: "*",
							Right:    &ast.IntegerLiteral{Value: 5},
						},
					},
				},
			},
		},
		{
			input: `5 > 4 == 3 < 4`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left:     &ast.IntegerLiteral{Value: 5},
								Operator: ">",
								Right:    &ast.IntegerLiteral{Value: 4},
							},
							Operator: "==",
							Right: &ast.InfixExpression{
								Left:     &ast.IntegerLiteral{Value: 3},
								Operator: "<",
								Right:    &ast.IntegerLiteral{Value: 4},
							},
						},
					},
				},
			},
		},
		{
			input: `5 > 4 != 3 < 4`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left:     &ast.IntegerLiteral{Value: 5},
								Operator: ">",
								Right:    &ast.IntegerLiteral{Value: 4},
							},
							Operator: "!=",
							Right: &ast.InfixExpression{
								Left:     &ast.IntegerLiteral{Value: 3},
								Operator: "<",
								Right:    &ast.IntegerLiteral{Value: 4},
							},
						},
					},
				},
			},
		},
		{
			input: `1 + (2 + 3) + 4`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left:     &ast.IntegerLiteral{Value: 1},
								Operator: "+",
								Right: &ast.InfixExpression{
									Left:     &ast.IntegerLiteral{Value: 2},
									Operator: "+",
									Right:    &ast.IntegerLiteral{Value: 3},
								},
							},
							Operator: "+",
							Right:    &ast.IntegerLiteral{Value: 4},
						},
					},
				},
			},
		},
		{
			input: `2 * (5 + 5)`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left:     &ast.IntegerLiteral{Value: 2},
							Operator: "*",
							Right: &ast.InfixExpression{
								Left:     &ast.IntegerLiteral{Value: 5},
								Operator: "+",
								Right:    &ast.IntegerLiteral{Value: 5},
							},
						},
					},
				},
			},
		},
		{
			input: `(5 + 5) * 2 * (5 + 5)`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left: &ast.InfixExpression{
									Left:     &ast.IntegerLiteral{Value: 5},
									Operator: "+",
									Right:    &ast.IntegerLiteral{Value: 5},
								},
								Operator: "*",
								Right:    &ast.IntegerLiteral{Value: 2},
							},
							Operator: "*",
							Right: &ast.InfixExpression{
								Left:     &ast.IntegerLiteral{Value: 5},
								Operator: "+",
								Right:    &ast.IntegerLiteral{Value: 5},
							},
						},
					},
				},
			},
		},
		{
			input: `!(5 + 5)`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.PrefixExpression{
							Operator: "!",
							Right: &ast.InfixExpression{
								Left:     &ast.IntegerLiteral{Value: 5},
								Operator: "+",
								Right:    &ast.IntegerLiteral{Value: 5},
							},
						},
					},
				},
			},
		},
		{
			input: `if (x < y) { x }`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.IfExpression{
							Condition: &ast.InfixExpression{
								Left:     &ast.Identifier{Value: "x"},
								Operator: "<",
								Right:    &ast.Identifier{Value: "y"},
							},
							Consequence: &ast.BlockStatement{
								Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expression: &ast.Identifier{
											Value: "x",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `if (x < y) { x } else { y }`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.IfExpression{
							Condition: &ast.InfixExpression{
								Left:     &ast.Identifier{Value: "x"},
								Operator: "<",
								Right:    &ast.Identifier{Value: "y"},
							},
							Consequence: &ast.BlockStatement{
								Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expression: &ast.Identifier{
											Value: "x",
										},
									},
								},
							},
							Alternative: &ast.BlockStatement{
								Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expression: &ast.Identifier{
											Value: "y",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `fn() {};`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.FunctionLiteral{
							Parameters: []*ast.Identifier{},
							Body: &ast.BlockStatement{
								Statements: []ast.Statement{},
							},
						},
					},
				},
			},
		},
		{
			input: `fn(x) {};`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.FunctionLiteral{
							Parameters: []*ast.Identifier{
								{Value: "x"},
							},
							Body: &ast.BlockStatement{
								Statements: []ast.Statement{},
							},
						},
					},
				},
			},
		},
		{
			input: `fn(x, y) { x + y; };`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.FunctionLiteral{
							Parameters: []*ast.Identifier{
								{Value: "x"},
								{Value: "y"},
							},
							Body: &ast.BlockStatement{
								Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expression: &ast.InfixExpression{
											Left: &ast.Identifier{
												Value: "x",
											},
											Operator: "+",
											Right: &ast.Identifier{
												Value: "y",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `add();`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Function:  &ast.Identifier{Value: "add"},
							Arguments: []ast.Expression{},
						},
					},
				},
			},
		},
		{
			input: `add(1);`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Function: &ast.Identifier{Value: "add"},
							Arguments: []ast.Expression{
								&ast.IntegerLiteral{Value: 1},
							},
						},
					},
				},
			},
		},
		{
			input: `add(1, 2 * 3, 4 + 5);`,
			expected: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.CallExpression{
							Function: &ast.Identifier{Value: "add"},
							Arguments: []ast.Expression{
								&ast.IntegerLiteral{Value: 1},
								&ast.InfixExpression{
									Left:     &ast.IntegerLiteral{Value: 2},
									Operator: "*",
									Right:    &ast.IntegerLiteral{Value: 3},
								},
								&ast.InfixExpression{
									Left:     &ast.IntegerLiteral{Value: 4},
									Operator: "+",
									Right:    &ast.IntegerLiteral{Value: 5},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, s := range setup {
		lex := lexer.New(s.input)
		p := New(lex)
		program := p.ParseProgram()
		testStatements(t, program.Statements, s.expected.Statements)
	}
}

func testStatements(
	t *testing.T,
	actual []ast.Statement,
	expected []ast.Statement,
) {
	if len(actual) != len(expected) {
		t.Fatalf(
			"number of statements mismatch. got=%v, expected=%v",
			len(actual),
			len(expected),
		)
	}
	for i := 0; i < len(actual); i++ {
		testStatement(t, actual[i], expected[i])
	}
}

func testStatement(
	t *testing.T,
	actual ast.Statement,
	expected ast.Statement,
) {
	switch e := expected.(type) {
	case *ast.ExpressionStatement:
		a, ok := actual.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"statement type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testExpressionStatement(t, a, e)
	case *ast.ReturnStatement:
		a, ok := actual.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf(
				"statement type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testReturnStatement(t, a, e)
	case *ast.LetStatement:
		a, ok := actual.(*ast.LetStatement)
		if !ok {
			t.Fatalf(
				"statement type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testLetStatement(t, a, e)
	}
}

func testExpressionStatement(
	t *testing.T,
	actual *ast.ExpressionStatement,
	expected *ast.ExpressionStatement,
) {
	testExpression(t, actual.Expression, expected.Expression)
}

func testExpression(
	t *testing.T,
	actual ast.Expression,
	expected ast.Expression,
) {
	switch e := expected.(type) {
	case *ast.IntegerLiteral:
		a, ok := actual.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf(
				"expression type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testIntegerLiteral(t, a, e)
	case *ast.Identifier:
		a, ok := actual.(*ast.Identifier)
		if !ok {
			t.Fatalf(
				"expression type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testIdentifier(t, a, e)
	case *ast.BooleanLiteral:
		a, ok := actual.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf(
				"expression type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testBooleanLiteral(t, a, e)
	case *ast.PrefixExpression:
		a, ok := actual.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf(
				"expression type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testPrefixExpression(t, a, e)
	case *ast.InfixExpression:
		a, ok := actual.(*ast.InfixExpression)
		if !ok {
			t.Fatalf(
				"expression type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testInfixExpression(t, a, e)
	case *ast.IfExpression:
		a, ok := actual.(*ast.IfExpression)
		if !ok {
			t.Fatalf(
				"expression type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testIfExpression(t, a, e)
	case *ast.FunctionLiteral:
		a, ok := actual.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf(
				"expression type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testFunctionLiteral(t, a, e)
	case *ast.CallExpression:
		a, ok := actual.(*ast.CallExpression)
		if !ok {
			t.Fatalf(
				"expression type mismatch. got=%T, expected=%T",
				actual,
				expected,
			)
		}
		testCallExpression(t, a, e)
	}

}

func testIntegerLiteral(
	t *testing.T,
	actual *ast.IntegerLiteral,
	expected *ast.IntegerLiteral,
) {
	if actual.Value != expected.Value {
		t.Fatalf(
			"integer literal value mismatch. got=%v, expected=%v",
			actual.Value,
			expected.Value,
		)
	}
}

func testIdentifier(
	t *testing.T,
	actual *ast.Identifier,
	expected *ast.Identifier,
) {
	if actual.Value != expected.Value {
		t.Fatalf(
			"identifier value mismatch. got=%v, expected=%v",
			actual.Value,
			expected.Value,
		)
	}
}

func testBooleanLiteral(
	t *testing.T,
	actual *ast.BooleanLiteral,
	expected *ast.BooleanLiteral,
) {
	if actual.Value != expected.Value {
		t.Fatalf(
			"boolean literal value mismatch. got=%v, expected=%v",
			actual.Value,
			expected.Value,
		)
	}
}

func testPrefixExpression(
	t *testing.T,
	actual *ast.PrefixExpression,
	expected *ast.PrefixExpression,
) {
	if actual.Operator != expected.Operator {
		t.Fatalf(
			"prefix operator mismatch. got=%v, expected=%v",
			actual.Operator,
			expected.Operator,
		)
	}
	testExpression(t, actual.Right, expected.Right)
}

func testInfixExpression(
	t *testing.T,
	actual *ast.InfixExpression,
	expected *ast.InfixExpression,
) {
	if actual.Operator != expected.Operator {
		t.Fatalf(
			"infix operator mismatch. got=%v, expected=%v",
			actual.Operator,
			expected.Operator,
		)
	}
	testExpression(t, actual.Left, expected.Left)
	testExpression(t, actual.Right, expected.Right)
}

func testIfExpression(
	t *testing.T,
	actual *ast.IfExpression,
	expected *ast.IfExpression,
) {
	testExpression(t, actual.Condition, expected.Condition)
	testBlockStatement(t, actual.Consequence, expected.Consequence)
	if expected.Alternative != nil {
		if actual.Alternative == nil {
			t.Fatal("if expression mismatch. expected alternative, but got nil")
		}
		testBlockStatement(t, actual.Alternative, expected.Alternative)
	} else if actual.Alternative != nil {
		t.Fatal("if expression mismatch. expected nil, but got alternative")
	}
}

func testBlockStatement(
	t *testing.T,
	actual *ast.BlockStatement,
	expected *ast.BlockStatement,
) {
	testStatements(t, actual.Statements, expected.Statements)
}

func testFunctionLiteral(
	t *testing.T,
	actual *ast.FunctionLiteral,
	expected *ast.FunctionLiteral,
) {
	testIdentifiers(t, actual.Parameters, expected.Parameters)
	testBlockStatement(t, actual.Body, expected.Body)
}

func testIdentifiers(
	t *testing.T,
	actual []*ast.Identifier,
	expected []*ast.Identifier,
) {
	if len(actual) != len(expected) {
		t.Fatalf(
			"number of identifiers mismatch. got=%v, expected=%v",
			len(actual),
			len(expected),
		)
	}
	for i := 0; i < len(actual); i++ {
		testIdentifier(t, actual[i], expected[i])
	}
}

func testCallExpression(
	t *testing.T,
	actual *ast.CallExpression,
	expected *ast.CallExpression,
) {
	testExpression(t, actual.Function, expected.Function)
	testExpressions(t, actual.Arguments, expected.Arguments)
}

func testExpressions(
	t *testing.T,
	actual []ast.Expression,
	expected []ast.Expression,
) {
	if len(actual) != len(expected) {
		t.Fatalf(
			"number of expressions mismatch. got=%v, expected=%v",
			len(actual),
			len(expected),
		)
	}
	for i := 0; i < len(actual); i++ {
		testExpression(t, actual[i], expected[i])
	}
}

func testReturnStatement(
	t *testing.T,
	actual *ast.ReturnStatement,
	expected *ast.ReturnStatement,
) {
	testExpression(t, actual.Value, expected.Value)
}

func testLetStatement(
	t *testing.T,
	actual *ast.LetStatement,
	expected *ast.LetStatement,
) {
	testIdentifier(t, actual.Name, expected.Name)
	testExpression(t, actual.Value, expected.Value)
}

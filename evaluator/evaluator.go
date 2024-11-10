package evaluator

import (
	"log"
	"reflect"

	"github.com/vincentlabelle/monkey/ast"
	"github.com/vincentlabelle/monkey/object"
)

func Eval(program *ast.Program, env *object.Environment) object.Object {
	var obj object.Object
	for _, statement := range program.Statements {
		switch s := statement.(type) {
		case *ast.ReturnStatement:
			return evalExpression(s.Value, env)
		case *ast.ExpressionStatement:
			obj = evalExpression(s.Expression, env)
			if rv, ok := obj.(*object.ReturnValue); ok {
				return rv.Value
			}
		case *ast.LetStatement:
			obj = evalLetStatement(s, env)
		default:
			message := "cannot evaluate program; unexpected statement type"
			log.Fatal(message)
		}
	}
	return obj
}

func evalExpression(
	expression ast.Expression,
	env *object.Environment,
) object.Object {
	var obj object.Object
	switch e := expression.(type) {
	case *ast.IntegerLiteral:
		obj = nativeToInteger(e.Value)
	case *ast.BooleanLiteral:
		obj = nativeToBoolean(e.Value)
	case *ast.StringLiteral:
		obj = nativeToString(e.Value)
	case *ast.Identifier:
		obj = evalIdentifier(e, env)
	case *ast.PrefixExpression:
		obj = evalPrefixExpression(e, env)
	case *ast.InfixExpression:
		obj = evalInfixExpression(e, env)
	case *ast.IfExpression:
		obj = evalIfExpression(e, env)
	case *ast.FunctionLiteral:
		obj = evalFunctionLiteral(e, env)
	case *ast.CallExpression:
		obj = evalCallExpression(e, env)
	case *ast.ArrayLiteral:
		obj = evalArrayLiteral(e, env)
	case *ast.IndexExpression:
		obj = evalIndexExpression(e, env)
	default:
		message := "cannot evaluate program; unexpected expression type"
		log.Fatal(message)
	}
	return obj
}

func nativeToInteger(native int) *object.Integer {
	return &object.Integer{Value: native}
}

func nativeToBoolean(native bool) *object.Boolean {
	if native {
		return object.TRUE
	}
	return object.FALSE
}

func nativeToString(native string) *object.String {
	return &object.String{Value: native}
}

func evalIdentifier(
	expression *ast.Identifier,
	env *object.Environment,
) object.Object {
	obj, ok := env.Get(expression.Value)
	if !ok {
		message := "cannot evaluate program; encountered undefined identifier"
		log.Fatal(message)
	}
	return obj
}

func evalPrefixExpression(
	expression *ast.PrefixExpression,
	env *object.Environment,
) object.Object {
	right := evalExpression(expression.Right, env)
	return innerEvalPrefixExpression(expression.Operator, right)
}

func innerEvalPrefixExpression(
	operator string,
	right object.Object,
) object.Object {
	var obj object.Object
	switch operator {
	case "-":
		obj = evalMinusPrefixExpression(right)
	case "!":
		obj = evalBangPrefixExpression(right)
	default:
		message := "cannot evaluate program; unexpected prefix operator"
		log.Fatal(message)
	}
	return obj
}

func evalMinusPrefixExpression(obj object.Object) *object.Integer {
	switch o := obj.(type) {
	case *object.Integer:
		return nativeToInteger(-o.Value)
	default:
		message := "cannot evaluate program; unexpected operand for - prefix"
		log.Fatal(message)
	}
	return nil
}

func evalBangPrefixExpression(obj object.Object) *object.Boolean {
	var new_ *object.Boolean
	switch o := obj.(type) {
	case *object.Boolean:
		new_ = nativeToBoolean(!o.Value)
	case *object.Null:
		new_ = object.TRUE
	default:
		new_ = object.FALSE
	}
	return new_
}

func evalInfixExpression(
	expression *ast.InfixExpression,
	env *object.Environment,
) object.Object {
	left := evalExpression(expression.Left, env)
	right := evalExpression(expression.Right, env)
	return innerEvalInfixExpression(left, expression.Operator, right)
}

func innerEvalInfixExpression(
	left object.Object,
	operator string,
	right object.Object,
) object.Object {
	var obj object.Object
	if isPtrToType(left, "Integer") && isPtrToType(right, "Integer") {
		l, r := left.(*object.Integer), right.(*object.Integer)
		obj = evalIntegerInfixExpression(l, operator, r)
	} else if isPtrToType(left, "String") && isPtrToType(right, "String") {
		l, r := left.(*object.String), right.(*object.String)
		obj = evalStringInfixExpression(l, operator, r)
	} else if operator == "==" {
		obj = nativeToBoolean(left == right)
	} else if operator == "!=" {
		obj = nativeToBoolean(left != right)
	} else if reflect.TypeOf(left) != reflect.TypeOf(right) {
		message := `cannot parse program; 
			operands with operator %v aren't of the same type`
		log.Fatalf(message, operator)
	} else {
		message := `cannot parse program; 
			unexpected operator for infix expression`
		log.Fatal(message)
	}
	return obj
}

func isPtrToType(obj object.Object, name string) bool {
	type_ := reflect.TypeOf(obj)
	return type_.Kind() == reflect.Pointer && type_.Elem().Name() == name
}

func evalIntegerInfixExpression(
	left *object.Integer,
	operator string,
	right *object.Integer,
) object.Object {
	var obj object.Object
	switch operator {
	case "+":
		obj = nativeToInteger(left.Value + right.Value)
	case "-":
		obj = nativeToInteger(left.Value - right.Value)
	case "*":
		obj = nativeToInteger(left.Value * right.Value)
	case "/":
		obj = nativeToInteger(left.Value / right.Value)
	case "<":
		obj = nativeToBoolean(left.Value < right.Value)
	case ">":
		obj = nativeToBoolean(left.Value > right.Value)
	case "==":
		obj = nativeToBoolean(left.Value == right.Value)
	case "!=":
		obj = nativeToBoolean(left.Value != right.Value)
	default:
		message := `cannot parse program; 
			unexpected operator for infix expression`
		log.Fatal(message)
	}
	return obj
}

func evalStringInfixExpression(
	left *object.String,
	operator string,
	right *object.String,
) *object.String {
	var obj *object.String
	switch operator {
	case "+":
		obj = nativeToString(left.Value + right.Value)
	default:
		message := `cannot parse program; 
			unexpected operator for infix expression`
		log.Fatal(message)
	}
	return obj
}

func evalIfExpression(
	expression *ast.IfExpression,
	env *object.Environment,
) object.Object {
	condition := evalIfExpressionCondition(expression.Condition, env)
	return evalIfExpressionBlock(condition, expression, env)
}

func evalIfExpressionCondition(
	expression ast.Expression,
	env *object.Environment,
) *object.Boolean {
	obj := evalExpression(expression, env)
	switch o := obj.(type) {
	case *object.Boolean:
		return o
	case *object.Null:
		return object.FALSE
	default:
		return object.TRUE
	}
}

func evalIfExpressionBlock(
	condition *object.Boolean,
	expression *ast.IfExpression,
	env *object.Environment,
) object.Object {
	if condition.Value {
		return evalBlockStatements(expression.Consequence, env)
	}
	if expression.Alternative != nil {
		return evalBlockStatements(expression.Alternative, env)
	}
	return object.NULL
}

func evalBlockStatements(
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	var obj object.Object
	for _, statement := range block.Statements {
		switch s := statement.(type) {
		case *ast.ReturnStatement:
			obj = evalExpression(s.Value, env)
			return &object.ReturnValue{Value: obj}
		case *ast.ExpressionStatement:
			obj = evalExpression(s.Expression, env)
			if rv, ok := obj.(*object.ReturnValue); ok {
				return rv
			}
		case *ast.LetStatement:
			obj = evalLetStatement(s, env)
		default:
			message := "cannot evaluate program; unexpected statement type"
			log.Fatal(message)
		}
	}
	return obj
}

func evalLetStatement(
	statement *ast.LetStatement,
	env *object.Environment,
) object.Object {
	env.Set(
		statement.Name.Value,
		evalExpression(statement.Value, env),
	)
	return nil
}

func evalFunctionLiteral(
	expression *ast.FunctionLiteral,
	env *object.Environment,
) *object.Function {
	return &object.Function{
		Parameters: expression.Parameters,
		Body:       expression.Body,
		Env:        env,
	}
}

func evalCallExpression(
	expression *ast.CallExpression,
	env *object.Environment,
) object.Object {
	function := evalExpression(expression.Function, env)
	arguments := evalExpressions(expression.Arguments, env)
	return innerEvalCallExpression(function, arguments)
}

func evalExpressions(
	expressions []ast.Expression,
	env *object.Environment,
) []object.Object {
	objs := []object.Object{}
	for _, expression := range expressions {
		obj := evalExpression(expression, env)
		objs = append(objs, obj)
	}
	return objs
}

func innerEvalCallExpression(
	function object.Object,
	arguments []object.Object,
) object.Object {
	var obj object.Object
	switch f := function.(type) {
	case *object.Function:
		obj = evalCallExpressionFunction(f, arguments)
	case *object.Builtin:
		obj = f.Fn(arguments...)
	default:
		message := `cannot evaluate program; 
			unexpected call expression function`
		log.Fatal(message)
	}
	return obj
}

func evalCallExpressionFunction(
	function *object.Function,
	arguments []object.Object,
) object.Object {
	inner := newInnerEnvironment(function, arguments)
	obj := evalBlockStatements(function.Body, inner)
	return unwrap(obj)
}

func newInnerEnvironment(
	function *object.Function,
	arguments []object.Object,
) *object.Environment {
	if len(arguments) != len(function.Parameters) {
		message := `cannot evaluate program; 
			incorrect number of arguments in call expression`
		log.Fatal(message)
	}
	return innerNewInnerEnvironment(function, arguments)
}

func innerNewInnerEnvironment(
	function *object.Function,
	arguments []object.Object,
) *object.Environment {
	inner := object.NewInnerEnvironment(function.Env)
	for i := 0; i < len(arguments); i++ {
		inner.Set(
			function.Parameters[i].Value,
			arguments[i],
		)
	}
	return inner
}

func unwrap(obj object.Object) object.Object {
	if rv, ok := obj.(*object.ReturnValue); ok {
		return rv.Value
	}
	return obj
}

func evalArrayLiteral(
	expression *ast.ArrayLiteral,
	env *object.Environment,
) *object.Array {
	elements := evalExpressions(expression.Elements, env)
	return &object.Array{Elements: elements}
}

func evalIndexExpression(
	expression *ast.IndexExpression,
	env *object.Environment,
) object.Object {
	left := evalIndexExpressionLeft(expression.Left, env)
	index := evalIndexExpressionIndex(expression.Index, env)
	return innerEvalIndexExpression(left, index)
}

func evalIndexExpressionLeft(
	expression ast.Expression,
	env *object.Environment,
) *object.Array {
	array := evalExpression(expression, env)
	a, ok := array.(*object.Array)
	if !ok {
		message := `cannot evaluate program; 
			unexpected left in index expression`
		log.Fatal(message)
	}
	return a
}

func evalIndexExpressionIndex(
	expression ast.Expression,
	env *object.Environment,
) *object.Integer {
	index := evalExpression(expression, env)
	i, ok := index.(*object.Integer)
	if !ok {
		message := `cannot evaluate program; 
			unexpected index in index expression`
		log.Fatal(message)
	}
	return i
}

func innerEvalIndexExpression(
	left *object.Array,
	index *object.Integer,
) object.Object {
	if index.Value >= len(left.Elements) || index.Value < 0 {
		return object.NULL
	}
	return left.Elements[index.Value]
}

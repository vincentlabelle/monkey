package evaluator

import (
	"log"

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
		obj = object.NativeToInteger(e.Value)
	case *ast.BooleanLiteral:
		obj = object.NativeToBoolean(e.Value)
	case *ast.StringLiteral:
		obj = object.NativeToString(e.Value)
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
	case *ast.HashLiteral:
		obj = evalHashLiteral(e, env)
	default:
		message := "cannot evaluate program; unexpected expression type"
		log.Fatal(message)
	}
	return obj
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
	return EvalPrefix(expression.Operator, right)
}

func evalInfixExpression(
	expression *ast.InfixExpression,
	env *object.Environment,
) object.Object {
	left := evalExpression(expression.Left, env)
	right := evalExpression(expression.Right, env)
	return EvalInfix(left, expression.Operator, right)
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
	return EvalTruthy(obj)
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
		message := "cannot evaluate program; " +
			"unexpected call expression function"
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
		message := "cannot evaluate program; " +
			"incorrect number of arguments in call expression"
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
	left := evalExpression(expression.Left, env)
	index := evalExpression(expression.Index, env)
	return EvalIndex(left, index)
}

func evalHashLiteral(
	expression *ast.HashLiteral,
	env *object.Environment,
) *object.Hash {
	pairs := map[object.HashKey]object.HashPair{}
	for key, value := range expression.Pairs {
		k := evalHashLiteralKey(key.Expression, env)
		v := evalExpression(value, env)
		pairs[k.HashKey()] = object.HashPair{Key: k, Value: v}
	}
	return &object.Hash{Pairs: pairs}
}

func evalHashLiteralKey(
	expression ast.Expression,
	env *object.Environment,
) object.Hashable {
	key := evalExpression(expression, env)
	return object.CastToHashable(key)
}

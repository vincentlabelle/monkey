package evaluator

import (
	"log"
	"reflect"

	"github.com/vincentlabelle/monkey/object"
)

func EvalPrefix(
	operator string,
	right object.Object,
) object.Object {
	var obj object.Object
	switch operator {
	case "-":
		obj = evalMinusPrefix(right)
	case "!":
		obj = evalBangPrefix(right)
	default:
		message := "cannot evaluate program; unexpected prefix operator"
		log.Fatal(message)
	}
	return obj
}

func evalMinusPrefix(obj object.Object) *object.Integer {
	switch o := obj.(type) {
	case *object.Integer:
		return object.NativeToInteger(-o.Value)
	default:
		message := "cannot evaluate program; unexpected operand for - prefix"
		log.Fatal(message)
	}
	return nil
}

func evalBangPrefix(obj object.Object) *object.Boolean {
	var new_ *object.Boolean
	switch o := obj.(type) {
	case *object.Boolean:
		new_ = object.NativeToBoolean(!o.Value)
	case *object.Null:
		new_ = object.TRUE
	default:
		new_ = object.FALSE
	}
	return new_
}

func EvalInfix(
	left object.Object,
	operator string,
	right object.Object,
) object.Object {
	var obj object.Object
	if object.IsType(left, "Integer") && object.IsType(right, "Integer") {
		l, r := left.(*object.Integer), right.(*object.Integer)
		obj = evalIntegerInfix(l, operator, r)
	} else if object.IsType(left, "String") && object.IsType(right, "String") {
		l, r := left.(*object.String), right.(*object.String)
		obj = evalStringInfix(l, operator, r)
	} else if operator == "==" {
		obj = object.NativeToBoolean(left == right)
	} else if operator == "!=" {
		obj = object.NativeToBoolean(left != right)
	} else if reflect.TypeOf(left) != reflect.TypeOf(right) {
		message := "cannot evaluate program; " +
			"operands with operator %v aren't of the same type"
		log.Fatalf(message, operator)
	} else {
		message := "cannot evaluate program; " +
			"unexpected operator for infix expression"
		log.Fatal(message)
	}
	return obj
}

func evalIntegerInfix(
	left *object.Integer,
	operator string,
	right *object.Integer,
) object.Object {
	var obj object.Object
	switch operator {
	case "+":
		obj = object.NativeToInteger(left.Value + right.Value)
	case "-":
		obj = object.NativeToInteger(left.Value - right.Value)
	case "*":
		obj = object.NativeToInteger(left.Value * right.Value)
	case "/":
		obj = object.NativeToInteger(left.Value / right.Value)
	case "<":
		obj = object.NativeToBoolean(left.Value < right.Value)
	case ">":
		obj = object.NativeToBoolean(left.Value > right.Value)
	case "==":
		obj = object.NativeToBoolean(left.Value == right.Value)
	case "!=":
		obj = object.NativeToBoolean(left.Value != right.Value)
	default:
		message := "cannot evaluate program; " +
			"unexpected operator for infix expression"
		log.Fatal(message)
	}
	return obj
}

func evalStringInfix(
	left *object.String,
	operator string,
	right *object.String,
) *object.String {
	var obj *object.String
	switch operator {
	case "+":
		obj = object.NativeToString(left.Value + right.Value)
	default:
		message := "cannot evaluate program; " +
			"unexpected operator for infix expression"
		log.Fatal(message)
	}
	return obj
}

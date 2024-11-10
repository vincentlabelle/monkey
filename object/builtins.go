package object

import (
	"fmt"
	"log"
)

func len_(args ...Object) Object {
	if len(args) != 1 {
		message := "cannot call built-in; one argument is expected"
		log.Fatal(message)
	}
	var obj *Integer
	switch a := args[0].(type) {
	case *String:
		obj = &Integer{Value: len(a.Value)}
	case *Array:
		obj = &Integer{Value: len(a.Elements)}
	default:
		message := "cannot call built-in; invalid argument"
		log.Fatal(message)
	}
	return obj
}

func first(args ...Object) Object {
	array := getUniqueArray(args)
	return innerFirst(array)
}

func getUniqueArray(args []Object) *Array {
	if len(args) != 1 {
		message := "cannot call built-in; one argument is expected"
		log.Fatal(message)
	}
	return getArray(args[0])
}

func getArray(arg Object) *Array {
	array, ok := arg.(*Array)
	if !ok {
		message := "cannot call built-in; argument must be an array"
		log.Fatal(message)
	}
	return array
}

func innerFirst(array *Array) Object {
	if len(array.Elements) == 0 {
		return NULL
	}
	return array.Elements[0]
}

func last(args ...Object) Object {
	array := getUniqueArray(args)
	return innerLast(array)
}

func innerLast(array *Array) Object {
	if len(array.Elements) == 0 {
		return NULL
	}
	return array.Elements[len(array.Elements)-1]
}

func rest(args ...Object) Object {
	array := getUniqueArray(args)
	return innerRest(array)
}

func innerRest(array *Array) Object {
	if len(array.Elements) == 0 {
		return NULL
	}
	elements := make([]Object, len(array.Elements)-1)
	copy(elements, array.Elements[1:])
	return &Array{Elements: elements}
}

func push(args ...Object) Object {
	if len(args) != 2 {
		message := "cannot call built-in; two arguments are expected"
		log.Fatal(message)
	}
	array := getArray(args[0])
	append_ := args[1]
	return innerPush(array, append_)
}

func innerPush(array *Array, append_ Object) Object {
	elements := make([]Object, len(array.Elements)+1)
	copy(elements, array.Elements)
	elements[len(elements)-1] = append_
	return &Array{Elements: elements}
}

func puts(args ...Object) Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NULL
}

package object

import (
	"fmt"

	"github.com/vincentlabelle/monkey/ast"
)

type Object interface {
	Inspect() string
}

type Integer struct {
	Value int
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%v", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}

type String struct {
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	return "fn(...) {...}"
}

type Builtin struct {
	Fn func(...Object) Object
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

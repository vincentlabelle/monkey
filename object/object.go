package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/vincentlabelle/monkey/ast"
	"github.com/vincentlabelle/monkey/code"
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

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: "Integer", Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	var value uint64 = 0
	if b.Value {
		value = 1
	}
	return HashKey{Type: "Boolean", Value: value}
}

type String struct {
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: "String", Value: h.Sum64()}
}

type Array struct {
	Elements []Object
}

func (a *Array) Inspect() string {
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Inspect() string {
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

type HashPair struct {
	Key   Object
	Value Object
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

type CompiledFunction struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFunction) Inspect() string {
	return "fn(...) {...}"
}

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) Inspect() string {
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

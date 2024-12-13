package vm

import (
	"github.com/vincentlabelle/monkey/code"
	"github.com/vincentlabelle/monkey/object"
)

type Frame struct {
	Closure        *object.Closure
	InsIndex       int
	BaseStackIndex int
}

func (f *Frame) Instructions() code.Instructions {
	return f.Closure.Fn.Instructions
}

func (f *Frame) NumLocals() int {
	return f.Closure.Fn.NumLocals
}

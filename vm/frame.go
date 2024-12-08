package vm

import (
	"github.com/vincentlabelle/monkey/object"
)

type Frame struct {
	Fn             *object.CompiledFunction
	InsIndex       int
	BaseStackIndex int
}

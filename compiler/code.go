package compiler

import (
	"github.com/vincentlabelle/monkey/code"
	"github.com/vincentlabelle/monkey/object"
)

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

package vm

import (
	"log"
	"slices"

	"github.com/vincentlabelle/monkey/code"
	"github.com/vincentlabelle/monkey/compiler"
	"github.com/vincentlabelle/monkey/evaluator"
	"github.com/vincentlabelle/monkey/object"
)

const (
	StackSize   = 2048
	GlobalsSize = 65536
)

type VM struct {
	code    *compiler.Bytecode
	globals []object.Object
	stack   []object.Object
	sp      int
}

func New(code *compiler.Bytecode) *VM {
	return &VM{
		code:    code,
		globals: make([]object.Object, GlobalsSize),
		stack:   make([]object.Object, StackSize),
	}
}

func (vm *VM) LastPopped() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() {
	ip := 0
	for ip < len(vm.code.Instructions) {
		remain := vm.code.Instructions[ip:]
		op, operands, width := code.Unmake(remain)
		ip += width
		switch op {
		case code.OpJump:
			ip = vm.getOperand(operands)
		case code.OpJumpIf:
			obj := vm.pop()
			if !vm.isTruthy(obj) {
				ip = vm.getOperand(operands)
			}
		default:
			vm.run(op, operands)
		}
	}
}

func (vm *VM) getOperand(operands []int) int {
	if len(operands) == 0 {
		message := "cannot run virtual machine; " +
			"unexpected number of operands encountered"
		log.Fatal(message)
	}
	return operands[0]
}

func (vm *VM) isTruthy(obj object.Object) bool {
	b := evaluator.EvalTruthy(obj)
	return b.Value
}

func (vm *VM) run(op code.Opcode, operands []int) {
	switch op {
	case code.OpConstant:
		vm.runOpConstant(operands)
	case code.OpTrue:
		vm.runOpTrue()
	case code.OpFalse:
		vm.runOpFalse()
	case code.OpNull:
		vm.runOpNull()
	case code.OpArray:
		vm.runOpArray(operands)
	case code.OpHash:
		vm.runOpHash(operands)
	case code.OpAdd,
		code.OpSub,
		code.OpMul,
		code.OpDiv,
		code.OpEqual,
		code.OpNotEqual,
		code.OpGreaterThan,
		code.OpLowerThan:
		vm.runInfixOperation(op)
	case code.OpBang, code.OpMinus:
		vm.runPrefixOperation(op)
	case code.OpIndex:
		vm.runOpIndex()
	case code.OpSetGlobal:
		vm.runOpSetGlobal(operands)
	case code.OpGetGlobal:
		vm.runOpGetGlobal(operands)
	case code.OpPop:
		vm.pop()
	default:
		message := "cannot run virtual machine; " +
			"unexpected Opcode encountered"
		log.Fatal(message)
	}
}

func (vm *VM) runOpConstant(operands []int) {
	operand := vm.getOperand(operands)
	constant := vm.code.Constants[operand]
	vm.push(constant)
}

func (vm *VM) push(obj object.Object) {
	if vm.sp >= StackSize {
		message := "cannot run virtual machine; stack overflow"
		log.Fatal(message)
	}
	vm.stack[vm.sp] = obj
	vm.sp++
}

func (vm *VM) runOpTrue() {
	vm.push(object.TRUE)
}

func (vm *VM) runOpFalse() {
	vm.push(object.FALSE)
}

func (vm *VM) runOpNull() {
	vm.push(object.NULL)
}

func (vm *VM) runOpArray(operands []int) {
	operand := vm.getOperand(operands)
	obj := &object.Array{Elements: vm.popNReverse(operand)}
	vm.push(obj)
}

func (vm *VM) popNReverse(n int) []object.Object {
	objs := vm.popN(n)
	slices.Reverse(objs)
	return objs
}

func (vm *VM) popN(n int) []object.Object {
	elements := []object.Object{}
	for i := 0; i < n; i++ {
		elements = append(elements, vm.pop())
	}
	return elements
}

func (vm *VM) runOpHash(operands []int) {
	operand := vm.getOperand(operands)
	obj := vm.innerRunOpHash(operand)
	vm.push(obj)
}

func (vm *VM) innerRunOpHash(operand int) *object.Hash {
	pairs := map[object.HashKey]object.HashPair{}
	for i := 0; i < operand; i++ {
		v, k := vm.pop(), vm.popHashable()
		pairs[k.HashKey()] = object.HashPair{Key: k, Value: v}
	}
	return &object.Hash{Pairs: pairs}
}

func (vm *VM) popHashable() object.Hashable {
	obj := vm.pop()
	return object.CastToHashable(obj)
}

func (vm *VM) runInfixOperation(op code.Opcode) {
	right, left := vm.pop(), vm.pop()
	operator := vm.getInfixOperator(op)
	obj := evaluator.EvalInfix(left, operator, right)
	vm.push(obj)
}

func (vm *VM) pop() object.Object {
	vm.sp--
	return vm.LastPopped()
}

func (vm *VM) getInfixOperator(op code.Opcode) string {
	operator, ok := code.InfixOperatorReverse[op]
	if !ok {
		message := "cannot run virtual machine; " +
			"unexpected Opcode encountered has infix operator"
		log.Fatal(message)
	}
	return operator
}

func (vm *VM) runPrefixOperation(op code.Opcode) {
	right := vm.pop()
	operator := vm.getPrefixOperator(op)
	obj := evaluator.EvalPrefix(operator, right)
	vm.push(obj)
}

func (vm *VM) getPrefixOperator(op code.Opcode) string {
	operator, ok := code.PrefixOperatorReverse[op]
	if !ok {
		message := "cannot run virtual machine; " +
			"unexpected Opcode encountered has prefix operator"
		log.Fatal(message)
	}
	return operator
}

func (vm *VM) runOpIndex() {
	index, left := vm.pop(), vm.pop()
	obj := evaluator.EvalIndex(left, index)
	vm.push(obj)
}

func (vm *VM) runOpSetGlobal(operands []int) {
	operand := vm.getOperand(operands)
	if operand > GlobalsSize {
		message := "cannot run virtual machine; globals overflow"
		log.Fatal(message)
	}
	vm.globals[operand] = vm.pop()
}

func (vm *VM) runOpGetGlobal(operands []int) {
	operand := vm.getOperand(operands)
	vm.push(vm.globals[operand])
}

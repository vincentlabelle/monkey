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
	GlobalsSize = 65536
	StackSize   = 2048
	FramesSize  = 1024
)

type VM struct {
	globals     []object.Object
	stack       []object.Object
	stackIndex  int
	frames      []*Frame
	framesIndex int
	constants   []object.Object
}

func New(code *compiler.Bytecode) *VM {
	vm := &VM{
		globals:   make([]object.Object, GlobalsSize),
		stack:     make([]object.Object, StackSize),
		frames:    make([]*Frame, FramesSize),
		constants: code.Constants,
	}
	vm.pushInitialFrame(code.Instructions)
	return vm
}

func (vm *VM) pushInitialFrame(instructions code.Instructions) {
	frame := &Frame{
		Closure: &object.Closure{
			Fn: &object.CompiledFunction{
				Instructions: instructions,
			},
		},
	}
	vm.pushFrame(frame)
}

func (vm *VM) LastPopped() object.Object {
	return vm.stack[vm.stackIndex]
}

func (vm *VM) Run() {
	frame := vm.currentFrame()
	for frame.InsIndex < len(frame.Instructions()) {
		remain := frame.Instructions()[frame.InsIndex:]
		op, operands, width := code.Unmake(remain)
		frame.InsIndex += width // Before run, because of jumps!!
		vm.run(op, operands)
		frame = vm.currentFrame()
	}
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
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
	case code.OpSetLocal:
		vm.runOpSetLocal(operands)
	case code.OpGetLocal:
		vm.runOpGetLocal(operands)
	case code.OpGetBuiltin:
		vm.runOpGetBuiltin(operands)
	case code.OpGetFree:
		vm.runOpGetFree(operands)
	case code.OpCall:
		vm.runOpCall(operands)
	case code.OpClosure:
		vm.runOpClosure(operands)
	case code.OpCurrentClosure:
		vm.runOpCurrentClosure()
	case code.OpReturnValue:
		vm.runOpReturnValue()
	case code.OpReturn:
		vm.runOpReturn()
	case code.OpJump:
		vm.runOpJump(operands)
	case code.OpJumpIf:
		vm.runOpJumpIf(operands)
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
	constant := vm.constants[operand]
	vm.push(constant)
}

func (vm *VM) push(obj object.Object) {
	if vm.stackIndex >= StackSize {
		message := "cannot run virtual machine; stack overflow"
		log.Fatal(message)
	}
	vm.stack[vm.stackIndex] = obj
	vm.stackIndex++
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
	vm.stackIndex--
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

func (vm *VM) runOpSetLocal(operands []int) {
	operand := vm.getOperand(operands)
	obj := vm.pop()
	vm.setLocal(obj, operand)
}

func (vm *VM) setLocal(obj object.Object, operand int) {
	frame := vm.currentFrame()
	vm.stack[frame.BaseStackIndex+operand] = obj
}

func (vm *VM) runOpGetLocal(operands []int) {
	operand := vm.getOperand(operands)
	obj := vm.getLocal(operand)
	vm.push(obj)
}

func (vm *VM) getLocal(operand int) object.Object {
	frame := vm.currentFrame()
	return vm.stack[frame.BaseStackIndex+operand]
}

func (vm *VM) runOpGetBuiltin(operands []int) {
	operand := vm.getOperand(operands)
	b := object.Builtins[operand]
	vm.push(b.Builtin)
}

func (vm *VM) runOpGetFree(operands []int) {
	operand := vm.getOperand(operands)
	frame := vm.currentFrame()
	obj := frame.Closure.Free[operand]
	vm.push(obj)
}

func (vm *VM) runOpCall(operands []int) {
	operand := vm.getOperand(operands)
	fn := vm.getFunction(operand)
	vm.dispatchCall(fn, operand)
}

func (vm *VM) getFunction(operand int) object.Object {
	return vm.stack[vm.stackIndex-operand-1]
}

func (vm *VM) dispatchCall(fn object.Object, operand int) {
	switch f := fn.(type) {
	case *object.Closure:
		vm.runClosure(f, operand)
	case *object.Builtin:
		vm.runBuiltinFunction(f, operand)
	default:
		message := "cannot run virtual machine; " +
			"unexpected object encountered has function in function call"
		log.Fatal(message)
	}
}

func (vm *VM) runClosure(closure *object.Closure, operand int) {
	vm.validateArguments(closure, operand)
	frame := &Frame{Closure: closure, BaseStackIndex: vm.stackIndex - operand}
	vm.pushFrame(frame)
}

func (vm *VM) validateArguments(closure *object.Closure, operand int) {
	if closure.Fn.NumParameters != operand {
		message := "cannot run virtual machine; " +
			"unexpected number of arguments in call to function"
		log.Fatal(message)
	}
}

func (vm *VM) pushFrame(frame *Frame) {
	vm.frames[vm.framesIndex] = frame
	vm.framesIndex++
	vm.stackIndex += frame.NumLocals()
}

func (vm *VM) runBuiltinFunction(fn *object.Builtin, operand int) {
	arguments := vm.popNReverse(operand)
	vm.pop() // Pop builtin!
	obj := fn.Fn(arguments...)
	vm.push(obj)
}

func (vm *VM) runOpClosure(operands []int) {
	index, count := vm.getTwoOperands(operands)
	obj := vm.getClosure(index, count)
	vm.push(obj)
}

func (vm *VM) getTwoOperands(operands []int) (int, int) {
	if len(operands) < 2 {
		message := "cannot run virtual machine; " +
			"unexpected number of operands encountered"
		log.Fatal(message)
	}
	return operands[0], operands[1]
}

func (vm *VM) getClosure(index int, count int) *object.Closure {
	obj := vm.getCompiledFunction(index)
	free := vm.popNReverse(count)
	return &object.Closure{Fn: obj, Free: free}
}

func (vm *VM) getCompiledFunction(index int) *object.CompiledFunction {
	obj, ok := vm.constants[index].(*object.CompiledFunction)
	if !ok {
		message := "cannot run virtual machine; " +
			"unexpected constant when expecting function"
		log.Fatal(message)
	}
	return obj
}

func (vm *VM) runOpCurrentClosure() {
	frame := vm.currentFrame()
	vm.push(frame.Closure)
}

func (vm *VM) runOpReturnValue() {
	obj := vm.pop()
	vm.popFrame()
	vm.push(obj)
}

func (vm *VM) popFrame() {
	frame := vm.currentFrame()
	vm.stackIndex = frame.BaseStackIndex - 1
	vm.framesIndex--
}

func (vm *VM) runOpReturn() {
	vm.popFrame()
	vm.push(object.NULL)
}

func (vm *VM) runOpJump(operands []int) {
	frame := vm.currentFrame()
	frame.InsIndex = vm.getOperand(operands)
}

func (vm *VM) runOpJumpIf(operands []int) {
	obj := vm.pop()
	if !vm.isTruthy(obj) {
		vm.runOpJump(operands)
	}
}

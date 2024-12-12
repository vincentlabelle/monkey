package code

import "log"

type Definition struct {
	Name          string
	Op            Opcode
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:    {"OpConstant", OpConstant, []int{2}},
	OpTrue:        {"OpTrue", OpTrue, []int{}},
	OpFalse:       {"OpFalse", OpFalse, []int{}},
	OpNull:        {"OpNull", OpNull, []int{}},
	OpArray:       {"OpArray", OpArray, []int{2}},
	OpHash:        {"OpHash", OpHash, []int{2}},
	OpAdd:         {"OpAdd", OpAdd, []int{}},
	OpSub:         {"OpSub", OpSub, []int{}},
	OpMul:         {"OpMul", OpMul, []int{}},
	OpDiv:         {"OpDiv", OpDiv, []int{}},
	OpEqual:       {"OpEqual", OpEqual, []int{}},
	OpNotEqual:    {"OpNotEqual", OpNotEqual, []int{}},
	OpGreaterThan: {"OpGreaterThan", OpGreaterThan, []int{}},
	OpLowerThan:   {"OpLowerThan", OpLowerThan, []int{}},
	OpMinus:       {"OpMinus", OpMinus, []int{}},
	OpBang:        {"OpBang", OpBang, []int{}},
	OpIndex:       {"OpIndex", OpIndex, []int{}},
	OpPop:         {"OpPop", OpPop, []int{}},
	OpJump:        {"OpJump", OpJump, []int{2}},
	OpJumpIf:      {"OpJumpIf", OpJumpIf, []int{2}},
	OpSetGlobal:   {"OpSetGlobal", OpSetGlobal, []int{2}},
	OpGetGlobal:   {"OpGetGlobal", OpGetGlobal, []int{2}},
	OpSetLocal:    {"OpSetLocal", OpSetLocal, []int{2}},
	OpGetLocal:    {"OpGetLocal", OpGetLocal, []int{2}},
	OpGetBuiltin:  {"OpGetBuiltin", OpGetBuiltin, []int{1}},
	OpCall:        {"OpCall", OpCall, []int{1}},
	OpReturnValue: {"OpReturnValue", OpReturnValue, []int{}},
	OpReturn:      {"OpReturn", OpReturn, []int{}},
}

func Lookup(op byte) *Definition {
	def, ok := definitions[Opcode(op)]
	if !ok {
		message := "cannot lookup opcode; opcode is undefined"
		log.Fatal(message)
	}
	return def
}

func getWidth(def *Definition) int {
	width := 1
	for _, w := range def.OperandWidths {
		width += w
	}
	return width
}

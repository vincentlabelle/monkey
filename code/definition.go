package code

import "log"

type Definition struct {
	Name          string
	Op            Opcode
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:    {"OpConstant", OpConstant, []int{2}},
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
	OpTrue:        {"OpTrue", OpTrue, []int{}},
	OpFalse:       {"OpFalse", OpFalse, []int{}},
	OpPop:         {"OpPop", OpPop, []int{}},
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

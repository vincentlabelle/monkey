package code

import "fmt"

type Instructions []byte

func (ins Instructions) String() string {
	s, i := "", 0
	for i < len(ins) {
		op, operands, width := Unmake(ins[i:])
		s += castInstruction(i, op, operands)
		i += width
	}
	return s
}

func castInstruction(i int, op Opcode, operands []int) string {
	return fmt.Sprintf(
		"%04d %v%s\n",
		i,
		definitions[op].Name,
		castOperands(operands),
	)
}

func castOperands(operands []int) string {
	s := ""
	for _, operand := range operands {
		s += fmt.Sprintf(" %v", operand)
	}
	return s
}

func Concatenate(instructions []Instructions) Instructions {
	concatenated := Instructions{}
	for _, ins := range instructions {
		concatenated = append(concatenated, ins...)
	}
	return concatenated
}

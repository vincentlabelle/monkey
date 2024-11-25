package code

import "fmt"

type Instructions []byte

func (instructions Instructions) String() string {
	s, i := "", 0
	for i < len(instructions) {
		op, operands, width := Unmake(instructions[i:])
		s += cast(i, op, operands)
		i += width
	}
	return s
}

func cast(i int, op Opcode, operands []int) string {
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

func Concatenate(pieces []Instructions) Instructions {
	concatenated := Instructions{}
	for _, instructions := range pieces {
		concatenated = append(concatenated, instructions...)
	}
	return concatenated
}

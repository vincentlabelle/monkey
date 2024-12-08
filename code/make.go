package code

import (
	"encoding/binary"
	"log"
)

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}
	return makeInstruction(def, operands)
}

func makeInstruction(def *Definition, operands []int) []byte {
	instruction := initializeInstruction(def)
	return addOperands(instruction, def, operands)
}

func initializeInstruction(def *Definition) []byte {
	width := getWidth(def)
	instruction := make([]byte, width)
	instruction[0] = byte(def.Op)
	return instruction
}

func addOperands(
	instruction []byte,
	def *Definition,
	operands []int,
) []byte {
	if len(def.OperandWidths) != len(operands) {
		message := "cannot make instruction; unexpected number of operands"
		log.Fatal(message)
	}
	return innerAddOperands(instruction, def, operands)
}

func innerAddOperands(
	instruction []byte,
	def *Definition,
	operands []int,
) []byte {
	remain := instruction[1:]
	for i := 0; i < len(operands); i++ {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(
				remain,
				uint16(operands[i]),
			)
		case 1:
			remain[0] = byte(operands[i])
		default:
			message := "cannot make instruction; unexpected operand width"
			log.Fatal(message)
		}
		remain = remain[width:]
	}
	return instruction
}

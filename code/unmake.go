package code

import (
	"encoding/binary"
	"log"
)

func Unmake(instruction []byte) (Opcode, []int, int) {
	if len(instruction) == 0 {
		message := "cannot unmake instruction; instruction is empty"
		log.Fatal(message)
	}
	return unmake(instruction)
}

func unmake(instruction []byte) (Opcode, []int, int) {
	def := Lookup(instruction[0])
	width := getWidth(def)
	validateWidth(width, instruction)
	operands := readOperands(def, instruction)
	return def.Op, operands, width
}

func validateWidth(width int, instruction []byte) {
	if len(instruction) < width {
		message := "cannot unmake instruction; unexpected instruction width"
		log.Fatal(message)
	}
}

func readOperands(def *Definition, instruction []byte) []int {
	operands := make([]int, len(def.OperandWidths))
	remain := instruction[1:]
	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(binary.BigEndian.Uint16(remain))
		case 1:
			operands[i] = int(uint8(remain[0]))
		default:
			message := "cannot unmake instruction; unexpected operand width"
			log.Fatal(message)
		}
		remain = remain[width:]
	}
	return operands
}

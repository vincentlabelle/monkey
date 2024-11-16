package code

import "testing"

func TestUnmake(t *testing.T) {
	setup := []struct {
		op       Opcode
		operands []int
		width    int
	}{
		{OpConstant, []int{65535}, 3},
		{OpAdd, []int{}, 1},
	}

	for _, s := range setup {
		instruction := Make(s.op, s.operands...)
		op, operands, width := Unmake(instruction)
		testOpcode(t, op, s.op)
		testOperands(t, operands, s.operands)
		testWidth(t, width, s.width)
	}
}

func testOpcode(t *testing.T, actual Opcode, expected Opcode) {
	if actual != expected {
		t.Fatalf("opcode mismatch. got=%v, expected=%v", actual, expected)
	}
}

func testOperands(t *testing.T, actual []int, expected []int) {
	if len(actual) != len(expected) {
		t.Fatalf(
			"number of operands mismatch. got=%v, expected=%v",
			len(actual),
			len(expected),
		)
	}
	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf(
				"operand mismatch. got=%v, expected=%v",
				actual[i],
				expected[i],
			)
		}
	}
}

func testWidth(t *testing.T, actual int, expected int) {
	if actual != expected {
		t.Fatalf("width mismatch. got=%v, expected=%v", actual, expected)
	}
}

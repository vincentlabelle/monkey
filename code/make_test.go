package code

import "testing"

func TestMake(t *testing.T) {
	setup := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{
			OpConstant,
			[]int{65534},
			[]byte{byte(OpConstant), 255, 254},
		},
		{
			OpAdd,
			[]int{},
			[]byte{byte(OpAdd)},
		},
		{
			OpCall,
			[]int{255},
			[]byte{byte(OpCall), 255},
		},
	}

	for _, s := range setup {
		instruction := Make(s.op, s.operands...)
		if len(instruction) != len(s.expected) {
			t.Fatalf(
				"number of bytes mismatch. got=%v, expected=%v",
				len(instruction),
				len(s.expected),
			)
		}
		for i, e := range s.expected {
			a := instruction[i]
			if a != e {
				t.Fatalf("byte mismatch. got=%v, expected=%v", a, e)
			}
		}
	}
}

package code

import "testing"

func TestInstructions(t *testing.T) {
	setup := []struct {
		pieces   []Instructions
		expected string
	}{
		{
			[]Instructions{
				Make(OpAdd),
				Make(OpConstant, 2),
				Make(OpConstant, 65535),
			},
			"0000 OpAdd\n0001 OpConstant 2\n0004 OpConstant 65535\n",
		},
	}

	for _, s := range setup {
		instructions := Concatenate(s.pieces)
		if instructions.String() != s.expected {
			t.Fatalf(
				"string mismatch. got=%q, expected=%q",
				instructions.String(),
				s.expected,
			)
		}
	}
}

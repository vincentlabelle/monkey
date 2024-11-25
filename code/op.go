package code

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpLowerThan
	OpMinus
	OpBang
	OpTrue
	OpFalse
	OpPop
	OpJump
	OpJumpIf
	OpNull
)

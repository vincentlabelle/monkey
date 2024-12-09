package code

type Opcode byte

const (
	OpConstant Opcode = iota
	OpTrue
	OpFalse
	OpNull
	OpArray
	OpHash
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
	OpIndex
	OpPop
	OpJump
	OpJumpIf
	OpSetGlobal
	OpGetGlobal
	OpSetLocal
	OpGetLocal
	OpGetBuiltin
	OpCall
	OpReturnValue
	OpReturn
)

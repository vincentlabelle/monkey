package code

var PrefixOperator = map[string]Opcode{
	"-": OpMinus,
	"!": OpBang,
}

var PrefixOperatorReverse = map[Opcode]string{
	OpMinus: "-",
	OpBang:  "!",
}

var InfixOperator = map[string]Opcode{
	"+":  OpAdd,
	"-":  OpSub,
	"*":  OpMul,
	"/":  OpDiv,
	"==": OpEqual,
	"!=": OpNotEqual,
	">":  OpGreaterThan,
	"<":  OpLowerThan,
}

var InfixOperatorReverse = map[Opcode]string{
	OpAdd:         "+",
	OpSub:         "-",
	OpMul:         "*",
	OpDiv:         "/",
	OpEqual:       "==",
	OpNotEqual:    "!=",
	OpGreaterThan: ">",
	OpLowerThan:   "<",
}

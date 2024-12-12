package symbol

type SymbolScope string

const (
	GlobalScope  = "GLOBAL"
	LocalScope   = "LOCAL"
	BuiltinScope = "BUILTIN"
)

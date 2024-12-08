package symbol

type SymbolTable struct {
	store map[string]Symbol
	outer *SymbolTable
}

func NewTable() *SymbolTable {
	return &SymbolTable{
		store: map[string]Symbol{},
	}
}

func NewInnerTable(outer *SymbolTable) *SymbolTable {
	table := NewTable()
	table.outer = outer
	return table
}

func (s *SymbolTable) Define(name string) Symbol {
	sym := Symbol{
		Name:  name,
		Scope: s.getScope(),
		Index: s.CountDefinitions(),
	}
	s.store[name] = sym
	return sym
}

func (s *SymbolTable) getScope() SymbolScope {
	if s.outer != nil {
		return LocalScope
	}
	return GlobalScope
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := s.store[name]
	if !ok && s.outer != nil {
		return s.outer.Resolve(name)
	}
	return sym, ok
}

func (s *SymbolTable) Outer() *SymbolTable {
	return s.outer
}

func (s *SymbolTable) CountDefinitions() int {
	return len(s.store)
}

package symbol

import "log"

type SymbolTable struct {
	store map[string]Symbol
}

func NewTable() *SymbolTable {
	return &SymbolTable{
		map[string]Symbol{},
	}
}

func (s *SymbolTable) Define(name string) Symbol {
	sym := Symbol{Name: name, Scope: GlobalScope, Index: len(s.store)}
	s.store[name] = sym
	return sym
}

func (s *SymbolTable) Resolve(name string) Symbol {
	sym, ok := s.store[name]
	if !ok {
		message := "cannot resolve; name is missing from symbol table"
		log.Fatal(message)
	}
	return sym
}

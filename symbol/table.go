package symbol

import "github.com/vincentlabelle/monkey/object"

type SymbolTable struct {
	store map[string]Symbol
	outer *SymbolTable
}

func NewTable() *SymbolTable {
	return &SymbolTable{
		store: map[string]Symbol{},
		outer: newBuiltinTable(),
	}
}

func newBuiltinTable() *SymbolTable {
	store := map[string]Symbol{}
	for i, b := range object.Builtins {
		store[b.Name] = Symbol{b.Name, BuiltinScope, i}
	}
	return &SymbolTable{store: store}
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
	if s.outer == nil {
		return BuiltinScope
	}
	if s.outer.outer == nil {
		return GlobalScope
	}
	return LocalScope
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

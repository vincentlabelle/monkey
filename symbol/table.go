package symbol

import "github.com/vincentlabelle/monkey/object"

type SymbolTable struct {
	store map[string]Symbol
	count int
	outer *SymbolTable
	free  []Symbol
}

func NewTable() *SymbolTable {
	outer := newBuiltinTable()
	return NewInnerTable(outer)
}

func newBuiltinTable() *SymbolTable {
	t := newTable()
	for i, b := range object.Builtins {
		t.store[b.Name] = Symbol{b.Name, BuiltinScope, i}
	}
	return t
}

func newTable() *SymbolTable {
	return &SymbolTable{
		store: map[string]Symbol{},
		free:  []Symbol{},
	}
}

func NewInnerTable(outer *SymbolTable) *SymbolTable {
	t := newTable()
	t.outer = outer
	return t
}

func (s *SymbolTable) Define(name string) Symbol {
	sym := Symbol{
		Name:  name,
		Scope: s.getScope(),
		Index: s.count,
	}
	s.store[name] = sym
	s.count++
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

func (s *SymbolTable) CountDefinitions() int {
	return s.count
}

func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	sym := Symbol{Name: name, Scope: FunctionScope, Index: 0}
	s.store[name] = sym
	return sym
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := s.store[name]
	if !ok && s.outer != nil {
		sym, ok = s.outer.Resolve(name)
		if ok && (sym.Scope == LocalScope || sym.Scope == FreeScope) {
			sym = s.Redefine(sym)
		}
	}
	return sym, ok
}

func (s *SymbolTable) Redefine(sym Symbol) Symbol {
	free := Symbol{Name: sym.Name, Scope: FreeScope, Index: len(s.free)}
	s.store[free.Name] = free
	s.free = append(s.free, sym)
	return free
}

func (s *SymbolTable) Outer() *SymbolTable {
	return s.outer
}

func (s *SymbolTable) Free() []Symbol {
	return s.free
}

package symbol

import "testing"

func TestDefine(t *testing.T) {
	global := NewTable()
	local := NewInnerTable(global)
	nested := NewInnerTable(local)

	setup := []struct {
		table    *SymbolTable
		expected []Symbol
	}{
		{
			global,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
			},
		},
		{
			local,
			[]Symbol{
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			nested,
			[]Symbol{
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, s := range setup {
		for _, e := range s.expected {
			a := s.table.Define(e.Name)
			if a != e {
				t.Fatalf("symbol mismatch. got=%v, expected=%v", a, e)
			}
		}
	}
}

func TestResolve(t *testing.T) {
	global := NewTable()
	global.Define("a")
	global.Define("b")

	local := NewInnerTable(global)
	local.Define("c")
	local.Define("d")

	nested := NewInnerTable(local)
	nested.Define("e")
	nested.Define("f")

	setup := []struct {
		table    *SymbolTable
		expected []Symbol
	}{
		{
			global,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
			},
		},
		{
			local,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			nested,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, s := range setup {
		for _, e := range s.expected {
			a, ok := s.table.Resolve(e.Name)
			if !ok {
				t.Fatalf("symbol missing. expected=%v", e.Name)
			}
			if a != e {
				t.Fatalf("symbol mismatch. got=%v, expected=%v", a, e)
			}
		}
	}
}

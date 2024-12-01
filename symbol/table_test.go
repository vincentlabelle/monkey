package symbol

import "testing"

func TestDefine(t *testing.T) {
	setup := []struct {
		name     string
		expected Symbol
	}{
		{"a", Symbol{Name: "a", Scope: GlobalScope, Index: 0}},
		{"b", Symbol{Name: "b", Scope: GlobalScope, Index: 1}},
	}

	table := NewTable()
	for _, s := range setup {
		actual := table.Define(s.name)
		if actual != s.expected {
			t.Fatalf("symbol mismatch. got=%v, expected=%v", actual, s.expected)
		}
	}
}

func TestResolve(t *testing.T) {
	setup := []struct {
		name     string
		expected Symbol
	}{
		{"a", Symbol{Name: "a", Scope: GlobalScope, Index: 0}},
		{"b", Symbol{Name: "b", Scope: GlobalScope, Index: 1}},
	}

	table := NewTable()
	for _, s := range setup {
		table.Define(s.name)
		actual := table.Resolve(s.name)
		if actual != s.expected {
			t.Fatalf("symbol mismatch. got=%v, expected=%v", actual, s.expected)
		}
	}
}

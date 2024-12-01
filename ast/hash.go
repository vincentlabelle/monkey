package ast

import (
	"cmp"
	"iter"
	"slices"
)

type HashKey struct {
	Index      int
	Expression Expression
}

func SortHashKeys(keys iter.Seq[HashKey]) []HashKey {
	return slices.SortedStableFunc(keys, func(x, y HashKey) int {
		return cmp.Compare(x.Index, y.Index)
	})
}

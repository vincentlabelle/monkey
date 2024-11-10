package object

type Hashable interface {
	Object
	HashKey() HashKey
}

type HashKey struct {
	Type  string
	Value uint64
}

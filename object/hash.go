package object

import "log"

type Hashable interface {
	Object
	HashKey() HashKey
}

func CastToHashable(obj Object) Hashable {
	hash, ok := obj.(Hashable)
	if !ok {
		message := "cannot cast to hashable; " +
			"unexpected object encountered"
		log.Fatal(message)
	}
	return hash
}

type HashKey struct {
	Type  string
	Value uint64
}

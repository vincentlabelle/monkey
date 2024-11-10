package object

import "testing"

func TestHashKeyWhenEqual(t *testing.T) {
	setup := []struct {
		one Hashable
		two Hashable
	}{
		{&String{Value: ""}, &String{Value: ""}},
		{&String{Value: "abc"}, &String{Value: "abc"}},
		{&Integer{Value: -1}, &Integer{Value: -1}},
		{&Integer{Value: 0}, &Integer{Value: 0}},
		{&Integer{Value: 1}, &Integer{Value: 1}},
		{&Boolean{Value: true}, &Boolean{Value: true}},
		{&Boolean{Value: false}, &Boolean{Value: false}},
	}

	for _, s := range setup {
		if s.one.HashKey() != s.two.HashKey() {
			t.Fatalf(
				"hash key mismatch. got=%v, expected=%v",
				s.one.HashKey(),
				s.two.HashKey(),
			)
		}
	}
}

func TestHashKeyWhenUnequal(t *testing.T) {
	setup := []struct {
		one Hashable
		two Hashable
	}{
		{&String{Value: "abc"}, &String{Value: ""}},
		{&Integer{Value: -1}, &Integer{Value: 0}},
		{&Boolean{Value: true}, &Boolean{Value: false}},
		{&Boolean{Value: true}, &Integer{Value: 1}},
		{&Boolean{Value: false}, &Integer{Value: 0}},
	}

	for _, s := range setup {
		if s.one.HashKey() == s.two.HashKey() {
			t.Fatalf(
				"hash key match. got=%v, expected=%v",
				s.one.HashKey(),
				s.two.HashKey(),
			)
		}
	}
}

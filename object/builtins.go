package object

import "log"

func len_(args ...Object) Object {
	if len(args) != 1 {
		message := "cannot call built-in len; one argument is expected"
		log.Fatal(message)
	}
	s, ok := args[0].(*String)
	if !ok {
		message := "cannot call built-in len; argument must be string"
		log.Fatal(message)
	}
	return &Integer{Value: len(s.Value)}
}

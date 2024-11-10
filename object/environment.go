package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	outer := newBuiltinEnvironment()
	return NewInnerEnvironment(outer)
}

func newBuiltinEnvironment() *Environment {
	store := map[string]Object{
		"len":   &Builtin{Fn: len_},
		"first": &Builtin{Fn: first},
		"last":  &Builtin{Fn: last},
		"rest":  &Builtin{Fn: rest},
		"push":  &Builtin{Fn: push},
	}
	return &Environment{store: store}
}

func NewInnerEnvironment(outer *Environment) *Environment {
	env := newEnvironment()
	env.outer = outer
	return env
}

func newEnvironment() *Environment {
	return &Environment{store: map[string]Object{}}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, obj Object) {
	e.store[name] = obj
}

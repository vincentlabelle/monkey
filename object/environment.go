package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	outer := newBuiltinEnvironment()
	return &Environment{store: map[string]Object{}, outer: outer}
}

func newBuiltinEnvironment() *Environment {
	store := map[string]Object{
		"len": &Builtin{Fn: len_},
	}
	return &Environment{store: store}
}

func NewInnerEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
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

package object

// Environment represents a programming language environment.
type Environment struct {
	store map[string]Object
}

// NewEnvironment returns a new environment.
func NewEnvironment() *Environment {
	return &Environment{
		store: map[string]Object{},
	}
}

// Get returns a named object from the environment.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set sets a named object in the environment.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

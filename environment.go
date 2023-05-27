package main

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	values := make(map[string]any)
	return &Environment{
		enclosing,
		values,
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) (any, error) {
	if value, ok := e.values[name.lexeme]; !ok {
		errMsg := "Undefined variable '" + name.lexeme + "'."
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		}
		return nil, &RuntimeError{name, errMsg}
	} else {
		return value, nil
	}
}

func (e *Environment) Assign(name Token, value any) error {
	if _, ok := e.values[name.lexeme]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Assign(name, value)
		}
		return &RuntimeError{name, "Undefined variable '" + name.lexeme + "'."}
	}
	e.values[name.lexeme] = value
	return nil
}

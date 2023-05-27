package main

type Environment struct {
	values map[string]any
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) (any, error) {
	if value, ok := e.values[name.lexeme]; !ok {
		errMsg := "Undefined variable '" + name.lexeme + "'."
		return nil, &RuntimeError{name, errMsg}
	} else {
		return value, nil
	}
}

func (e *Environment) Assign(name Token, value any) error {
	if _, ok := e.values[name.lexeme]; !ok {
		return &RuntimeError{name, "Undefined variable '" + name.lexeme + "'."}
	}
	e.values[name.lexeme] = value
	return nil
}

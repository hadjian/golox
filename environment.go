package main

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

var NewEnvironment func(enclosing *Environment) *Environment = func(enclosing *Environment) *Environment {
	values := make(map[string]any)
	return &Environment{
		enclosing,
		values,
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) any {
	if value, ok := e.values[name.lexeme]; !ok {
		errMsg := "Undefined variable '" + name.lexeme + "'."
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		}
		panic(RuntimeError{name, errMsg})
	} else {
		return value
	}
}

func (e *Environment) GetAt(distance int, name string) any {
	return e.Ancestor(distance).values[name]
}

func (e *Environment) Ancestor(distance int) Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return *env
}

func (e *Environment) Assign(name Token, value any) any {
	if _, ok := e.values[name.lexeme]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Assign(name, value)
		}
		panic(RuntimeError{name, "Undefined variable '" + name.lexeme + "'."})
	}
	e.values[name.lexeme] = value
	return nil
}

func (e *Environment) AssignAt(distance int, name Token, value any) {
	e.Ancestor(distance).values[name.lexeme] = value
}

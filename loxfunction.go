package main

type LoxFunction struct {
	declaration Function
}

func (l LoxFunction) Call(i *Interpreter, args []any) any {
	env := NewEnvironment(i.globals)
	for i := 0; i < len(l.declaration.params); i++ {
		env.Define(l.declaration.params[i].lexeme, args[i])
	}
	i.executeBlock(l.declaration.body, env)
	return nil
}

func (l LoxFunction) Arity() int {
	return len(l.declaration.params)
}

func (l LoxFunction) String() string {
	return "<fn " + l.declaration.name.lexeme + ">"
}

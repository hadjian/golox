package main

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []any) any
	Arity() int
}

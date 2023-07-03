package main

import (
	"time"
)

type LoxTime struct {
}

func (l *LoxTime) Arity() int {
	return 0
}

func (l *LoxTime) Call(i Interpreter, args []any) any {
	return time.Now()
}

func (l *LoxTime) String() string {
	return "<native fn>"
}

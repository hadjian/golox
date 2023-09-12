package main

import (
	"fmt"
	"strings"
)

type RuntimeError struct {
	token Token
	msg   string
}

type ReturnValue struct {
	value any
}

func (re RuntimeError) Error() string {
	return fmt.Sprintf("%v \n[line %v]", re.msg, re.token.line)
}

type Interpreter struct {
	environment *Environment
	globals     *Environment
	locals      map[Expr]int
}

func NewInterpreter() *Interpreter {
	env := NewEnvironment(nil)
	i := &Interpreter{
		environment: env,
		globals:     env,
		locals:      map[Expr]int{},
	}

	i.globals.Define("clock", LoxTime{})
	return i
}

func (i *Interpreter) Interpret(stmts []Stmt) {
	defer func() {
		if err := recover(); err != nil {
			runtimeError(err.(error))
		}
	}()
	for _, stmt := range stmts {
		i.Execute(stmt)
	}
}

func (i *Interpreter) Execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) VisitBlock(b *Block) {
	i.executeBlock(b.statements, NewEnvironment(i.environment))
}

func (i *Interpreter) executeBlock(stmts []Stmt, env *Environment) {
	previous := i.environment
	i.environment = env
	defer func() {
		i.environment = previous
	}()
	for _, stmt := range stmts {
		i.Execute(stmt)
	}
}

func (i *Interpreter) Evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitBinary(b *Binary) any {
	left := i.Evaluate(b.Left)
	right := i.Evaluate(b.Right)

	switch b.Operator.tType {
	case BANG_EQUAL:
		return left != right
	case EQUAL_EQUAL:
		return left == right
	case GREATER:
		left, right := i.checkNumberOperands(b.Operator, left, right)
		return left > right
	case GREATER_EQUAL:
		left, right := i.checkNumberOperands(b.Operator, left, right)
		return left >= right
	case LESS:
		left, right := i.checkNumberOperands(b.Operator, left, right)
		return left < right
	case LESS_EQUAL:
		left, right := i.checkNumberOperands(b.Operator, left, right)
		return left <= right
	case MINUS:
		left, right := i.checkNumberOperands(b.Operator, left, right)
		return left - right
	case PLUS:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left + right
			}
		}
		if left, ok := left.(string); ok {
			if right, ok := right.(string); ok {
				return left + right
			}
		}
		msg := "Operands must be two numbers or two strings."
		panic(RuntimeError{b.Operator, msg})
	case SLASH:
		left, right := i.checkNumberOperands(b.Operator, left, right)
		return left / right
	case STAR:
		left, right := i.checkNumberOperands(b.Operator, left, right)
		return left * right
	}
	return nil
}

func (i *Interpreter) VisitCallExpr(c *Call) any {
	callee := i.Evaluate(c.callee)

	var arguments []any
	for _, argument := range c.arguments {
		arguments = append(arguments, i.Evaluate(argument))
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		panic(RuntimeError{c.paren, "Can only call functions and classes"})
	}

	if len(arguments) != function.Arity() {
		msg := "Expected %v arguments but got %v\n"
		msg = fmt.Sprintf(msg, function.Arity(), len(arguments))
		panic(RuntimeError{c.paren, msg})
	}
	return function.Call(i, arguments)
}

func (i *Interpreter) VisitGrouping(g *Grouping) any {
	return i.Evaluate(g.Expression)
}

func (i *Interpreter) VisitLiteral(l *Literal) any {
	return l.Value
}

func (i *Interpreter) VisitLogical(l *Logical) any {
	left := i.Evaluate(l.left)
	if l.operator.tType == OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}
	return i.Evaluate(l.right)
}

func (i *Interpreter) VisitVariableExpr(expr *Variable) any {
	return i.LookUpVariable(expr.name, expr)
}

func (i *Interpreter) LookUpVariable(name Token, expr Expr) any {
	distance, ok := i.locals[expr]
	if ok != false {
		return i.environment.GetAt(distance, name.lexeme)
	} else {
		return i.globals.Get(name)
	}
}

func (i *Interpreter) VisitVarStmt(stmt *Var) {
	var value any
	if stmt.initializer != nil {
		value = i.Evaluate(stmt.initializer)
	}
	i.environment.Define(stmt.name.lexeme, value)
}

func (i *Interpreter) VisitAssign(a *Assign) any {
	value := i.Evaluate(a.value)

	distance, ok := i.locals[a]
	if ok != false {
		i.environment.AssignAt(distance, a.name, value)
	} else {
		i.globals.Assign(a.name, value)
	}

	return value
}

func (i *Interpreter) VisitUnary(u *Unary) any {
	right := i.Evaluate(u.Right)
	switch u.Operator.tType {
	case BANG:
		return !i.isTruthy(right)
	case MINUS:
		number := i.checkNumberOperand(u.Operator, right)
		return -number
	}
	return nil
}

func (i *Interpreter) VisitExpressionStmt(e *Expression) {
	i.Evaluate(e.expr)
}

func (i *Interpreter) VisitFunction(stmt *Function) {
	function := LoxFunction{*stmt, *i.environment}
	i.environment.Define(stmt.name.lexeme, function)
}

func (i *Interpreter) VisitWhile(w *While) {
	loop := i.Evaluate(w.condition)
	for i.isTruthy(loop) {
		i.Execute(w.body)
		loop = i.Evaluate(w.condition)
	}
}

func (i *Interpreter) VisitIf(f *If) {
	condition := i.Evaluate(f.condition)
	if i.isTruthy(condition) {
		i.Execute(f.thenBranch)
	} else if f.elseBranch != nil {
		i.Execute(f.elseBranch)
	}
}

func (i *Interpreter) VisitPrint(p *Print) {
	value := i.Evaluate(p.expr)
	fmt.Println(i.stringify(value))
}

func (i *Interpreter) VisitReturn(r *Return) {
	var value any
	if r.value != nil {
		value = i.Evaluate(r.value)
	}
	panic(ReturnValue{value})
}

func (i *Interpreter) isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if v, ok := value.(bool); ok {
		return v
	}
	return true
}

func (i *Interpreter) checkNumberOperand(op Token, value any) float64 {
	if number, ok := value.(float64); !ok {
		panic(RuntimeError{op, "Operand must be a number."})
	} else {
		return number
	}
}

func (i *Interpreter) checkNumberOperands(op Token, left, right any) (float64, float64) {
	l, ok1 := left.(float64)
	r, ok2 := right.(float64)
	if !ok1 || !ok2 {
		panic(RuntimeError{op, "Operands must be numbers."})
	}
	return l, r
}

func (i *Interpreter) stringify(value any) string {
	var text string
	if runes, ok := value.([]rune); ok {
		value = string(runes)
	}
	text = fmt.Sprintf("%v", value)
	if _, ok := value.(float64); ok {
		if strings.HasSuffix(text, ".0") {
			text = text[:len(text)-2]
		}
	}
	return text
}

func (i *Interpreter) Resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

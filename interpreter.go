package main

import (
	"errors"
	"fmt"
	"strings"
)

type RuntimeError struct {
	token Token
	msg   string
}

func (re *RuntimeError) Error() string {
	return re.msg
}

type Interpreter struct {
}

func (i *Interpreter) Interpret(expr Expr) {
	value, err := i.Evaluate(expr)
	if err != nil {
		runtimeError(err)
	}
	fmt.Println(i.stringify(value))
}

func (i *Interpreter) Evaluate(expr Expr) (any, error) {
	value, err := expr.Accept(i)
	if err != nil {
		return nil, errors.New("Error evaluating expression")
	}
	return value, nil
}

func (i *Interpreter) VisitBinary(b *Binary) (any, error) {
	var left, right any
	var err error
	if left, err = i.Evaluate(b.Left); err != nil {
		return nil, err
	}
	if right, err = i.Evaluate(b.Right); err != nil {
		return nil, err
	}

	switch b.Operator.tType {
	case BANG_EQUAL:
		return left != right, nil
	case EQUAL_EQUAL:
		return left == right, nil
	case GREATER:
		left, right, err := i.checkNumberOperands(b.Operator, left, right)
		return left > right, err
	case GREATER_EQUAL:
		left, right, err := i.checkNumberOperands(b.Operator, left, right)
		return left >= right, err
	case LESS:
		left, right, err := i.checkNumberOperands(b.Operator, left, right)
		return left < right, err
	case LESS_EQUAL:
		left, right, err := i.checkNumberOperands(b.Operator, left, right)
		return left <= right, err
	case MINUS:
		left, right, err := i.checkNumberOperands(b.Operator, left, right)
		return left - right, err
	case PLUS:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left + right, nil
			}
		}
		if left, ok := left.(string); ok {
			if right, ok := right.(string); ok {
				return left + right, nil
			}
		}
		msg := "Operands must be two numbers or two strings."
		return nil, &RuntimeError{b.Operator, msg}
	case SLASH:
		left, right, err := i.checkNumberOperands(b.Operator, left, right)
		return left / right, err
	case STAR:
		left, right, err := i.checkNumberOperands(b.Operator, left, right)
		return left * right, err
	}
	return nil, nil
}

func (i *Interpreter) VisitGrouping(g *Grouping) (any, error) {
	return i.Evaluate(g.Expression)
}

func (i *Interpreter) VisitLiteral(l *Literal) (any, error) {
	return l.Value, nil
}

func (i *Interpreter) VisitUnary(u *Unary) (any, error) {
	var right interface{}
	var err error
	if right, err = i.Evaluate(u.Right); err != nil {
		return nil, err
	}
	switch u.Operator.tType {
	case BANG:
		return !i.isTruthy(right), nil
	case MINUS:
		if right, err := i.checkNumberOperand(u.Operator, right); err != nil {
			return right, err
		} else {
			return -right, nil
		}
	}
	return nil, nil
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

func (i *Interpreter) checkNumberOperand(op Token, value any) (float64, error) {
	if number, ok := value.(float64); !ok {
		return number, &RuntimeError{op, "Operand must be a number."}
	} else {
		return number, nil
	}
}

func (i *Interpreter) checkNumberOperands(op Token, left, right any) (float64, float64, error) {
	l, ok1 := left.(float64)
	r, ok2 := right.(float64)
	if !ok1 || !ok2 {
		return l, r, &RuntimeError{op, "Operands must be numbers."}
	}
	return l, r, nil
}

func (i *Interpreter) stringify(value any) string {
	text := fmt.Sprintf("%v", value)
	if _, ok := value.(float64); ok {
		if strings.HasSuffix(text, ".0") {
			text = text[:len(text)-2]
		}
	}
	return text
}

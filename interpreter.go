package main

type Interpreter struct {
}

func (i *Interpreter) Evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitBinary(b *Binary) any {

	return nil
}

func (i *Interpreter) VisitGrouping(g *Grouping) any {
	return i.Evaluate(g.Expression)

}

func (i *Interpreter) VisitLiteral(l *Literal) any {
	return l.Value
}

func (i *Interpreter) VisitUnary(u *Unary) any {
	right := i.Evaluate(u.Right)
	switch u.Operator.tType {
	case BANG:
		return !i.isTruthy(right)
	case MINUS:
		return -right.(float64)
	}
	return nil
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

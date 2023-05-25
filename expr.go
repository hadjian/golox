package main

type ExprVisitor interface {
	VisitBinary(b *Binary) (any, error)
	VisitGrouping(g *Grouping) (any, error)
	VisitLiteral(l *Literal) (any, error)
	VisitVariableExpr(v *Variable) (any, error)
	VisitUnary(u *Unary) (any, error)
}

type Expr interface {
	Accept(v ExprVisitor) (any, error)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(v ExprVisitor) (any, error) {
	return v.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(v ExprVisitor) (any, error) {
	return v.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(v ExprVisitor) (any, error) {
	return v.VisitLiteral(l)
}

type Variable struct {
	name Token
}

func (ev *Variable) Accept(v ExprVisitor) (any, error) {
	return v.VisitVariableExpr(ev)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v ExprVisitor) (any, error) {
	return v.VisitUnary(u)
}

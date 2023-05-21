package main

type Visitor interface {
	VisitBinary(b *Binary) (any, error)
	VisitGrouping(g *Grouping) (any, error)
	VisitLiteral(l *Literal) (any, error)
	VisitUnary(u *Unary) (any, error)
}

type Expr interface {
	Accept(v Visitor) (any, error)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(v Visitor) (any, error) {
	return v.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(v Visitor) (any, error) {
	return v.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(v Visitor) (any, error) {
	return v.VisitLiteral(l)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v Visitor) (any, error) {
	return v.VisitUnary(u)
}

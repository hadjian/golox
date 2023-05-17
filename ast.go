package main

type Visitor interface {
	VisitBinary(b *Binary) any
	VisitGrouping(g *Grouping) any
	VisitLiteral(l *Literal) any
	VisitUnary(u *Unary) any
}

type Expr interface {
	Accept(v Visitor) any
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(v Visitor) any {
	return v.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(v Visitor) any {
	return v.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(v Visitor) any {
	return v.VisitLiteral(l)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v Visitor) any {
	return v.VisitUnary(u)
}

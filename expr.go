package main

type ExprVisitor interface {
	VisitAssign(b *Assign) any
	VisitBinary(b *Binary) any
	VisitCallExpr(c *Call) any
	VisitGrouping(g *Grouping) any
	VisitLiteral(l *Literal) any
	VisitLogical(l *Logical) any
	VisitVariableExpr(v *Variable) any
	VisitUnary(u *Unary) any
}

type Expr interface {
	Accept(v ExprVisitor) any
}

type Assign struct {
	name  Token
	value Expr
}

func (a *Assign) Accept(v ExprVisitor) any {
	return v.VisitAssign(a)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(v ExprVisitor) any {
	return v.VisitBinary(b)
}

type Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func (c *Call) Accept(v ExprVisitor) any {
	return v.VisitCallExpr(c)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(v ExprVisitor) any {
	return v.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(v ExprVisitor) any {
	return v.VisitLiteral(l)
}

type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

func (l *Logical) Accept(v ExprVisitor) any {
	return v.VisitLogical(l)
}

type Variable struct {
	name Token
}

func (ev *Variable) Accept(v ExprVisitor) any {
	return v.VisitVariableExpr(ev)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v ExprVisitor) any {
	return v.VisitUnary(u)
}

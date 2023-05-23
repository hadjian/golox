package main

type StmtVisitor interface {
	VisitExpression(stmt Stmt) (any, error)
	VisitPrint(stmt Stmt) (any, error)
}

type Stmt interface {
	Accept(v StmtVisitor) (any, error)
}

type Expression struct {
	expr Expr
}

func (e *Expression) Accept(v StmtVisitor) (any, error) {
	return v.VisitExpression(e)
}

type Print struct {
	expr Expr
}

func (p *Print) Accept(v StmtVisitor) (any, error) {
	return v.VisitPrint(p)
}

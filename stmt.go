package main

type StmtVisitor interface {
	VisitBlock(b *Block) error
	VisitExpression(e *Expression) error
	VisitPrint(p *Print) error
	VisitVarStmt(v *Var) error
}

type Stmt interface {
	Accept(v StmtVisitor) error
}

type Block struct {
	statements []Stmt
}

func (b *Block) Accept(v StmtVisitor) error {
	return v.VisitBlock(b)
}

type Expression struct {
	expr Expr
}

func (e *Expression) Accept(v StmtVisitor) error {
	return v.VisitExpression(e)
}

type Print struct {
	expr Expr
}

func (p *Print) Accept(v StmtVisitor) error {
	return v.VisitPrint(p)
}

type Var struct {
	name        Token
	initializer Expr
}

func (vr *Var) Accept(v StmtVisitor) error {
	return v.VisitVarStmt(vr)
}

package main

type StmtVisitor interface {
	VisitBlock(b *Block) error
	VisitExpressionStmt(e *Expression) error
	VisitIf(f *If) error
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
	return v.VisitExpressionStmt(e)
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (f *If) Accept(v StmtVisitor) error {
	return v.VisitIf(f)
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

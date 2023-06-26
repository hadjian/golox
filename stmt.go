package main

type StmtVisitor interface {
	VisitBlock(b *Block)
	VisitExpressionStmt(e *Expression)
	VisitWhile(w *While)
	VisitIf(f *If)
	VisitPrint(p *Print)
	VisitVarStmt(v *Var)
}

type Stmt interface {
	Accept(v StmtVisitor)
}

type Block struct {
	statements []Stmt
}

func (b *Block) Accept(v StmtVisitor) {
	v.VisitBlock(b)
}

type Expression struct {
	expr Expr
}

func (e *Expression) Accept(v StmtVisitor) {
	v.VisitExpressionStmt(e)
}

type While struct {
	condition Expr
	body      Stmt
}

func (w *While) Accept(v StmtVisitor) {
	v.VisitWhile(w)
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (f *If) Accept(v StmtVisitor) {
	v.VisitIf(f)
}

type Print struct {
	expr Expr
}

func (p *Print) Accept(v StmtVisitor) {
	v.VisitPrint(p)
}

type Var struct {
	name        Token
	initializer Expr
}

func (vr *Var) Accept(v StmtVisitor) {
	v.VisitVarStmt(vr)
}

package main

type StmtVisitor interface {
	VisitBlock(b *Block)
	VisitExpressionStmt(e *Expression)
	VisitFunction(f *Function)
	VisitWhile(w *While)
	VisitIf(f *If)
	VisitPrint(p *Print)
	VisitReturn(r *Return)
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

type Function struct {
	name   Token
	params []Token
	body   []Stmt
}

func (f *Function) Accept(v StmtVisitor) {
	v.VisitFunction(f)
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

type Return struct {
	keyword Token
	value   Expr
}

func (r *Return) Accept(v StmtVisitor) {
	v.VisitReturn(r)
}

type Var struct {
	name        Token
	initializer Expr
}

func (vr *Var) Accept(v StmtVisitor) {
	v.VisitVarStmt(vr)
}

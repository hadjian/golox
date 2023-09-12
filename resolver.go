package main

import "github.com/hadjian/golox/util"

type Resolver struct {
	interpreter Interpreter
	scopes      util.Stack
}

func NewResolver(i Interpreter) Resolver {
	return Resolver{i, util.Stack{}}
}

func (r *Resolver) resolveStmts(stmts []Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes.Push(map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) VisitBlock(b *Block) {
	r.beginScope()
	r.resolveStmts(b.statements)
	r.endScope()
}

func (r *Resolver) VisitVarStmt(v *Var) {
	r.declare(v.name)
	if v.initializer != nil {
		r.resolveExpr(v.initializer)
	}
	r.define(v.name)
}

func (r *Resolver) declare(name Token) {
	if r.scopes.IsEmpty() {
		return
	}
	scope := r.scopes.Peek().(map[string]bool)
	scope[name.lexeme] = false
}

func (r *Resolver) define(name Token) {
	if r.scopes.IsEmpty() {
		return
	}
	r.scopes.Peek().(map[string]bool)[name.lexeme] = true
}

func (r *Resolver) VisitVariableExpr(expr *Variable) any {
	if env, valid := r.scopes.Peek().(map[string]bool); valid {
		if defined, ok := env[expr.name.lexeme]; !defined && ok {
			errToken(expr.name, "Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		if _, ok := r.scopes.Get(i).(map[string]bool)[name.lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.Size()-1-i)
		}
	}
}

func (r *Resolver) VisitAssign(expr *Assign) any {
	r.resolveExpr(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) VisitFunction(stmt *Function) {
	r.declare(stmt.name)
	r.define(stmt.name)
	r.resolveFunction(*stmt)
}

func (r *Resolver) resolveFunction(fn Function) {
	r.beginScope()
	for _, param := range fn.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmts(fn.body)
	r.endScope()
}

func (r *Resolver) VisitExpressionStmt(stmt *Expression) {
	r.resolveExpr(stmt.expr)
}

func (r *Resolver) VisitIf(stmt *If) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStmt(stmt.elseBranch)
	}
}

func (r *Resolver) VisitPrint(stmt *Print) {
	r.resolveExpr(stmt.expr)
}

func (r *Resolver) VisitReturn(stmt *Return) {
	if stmt.value != nil {
		r.resolveExpr(stmt.value)
	}
}

func (r *Resolver) VisitWhile(stmt *While) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
}

func (r *Resolver) VisitBinary(expr *Binary) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *Call) any {
	r.resolveExpr(expr.callee)
	for _, argument := range expr.arguments {
		r.resolveExpr(argument)
	}
	return nil
}

func (r *Resolver) VisitGrouping(expr *Grouping) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteral(expr *Literal) any {
	return nil
}

func (r *Resolver) VisitLogical(expr *Logical) any {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

func (r *Resolver) VisitUnary(expr *Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}

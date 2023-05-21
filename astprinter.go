package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (a *AstPrinter) Print(e Expr) string {
	result, _ := e.Accept(a)
	return result.(string)
}

func (a *AstPrinter) VisitBinary(b *Binary) (any, error) {
	return a.parenthesize(b.Operator.lexeme, b.Left, b.Right)
}

func (a *AstPrinter) VisitGrouping(g *Grouping) (any, error) {
	return a.parenthesize("group", g.Expression)

}

func (a *AstPrinter) VisitLiteral(l *Literal) (any, error) {
	if l.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", l.Value), nil
}

func (a *AstPrinter) VisitUnary(u *Unary) (any, error) {
	return a.parenthesize(u.Operator.lexeme, u.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) (any, error) {
	builder := strings.Builder{}
	builder.WriteString("(" + name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(a.Print(expr))
	}
	builder.WriteString(")")
	return builder.String(), nil
}

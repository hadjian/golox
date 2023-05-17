package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (a *AstPrinter) Print(e Expr) string {
	return e.Accept(a).(string)
}

func (a *AstPrinter) VisitBinary(b *Binary) any {
	return a.parenthesize(b.Operator.lexeme, b.Left, b.Right)
}

func (a *AstPrinter) VisitGrouping(g *Grouping) any {
	return a.parenthesize("group", g.Expression)

}

func (a *AstPrinter) VisitLiteral(l *Literal) any {
	if l.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", l.Value)
}

func (a *AstPrinter) VisitUnary(u *Unary) any {
	return a.parenthesize(u.Operator.lexeme, u.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) any {
	builder := strings.Builder{}
	builder.WriteString("(" + name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(a.Print(expr))
	}
	builder.WriteString(")")
	return builder.String()
}

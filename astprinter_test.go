package main

import (
	"fmt"
	"testing"
)

func TestAstPrinter(t *testing.T) {
	astPrinter := AstPrinter{}
	unary := Unary{}
	unary.Operator = Token{
		MINUS,
		"-",
		nil,
		1,
	}
	unary.Right = &Literal{
		123,
	}
	group := Grouping{}
	group.Expression = &Literal{
		45.67,
	}
	star := Token{
		STAR,
		"*",
		nil,
		1,
	}
	binary := Binary{
		&unary,
		star,
		&group,
	}
	fmt.Println(astPrinter.Print(&binary))
}

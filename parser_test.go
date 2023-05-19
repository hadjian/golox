package main

import (
	"testing"
)

func TestParser(t *testing.T) {
	src := "32/2+33==1+1+34-2*(-2)==1"
	expected := "(== (== (+ (/ 32 2) 33) (- (+ (+ 1 1) 34) (* 2 (- 2)))) 1)"
	scanner := NewScanner(src)
	tokens := scanner.scanTokens()
	parser := NewParser(tokens)
	expr := parser.parse()
	result := (&AstPrinter{}).Print(expr)
	if expected != result {
		t.Errorf("expression: " + src)
		t.Errorf("result:     " + result)
		t.Errorf("expected:   " + expected)
	}
}

package main

import (
	"fmt"
	"testing"
)

func TestStringer(t *testing.T) {
	want := "LEFT_PAREN"
	var token Token
	got := fmt.Sprintf("%v", token.tType)
	if got != want {
		t.Error("Printing TokenType as string failed: " + got + " / " + want)
	}
}

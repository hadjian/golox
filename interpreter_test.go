package main

import "testing"

func TestIsTruthy(t *testing.T) {
	i := Interpreter{}
	tests := map[any]bool{
		0:       true,
		1:       true,
		45:      true,
		nil:     false,
		"hello": true,
		'r':     true,
		false:   false,
	}
	for input, expected := range tests {
		if out := i.isTruthy(input); out != expected {
			t.Errorf("output  : %v == %v\n", input, out)
			t.Errorf("expected: %v == %v\n", input, expected)
			t.Errorf("\n")
		}
	}
}

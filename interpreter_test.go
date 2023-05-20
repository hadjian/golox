package main

import (
	"reflect"
	"testing"
)

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

func TestAddition(t *testing.T) {
	tests := []struct {
		left      any
		operator  Token
		right     any
		valueType reflect.Type
		expected  any
	}{
		{
			4.0,
			Token{tType: PLUS, lexeme: "+"},
			5.0,
			reflect.TypeOf(float64(0.0)),
			9.0,
		},
		{
			"Hello ",
			Token{tType: PLUS, lexeme: "+"},
			"World!",
			reflect.TypeOf(string("")),
			"Hello World!",
		},
		{
			4.0,
			Token{tType: STAR, lexeme: "*"},
			9.0,
			reflect.TypeOf(float64(0)),
			36.0,
		},
		{
			16.0,
			Token{tType: SLASH, lexeme: "/"},
			4.0,
			reflect.TypeOf(float64(0)),
			4.0,
		},
		{
			16.0,
			Token{tType: LESS, lexeme: "<"},
			4.0,
			reflect.TypeOf(bool(true)),
			false,
		},
		{
			4.0,
			Token{tType: LESS, lexeme: "<"},
			16.0,
			reflect.TypeOf(bool(true)),
			true,
		},
	}

	i := Interpreter{}
	for _, test := range tests {
		b := Binary{}
		b.Left = &Literal{test.left}
		b.Operator = test.operator
		b.Right = &Literal{test.right}
		expr := (&AstPrinter{}).Print(&b)

		output := i.Evaluate(&b)
		outValue := reflect.ValueOf(output)
		if outValue.Type().ConvertibleTo(test.valueType) {
			output := outValue.Convert(test.valueType).Interface()
			if output != test.expected {
				t.Errorf("Evaluation failed for binary %s\n", expr)
				t.Errorf("%v != %v", output, test.expected)
			}
		} else {
			t.Errorf("Evaluation failed for binary %s\n", expr)
			t.Errorf("Could not cast ouput %v to float64", output)
		}
	}
}

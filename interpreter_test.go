package main

import (
	"fmt"
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
			5.0,
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

		output, _ := i.Evaluate(&b)
		outValue := reflect.ValueOf(output)
		if outValue.Type().ConvertibleTo(test.valueType) {
			output := outValue.Convert(test.valueType).Interface()
			if output != test.expected {
				t.Errorf("Evaluation failed for binary %s\n", b)
				t.Errorf("%v != %v", output, test.expected)
			}
		} else {
			t.Errorf("Evaluation failed for binary %s\n", b)
			t.Errorf("Could not cast ouput %v to float64", output)
		}
	}
}

func TestStringify(t *testing.T) {
	i := Interpreter{}
	fmt.Println(i.stringify(int(134235.0)))
}

func TestScope(t *testing.T) {
	tests := []struct {
		src      string
		expected []map[string]any
	}{
		{
			"var a = 3; { var a = 2; a = a*2; } a = a + 1;",
			[]map[string]any{
				{
					"a": 4.0,
				},
				{
					"a": 4.0,
				},
			},
		},
		{
			"var a = 1; { a = a+3; var b = 1; } a = a*3;",
			[]map[string]any{
				{
					"a": 12.0,
				},
				{
					"b": 1.0,
				},
			},
		},
	}
	var createdEnvironments []*Environment
	previous := NewEnvironment
	NewEnvironment = func(enclosing *Environment) *Environment {
		newEnv := previous(enclosing)
		createdEnvironments = append(createdEnvironments, newEnv)
		return newEnv
	}

	for _, test := range tests {
		scanner := Scanner{}
		parser := Parser{}
		scanner.Source = []rune(test.src)
		tokens := scanner.scanTokens()
		parser.Tokens = tokens
		stmts := parser.parse()
		NewInterpreter().Interpret(stmts)
		numCreated := len(createdEnvironments)
		numExpected := len(test.expected)
		if numCreated != numExpected {
			msg := "Expected %d environments, but %d were created."
			t.Errorf(msg, numExpected, numCreated)
		}
		for i, env := range createdEnvironments {
			expectedEnv := test.expected[i]
			for varname, value := range env.values {
				expectedValue := expectedEnv[varname]
				expectedType := reflect.TypeOf(expectedValue)
				actualType := reflect.TypeOf(value)
				if expectedType != actualType {
					msg := "Expected type %v for var %s in env %d, got %v\n"
					t.Errorf(msg, expectedType, varname, i, actualType)
					continue
				}
				if reflect.ValueOf(value).Interface() != reflect.ValueOf(expectedValue).Interface() {
					msg := "Expected %s=%v in env %d, got %v\n"
					t.Errorf(msg, varname, expectedEnv[varname], i, value)
				}
			}
		}
		createdEnvironments = []*Environment{}
	}
	NewEnvironment = previous
}

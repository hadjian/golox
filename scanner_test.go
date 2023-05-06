package main

import (
	"testing"
)

func TestMatchToken(t *testing.T) {
	tests := []struct {
		expected rune
		true     TokenType
		false    TokenType
		want     TokenType
	}{
		{'!', BANG, BANG_EQUAL, BANG},
		{'=', EQUAL, EQUAL_EQUAL, EQUAL},
	}

	s := NewScanner("!=")
	for _, test := range tests {
		got := s.matchToken(test.expected, test.true, test.false)
		if got != test.want {
			t.Errorf("matchToken(%s)=%s", test.true, got)
		}
	}
}

func TestScanner(t *testing.T) {
	code := "!====\n// laskdjflsadjf \n/>."
	expected := []TokenType{
		BANG_EQUAL,
		EQUAL_EQUAL,
		EQUAL,
		SLASH,
		GREATER,
		DOT,
		EOF,
	}
	s := NewScanner(code)
	s.scanTokens()
	if len(s.tokens) == 0 {
		t.Errorf("No tokens were scanned.")
	}
	for i, token := range s.tokens {
		if token.tType != expected[i] {
			t.Errorf("tokens[%d] = %s, expected %s", i, token.tType, expected[i])
		}

	}
}

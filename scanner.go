package main

import (
	"fmt"
	"strconv"
)

var keywords map[string]TokenType

func init() {
	keywords = map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
}

type Scanner struct {
	Source  []rune
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Source: []rune(source),
		line:   1,
	}
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) scanToken() {
	r := s.advance()
	switch r {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		s.addToken(s.matchToken('=', BANG_EQUAL, BANG))
	case '=':
		s.addToken(s.matchToken('=', EQUAL_EQUAL, EQUAL))
	case '<':
		s.addToken(s.matchToken('=', LESS_EQUAL, LESS))
	case '>':
		s.addToken(s.matchToken('=', GREATER_EQUAL, GREATER))
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
		fallthrough
	case '\r':
		fallthrough
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(r) {
			s.number()
		} else if s.isAlpha(r) {
			s.identifier()
		} else {
			msg := fmt.Sprintf("Unexpected character %c", r)
			errLine(s.line, msg)
		}
	}
}

func (s *Scanner) advance() rune {
	r := s.Source[s.current]
	s.current++
	return r
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t TokenType, literal any) {
	text := s.Source[s.start:s.current]
	s.tokens = append(s.tokens, Token{t, string(text), literal, s.line})
}

func (s *Scanner) matchToken(expected rune, true, false TokenType) TokenType {
	matchedToken := true
	if !s.match(expected) {
		matchedToken = false
	}
	return matchedToken
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.Source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\u0000'
	}
	return s.Source[s.current]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		errLine(s.line, "Unterminated string.")
		return
	}

	s.advance()

	value := string(s.Source[s.start+1 : s.current-1])
	s.addTokenWithLiteral(STRING, value)
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.Source) {
		return '\u0000'
	}
	return s.Source[s.current+1]
}

func (s *Scanner) isDigit(r rune) bool {
	return (r >= '0' && r <= '9')
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}

	}

	substring := string(s.Source[s.start:s.current])
	value, _ := strconv.ParseFloat(substring, 64)
	s.addTokenWithLiteral(NUMBER, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.Source[s.start:s.current]
	tokenType, ok := keywords[string(text)]
	if ok == false {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType)
}

func (s *Scanner) isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r == '_')
}

func (s *Scanner) isAlphaNumeric(r rune) bool {
	return s.isAlpha(r) || s.isDigit(r)
}

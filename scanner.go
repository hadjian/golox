package main

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
	default:
		err(s.line, "Unexpected character.")
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

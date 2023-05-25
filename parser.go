// This is the parser implementation of the Lox language
// It implements the following grammar
//
//			  program    -> statement* EOF;
//			  decl       -> varDecl | statement;
//			  varDecl    -> "var" IDENTIFIER ( "=" expression )? ";";
//			  statement  -> exprStmt | printStmt;
//			  exprStmt   -> expression ";";
//			  printStmt  -> "print" expression ";";
//			  expression -> equality;
//			  equality   -> comparison ( ( "!=" | "==" ) ) comparison )*;
//			  comparison -> term ( ( ">" | "<" | ">=" | "<=" ) term)*;
//			  term       -> factor ( ( "+" | "-" ) factor)*;
//				factor     -> unary ( ( "/" | "*" ) unary )*;
//				unary      -> ( "!" | "-") unary | primary;
//				primary    -> NUMBER     |
//		                  STRING     |
//						          "true"     |
//						          "false"    |
//						          "nil"      |
//	                   IDENTIFIER |
//						          "("expression")";
package main

type ParseError struct {
	msg string
}

func (p *ParseError) Error() string {
	return p.msg
}

type Parser struct {
	Tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens: tokens,
	}
}

func (p *Parser) parse() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		if stmt, err := p.varDeclaration(); err != nil {
			p.synchronize()
			return nil
		} else {
			return stmt
		}
	}
	if stmt, err := p.statement(); err != nil {
		p.synchronize()
		return nil
	} else {
		return stmt
	}
}

func (p *Parser) varDeclaration() (Stmt, error) {
	errMsg := "Expected identifier after 'var'."
	var varID Token
	var err error
	if varID, err = p.consume(IDENTIFIER, errMsg); err != nil {
		return nil, err
	}
	// Check if there is an initializer expression
	var initializer Expr
	if p.match(EQUAL) {
		if initializer, err = p.expression(); err != nil {
			return nil, err

		}
	}
	errMsg = "Expected ';' after variable declaration."
	if _, err = p.consume(SEMICOLON, errMsg); err != nil {
		return nil, err
	}
	return &Var{varID, initializer}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	if value, err := p.expression(); err != nil {
		return nil, err
	} else {
		p.consume(SEMICOLON, "Expect ';' after value.")
		return &Print{value}, nil
	}
}

func (p *Parser) expressionStatement() (Stmt, error) {
	if expr, err := p.expression(); err != nil {
		return nil, err
	} else {
		p.consume(SEMICOLON, "Expected ';' after expression.")
		return &Expression{expr}, nil
	}
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	var expr Expr
	var err error

	if expr, err = p.comparison(); err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		op := p.previous()
		var right Expr
		if right, err = p.comparison(); err != nil {
			break
		}
		expr = &Binary{expr, op, right}
	}

	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	var err error
	var expr Expr
	if expr, err = p.term(); err != nil {
		return nil, err
	}
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		var right Expr
		op := p.previous()
		if right, err = p.term(); err != nil {
			return nil, err
		}
		expr = &Binary{expr, op, right}
	}
	return expr, nil
}

// match is like check, but also advances the parser.
//
// Use if one of a set of token types is expected, mostly when the
// production rule contains alternative terminals, e.g. ( "!" | "-").
//
// TODO: why doesn't this return the token? Wouldn't this make previous
// unnecessary?
func (p *Parser) match(types ...TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) term() (Expr, error) {
	var expr Expr
	var err error

	if expr, err = p.factor(); err != nil {
		return nil, err
	}

	for p.match(PLUS, MINUS) {
		op := p.previous()
		var right Expr
		if right, err = p.factor(); err != nil {
			return nil, err
		}
		expr = &Binary{expr, op, right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	var expr Expr
	var err error

	if expr, err = p.unary(); err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		op := p.previous()
		var right Expr
		if right, err = p.unary(); err != nil {
			return nil, err
		}
		expr = &Binary{expr, op, right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(MINUS, BANG) {
		var err error
		var expr Expr

		op := p.previous()
		if expr, err = p.unary(); err != nil {
			return nil, err
		}
		return &Unary{op, expr}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	switch {
	case p.match(FALSE):
		return &Literal{false}, nil
	case p.match(TRUE):
		return &Literal{true}, nil
	case p.match(NIL):
		return &Literal{nil}, nil
	case p.match(NUMBER, STRING):
		return &Literal{p.previous().literal}, nil
	case p.match(IDENTIFIER):
		return &Variable{p.previous()}, nil
	case p.match(LEFT_PAREN):
		var expr Expr
		var err error
		if expr, err = p.expression(); err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return expr, nil
	}
	err := p.err(p.peek(), "Expect expression.")
	return nil, &err
}

// consume checks if the current token is of the expected type and
// returns it or prints the error message and returns a ParseError.
//
// This method should be used, if a certain token must appear next, like
// a closing parenthesis.
func (p *Parser) consume(typ TokenType, message string) (Token, error) {
	if p.check(typ) {
		return p.advance(), nil
	}
	err := p.err(p.peek(), message)
	return Token{}, &err
}

func (p *Parser) err(token Token, message string) ParseError {
	errToken(token, message)
	return ParseError{
		message,
	}
}

// check checks if the current token is of type typ
//
// Used in match(), where the terminal might appear and in consume(),
// where a terminal must appear. The latter one throws an error, if it
// doesn't.
func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tType == typ
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tType == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().tType == SEMICOLON {
			return
		}
		switch p.peek().tType {
		case CLASS, FOR, FUN, IF, PRINT, RETURN, VAR, WHILE:
			return
		}
		p.advance()
	}
}

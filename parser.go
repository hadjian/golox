// This is the parser implementation of the Lox language
// It implements the following grammar
//
// program    -> decl* EOF;
// decl       -> varDecl | statement;
// varDecl    -> "var" IDENTIFIER ( "=" expression )? ";";
// statement  -> exprStmt | forStmt | ifStmt | printStmt | whileStmt | block;
// whileStmt  -> "while" "(" expression ")" statement;
// ifStmt     -> "if" "(" expression ")" statement ("else" statement )?;
// block      -> "{" declaration* "}";
// exprStmt   -> expression ";";
// forStmt    -> "for" "(" ( varDecl | exprStmt | ";" )
//
//		             expression? ";"
//	               expression? ")" statement;
//
// printStmt  -> "print" expression ";";
// expression -> assignment;
// assignment -> IDENTIFIER "=" assignment | logic_or;
// logic_or   -> logic_and ( "or" logic_and)*;
// logic_and  -> equality ( "and" equality)*;
// equality   -> comparison ( ( "!=" | "==" ) ) comparison )*;
// comparison -> term ( ( ">" | "<" | ">=" | "<=" ) term)*;
// term       -> factor ( ( "+" | "-" ) factor)*;
// factor     -> unary ( ( "/" | "*" ) unary )*;
// unary      -> ( "!" | "-") unary | call;
// call       -> primary ( "(" arguments? ")" )*;
// arguments  -> expression ( "," expression )*;
// primary    -> NUMBER   |
//
//		         STRING     |
//						 "true"     |
//						 "false"    |
//						 "nil"      |
//	           IDENTIFIER |
//						 "("expression")";
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

func (p *Parser) declaration() (stmt Stmt) {
	defer func() {
		if err := recover(); err != nil {
			p.synchronize()
			stmt = nil
		}
	}()
	if p.match(VAR) {
		stmt := p.varDeclaration()
		return stmt
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	errMsg := "Expected identifier after 'var'."
	varID := p.consume(IDENTIFIER, errMsg)
	// Check if there is an initializer expression
	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	errMsg = "Expected ';' after variable declaration."
	p.consume(SEMICOLON, errMsg)
	return &Var{varID, initializer}
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(LEFT_BRACE) {
		return p.block()
	}
	if p.match(WHILE) {
		return p.whileStmt()
	}
	if p.match(FOR) {
		return p.forStmt()
	}
	if p.match(IF) {
		return p.ifStmt()
	}
	return p.expressionStmt()
}

func (p *Parser) block() Stmt {
	var statements []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	msg := "Expected '}' after block."
	p.consume(RIGHT_BRACE, msg)
	return &Block{statements}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return &Print{value}
}

func (p *Parser) expressionStmt() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expected ';' after expression.")
	return &Expression{expr}
}

func (p *Parser) whileStmt() Stmt {
	p.consume(LEFT_PAREN, "Expected '(' after while statement.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expected ')' after while condition.")
	body := p.statement()
	return &While{condition, body}
}

func (p *Parser) forStmt() Stmt {
	p.consume(LEFT_PAREN, "Expected '(' after 'for'.")
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStmt()
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}

	p.consume(RIGHT_PAREN, "Expect ')' after for clause.")

	body := p.statement()

	if increment != nil {
		body = &Block{
			[]Stmt{body, &Expression{increment}},
		}
	}

	if condition == nil {
		condition = &Literal{true}
	}
	body = &While{condition, body}

	if initializer != nil {
		body = &Block{
			[]Stmt{initializer, body},
		}
	}

	return body
}

func (p *Parser) ifStmt() Stmt {
	p.consume(LEFT_PAREN, "Expected opening '(' after 'if'.")
	expr := p.expression()
	p.consume(RIGHT_PAREN, "Expected closing ')' after 'if' expression.")
	stmt := p.statement()
	var elseStmt Stmt
	if p.match(ELSE) {
		elseStmt = p.statement()
	}
	return &If{expr, stmt, elseStmt}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.logic_or()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if expr, ok := expr.(*Variable); !ok {
			panic(p.err(equals, "Invalid assignment target."))
		} else {
			return &Assign{expr.name, value}
		}
	}
	return expr
}

func (p *Parser) logic_or() Expr {
	expr := p.logic_and()
	for p.match(OR) {
		operator := p.previous()
		right := p.logic_and()
		expr = &Logical{expr, operator, right}
	}
	return expr
}

func (p *Parser) logic_and() Expr {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = &Logical{expr, operator, right}
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		op := p.previous()
		right := p.comparison()
		expr = &Binary{expr, op, right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right := p.term()
		expr = &Binary{expr, op, right}
	}
	return expr
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

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(PLUS, MINUS) {
		op := p.previous()
		right := p.factor()
		expr = &Binary{expr, op, right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		op := p.previous()
		right := p.unary()
		expr = &Binary{expr, op, right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(MINUS, BANG) {
		op := p.previous()
		expr := p.unary()
		return &Unary{op, expr}
	}
	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for true {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	var arguments []Expr
	if !p.check(RIGHT_PAREN) {
		for ok := true; ok; p.match(COMMA) {
			if len(arguments) >= 255 {
				p.err(p.peek(), "Can't have more than 255 arguments")
			}
			arguments = append(arguments, p.expression())
		}
	}

	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")

	return &Call{callee, paren, arguments}
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(FALSE):
		return &Literal{false}
	case p.match(TRUE):
		return &Literal{true}
	case p.match(NIL):
		return &Literal{nil}
	case p.match(NUMBER, STRING):
		return &Literal{p.previous().literal}
	case p.match(IDENTIFIER):
		return &Variable{p.previous()}
	case p.match(LEFT_PAREN):
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression")
		return expr
	}
	panic(p.err(p.peek(), "Expect expression."))
}

// consume checks if the current token is of the expected type and
// returns it or prints the error message and returns a ParseError.
//
// This method should be used, if a certain token must appear next, like
// a closing parenthesis.
func (p *Parser) consume(typ TokenType, message string) Token {
	if p.check(typ) {
		return p.advance()
	}
	panic(p.err(p.peek(), message))
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

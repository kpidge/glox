package lox

type Parser struct {
	tokens []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

// Wrapper for parser errors
type parserError struct { error }

func (e parserError) Error() string {
	return "Encountered error during parsing"
}

func (p *Parser) Parse() (statements []Stmt, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Determine that we are recovering from a 
			// parserError
			if _, ok := r.(parserError); ok {
				p.synchronise()
			} else {
				panic(r)
			}
		}

	}()
	// Default cap necessary?
	statements = make([]Stmt, 0, 100)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements, nil
}

func (p *Parser) declaration() Stmt {
	if p.match(FUN) {
		function := p.function("function")
		return function
	}
	if p.match(VAR) {
		return p.variableDecl()
	}
	return p.statement()
}

func (p *Parser) function(kind string) *FunctionStmt{
	name := p.consume(IDENTIFIER, "Expect " + kind + " name")
	p.consume(LEFT_PAREN, "Expect '(' after " + kind + " name")
	var params []Token
	if !p.check(RIGHT_PAREN) {
		params = append(params, p.consume(IDENTIFIER, "Expect parameter name"))
		for p.match(COMMA) {
			if len(params) >= 255 {
				ErrorOnToken(p.peek(), "Can't have more than 255 parameters")
			}
			params = append(params, p.consume(IDENTIFIER, "Expect parameter name"))
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters")
	p.consume(LEFT_BRACE, "Expect '{' before " + kind + " body")
	body := p.block()
	return &FunctionStmt{name: name, params: params, body: body}
}

func (p *Parser) variableDecl() *VarStmt {
	name := p.consume(IDENTIFIER, "Expect identifier after 'var' keyword")

	var init Expr
	if p.match(EQUAL) {
		init = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after variable declaration")

	return &VarStmt{Name: name, Initialiser: init}
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(LEFT_BRACE) {
		return &BlockStmt{statements: p.block()}
	}

	// Fallthrough case, as difficult to detect based on token
	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value")
	return &PrintStmt{Expr: expr}
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' following 'if'")
	expr := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' following condition")
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return &IfStmt{expr: expr, thenBranch: thenBranch, elseBranch: elseBranch}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' following 'while'")
	expr := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' following condition")
	body := p.statement()
	return &WhileStmt{expr: expr, body: body}
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' following 'for'")
	var init Stmt
	var cond Expr
	var inc Expr
	if p.match(VAR) {
		init = p.variableDecl()
	} else {
		if !p.match(SEMICOLON) {
			init = p.expressionStatement()
		}
	}
	if !p.check(SEMICOLON) {
		cond = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after condition")
	if !p.check(RIGHT_PAREN) {
		inc = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clause")

	// Consume the for block statement
	body := p.statement()
	if inc != nil {
		body = &BlockStmt{statements: []Stmt{body, &ExpressionStmt{Expr: inc}}}
	}
	if cond == nil {
		cond = &Literal{Value: true}
	}
	body = &WhileStmt{expr: cond, body: body}
	if init != nil {
		body = &BlockStmt{statements: []Stmt{init, body}}
	}

	return body
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		value = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after return value")
	return &ReturnStmt{keyword: keyword, value: value}
}

func (p *Parser) block() []Stmt {
	var statements []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expect '}' after block")
	return statements
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression")
	return &ExpressionStmt{Expr: expr}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.logicalOr()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if v, ok := expr.(*Variable); ok {
			return &Assign{Name: v.Name, Value: value}
		}
		ErrorOnToken(equals, "Invalid assignment target")
	}
	return expr
}

func (p *Parser) logicalOr() Expr {
	expr := p.logicalAnd()
	for p.match(OR) {
		op := p.previous()
		right := p.logicalAnd()
		return &Logical{op:op, left: expr, right: right}
	}
	return expr
}

func (p *Parser) logicalAnd() Expr {
	expr := p.equality()
	for p.match(AND) {
		op := p.previous()
		right := p.equality()
		return &Logical{op:op, left: expr, right: right}
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		op := p.previous()
		right := p.comparison()
		expr = &Binary{Left: expr, Op: op, Right: right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right := p.term()
		expr = &Binary{Left: expr, Op: op, Right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		op := p.previous()
		right := p.factor()
		expr = &Binary{Left: expr, Op: op, Right: right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(STAR, SLASH) {
		op := p.previous()
		right := p.unary()
		expr = &Binary{Left: expr, Op: op, Right: right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		op := p.previous()
		right := p.unary()
		return &Unary{Op: op, Right: right}
	}
	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(LEFT_PAREN) {
			return p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	var args []Expr
	if !p.check(RIGHT_PAREN) {
		args = append(args, p.expression())
		for p.match(COMMA) {
			if len(args) >= 255 {
				ErrorOnToken(p.peek(), "Can't have more than 255 arguments")
			}
			args = append(args, p.expression())
		}
	}
	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments")
	return &Call{paren: paren, callee: callee, arguments: args}
}

func (p *Parser) primary() Expr {
	if p.match(TRUE) { return &Literal{Value: true} }
	if p.match(FALSE) { return &Literal{Value: false} }
	if p.match(NIL) { return &Literal{Value: nil} }
	if p.match(NUMBER, STRING) {
		return &Literal{Value: p.previous().Literal}
	}
	if p.match(IDENTIFIER) {
		return &Variable{Name: p.previous()}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &Grouping{Expr: expr}
	}
	panic(p.parserError(p.peek(), "Expect expression."))
}

func (p *Parser) match(comps ...TokenType) bool {
	for _, t := range comps {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(ttype TokenType, msg string) Token {
	if p.check(ttype) { return p.advance() }

	panic(p.parserError(p.peek(), msg))
}

func (p *Parser) check(ttype TokenType) bool {
	if p.isAtEnd() { return false }
	return p.peek().Type == ttype
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() { p.current += 1 }
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.tokens[p.current].Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func(p *Parser) previous() Token {
	return p.tokens[p.current - 1]
}

// Synchronise sorts out the state of the parser when recovering
// from a parser error, by advancing to the next statement
func (p *Parser) synchronise() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}

func(p *Parser) parserError(token Token, msg string) error {
	ErrorOnToken(token, msg)
	return parserError{}
}

package lox


type Parser struct {
	tokens []Token
	current int
}

// Wrapper for parser errors
type parserError struct { error }

func (e parserError) Error() string {
	return "Encountered error during parsing"
}

func (p *Parser) Parse() (expr Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Determine that we are recovering from a 
			// parserError
			if pe, ok := r.(parserError); ok {
				err = pe.error
			} else {
				panic(r)
			}
		}

	}()
	return p.expression(), nil
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) expression() Expr {
	return p.equality()
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
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(TRUE) { return &Literal{Value: true} }
	if p.match(FALSE) { return &Literal{Value: false} }
	if p.match(NIL) { return &Literal{Value: nil} }
	if p.match(NUMBER, STRING) {
		return &Literal{Value: p.previous().Literal}
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

func(p *Parser) parserError(token Token, msg string) error {
	ErrorOnToken(token, msg)
	return parserError{}
}

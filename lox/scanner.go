package lox

type Scanner struct {
	Source string
	tokens []Token
	start int
	current int
	line int
}

func NewScanner(source string) *Scanner {
	return &Scanner{Source: source, start: 0, current: 0, line: 1}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{Type: EOF, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
		case '(': s.addToken(LEFT_PAREN)
		case ')': s.addToken(RIGHT_PAREN)
		case '{': s.addToken(LEFT_BRACE)
		case '}': s.addToken(RIGHT_BRACE)
		case ',': s.addToken(COMMA)
		case '.': s.addToken(DOT)
		case '-': s.addToken(MINUS)
		case '+': s.addToken(PLUS)
		case ';': s.addToken(SEMICOLON)
		case '*': s.addToken(STAR)
		default: Error(s.line, "Unexpected character")
	}
}


func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) advance() rune {
	current := s.current
	s.current += 1
	return rune(s.Source[current])
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t TokenType, literal any) {
	text := s.Source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Type: t, Lexeme: text, Literal: literal, Line: s.line})
}

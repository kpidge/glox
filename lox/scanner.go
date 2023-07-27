package lox

import (
	"fmt"
	"strconv"
)

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
		case '!':
			token := BANG
			if s.match('=') {
				token = BANG_EQUAL
			}
			s.addToken(token)
		case '=':
			token := EQUAL
			if s.match('=') {
				token = EQUAL_EQUAL
			}
			s.addToken(token)
		case '<':
			token := LESS
			if s.match('=') {
				token = LESS_EQUAL
			}
			s.addToken(token)
		case '>':
			token := GREATER
			if s.match('=') {
				token = GREATER_EQUAL
			}
			s.addToken(token)
		case '/':
			if s.match('/') {
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				s.addToken(SLASH)
			}
			// Either this, or use s.match instead of peek above
			fallthrough
		case ' ':
		case '\r':
		case '\t':
		case '\n': s.line += 1
		case '"': s.string()
		default: // This is checked here to save adding a case for each digit
			if isDigit(c) {
				s.number()
			} else {
				Error(s.line, "Unexpected character")
			}
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

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() { return false }
	if rune(s.Source[s.current]) != expected { return false }

	s.current += 1
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() { return '\u0000' }
	return rune(s.Source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current + 1 >= len(s.Source) { return '\u0000' }
	return rune(s.Source[s.current+1])
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' { s.line += 1 }
		s.advance()
	}

	if s.isAtEnd() {
		Error(s.line, "Unterminated string")
		return
	}

	s.advance() // Closing "
	value := s.Source[s.start+1:s.current-1]
	s.addTokenWithLiteral(STRING, value)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) { s.advance() }

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) { s.advance() }
	}
	f64, err := strconv.ParseFloat(s.Source[s.start:s.current], 64)
	if err != nil {
		Error(s.line, fmt.Sprint(err))
	} else {
		s.addTokenWithLiteral(NUMBER, f64)
	}

}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t TokenType, literal any) {
	text := s.Source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Type: t, Lexeme: text, Literal: literal, Line: s.line})
}

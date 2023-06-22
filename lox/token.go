package lox

import "fmt"

type Token struct {
	Type TokenType
	Lexeme string
	Literal any
	Line int
}

type Repr interface {
	ToString()
}

func (t *Token) ToString() string {
	return fmt.Sprint(t.Type, " " + t.Lexeme + " ", t.Literal)
}

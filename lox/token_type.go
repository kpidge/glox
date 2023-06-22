package lox

type TokenType int

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN TokenType = iota
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE TokenType = iota
	COMMA TokenType = iota
	DOT TokenType = iota
	MINUS TokenType = iota
	PLUS TokenType = iota
	SEMICOLON TokenType = iota
	SLASH TokenType = iota
	STAR TokenType = iota

	// One or two character tokens.
	BANG TokenType = iota
	BANG_EQUAL TokenType = iota
	EQUAL TokenType = iota
	EQUAL_EQUAL TokenType = iota
	GREATER TokenType = iota
	GREATER_EQUAL TokenType = iota
	LESS TokenType = iota
	LESS_EQUAL TokenType = iota

	// Literals. TokenType = iota
	IDENTIFIER TokenType = iota
	STRING TokenType = iota
	NUMBER TokenType = iota

	// Keywords. TokenType = iota
	AND TokenType = iota
	CLASS TokenType = iota
	ELSE TokenType = iota
	FALSE TokenType = iota
	FUN TokenType = iota
	FOR TokenType = iota
	IF TokenType = iota
	NIL TokenType = iota
	OR TokenType = iota
	PRINT TokenType = iota
	RETURN TokenType = iota
	SUPER TokenType = iota
	THIS TokenType = iota
	TRUE TokenType = iota
	VAR TokenType = iota
	WHILE TokenType = iota
	EOF TokenType = iota
)

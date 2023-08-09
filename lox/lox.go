package lox

import (
	"fmt"
	"os"
	"bufio"
)

var hadError = false

func RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		os.Exit(1)
	}
	run(string(bytes))
	if hadError {
		os.Exit(65)
	}
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(1)
		}
		run(line)
		hadError = false
	}
}

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := NewParser(tokens)
	expr, err := parser.Parse()

	if hadError || err != nil { panic("Error while parsing") }

	astPrinter := ASTPrinter{}
	astPrinter.PrintAST(expr)
}

func Error(line int, msg string) {
	report(line, "", msg)
}

func report(line int, where string, msg string) {
	fmt.Println(fmt.Sprint("[line ", line, "] Error", where, ": ", msg))
	hadError = true
}

func ErrorOnToken(token Token, msg string) {
	if token.Type == EOF {
		report(token.Line, " at end", msg)
	} else {
		report(token.Line, " at '" + token.Lexeme + "'", msg)
	}
}

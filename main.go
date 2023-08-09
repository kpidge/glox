package main

import (
	"fmt"
	"glox/lox"
	"os"
)

func main() {

	if len(os.Args) > 2 {
		fmt.Fprint(os.Stderr, "Usage: glox [script]\n")
		os.Exit(64)
	} else if args := os.Args; len(args) == 2 {
		lox.RunFile(args[1])
	} else {
		lox.RunPrompt()
	}
}


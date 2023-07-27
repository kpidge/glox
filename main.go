package main

import (
	"fmt"
	"glox/lox"
	"os"
)

func main() {
	fmt.Println("Hello world")
	fmt.Println(os.Args)
	if len(os.Args) > 2 {
		fmt.Fprint(os.Stderr, "Usage: glox [script]\n")
		os.Exit(64)
	} else if args := os.Args; len(args) == 2 {
		fmt.Println("Scanning file")
		lox.RunFile(args[1])
	} else {
		lox.RunPrompt()
	}
}


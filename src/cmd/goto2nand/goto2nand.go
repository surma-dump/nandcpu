package main

import (
	"nandcpu"
	"os"
	"bufio"
	"fmt"
	)

func main() {
	lexer := nandcpu.NewLexer(bufio.NewReader(os.Stdin))
	for {
		t, e := lexer.GetNextToken()
		if e != nil {
			panic(e)
		}
		fmt.Printf("Token: \"%s\" (Type 0x%x)\n", t.Match, t.Type)
	}
}
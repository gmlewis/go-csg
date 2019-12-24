// Package repl implements a Read, Eval, Print, Loop for our language.
package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/gmlewis/go-monkey/lexer"
	"github.com/gmlewis/go-monkey/token"
)

// PROMPT is the repl prompt.
const PROMPT = ">> "

// Start starts the repl.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		le := lexer.New(line)

		for tok := le.NextToken(); tok.Type != token.EOF; tok = le.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}

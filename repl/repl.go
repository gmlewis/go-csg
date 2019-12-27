// Package repl implements a Read, Eval, Print, Loop for our language.
package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/gmlewis/go-monkey/evaluator"
	"github.com/gmlewis/go-monkey/lexer"
	"github.com/gmlewis/go-monkey/parser"
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
		p := parser.New(le)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

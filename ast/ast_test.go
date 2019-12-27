package ast

import (
	"testing"

	"github.com/gmlewis/go-monkey/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if got, want := program.String(), "let myVar = anotherVar;"; got != want {
		t.Errorf("program.String() = %v, want %v", got, want)
	}
}

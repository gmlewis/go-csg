package lexer

import (
	"fmt"
	"testing"

	"github.com/gmlewis/go-monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.T
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	le := New(input)

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			tok := le.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tokentype = %q, want %q", tok.Type, tt.expectedType)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("literal = %q, want %q", tok.Literal, tt.expectedLiteral)
			}
		})
	}
}

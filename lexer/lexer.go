// Package lexer implements the language lexer.
package lexer

import (
	"log"

	"github.com/gmlewis/go-monkey/token"
)

// Lexer represents the lexer.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New returns a new Lexer.
func New(input string) *Lexer {
	le := &Lexer{input: input}
	le.readChar()
	return le
}

func (le *Lexer) readChar() {
	if le.readPosition >= len(le.input) {
		le.ch = 0
	} else {
		le.ch = le.input[le.readPosition]
	}
	le.position = le.readPosition
	le.readPosition++
}

// NextToken returns the next token.
func (le *Lexer) NextToken() token.Token {
	var tok token.Token
	switch le.ch {
	case '=':
		tok = newToken(token.ASSIGN, le.ch)
	case ';':
		tok = newToken(token.SEMICOLON, le.ch)
	case '(':
		tok = newToken(token.LPAREN, le.ch)
	case ')':
		tok = newToken(token.RPAREN, le.ch)
	case ',':
		tok = newToken(token.COMMA, le.ch)
	case '+':
		tok = newToken(token.PLUS, le.ch)
	case '{':
		tok = newToken(token.LBRACE, le.ch)
	case '}':
		tok = newToken(token.RBRACE, le.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		log.Fatalf("unknown token %c", le.ch)
	}
	le.readChar()
	return tok
}

func newToken(tokenType token.T, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

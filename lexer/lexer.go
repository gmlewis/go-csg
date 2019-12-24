// Package lexer implements the language lexer.
package lexer

import (
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

	le.skipWhitespace()

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
	case '-':
		tok = newToken(token.MINUS, le.ch)
	case '!':
		tok = newToken(token.BANG, le.ch)
	case '/':
		tok = newToken(token.SLASH, le.ch)
	case '*':
		tok = newToken(token.ASTERISK, le.ch)
	case '<':
		tok = newToken(token.LT, le.ch)
	case '>':
		tok = newToken(token.GT, le.ch)
	case '{':
		tok = newToken(token.LBRACE, le.ch)
	case '}':
		tok = newToken(token.RBRACE, le.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(le.ch) {
			tok.Literal = le.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}
		if isDigit(le.ch) {
			tok.Type = token.INT
			tok.Literal = le.readNumber()
			return tok
		}
		tok = newToken(token.ILLEGAL, le.ch)
	}
	le.readChar()
	return tok
}

func newToken(tokenType token.T, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (le *Lexer) readIdentifier() string {
	position := le.position
	for isLetter(le.ch) {
		le.readChar()
	}
	return le.input[position:le.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (le *Lexer) skipWhitespace() {
	for le.ch == ' ' || le.ch == '\t' || le.ch == '\n' || le.ch == '\r' {
		le.readChar()
	}
}

func (le *Lexer) readNumber() string {
	position := le.position
	for isDigit(le.ch) {
		le.readChar()
	}
	return le.input[position:le.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

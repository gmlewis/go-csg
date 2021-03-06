// Package lexer implements the language lexer.
package lexer

import (
	"github.com/gmlewis/go-csg/token"
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
		if le.peekChar() == '=' {
			le.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, le.ch)
		}
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
		if le.peekChar() == '=' {
			le.readChar()
			tok = token.Token{Type: token.NOTEQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, le.ch)
		}
	case '/':
		if le.peekChar() == '/' {
			le.readChar()
			tok.Type = token.LINECOMMENT
			tok.Literal = le.readLineComment()
		} else {
			tok = newToken(token.SLASH, le.ch)
		}
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
	case '"':
		tok.Type = token.STRING
		tok.Literal = le.readString()
	case '[':
		tok = newToken(token.LBRACKET, le.ch)
	case ']':
		tok = newToken(token.RBRACKET, le.ch)
	case ':':
		tok = newToken(token.COLON, le.ch)
	// case '$':
	// 	tok = newToken(token.DOLLAR, le.ch)
	default:
		if isLetter(le.ch) || le.ch == '$' {
			tok.Literal = le.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}
		if isDigit(le.ch) {
			tok.Literal, tok.Type = le.readNumber()
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
	for isLetter(le.ch) || le.ch == '$' || isDigit(le.ch) {
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

func (le *Lexer) readNumber() (string, token.T) {
	position := le.position
	var tokType token.T = token.INT
	for isDigit(le.ch) || le.ch == '.' || le.ch == 'e' || le.ch == '-' {
		if le.ch == '.' || le.ch == 'e' {
			tokType = token.FLOAT
		}
		le.readChar()
	}
	return le.input[position:le.position], tokType
}

func (le *Lexer) readString() string {
	position := le.position + 1
	for {
		le.readChar()
		if le.ch == '"' || le.ch == 0 {
			break
		}
	}
	return le.input[position:le.position]
}

func (le *Lexer) readLineComment() string {
	position := le.position + 1
	for {
		le.readChar()
		if le.ch == '\n' || le.ch == 0 {
			break
		}
	}
	return le.input[position:le.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (le *Lexer) peekChar() byte {
	if le.readPosition >= len(le.input) {
		return 0
	}
	return le.input[le.readPosition]
}

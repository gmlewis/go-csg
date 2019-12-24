// Package parser implements a parser for our programming language.
package parser

import (
	"github.com/gmlewis/go-monkey/ast"
	"github.com/gmlewis/go-monkey/lexer"
	"github.com/gmlewis/go-monkey/token"
)

// Parser represents our language parser.
type Parser struct {
	le *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

// New returns a new Parser.
func New(le *lexer.Lexer) *Parser {
	p := &Parser{le: le}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.le.NextToken()
}

// ParseProgram parses a program.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.T) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.T) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.T) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

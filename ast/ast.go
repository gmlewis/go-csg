// Package ast defines the Abstract Syntax Tree for our language.
package ast

import "github.com/gmlewis/go-monkey/token"

// Node represents a node in our abstract syntax tree.
type Node interface {
	TokenLiteral() string
}

// Statement represents a statement in our abstract syntax tree.
type Statement interface {
	Node
	statementNode()
}

// Expression respresents an expression in our abstract syntax tree.
type Expression interface {
	Node
	expressionNode()
}

// Program represents a program using our abstract syntax tree.
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the first token literal of the program.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// LetStatement represents a 'let' statement.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral returns the token literal.
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier represents an identifier.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns the token literal.
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// ReturnStatement respresents a 'return' statement.
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the token literal.
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// Package ast defines the Abstract Syntax Tree for our language.
package ast

import (
	"bytes"

	"github.com/gmlewis/go-monkey/token"
)

// Node represents a node in our abstract syntax tree.
type Node interface {
	String() string
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

// String returns the string representation of the Node.
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
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

// String returns the string representation of the Node.
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// TokenLiteral returns the token literal.
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier represents an identifier.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// String returns the string representation of the Node.
func (i *Identifier) String() string { return i.Value }

// TokenLiteral returns the token literal.
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// ReturnStatement respresents a 'return' statement.
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// String returns the string representation of the Node.
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// TokenLiteral returns the token literal.
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// ExpressionStatement represents an expression statement.
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// String returns the string representation of the Node.
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

// TokenLiteral returns the token literal.
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

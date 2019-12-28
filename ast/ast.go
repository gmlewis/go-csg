// Package ast defines the Abstract Syntax Tree for our language.
package ast

import (
	"bytes"
	"strings"

	"github.com/gmlewis/go-openscad/token"
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

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// String returns the string representation of the Node.
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// TokenLiteral returns the token literal.
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// StringLiteral represents a string literal.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// String returns the string representation of the Node.
func (sl *StringLiteral) String() string { return sl.Token.Literal }

// TokenLiteral returns the token literal.
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// PrefixExpression represents a prefix expression.
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// String returns the string representation of the Node.
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// TokenLiteral returns the token literal.
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// InfixExpression represents an infix expression.
type InfixExpression struct {
	Token    token.Token // The infix token, e.g. ==
	Left     Expression
	Operator string
	Right    Expression
}

func (pe *InfixExpression) expressionNode() {}

// String returns the string representation of the Node.
func (pe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(" " + pe.Operator + " ")
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// TokenLiteral returns the token literal.
func (pe *InfixExpression) TokenLiteral() string { return pe.Token.Literal }

// BooleanLiteral represents a boolean literal.
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}

// String returns the string representation of the Node.
func (bl *BooleanLiteral) String() string { return bl.Token.Literal }

// TokenLiteral returns the token literal.
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }

// IfExpression represents an "if" expression.
type IfExpression struct {
	Token       token.Token // The "if" token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// String returns the string representation of the Node.
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// TokenLiteral returns the token literal.
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// BlockStatement represents a block statement.
type BlockStatement struct {
	Token      token.Token // The "if" token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {}

// String returns the string representation of the Node.
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// TokenLiteral returns the token literal.
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// FunctionLiteral represents a function literal.
type FunctionLiteral struct {
	Token      token.Token // The "function" token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// String returns the string representation of the Node.
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// TokenLiteral returns the token literal.
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// CallExpression represents a call expression.
type CallExpression struct {
	Token     token.Token // The "(" token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// String returns the string representation of the Node.
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// TokenLiteral returns the token literal.
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// ArrayLiteral represents a prefix expression.
type ArrayLiteral struct {
	Token    token.Token // The '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

// String returns the string representation of the Node.
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// TokenLiteral returns the token literal.
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

// IndexExpression represents a prefix expression.
type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// String returns the string representation of the Node.
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// TokenLiteral returns the token literal.
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

// HashLiteral represents a prefix expression.
type HashLiteral struct {
	Token token.Token // The '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

// String returns the string representation of the Node.
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	var pairs []string
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// TokenLiteral returns the token literal.
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }

package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gmlewis/go-csg/token"
)

// NamedArgument represents a CSG named argument.
type NamedArgument struct {
	Token token.Token
	Name  Expression
	Value Expression
}

func (cp *NamedArgument) expressionNode() {}

// String returns the string representation of the Node.
func (cp *NamedArgument) String() string {
	return fmt.Sprintf("%v = %v", cp.Name.String(), cp.Value.String())
}

// TokenLiteral returns the token literal.
func (cp *NamedArgument) TokenLiteral() string { return cp.Token.Literal }

// CirclePrimitive represents a CSG primitive.
type CirclePrimitive struct {
	Token     token.Token
	Arguments []Expression
}

func (cp *CirclePrimitive) expressionNode() {}

func primitiveString(tokenLiteral string, args []Expression) string {
	var out bytes.Buffer

	var params []string
	for _, p := range args {
		params = append(params, p.String())
	}

	out.WriteString(tokenLiteral)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	return out.String()
}

// String returns the string representation of the Node.
func (cp *CirclePrimitive) String() string {
	return primitiveString(cp.TokenLiteral(), cp.Arguments)
}

// TokenLiteral returns the token literal.
func (cp *CirclePrimitive) TokenLiteral() string { return cp.Token.Literal }

// CubePrimitive represents a CSG primitive.
type CubePrimitive struct {
	Token     token.Token
	Arguments []Expression
}

func (cp *CubePrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (cp *CubePrimitive) String() string {
	return primitiveString(cp.TokenLiteral(), cp.Arguments)
}

// TokenLiteral returns the token literal.
func (cp *CubePrimitive) TokenLiteral() string { return cp.Token.Literal }

// CylinderPrimitive represents a CSG primitive.
type CylinderPrimitive struct {
	Token     token.Token
	Arguments []Expression
}

func (cp *CylinderPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (cp *CylinderPrimitive) String() string {
	return primitiveString(cp.TokenLiteral(), cp.Arguments)
}

// TokenLiteral returns the token literal.
func (cp *CylinderPrimitive) TokenLiteral() string { return cp.Token.Literal }

// GroupPrimitive represents a CSG primitive.
type GroupPrimitive struct {
	Token     token.Token
	Arguments []Expression
}

func (gp *GroupPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (gp *GroupPrimitive) String() string {
	return primitiveString(gp.TokenLiteral(), gp.Arguments)
}

// TokenLiteral returns the token literal.
func (gp *GroupPrimitive) TokenLiteral() string { return gp.Token.Literal }

// PolygonPrimitive represents a CSG primitive.
type PolygonPrimitive struct {
	Token     token.Token
	Arguments []Expression
}

func (pp *PolygonPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (pp *PolygonPrimitive) String() string {
	return primitiveString(pp.TokenLiteral(), pp.Arguments)
}

// TokenLiteral returns the token literal.
func (pp *PolygonPrimitive) TokenLiteral() string { return pp.Token.Literal }

// SpherePrimitive represents a CSG primitive.
type SpherePrimitive struct {
	Token     token.Token
	Arguments []Expression
}

func (sp *SpherePrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (sp *SpherePrimitive) String() string {
	return primitiveString(sp.TokenLiteral(), sp.Arguments)
}

// TokenLiteral returns the token literal.
func (sp *SpherePrimitive) TokenLiteral() string { return sp.Token.Literal }

// SquarePrimitive represents a CSG primitive.
type SquarePrimitive struct {
	Token     token.Token
	Arguments []Expression
}

func (sp *SquarePrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (sp *SquarePrimitive) String() string {
	return primitiveString(sp.TokenLiteral(), sp.Arguments)
}

// TokenLiteral returns the token literal.
func (sp *SquarePrimitive) TokenLiteral() string { return sp.Token.Literal }

// TextPrimitive represents a CSG primitive.
type TextPrimitive struct {
	Token     token.Token
	Arguments []Expression
}

func (tp *TextPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (tp *TextPrimitive) String() string {
	return primitiveString(tp.TokenLiteral(), tp.Arguments)
}

// TokenLiteral returns the token literal.
func (tp *TextPrimitive) TokenLiteral() string { return tp.Token.Literal }

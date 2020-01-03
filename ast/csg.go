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

// PolyhedronPrimitive represents a CSG primitive.
type PolyhedronPrimitive struct {
	Token     token.Token
	Arguments []Expression
}

func (pp *PolyhedronPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (pp *PolyhedronPrimitive) String() string {
	return primitiveString(pp.TokenLiteral(), pp.Arguments)
}

// TokenLiteral returns the token literal.
func (pp *PolyhedronPrimitive) TokenLiteral() string { return pp.Token.Literal }

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

// UndefLiteral represents a undef literal.
type UndefLiteral struct {
	Token token.Token
}

func (tp *UndefLiteral) expressionNode() {}

// String returns the string representation of the Node.
func (tp *UndefLiteral) String() string {
	return tp.TokenLiteral()
}

// TokenLiteral returns the token literal.
func (tp *UndefLiteral) TokenLiteral() string { return tp.Token.Literal }

func blockPrimitiveString(tokenLiteral string, args []Expression, block *BlockStatement) string {
	var out bytes.Buffer

	out.WriteString(primitiveString(tokenLiteral, args))
	if block != nil {
		out.WriteString(" { ")
		out.WriteString(block.String())
		out.WriteString(" }")
	}

	return out.String()
}

// ColorBlockPrimitive represents a CSG block primitive.
type ColorBlockPrimitive struct {
	Token     token.Token
	Arguments []Expression
	Body      *BlockStatement
}

func (cbp *ColorBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (cbp *ColorBlockPrimitive) String() string {
	return blockPrimitiveString(cbp.TokenLiteral(), cbp.Arguments, cbp.Body)
}

// TokenLiteral returns the token literal.
func (cbp *ColorBlockPrimitive) TokenLiteral() string { return cbp.Token.Literal }

// DifferenceBlockPrimitive represents a CSG block primitive.
type DifferenceBlockPrimitive struct {
	Token token.Token
	Body  *BlockStatement
}

func (dbp *DifferenceBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (dbp *DifferenceBlockPrimitive) String() string {
	return blockPrimitiveString(dbp.TokenLiteral(), nil, dbp.Body)
}

// TokenLiteral returns the token literal.
func (dbp *DifferenceBlockPrimitive) TokenLiteral() string { return dbp.Token.Literal }

// GroupBlockPrimitive represents a CSG block primitive.
type GroupBlockPrimitive struct {
	Token token.Token
	Body  *BlockStatement
}

func (gbp *GroupBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (gbp *GroupBlockPrimitive) String() string {
	return blockPrimitiveString(gbp.TokenLiteral(), nil, gbp.Body)
}

// TokenLiteral returns the token literal.
func (gbp *GroupBlockPrimitive) TokenLiteral() string { return gbp.Token.Literal }

// HullBlockPrimitive represents a CSG block primitive.
type HullBlockPrimitive struct {
	Token token.Token
	Body  *BlockStatement
}

func (hbp *HullBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (hbp *HullBlockPrimitive) String() string {
	return blockPrimitiveString(hbp.TokenLiteral(), nil, hbp.Body)
}

// TokenLiteral returns the token literal.
func (hbp *HullBlockPrimitive) TokenLiteral() string { return hbp.Token.Literal }

// IntersectionBlockPrimitive represents a CSG block primitive.
type IntersectionBlockPrimitive struct {
	Token token.Token
	Body  *BlockStatement
}

func (ibp *IntersectionBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (ibp *IntersectionBlockPrimitive) String() string {
	return blockPrimitiveString(ibp.TokenLiteral(), nil, ibp.Body)
}

// TokenLiteral returns the token literal.
func (ibp *IntersectionBlockPrimitive) TokenLiteral() string { return ibp.Token.Literal }

// LinearExtrudeBlockPrimitive represents a CSG block primitive.
type LinearExtrudeBlockPrimitive struct {
	Token     token.Token
	Arguments []Expression
	Body      *BlockStatement
}

func (lebp *LinearExtrudeBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (lebp *LinearExtrudeBlockPrimitive) String() string {
	return blockPrimitiveString(lebp.TokenLiteral(), lebp.Arguments, lebp.Body)
}

// TokenLiteral returns the token literal.
func (lebp *LinearExtrudeBlockPrimitive) TokenLiteral() string { return lebp.Token.Literal }

// MinkowskiBlockPrimitive represents a CSG block primitive.
type MinkowskiBlockPrimitive struct {
	Token     token.Token
	Arguments []Expression
	Body      *BlockStatement
}

func (mbp *MinkowskiBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (mbp *MinkowskiBlockPrimitive) String() string {
	return blockPrimitiveString(mbp.TokenLiteral(), mbp.Arguments, mbp.Body)
}

// TokenLiteral returns the token literal.
func (mbp *MinkowskiBlockPrimitive) TokenLiteral() string { return mbp.Token.Literal }

// MultmatrixBlockPrimitive represents a CSG block primitive.
type MultmatrixBlockPrimitive struct {
	Token     token.Token
	Arguments []Expression
	Body      *BlockStatement
}

func (mbp *MultmatrixBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (mbp *MultmatrixBlockPrimitive) String() string {
	return blockPrimitiveString(mbp.TokenLiteral(), mbp.Arguments, mbp.Body)
}

// TokenLiteral returns the token literal.
func (mbp *MultmatrixBlockPrimitive) TokenLiteral() string { return mbp.Token.Literal }

// ProjectionBlockPrimitive represents a CSG block primitive.
type ProjectionBlockPrimitive struct {
	Token     token.Token
	Arguments []Expression
	Body      *BlockStatement
}

func (pbp *ProjectionBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (pbp *ProjectionBlockPrimitive) String() string {
	return blockPrimitiveString(pbp.TokenLiteral(), pbp.Arguments, pbp.Body)
}

// TokenLiteral returns the token literal.
func (pbp *ProjectionBlockPrimitive) TokenLiteral() string { return pbp.Token.Literal }

// RotateExtrudeBlockPrimitive represents a CSG block primitive.
type RotateExtrudeBlockPrimitive struct {
	Token     token.Token
	Arguments []Expression
	Body      *BlockStatement
}

func (rebp *RotateExtrudeBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (rebp *RotateExtrudeBlockPrimitive) String() string {
	return blockPrimitiveString(rebp.TokenLiteral(), rebp.Arguments, rebp.Body)
}

// TokenLiteral returns the token literal.
func (rebp *RotateExtrudeBlockPrimitive) TokenLiteral() string { return rebp.Token.Literal }

// UnionBlockPrimitive represents a CSG block primitive.
type UnionBlockPrimitive struct {
	Token token.Token
	Body  *BlockStatement
}

func (ubp *UnionBlockPrimitive) expressionNode() {}

// String returns the string representation of the Node.
func (ubp *UnionBlockPrimitive) String() string {
	return blockPrimitiveString(ubp.TokenLiteral(), nil, ubp.Body)
}

// TokenLiteral returns the token literal.
func (ubp *UnionBlockPrimitive) TokenLiteral() string { return ubp.Token.Literal }

package parser

import (
	"github.com/gmlewis/go-csg/ast"
	"github.com/gmlewis/go-csg/token"
)

func (p *Parser) parseNamedArgument(left ast.Expression) ast.Expression {
	exp := &ast.NamedArgument{
		Token: p.curToken,
		Name:  left,
	}

	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)

	return exp
}

func (p *Parser) parsePrimitiveArguments() ([]ast.Expression, bool) {
	if !p.expectPeek(token.LPAREN) {
		return nil, false
	}

	args := p.parseExpressionList(token.RPAREN)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return args, true
}

func (p *Parser) parseCirclePrimitive() ast.Expression {
	prim := &ast.CirclePrimitive{Token: p.curToken}

	args, ok := p.parsePrimitiveArguments()
	if !ok {
		return nil
	}

	prim.Arguments = args

	return prim
}

func (p *Parser) parseCubePrimitive() ast.Expression {
	prim := &ast.CubePrimitive{Token: p.curToken}

	args, ok := p.parsePrimitiveArguments()
	if !ok {
		return nil
	}

	prim.Arguments = args

	return prim
}

func (p *Parser) parseCylinderPrimitive() ast.Expression {
	prim := &ast.CylinderPrimitive{Token: p.curToken}

	args, ok := p.parsePrimitiveArguments()
	if !ok {
		return nil
	}

	prim.Arguments = args

	return prim
}

// func (p *Parser) parseGroupPrimitive() ast.Expression {
// 	prim := &ast.GroupPrimitive{Token: p.curToken}

// 	args, ok := p.parsePrimitiveArguments()
// 	if !ok {
// 		return nil
// 	}

// 	prim.Arguments = args

// 	return prim
// }

func (p *Parser) parsePolygonPrimitive() ast.Expression {
	prim := &ast.PolygonPrimitive{Token: p.curToken}

	args, ok := p.parsePrimitiveArguments()
	if !ok {
		return nil
	}

	prim.Arguments = args

	return prim
}

func (p *Parser) parseSpherePrimitive() ast.Expression {
	prim := &ast.SpherePrimitive{Token: p.curToken}

	args, ok := p.parsePrimitiveArguments()
	if !ok {
		return nil
	}

	prim.Arguments = args

	return prim
}

func (p *Parser) parseSquarePrimitive() ast.Expression {
	prim := &ast.SquarePrimitive{Token: p.curToken}

	args, ok := p.parsePrimitiveArguments()
	if !ok {
		return nil
	}

	prim.Arguments = args

	return prim
}

// func (p *Parser) parseTextPrimitive() ast.Expression {
// 	prim := &ast.TextPrimitive{Token: p.curToken}

// 	args, ok := p.parsePrimitiveArguments()
// 	if !ok {
// 		return nil
// 	}

// 	prim.Arguments = args

// 	return prim
// }

func (p *Parser) parseBlockPrimitive() (args []ast.Expression, block *ast.BlockStatement, ok bool) {
	if !p.expectPeek(token.LPAREN) {
		return nil, nil, false
	}

	args = p.parseExpressionList(token.RPAREN)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	} else {
		if !p.expectPeek(token.LBRACE) {
			return nil, nil, false
		}

		block = p.parseBlockStatement()
	}

	return args, block, true
}

func (p *Parser) parseColorBlockPrimitive() ast.Expression {
	blockPrim := &ast.ColorBlockPrimitive{Token: p.curToken}

	args, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Arguments = args
	blockPrim.Body = body

	return blockPrim
}

func (p *Parser) parseDifferenceBlockPrimitive() ast.Expression {
	blockPrim := &ast.DifferenceBlockPrimitive{Token: p.curToken}

	_, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Body = body

	return blockPrim
}

func (p *Parser) parseGroupBlockPrimitive() ast.Expression {
	blockPrim := &ast.GroupBlockPrimitive{Token: p.curToken}

	_, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Body = body

	return blockPrim
}

func (p *Parser) parseHullBlockPrimitive() ast.Expression {
	blockPrim := &ast.HullBlockPrimitive{Token: p.curToken}

	_, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Body = body

	return blockPrim
}

func (p *Parser) parseIntersectionBlockPrimitive() ast.Expression {
	blockPrim := &ast.IntersectionBlockPrimitive{Token: p.curToken}

	_, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Body = body

	return blockPrim
}

func (p *Parser) parseLinearExtrudeBlockPrimitive() ast.Expression {
	blockPrim := &ast.LinearExtrudeBlockPrimitive{Token: p.curToken}

	args, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Arguments = args
	blockPrim.Body = body

	return blockPrim
}

func (p *Parser) parseMinkowskiBlockPrimitive() ast.Expression {
	blockPrim := &ast.MinkowskiBlockPrimitive{Token: p.curToken}

	args, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Arguments = args
	blockPrim.Body = body

	return blockPrim
}

func (p *Parser) parseMultmatrixBlockPrimitive() ast.Expression {
	blockPrim := &ast.MultmatrixBlockPrimitive{Token: p.curToken}

	args, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Arguments = args
	blockPrim.Body = body

	return blockPrim
}

func (p *Parser) parseProjectionBlockPrimitive() ast.Expression {
	blockPrim := &ast.ProjectionBlockPrimitive{Token: p.curToken}

	args, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Arguments = args
	blockPrim.Body = body

	return blockPrim
}

func (p *Parser) parseRotateExtrudeBlockPrimitive() ast.Expression {
	blockPrim := &ast.RotateExtrudeBlockPrimitive{Token: p.curToken}

	args, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Arguments = args
	blockPrim.Body = body

	return blockPrim
}

func (p *Parser) parseUnionBlockPrimitive() ast.Expression {
	blockPrim := &ast.UnionBlockPrimitive{Token: p.curToken}

	_, body, ok := p.parseBlockPrimitive()
	if !ok {
		return nil
	}

	blockPrim.Body = body

	return blockPrim
}

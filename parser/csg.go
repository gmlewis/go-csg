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

func (p *Parser) parseGroupPrimitive() ast.Expression {
	prim := &ast.GroupPrimitive{Token: p.curToken}

	args, ok := p.parsePrimitiveArguments()
	if !ok {
		return nil
	}

	prim.Arguments = args

	return prim
}

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

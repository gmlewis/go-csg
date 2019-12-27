package parser

import (
	"fmt"
	"testing"

	"github.com/gmlewis/go-monkey/ast"
	"github.com/gmlewis/go-monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	le := lexer.New(input)
	p := New(le)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements len = %v, want 3", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			stmt := program.Statements[i]
			if !testLetStatements(t, stmt, tt.expectedIdentifier) {
				return
			}
		})
	}
}

func testLetStatements(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral = %q, want 'let'", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s = %T, want *ast.LetStatement", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value = %v, want %v", letStmt.Name.Value, name)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral = %v, want %v", letStmt.Name.TokenLiteral(), name)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %v errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	le := lexer.New(input)
	p := New(le)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements len = %v, want 3", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"5"},
		{"10"},
		{"993322"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			stmt := program.Statements[i]
			if !testReturnStatements(t, stmt, tt.expectedIdentifier) {
				return
			}
		})
	}
}

func testReturnStatements(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "return" {
		t.Errorf("s.TokenLiteral = %q, want 'return'", s.TokenLiteral())
		return false
	}

	returnStmt, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("s = %T, want *ast.returnStatement", s)
		return false
	}

	// if returnStmt.ReturnValue.TokenLiteral() != name {
	// 	t.Errorf("returnStmt.ReturnValue = %v, want %v", returnStmt.ReturnValue, name)
	// 	return false
	// }

	if returnStmt.TokenLiteral() != "return" {
		t.Errorf("returnStmt.Name.TokenLiteral = %v, want 'return'", returnStmt.TokenLiteral())
		return false
	}

	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	le := lexer.New(input)
	p := New(le)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program len(statements) = %v, want 1", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] = %T, want *ast.ExpressionsStatment", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp = %T, want *ast.Identifier", stmt.Expression)
	}
	if want := "foobar"; ident.Value != want {
		t.Errorf("ident.Value = %v, want %v", ident.Value, want)
	}
	if got, want := ident.TokenLiteral(), "foobar"; got != want {
		t.Errorf("ident.TokenLiteral = %v, want %v", got, want)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	le := lexer.New(input)
	p := New(le)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program len(statements) = %v, want 1", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] = %T, want *ast.ExpressionsStatment", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp = %T, want *ast.IntegerLiteral", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value = %v, want %v", literal.Value, 5)
	}
	if got, want := literal.TokenLiteral(), "5"; got != want {
		t.Errorf("literal.TokenLiteral = %v, want %v", got, want)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		val      int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			le := lexer.New(tt.input)
			p := New(le)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program len(statements) = %v, want 1", len(program.Statements))
			}
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] = %T, want *ast.ExpressionsStatment", program.Statements[0])
			}
			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("exp = %T, want *ast.PrefixExpression", stmt.Expression)
			}
			if exp.Operator != tt.operator {
				t.Errorf("exp.Operator = %v, want %v", exp.Operator, tt.operator)
			}
			if !testIntegerLiteral(t, exp.Right, tt.val) {
				return
			}
		})
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il = %T, want *ast.IntegerLiteral", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value = %v, want %v", integ.Value, value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("integ.TokenLiteral = %v, want %v", integ.TokenLiteral(), value)
		return false
	}

	return true
}

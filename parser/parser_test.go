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

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		leftVal  int64
		operator string
		rightVal int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
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
			exp, ok := stmt.Expression.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("exp = %T, want *ast.InfixExpression", stmt.Expression)
			}
			if !testIntegerLiteral(t, exp.Left, tt.leftVal) {
				return
			}
			if exp.Operator != tt.operator {
				t.Errorf("exp.Operator = %v, want %v", exp.Operator, tt.operator)
			}
			if !testIntegerLiteral(t, exp.Right, tt.rightVal) {
				return
			}
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			le := lexer.New(tt.input)
			p := New(le)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			got := program.String()
			if got != tt.want {
				t.Errorf("string = %v, want %v", got, tt.want)
			}
		})
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp = %T, want *ast.Identifier", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value = %v, want %v", ident.Value, value)
		return false
	}

	if ident.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("ident.TokenLiteral = %v, want %v", ident.TokenLiteral(), value)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp (%T) not handled", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp = %T, want *ast.InfixExpression", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("operator = %v, want %v", opExp.Operator, operator)
		return false
	}

	return testLiteralExpression(t, opExp.Right, right)
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

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

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp = %T, want *ast.IfExpression", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence len(statements) = %v, want 1", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence Statemetns[0] = %T, want *ast.ExpressionStatement", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative = %v, want nil", exp.Alternative)
	}
}

func TestFunctionLiteralExpression(t *testing.T) {
	input := `fn(x, y) { x + y; }`

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

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("exp = %T, want *ast.FunctionLiteral", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("len(parameters) = %v, want 2", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("len(Body.Statements) = %v, want 1", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("bodyStmt = %T, want *ast.ExpressionStatement", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"fn() {};", nil},
		{"fn(x) {}; ", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			le := lexer.New(tt.input)
			p := New(le)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			function := stmt.Expression.(*ast.FunctionLiteral)

			if len(function.Parameters) != len(tt.want) {
				t.Errorf("len(params) = %v, want %v", len(function.Parameters), len(tt.want))
			}

			for j, ident := range tt.want {
				testLiteralExpression(t, function.Parameters[j], ident)
			}
		})
	}
}

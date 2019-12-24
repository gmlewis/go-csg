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

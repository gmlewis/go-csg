package lexer

import (
	"fmt"
	"testing"

	"github.com/gmlewis/go-openscad/token"
)

func TestNextToken(t *testing.T) {
	input := `sphere(20);
$fa=0.5; // default minimum facet angle is now 0.5
cylinder(r=3,h=10);
cube([20,10,5]);

union()
{
  cube([20,20,20], center=true);
  sphere(14);
}

let five = 5;
let ten = 10;

let add = function(x, y) {
	x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
`

	tests := []struct {
		expectedType    token.T
		expectedLiteral string
	}{
		{token.SPHERE, "sphere"},
		{token.LPAREN, "("},
		{token.INT, "20"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.DOLLAR, "$"},
		{token.IDENT, "fa"},
		{token.ASSIGN, "="},
		{token.FLOAT, "0.5"},
		{token.SEMICOLON, ";"},
		{token.LINECOMMENT, " default minimum facet angle is now 0.5"},

		{token.CYLINDER, "cylinder"},
		{token.LPAREN, "("},
		{token.IDENT, "r"},
		{token.ASSIGN, "="},
		{token.INT, "3"},
		{token.COMMA, ","},
		{token.IDENT, "h"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.CUBE, "cube"},
		{token.LPAREN, "("},
		{token.LBRACKET, "["},
		{token.INT, "20"},
		{token.COMMA, ","},
		{token.INT, "10"},
		{token.COMMA, ","},
		{token.INT, "5"},
		{token.RBRACKET, "]"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.UNION, "union"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.CUBE, "cube"},
		{token.LPAREN, "("},
		{token.LBRACKET, "["},
		{token.INT, "20"},
		{token.COMMA, ","},
		{token.INT, "20"},
		{token.COMMA, ","},
		{token.INT, "20"},
		{token.RBRACKET, "]"},
		{token.COMMA, ","},
		{token.IDENT, "center"},
		{token.ASSIGN, "="},
		{token.TRUE, "true"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.SPHERE, "sphere"},
		{token.LPAREN, "("},
		{token.INT, "14"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "function"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOTEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	le := New(input)

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			tok := le.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tokentype = %q, want %q", tok.Type, tt.expectedType)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("literal = %q, want %q", tok.Literal, tt.expectedLiteral)
			}
		})
	}
}

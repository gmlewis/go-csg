package evaluator

import (
	"fmt"
	"testing"

	"github.com/gmlewis/go-monkey/lexer"
	"github.com/gmlewis/go-monkey/object"
	"github.com/gmlewis/go-monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)
			testIntegerObject(t, got, tt.want)
		})
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)
			testBooleanObject(t, got, tt.want)
		})
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)
			testBooleanObject(t, got, tt.want)
		})
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)
			if integer, ok := tt.want.(int); ok {
				testIntegerObject(t, got, int64(integer))
			} else {
				testNullObject(t, got)
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if ( 10 > 1 ) { if ( 10 > 1 ) { return 10 ; } return 1 ; }`, 10},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)
			testIntegerObject(t, got, tt.want)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{` if (10 > 1) { if (10 > 1) { return true + false; } return 1; } `, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)

			errObj, ok := got.(*object.Error)
			if !ok {
				t.Fatalf("got = %T (%+v), want *object.Error", got, got)
			}

			if errObj.Message != tt.want {
				t.Errorf("error = %v, want %v", errObj.Message, tt.want)
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)
			testIntegerObject(t, got, tt.want)
		})
	}
}

func testEval(input string) object.Object {
	le := lexer.New(input)
	p := parser.New(le)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, want int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object = %T, want *object.Integer", obj)
		return false
	}

	if result.Value != want {
		t.Errorf("value = %v, want %v", result.Value, want)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, want bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object = %T, want *object.Boolean", obj)
		return false
	}

	if result.Value != want {
		t.Errorf("value = %v, want %v", result.Value, want)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != Null {
		t.Errorf("got = %T (%+v), want Null", obj, obj)
		return false
	}

	return true
}

package evaluator

import (
	"fmt"
	"testing"

	"github.com/gmlewis/go-openscad/lexer"
	"github.com/gmlewis/go-openscad/object"
	"github.com/gmlewis/go-openscad/parser"
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
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
		{`{"name": "Monkey"}[function(x) { x }];`, "unusable as hash key: FUNCTION"},
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

func TestFunctionObject(t *testing.T) {
	input := "function(x) { x + 2; };"

	got := testEval(input)
	fn, ok := got.(*object.Function)
	if !ok {
		t.Fatalf("got = %T (%+v), want *object.Function", got, got)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("len(params) = %v, want 1", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("params[0] = %v, want x", fn.Parameters[0].String())
	}

	if want := "(x + 2)"; fn.Body.String() != want {
		t.Fatalf("body = %v, want %v", fn.Body.String(), want)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	got := testEval(input)
	str, ok := got.(*object.String)
	if !ok {
		t.Fatalf("got = %T (%+v), want *object.String", got, got)
	}

	if want := "Hello World!"; str.Value != want {
		t.Fatalf("value = %v, want %v", str.Value, want)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	got := testEval(input)
	str, ok := got.(*object.String)
	if !ok {
		t.Fatalf("got = %T (%+v), want *object.String", got, got)
	}

	if want := "Hello World!"; str.Value != want {
		t.Fatalf("value = %v, want %v", str.Value, want)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"let identity = function(x) { x; }; identity(5);", 5},
		{"let identity = function(x) { return x; }; identity(5);", 5},
		{"let double = function(x) { x * 2; }; double(5);", 10},
		{"let add = function(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = function(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"function(x) { x; }(5)", 5},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)
			testIntegerObject(t, got, tt.want)
		})
	}
}

func TestClosures(t *testing.T) {
	input := ` let newAdder = function(x) { function(y) { x + y }; }; let addTwo = newAdder(2); addTwo(2);`

	got := testEval(input)
	testIntegerObject(t, got, 4)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)
			switch want := tt.want.(type) {
			case int:
				testIntegerObject(t, got, int64(want))
			case string:
				errObj, ok := got.(*object.Error)
				if !ok {
					t.Fatalf("got = %T (%+v), want *object.Error", got, got)
				}
				if errObj.Message != want {
					t.Errorf("got = %v, want %v", errObj.Message, want)
				}
			}
		})
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	got := testEval(input)
	result, ok := got.(*object.Array)
	if !ok {
		t.Fatalf("got = %T (%+v), want *object.Array", got, got)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("len(result.Elements) = %v, want 3", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)
			switch want := tt.want.(type) {
			case int:
				testIntegerObject(t, got, int64(want))
			default:
				testNullObject(t, got)
			}
		})
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two"; { "one": 10 - 9, two: 1 + 1, "thr" + "ee": 6 / 2, 4: 4, true: 5, false: 6 }`

	got := testEval(input)
	result, ok := got.(*object.Hash)
	if !ok {
		t.Fatalf("got = %T (%+v), want *object.Hash", got, got)
	}

	want := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		True.HashKey():                             5,
		False.HashKey():                            6,
	}

	if len(result.Pairs) != len(want) {
		t.Fatalf("len(result.Pairs) = %v, want %v", len(result.Pairs), len(want))
	}

	for wantKey, wantValue := range want {
		pair, ok := result.Pairs[wantKey]
		if !ok {
			t.Fatalf("no pair for key %v", wantKey)
		}

		testIntegerObject(t, pair.Value, wantValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := testEval(tt.input)
			switch want := tt.want.(type) {
			case int:
				testIntegerObject(t, got, int64(want))
			default:
				testNullObject(t, got)
			}
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
		t.Errorf("object = %T (%+v), want *object.Integer", obj, obj)
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

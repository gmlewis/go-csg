// Package evaluator implements the language evaluator.
package evaluator

import (
	"github.com/gmlewis/go-monkey/ast"
	"github.com/gmlewis/go-monkey/object"
)

// Singleton object types.
var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

// Eval evaluates the AST node and returns the evaluated object.
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.BooleanLiteral:
		if node.Value {
			return True
		}
		return False

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	}

	return nil
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	}
	return False
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerT {
		return Null
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerT && right.Type() == object.IntegerT:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		if left == right {
			return True
		}
		return False
	case operator == "!=":
		if left != right {
			return True
		}
		return False
	}
	return Null
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}

	case "<":
		if leftVal < rightVal {
			return True
		}
		return False
	case ">":
		if leftVal > rightVal {
			return True
		}
		return False
	case "==":
		if leftVal == rightVal {
			return True
		}
		return False
	case "!=":
		if leftVal != rightVal {
			return True
		}
		return False
	}
	return Null
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	}
	return Null
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

// Package evaluator implements the language evaluator.
package evaluator

import (
	"fmt"

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
		return evalProgram(node)

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

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
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
		return newError("unknown operator: -%v", right.Type())
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
	case left.Type() != right.Type():
		return newError("type mismatch: %v %v %v", left.Type(), operator, right.Type())
	}
	return newError("unknown operator: %v %v %v", left.Type(), operator, right.Type())
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

	return newError("unknown operator: %v %v %v", left.Type(), operator, right.Type())
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	}
	return newError("unknown operator: %v%v", operator, right.Type())
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			if rt := result.Type(); rt == object.ReturnValueT || rt == object.ErrorT {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	switch {
	case isTruthy(condition):
		return Eval(ie.Consequence)
	case ie.Alternative != nil:
		return Eval(ie.Alternative)
	}

	return Null
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case Null:
		return false
	case True:
		return true
	case False:
		return false
	}

	return true
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

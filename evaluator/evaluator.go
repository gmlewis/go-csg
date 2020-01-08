// Package evaluator implements the language evaluator.
package evaluator

import (
	"fmt"
	"log"

	"github.com/gmlewis/go-csg/ast"
	"github.com/gmlewis/go-csg/object"
)

// Singleton object types.
var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

// Eval evaluates the AST node and returns the evaluated object.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// Expressions
	case *ast.BooleanLiteral:
		if node.Value {
			return True
		}
		return False

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	// CSG...

	case *ast.CubePrimitive:
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return &object.CubePrimitive{Arguments: args}

	case *ast.CylinderPrimitive:
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return &object.CylinderPrimitive{Arguments: args}

	case *ast.GroupBlockPrimitive:
		body := Eval(node.Body, env)
		return &object.GroupBlockPrimitive{Body: body}

	case *ast.NamedArgument:
		value := Eval(node.Value, env)
		return &object.NamedArgument{Name: node.Name.String(), Value: value}

	case *ast.PolygonPrimitive:
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return &object.PolygonPrimitive{Arguments: args}

	case *ast.SpherePrimitive:
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return &object.SpherePrimitive{Arguments: args}

	case *ast.SquarePrimitive:
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return &object.SquarePrimitive{Arguments: args}

	case *ast.UndefLiteral:
		return Null

	default:
		log.Fatalf("unhandled AST Node type %T (%+v)", node, node)
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
	switch {
	case right.Type() == object.IntegerT:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	case right.Type() == object.FloatT:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	default:
		return newError("unknown operator: -%v", right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerT && right.Type() == object.IntegerT:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.StringT && right.Type() == object.StringT:
		return evalStringInfixExpression(operator, left, right)
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

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
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

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			if rt := result.Type(); rt == object.ReturnValueT || rt == object.ErrorT {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	switch {
	case isError(condition):
		return condition
	case isTruthy(condition):
		return Eval(ie.Consequence, env)
	case ie.Alternative != nil:
		return Eval(ie.Alternative, env)
	}

	return Null
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayT && index.Type() == object.IntegerT:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HashT:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %v", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return Null
	}

	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %v", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return Null
	}

	return pair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := map[object.HashKey]object.HashPair{}

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %v", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %v", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if v, ok := obj.(*object.ReturnValue); ok {
		return v.Value
	}
	return obj
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

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorT
	}
	return false
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

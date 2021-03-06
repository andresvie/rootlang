package evaluator

import (
	"rootlang/ast"
	"rootlang/object"
	"rootlang/builtin"
	"fmt"
	"strings"
)

func CallMainFunction(function *object.Function, builtinSymbols *builtin.Builtin) object.Object {
	return applyArgumentsToFunctionAndCall(function, []object.Object{}, builtinSymbols)
}
func Eval(node ast.Node, environment *object.Environment, builtinSymbols *builtin.Builtin) object.Object {

	switch nodeType := node.(type) {
	case *ast.Program:
		return evalProgram(nodeType, environment, builtinSymbols)
	case *ast.BlockStatement:
		return evalStatement(nodeType.Statements, environment, builtinSymbols)
	case *ast.ExpressionStatement:
		return Eval(nodeType.Exp, environment, builtinSymbols)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: nodeType.Value}
	case *ast.BoolExpression:
		return nativeToBooleanObject(nodeType.Value == "true")
	case *ast.StringExpression:
		return nativeStringToObject(nodeType.Value)
	case *ast.ImportStatement:
		builtinModule, ok := builtinSymbols.GetObject(nodeType.Name.Value)
		if ok {
			environment.SetVar(nodeType.Name.Value, builtinModule)
			return nil
		}
		resultImport := importModule(nodeType, builtinSymbols)
		if isError(resultImport) {
			return resultImport
		}
		module := resultImport.(*object.Module)
		environment.SetVar(module.Name, module)
		return nil
	case *ast.LetStatement:
		valueExpression := Eval(nodeType.Value, environment, builtinSymbols)
		if isError(valueExpression) {
			return valueExpression
		}
		environment.SetVar(nodeType.Name.Value, valueExpression)
		return nil
	case *ast.Identifier:
		value, ok := environment.GetVar(nodeType.Value)
		if !ok {
			value, ok = builtinSymbols.GetObject(nodeType.Value)
			if !ok {
				return &object.ErrorObject{Error: fmt.Sprintf("%s was not declare", nodeType.Value)}
			}
		}
		return value
	case *ast.FunctionExpression:
		return &object.Function{Params: nodeType.Params, Body: nodeType.Block, Env: environment.ExtendNewEnvironment()}
	case *ast.CallFunctionExpression:
		params := evalExpressions(nodeType.Arguments, environment, builtinSymbols)
		if len(params) == 1 && isError(params[0]) {
			return params[0]
		}
		value := Eval(nodeType.Function, environment, builtinSymbols)
		if isError(value) {
			return value
		}
		switch valueFunction := value.(type) {
		case *object.Function:
			return applyArgumentsToFunctionAndCall(valueFunction, params, builtinSymbols)
		case *builtin.BuiltinFunction:
			return valueFunction.Function(environment, builtinSymbols, Eval, params...)
		default:
			return newError(fmt.Sprintf("expected function %s", value.Inspect()))
		}
		function, ok := value.(*object.Function)
		if !ok {
			return newError(fmt.Sprintf("expected function %s", value.Inspect()))
		}
		return applyArgumentsToFunctionAndCall(function, params, builtinSymbols)
	case *ast.PrefixExpression:
		rightExpression := Eval(nodeType.RightExpression, environment, builtinSymbols)
		if isError(rightExpression) {
			return rightExpression
		}
		return evalPrefixExpression(nodeType.Operator, rightExpression)
	case *ast.ReturnStatement:
		value := Eval(nodeType.Value, environment, builtinSymbols)
		if isError(value) {
			return value
		}
		return &object.ReturnObject{Value: value}
	case *ast.IfExpression:
		condition := Eval(nodeType.Condition, environment, builtinSymbols)
		if isError(condition) {
			return condition
		}
		if evalTruthValue(condition) {
			return Eval(nodeType.ConditionalBlock, environment, builtinSymbols)
		} else if nodeType.AlternativeBlock != nil {
			return Eval(nodeType.AlternativeBlock, environment, builtinSymbols)
		} else {
			return nil
		}
	case *ast.InfixExpression:
		leftExpression := Eval(nodeType.LeftExpression, environment, builtinSymbols)
		if isError(leftExpression) {
			return leftExpression
		}
		if nodeType.Operator == "::" && leftExpression.Type() != object.MODULE_OBJ {
			return newError("module was expected")
		}

		if nodeType.Operator == "::" && leftExpression.Type() == object.MODULE_OBJ {
			return moduleEvaluation(leftExpression.(*object.Module), nodeType.RightExpression, environment, builtinSymbols)
		}
		rightExpression := Eval(nodeType.RightExpression, environment, builtinSymbols)

		if isError(rightExpression) {
			return rightExpression
		}
		return evalInfixExpression(nodeType.Operator, rightExpression, leftExpression)
	}

	return nil
}

func moduleEvaluation(module *object.Module, expression ast.Expression, environment *object.Environment, builtinSymbols *builtin.Builtin) object.Object {
	switch nodeType := expression.(type) {
	case *ast.CallFunctionExpression:
		params := evalExpressions(nodeType.Arguments, environment, builtinSymbols)
		if len(params) == 1 && isError(params[0]) {
			return params[0]
		}
		value := Eval(nodeType.Function, module.Env, builtinSymbols)
		if isError(value) {
			return value
		}
		if value.Type() != object.FUNCTION_OBJ && value.Type() != object.BUILTIN_FUNCTION_OBJ {
			return newError(fmt.Sprintf("expected function and got %s", value.Type()))
		}
		switch functionType := value.(type) {
		case *object.Function:
			return applyArgumentsToFunctionAndCall(functionType, params, builtinSymbols)
		case *builtin.BuiltinFunction:
			return functionType.Function(environment, builtinSymbols, Eval, params...)
		default:
			return newError("expression not expected on module")
		}

	case *ast.Identifier:
		value, ok := module.Env.GetVar(nodeType.Value)
		if !ok {
			return newError(fmt.Sprintf("symbol %s not found in module %s", nodeType.Value, module.Name))
		}
		return value
	default:
		return newError("expression not expected on module")
	}

}

func applyArgumentsToFunctionAndCall(function *object.Function, params []object.Object, builtinSymbols *builtin.Builtin) object.Object {

	if len(params) > len(function.Params) {
		return newError(fmt.Sprintf("this function takes at least %d arguments (%d given)", len(function.Params), len(params)))
	}
	newEnvironment := applyArguments(function, params)
	if len(function.Params) == len(params) {
		returnValue := Eval(function.Body, newEnvironment, builtinSymbols)
		if returnValue.Type() == object.RETURN_OBJ {
			return returnValue.(*object.ReturnObject).Value
		}
		return returnValue
	}
	return function.Clone(function.Params[len(params):], newEnvironment)
}

func applyArguments(function *object.Function, params []object.Object) *object.Environment {
	newEnvironment := function.Env.ExtendNewEnvironment()
	for i := 0; i < len(params); i++ {
		newEnvironment.SetVar(function.Params[i].Value, params[i])
	}
	return newEnvironment
}

func evalExpressions(expressions []ast.Expression, environment *object.Environment, builtinSymbols *builtin.Builtin) []object.Object {
	expressionsObjects := make([]object.Object, 0)
	for _, expression := range expressions {
		expressionObject := Eval(expression, environment, builtinSymbols)
		if isError(expressionObject) {
			return []object.Object{expressionObject}
		}
		expressionsObjects = append(expressionsObjects, expressionObject)
	}
	return expressionsObjects
}

func evalProgram(program *ast.Program, environment *object.Environment, builtinSymbols *builtin.Builtin) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, environment, builtinSymbols)
		if result == nil {
			continue
		}
		if returnValue, ok := result.(*object.ReturnObject); ok {
			return returnValue.Value
		}
		if error, ok := result.(*object.ErrorObject); ok {
			return error
		}
	}
	return result
}

func evalTruthValue(value object.Object) bool {
	switch valueType := value.(type) {
	case *object.Boolean:
		return valueType.Value
	case *object.Integer:
		return valueType.Value != 0
	default:
		return false
	}
}

func evalPrefixExpression(operator string, rightValue object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(rightValue)
	case "-":
		return evalMinusOperator(rightValue)
	default:
		return newError(fmt.Sprintf("unknow operator for %s%s", operator, rightValue.Inspect()))
	}
}

func newError(error string) *object.ErrorObject {
	return &object.ErrorObject{Error: error}
}

func evalInfixExpression(operator string, rightValue, leftValue object.Object) object.Object {

	if rightValue.Type() == object.INTEGER_OBJ && leftValue.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, rightValue, leftValue);
	}
	if (rightValue.Type() == object.STRING_OBJ || leftValue.Type() == object.STRING_OBJ) && operator == "+" {
		return nativeStringToObject(fmt.Sprintf("%s%s", leftValue.Inspect(), rightValue.Inspect()));
	}
	if (rightValue.Type() == object.STRING_OBJ || leftValue.Type() == object.STRING_OBJ) && operator == "!=" {
		return nativeToBooleanObject(strings.Compare(leftValue.Inspect(), rightValue.Inspect()) != 0);
	}
	switch operator {
	case "==":
		return nativeToBooleanObject(leftValue == rightValue)
	case "!=":

		return nativeToBooleanObject(leftValue != rightValue)
	default:
		return &object.ErrorObject{Error: fmt.Sprintf("unknow operator for %s %s %s", leftValue.Inspect(), operator, rightValue.Inspect())}
	}

}

func nativeStringToObject(text string) object.Object {
	return &object.String{Value: text}
}

func evalIntegerInfixExpression(operator string, rightValue, leftValue object.Object) object.Object {
	rightIntegerValue := rightValue.(*object.Integer)
	leftIntegerValue := leftValue.(*object.Integer)
	switch operator {
	case "+":
		return &object.Integer{Value: leftIntegerValue.Value + rightIntegerValue.Value }
	case "-":
		return &object.Integer{Value: leftIntegerValue.Value - rightIntegerValue.Value }
	case "/":
		return &object.Integer{Value: leftIntegerValue.Value / rightIntegerValue.Value }
	case "*":
		return &object.Integer{Value: leftIntegerValue.Value * rightIntegerValue.Value }
	case "%":
		return &object.Integer{Value: leftIntegerValue.Value % rightIntegerValue.Value }
	case "==":
		return nativeToBooleanObject(leftIntegerValue.Value == rightIntegerValue.Value)
	case "!=":
		return nativeToBooleanObject(leftIntegerValue.Value != rightIntegerValue.Value)
	case ">":
		return nativeToBooleanObject(leftIntegerValue.Value > rightIntegerValue.Value)
	case "<":
		return nativeToBooleanObject(leftIntegerValue.Value < rightIntegerValue.Value)
	default:
		return object.NULL
	}
}
func isError(error object.Object) bool {
	return error != nil && error.Type() == object.ERROR_OBJ
}

func evalMinusOperator(rightValue object.Object) object.Object {
	if rightValue.Type() != object.INTEGER_OBJ {
		return &object.ErrorObject{Error: fmt.Sprintf("unknow operator for -%s", rightValue.Inspect())}
	}
	value := rightValue.(*object.Integer)
	return &object.Integer{Value: -value.Value}
}

func evalBangOperator(rightValue object.Object) object.Object {
	integerObject, ok := rightValue.(*object.Integer)
	if ok && integerObject.Value == 0 {
		return object.TRUE
	}
	switch rightValue {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE

	default:
		return object.FALSE
	}
}

func evalStatement(statements []ast.Statement, environment *object.Environment, builtinSymbols *builtin.Builtin) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement, environment, builtinSymbols)
		if result != nil && (result.Type() == object.RETURN_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}

	}
	return result
}

func nativeToBooleanObject(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

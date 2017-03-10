package builtin

import (
	"rootlang/object"
	"rootlang/ast"
	"fmt"
)

func __callFunction(function *object.Function, env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, param object.Object) [][]object.Object {
	elements := make([][]object.Object, 0)
	switch paramTye := param.(type) {
	case *object.List:
		for _, paramElement := range paramTye.Elements {
			returnValues := applyArgumentsToFunctionAndCall(function, []object.Object{paramElement}, b, eval)
			if returnValues.Type() == object.ERROR_OBJ {
				return [][]object.Object{[]object.Object{returnValues}};
			}
			elements = append(elements, []object.Object{returnValues, paramElement})
		}
	default:
		returnValues := applyArgumentsToFunctionAndCall(function, []object.Object{paramTye}, b, eval)
		elements = append(elements, []object.Object{returnValues, paramTye})
	}
	return elements
}

func applyArgumentsToFunctionAndCall(function *object.Function, params []object.Object, builtinSymbols *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object) object.Object {

	if len(params) > len(function.Params) {
		return &object.ErrorObject{Error: (fmt.Sprintf("this function takes at least %d arguments (%d given)", len(function.Params), len(params)))}
	}
	newEnvironment := applyArguments(function, params)
	if len(function.Params) == len(params) {
		returnValue := eval(function.Body, newEnvironment, builtinSymbols)
		if returnValue != nil && returnValue.Type() == object.RETURN_OBJ {
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

func isErrorObject(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}

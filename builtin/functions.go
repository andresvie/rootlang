package builtin

import (
	"rootlang/object"
	"fmt"
	"rootlang/ast"
	"math"
	"bytes"
)

func _zip(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	elements := make([]object.Object, 0)
	if len(params) < 1 {
		return &object.ErrorObject{Error: fmt.Sprintf("zip expect more than 1 params and got %d", len(params))}
	}
	if !_allParamsAreList(params) {
		return &object.ErrorObject{Error: fmt.Sprintf("zip all arguments expected to be list", len(params))}
	}
	minIndex := _getListMinIndex(params)
	for i := uint64(0); i < minIndex; i++ {
		zipArray := make([]object.Object, 0)
		for j := 0; j < len(params); j++ {
			list := params[j].(*object.List)
			zipArray = append(zipArray, list.Elements[i])
		}
		elements = append(elements, &object.List{Elements: zipArray})
	}
	return &object.List{Elements: elements}
}

func _print(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	text := bytes.NewBufferString("")
	for _, param := range params {
		text.WriteString(param.Inspect())
	}
	fmt.Println(text)
	return &object.String{Value: text.String()}
}

func _getListMinIndex(params []object.Object) uint64 {
	minIndex := uint64(math.MaxUint64)
	for i := 0; i < len(params); i++ {
		list := params[0].(*object.List)
		minIndex = _min(minIndex, uint64(len(list.Elements)))
	}
	return minIndex
}

func _min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

func _allParamsAreList(params []object.Object) bool {
	for _, element := range params {
		if element.Type() != object.LIST_OBJ {
			return false
		}
	}
	return true
}

func _map(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	elements := make([]object.Object, 0)
	if len(params) < 1 {
		return &object.ErrorObject{Error: fmt.Sprintf("map expect more than 1 params and got %d", len(params))}
	}
	function, ok := params[0].(*object.Function)
	if !ok {
		return &object.ErrorObject{Error: fmt.Sprintf("map first params should be function and got %s", params[0].Type())}
	}
	for _, objectParam := range params[1:] {
		returnValues := __callFunction(function, env, b, eval, objectParam)
		if len(returnValues) == 1 && len(returnValues[0]) == 1 && returnValues[0][0].Type() == object.ERROR_OBJ {
			return returnValues[0][0]
		}
		elements = append(elements, _get_only_values(returnValues)...)
	}
	return &object.List{Elements: elements}
}

func _reduce(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	if len(params) < 2 {
		return &object.ErrorObject{Error: fmt.Sprintf("reduce expect more than 1 params and got %d", len(params))}
	}
	function, ok := params[0].(*object.Function)
	if !ok {
		return &object.ErrorObject{Error: fmt.Sprintf("reduce first params should be function and got %s", params[0].Type())}
	}
	if len(function.Params) != 2 {
		return &object.ErrorObject{Error: fmt.Sprintf("reduce function should recive 2 params and recive %d", len(function.Params))}
	}
	if len(params[1:]) > 2 {
		return &object.ErrorObject{Error: fmt.Sprintf("reduce function should has max 2 arguments the list and initizial value and got %d", len(params[1:]))}
	}
	list, ok := params[1].(*object.List)
	if !ok {
		return &object.ErrorObject{Error: fmt.Sprintf("the second arguments is expected to be a list")}
	}
	var initialValue object.Object = nil
	var reduceParams []object.Object = list.Elements
	if len(params) == 3 {
		initialValue = params[2]
	} else {
		if len(reduceParams) == 0 {
			return &object.ErrorObject{Error: fmt.Sprintf("you provide empty list and not initial value, please dont be a fucking ass hole")}
		}
		initialValue = reduceParams[0]
		reduceParams = reduceParams[1:]
	}
	for _, objectParam := range reduceParams {
		initialValue = applyArgumentsToFunctionAndCall(function, []object.Object{initialValue, objectParam}, b, eval)
		if initialValue.Type() == object.ERROR_OBJ {
			return initialValue
		}
	}
	return initialValue
}

func _get_only_values(returnValues [][]object.Object) []object.Object {
	values := make([]object.Object, 0)
	for _, value := range returnValues {

		values = append(values, value[0])
	}
	return values
}

func _filter(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	elements := make([]object.Object, 0)
	if len(params) < 1 {
		return &object.ErrorObject{Error: fmt.Sprintf("filter expect more than 1 params and got %d", len(params))}
	}
	function, ok := params[0].(*object.Function)
	if !ok {
		return &object.ErrorObject{Error: fmt.Sprintf("filter first params should be function and got %s", params[0].Type())}
	}
	for _, objectParam := range params[1:] {
		returnValues := __callFunction(function, env, b, eval, objectParam)
		if len(returnValues) == 1 && len(returnValues[0]) == 1 && returnValues[0][0].Type() == object.ERROR_OBJ {
			return returnValues[0][0]
		}
		elements = append(elements, filterValues(returnValues)...)
	}
	return &object.List{Elements: elements}
}

func filterValues(values [][]object.Object) []object.Object {
	filteredValues := make([]object.Object, 0)
	for _, value := range values {
		if !evalTruthValue(value[0]) {
			continue
		}
		filteredValues = append(filteredValues, value[1])
	}
	return filteredValues
}

func evalTruthValue(value object.Object) bool {
	switch valueType := value.(type) {
	case *object.Boolean:
		return valueType.Value
	case *object.Integer:
		return valueType.Value != 0
	case *object.String:
		return len(valueType.Value) != 0
	default:
		return false
	}
}

func _list(_ *object.Environment, _ *Builtin, _ func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	return &object.List{Elements: params}
}

func _append(_ *object.Environment, _ *Builtin, _ func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	if len(params) < 1 {
		return &object.ErrorObject{Error: fmt.Sprintf("append expect more than 1 params and got %d", len(params))}
	}
	list, ok := params[0].(*object.List)
	if !ok {
		return &object.ErrorObject{Error: fmt.Sprintf("first params expected to be a list and got %s", params[0].Type())}
	}
	for _, element := range params[1:] {
		list.Elements = append(list.Elements, element)
	}
	return list
}

func _len(_ *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {

	if len(params) != 1 {
		return &object.ErrorObject{Error: fmt.Sprintf("len only revice 1 params and got %d", len(params))}
	}
	value := params[0]
	switch valueType := value.(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(valueType.Value))}
	case *object.List:
		return &object.Integer{Value: int64(len(valueType.Elements))}
	default:
		return &object.ErrorObject{Error: fmt.Sprintf("expected string or list type and got %s", value.Type())}
	}

}

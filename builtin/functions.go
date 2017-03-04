package builtin

import (
  "rootlang/object"
  "fmt"
  "rootlang/ast"
)

const (
  LEN    = "len"
  LIST   = "list"
  APPEND = "append"
  MAP    = "map"
  FILTER = "filter"
)

type function func(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object

type BuiltinFunction struct {
  Name     string
  Function function
  Params   []object.Object
}

func (builtinFunction *BuiltinFunction) Type() object.ObjectType {
  return object.BUILTIN_FUNCTION_OBJ
}

func (builtinFunction *BuiltinFunction) Inspect() string {
  return builtinFunction.Name
}

type Builtin struct {
  symbols map[string]object.Object
}

func New() *Builtin {
  symbols := registerSymbols()
  return &Builtin{symbols:symbols}
}

func (b *Builtin) GetObject(name string) (object.Object, bool) {
  value, ok := b.symbols[name]
  return value, ok
}

func registerSymbols() map[string]object.Object {
  symbols := make(map[string]object.Object)
  symbols[LEN] = getBuiltinFunction(_len, LEN)
  symbols[LIST] = getBuiltinFunction(_list, LIST)
  symbols[APPEND] = getBuiltinFunction(_append, APPEND)
  symbols[MAP] = getBuiltinFunction(_map, MAP)
  symbols[FILTER] = getBuiltinFunction(_filter, FILTER)
  return symbols
}

func getBuiltinFunction(f function, symbol string) *BuiltinFunction {
  return &BuiltinFunction{Name:symbol, Function:f, Params:make([]object.Object, 0)}
}

func _map(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
  elements := make([]object.Object, 0)
  if len(params) < 1 {
    return &object.ErrorObject{Error:fmt.Sprintf("map expect more than 1 params and got %d", len(params))}
  }
  function, ok := params[0].(*object.Function)
  if !ok {
    return &object.ErrorObject{Error:fmt.Sprintf("map first params should be function and got %s", params[0].Type())}
  }
  for _, objectParam := range params[1:] {
    returnValues := __callFunction(function, env, b, eval, objectParam)
    if len(returnValues) == 1 && len(returnValues[0]) == 1 && returnValues[0][0].Type() == object.ERROR_OBJ {
      return returnValues[0][0]
    }
    elements = append(elements, _get_only_values(returnValues)...)
  }
  return &object.List{Elements:elements}
}

func _get_only_values(returnValues [][]object.Object) []object.Object{
  values := make([]object.Object, 0)
  for _, value := range returnValues {

    values = append(values, value[0])
  }
  return values
}

func _filter(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
  elements := make([]object.Object, 0)
  if len(params) < 1 {
    return &object.ErrorObject{Error:fmt.Sprintf("filter expect more than 1 params and got %d", len(params))}
  }
  function, ok := params[0].(*object.Function)
  if !ok {
    return &object.ErrorObject{Error:fmt.Sprintf("filter first params should be function and got %s", params[0].Type())}
  }
  for _, objectParam := range params[1:] {
    returnValues := __callFunction(function, env, b, eval, objectParam)
    if len(returnValues) == 1 && len(returnValues[0]) == 1 &&  returnValues[0][0].Type() == object.ERROR_OBJ {
      return returnValues[0][0]
    }
    elements = append(elements, filterValues(returnValues)...)
  }
  return &object.List{Elements:elements}
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
    return &object.ErrorObject{Error:(fmt.Sprintf("this function takes at least %d arguments (%d given)", len(function.Params), len(params)))}
  }
  newEnvironment := applyArguments(function, params)
  if len(function.Params) == len(params) {
    returnValue := eval(function.Body, newEnvironment, builtinSymbols)
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

func _list(_ *object.Environment, _ *Builtin, _ func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
  return &object.List{Elements:params}
}

func _append(_ *object.Environment, _ *Builtin, _ func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
  if len(params) < 1 {
    return &object.ErrorObject{Error:fmt.Sprintf("append expect more than 1 params and got %d", len(params))}
  }
  list, ok := params[0].(*object.List)
  if !ok {
    return &object.ErrorObject{Error:fmt.Sprintf("first params expected to be a list and got %s", params[0].Type())}
  }
  for _, element := range params[1:] {
    list.Elements = append(list.Elements, element)
  }
  return list
}

func _len(_ *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {

  if len(params) != 1 {
    return &object.ErrorObject{Error:fmt.Sprintf("len only revice 1 params and got %d", len(params))}
  }
  value := params[0]
  switch valueType := value.(type) {
  case *object.String:
    return &object.Integer{Value:int64(len(valueType.Value))}
  case *object.List:
    return &object.Integer{Value:int64(len(valueType.Elements))}
  default:
    return &object.ErrorObject{Error:fmt.Sprintf("expected string or list type and got %s", value.Type())}
  }

}

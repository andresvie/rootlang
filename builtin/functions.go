package builtin

import (
  "rootlang/object"
  "fmt"
)
const(LEN = "len")

type function func(env * object.Environment,b *Builtin , params ...object.Object)object.Object

type BuiltinFunction struct {
  Name  string
  Function function
  Params  []object.Object
}

func (builtinFunction *BuiltinFunction) Type() object.ObjectType {
  return object.BUILTIN_FUNCTION_OBJ
}

func (builtinFunction *BuiltinFunction) Inspect() string {
  return builtinFunction.Name
}


type Builtin struct{
  symbols map[string] object.Object
}

func New()*Builtin{
  symbols := registerSymbols()
  return &Builtin{symbols:symbols}
}

func (b *Builtin) GetObject(name string) (object.Object, bool){
  value, ok := b.symbols[name]
  return value,ok
}

func registerSymbols() map[string] object.Object{
  symbols := make(map[string] object.Object)
  symbols[LEN] = getBuiltinFunction(_len, LEN)
  return symbols
}

func getBuiltinFunction(f function, symbol string) *BuiltinFunction{
  return &BuiltinFunction{Name:symbol, Function:f, Params:make([]object.Object, 0)}
}


func _len(_ * object.Environment,b *Builtin,params ...object.Object)object.Object{

  if len(params)!=1 {
    return &object.ErrorObject{Error:fmt.Sprintf("len only revice 1 params and got %d", len(params))}
  }
  value := params[0]
  switch valueType:= value.(type) {
  case *object.String:
    return &object.Integer{Value:int64(len(valueType.Value))}
  default:
    return &object.ErrorObject{Error:fmt.Sprintf("expected string type and got %s", value.Type())}
  }

}

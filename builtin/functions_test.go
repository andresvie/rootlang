package builtin

import (
  "testing"
  "rootlang/object"
)

func TestLenBuildExpression(t *testing.T) {
  b := New()
  lenFunction,ok := b.GetObject("len")
  test := &object.String{Value:"Carlos"}
  valueExpected := 6
  if !ok{
    t.Errorf("len functions is expected on the registers")
    return
  }
  if lenFunction.Type() != object.BUILTIN_FUNCTION_OBJ{
    t.Errorf("len functions should be type BUILTINT_FUNCTION_OBJECT and got %s", lenFunction.Type())
    return
  }
  l, ok:= lenFunction.(*BuiltinFunction)
  if !ok{
    t.Errorf("len functions should be type BUILTINT_FUNCTION_OBJECT")
    return
  }
  returnValue := l.Function(object.NewEnvironment(), b, test)
  value, ok:= returnValue.(*object.Integer)
  if !ok{
    t.Errorf("len functions should return Integer Object and got %s", returnValue.Type())
    return
  }

  if value.Value != int64(valueExpected){
    t.Errorf("value expected was %d and got %d", valueExpected, value.Value)
    return
  }


}

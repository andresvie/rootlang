package evaluator

import (
  "testing"
  "rootlang/lexer"
  "rootlang/parser"
  "rootlang/object"
)

func TestIntegerEvaluator(t *testing.T) {
  input := `5`
  l := lexer.New(input)
  programParser := parser.New(l)
  program := programParser.ParseProgram()
  returnValue := Eval(program, object.NewEnvironment())
  objectInteger, ok := returnValue.(*object.Integer)
  if !ok {
    t.Error("should return object integer")
    return
  }
  if objectInteger.Value != 5 {
    t.Error("should be has 5 integer value")
    return
  }
}

func TestFalseBooleanEvaluator(t *testing.T) {
  input := `false`
  l := lexer.New(input)
  programParser := parser.New(l)
  program := programParser.ParseProgram()
  returnValue := Eval(program, object.NewEnvironment())
  objectBoolean, ok := returnValue.(*object.Boolean)
  if !ok {
    t.Error("should return object boolean")
    return
  }
  if objectBoolean.Value {
    t.Error("should be has false boolean value")
    return
  }
}

func TestBooleanEvaluator(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    {"!true", false},
    {"!false", true},
    {"!5", false},
    {"!!true", true},
    {"!!false", false},
    {"!!5", true},
    {"!0", true},
    {"5 == 5", true},
    {"5 > 5", false},
    {"6 != 5", true},
    {"5 != 5", false},
    {"2 < 3", true},
    {"2 > 3", false},
    {"(2 > 3) == true", false},
    {"(2 < 3) == true", true},

  }
  for _, test := range tests {
    l := lexer.New(test.input)
    programParser := parser.New(l)
    program := programParser.ParseProgram()
    returnValue := Eval(program, object.NewEnvironment())
    objectBoolean, ok := returnValue.(*object.Boolean)
    if !ok {
      t.Errorf("should return object boolean %s", test.input)
      return
    }
    if test.expected != objectBoolean.Value {
      t.Errorf("should be has %t and got %t, %s", test.expected, objectBoolean.Value, test.input)
      return
    }
  }

}

func TestIntegerExpressionEvaluator(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    {"5", 5},
    {"-10", -10},
    {"-10 + 5", -5},
    {"10 + 5", 15},
    {"-10 - 5", -15},
    {"10 * 5", 50},
    {"10 * -6", -60},
  }
  for _, test := range tests {
    l := lexer.New(test.input)
    programParser := parser.New(l)
    program := programParser.ParseProgram()
    returnValue := Eval(program, object.NewEnvironment())
    objectInteger, ok := returnValue.(*object.Integer)
    if !ok {
      t.Error("should return integer object")
      return
    }
    if test.expected != objectInteger.Value {
      t.Errorf("should has %d and got %d %s", test.expected, objectInteger.Value, test.input)
      return
    }
  }

}

func TestIfExpressionEvaluator(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"if(true){20}", 20},
    {"if(false){20}", nil},
    {"if(false){20}else{10}", 10},
    {"if(true){}else{10}", nil},
    {"if(2<3){40+20}else{10}", 60},
    {"if(2>3){40+20}else{50}", 50},
  }
  for _, test := range tests {
    l := lexer.New(test.input)
    programParser := parser.New(l)
    program := programParser.ParseProgram()
    returnValue := Eval(program, object.NewEnvironment())
    if test.expected == nil && returnValue == nil {
      continue
    }
    objectInteger, ok := returnValue.(*object.Integer)

    if !ok {
      t.Errorf("should return integer object %s", test.input)
      return
    }
    if int64(test.expected.(int)) != objectInteger.Value {
      t.Errorf("should has %d and got %d %s", test.expected, objectInteger.Value, test.input)
      return
    }
  }

}

func TestReturnExpressionEvaluator(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    {"9;return 20;10", 20},
    {"return 9;return 20;10", 9},
    {"if (20>1){if(20>1){return 11;} return 12;}", 11},
  }
  for _, test := range tests {
    l := lexer.New(test.input)
    programParser := parser.New(l)
    program := programParser.ParseProgram()
    returnValue := Eval(program, object.NewEnvironment())
    objectInteger, ok := returnValue.(*object.Integer)
    if !ok {
      t.Errorf("should return integer object %s", test.input)
      return
    }
    if test.expected != objectInteger.Value {
      t.Errorf("should has %d and got %d %s", test.expected, objectInteger.Value, test.input)
      return
    }
  }

}

func TestErrorExpression(t *testing.T) {
  tests := []struct {
    input    string
    expected string
  }{
    {"false + false;return 1 + 1;10", "unknow operator for false + false"},
  }
  for _, test := range tests {
    l := lexer.New(test.input)
    programParser := parser.New(l)
    program := programParser.ParseProgram()
    returnValue := Eval(program, object.NewEnvironment())
    errorObject, ok := returnValue.(*object.ErrorObject)
    if !ok {
      t.Errorf("should return error object %s", test.input)
      return
    }
    if test.expected != errorObject.Error {
      t.Errorf("should has %d and got %d %s", test.expected, errorObject.Error, test.input)
      return
    }
  }
}

func TestLetExpression(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    {"let a=20;a;", 20},
  }
  for _, test := range tests {
    l := lexer.New(test.input)
    programParser := parser.New(l)
    program := programParser.ParseProgram()
    returnValue := Eval(program, object.NewEnvironment())
    integerObject, ok := returnValue.(*object.Integer)
    if !ok {
      t.Errorf("should integer object %s", test.input)
      return
    }
    if test.expected != integerObject.Value {
      t.Errorf("should has %d and got %d %s", test.expected, integerObject.Value, test.input)
      return
    }
  }
}

func TestFunctionExpression(t *testing.T) {
  tests := []struct {
    input    string
    expected string
    params   []string
  }{
    {"(x,y)=>{let x=y;return x+y;}", "(x,y)=>{let x = y;return (x + y);}", []string{"x", "y"}},
    {"(x)=>{return x;}", "(x)=>{return x;}", []string{"x"}},
    {"()=>{return 20;}", "()=>{return 20;}", []string{}},
  }
  for _, test := range tests {
    l := lexer.New(test.input)
    programParser := parser.New(l)
    program := programParser.ParseProgram()
    returnValue := Eval(program, object.NewEnvironment())
    functionObject, ok := returnValue.(*object.Function)
    if !ok {
      t.Errorf("should function object %s", test.input)
      return
    }
    if len(functionObject.Params) != len(test.params) {
      t.Errorf("should have %s and got %d", len(functionObject.Params), len(test.params))
      return
    }
    for i := 0; i < len(test.params); i++ {
      if test.params[i] != functionObject.Params[i].Value {
        t.Errorf("param expected is  %s and got %d", test.params[i], functionObject.Params[i].Value)
        return
      }
    }
    if test.expected != functionObject.Inspect() {
      t.Errorf("should have %s and got %s --> %s", test.expected, functionObject.Inspect(), test.input)
      return
    }
  }
}

func TestFunctionCallExpression(t *testing.T) {
  tests := []struct {
    input string
    value int64
  }{
   {"((x,y)=>{return x+y;})(10, 5);", 15},
    {"let add = (x,y)=>{return x+y;}; add(5,10);", 15},
    {"let x = 10;let add = (x,y)=>{return x+y;}; add(5,10);", 15},
    {"let z = 10;let add = (x,y)=>{return x+y+z;}; add(5,10);", 25},
    {"let z = (x,y)=>{let w = ()=>{return x+y;};return w;}; let b= z(10, 15); b();", 25},
    {"let z = (x,y)=>{ return x + y;}; let b= ()=>{return 2;}; z(23, b());", 25},
  }
  for _, test := range tests {
    l := lexer.New(test.input)
    programParser := parser.New(l)
    program := programParser.ParseProgram()
    returnValue := Eval(program, object.NewEnvironment())
    integerObject, ok := returnValue.(*object.Integer)
    if !ok {
      t.Errorf("should Integer object %s", test.input)
      return
    }
    if test.value != integerObject.Value {
      t.Errorf("should have %d and got %d %s", test.value, integerObject.Value, test.input)
      return
    }
  }



}


func TestStringObject(t *testing.T) {
  tests := []struct {
    input string
    value string
  }{
    {`"carlos viera"`, "carlos viera"},
    {`"carlos viera" + " hola mundo"`, "carlos viera hola mundo"},
    {`"carlos viera " + (3+5)`, "carlos viera 8"},
  }
  for _, test := range tests {
    l := lexer.New(test.input)
    programParser := parser.New(l)
    program := programParser.ParseProgram()
    returnValue := Eval(program, object.NewEnvironment())
    stringObject, ok := returnValue.(*object.String)
    if !ok {
      t.Errorf("should String object %s", test.input)
      return
    }
    if test.value != stringObject.Value {
      t.Errorf("should have %s and got %s %s", test.value, stringObject.Value, test.input)
      return
    }
  }
}


func TestClosure(t *testing.T) {
  tests := []struct {
    input string
    value int64
  }{
    {`let test = (x,y)=>{return x+y;};let test1 = test(1);test1(2);`, 3},
  }
  for _, test := range tests {
    l := lexer.New(test.input)
    programParser := parser.New(l)
    program := programParser.ParseProgram()
    returnValue := Eval(program, object.NewEnvironment())
    integerObject, ok := returnValue.(*object.Integer)
    if !ok {
      t.Errorf("should Integer object %s", test.input)
      return
    }
    if test.value != integerObject.Value {
      t.Errorf("should have %d and got %d %s", test.value, integerObject.Value, test.input)
      return
    }
  }
}

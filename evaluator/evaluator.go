package evaluator

import (
  "rootlang/ast"
  "rootlang/object"
  "fmt"
)

func Eval(node ast.Node, environment *object.Environment) object.Object {

  switch nodeType := node.(type) {
  case *ast.Program:
    return evalProgram(nodeType, environment)
  case *ast.BlockStatement:
    return evalStatement(nodeType.Statements, environment)
  case *ast.ExpressionStatement:
    return Eval(nodeType.Exp, environment)
  case *ast.IntegerLiteral:
    return &object.Integer{Value: nodeType.Value}
  case *ast.BoolExpression:
    return nativeToBooleanObject(nodeType.Value == "true")
  case *ast.StringExpression:
    return nativeStringToObject(nodeType.Value)
  case *ast.LetStatement:
    valueExpression := Eval(nodeType.Value, environment)
    if isError(valueExpression) {
      return valueExpression
    }
    environment.SetVar(nodeType.Name.Value, valueExpression)
    return nil
  case *ast.Identifier:
    value, ok := environment.GetVar(nodeType.Value)
    if !ok {
      return &object.ErrorObject{Error:fmt.Sprintf("%s was not declare", nodeType.Value)}
    }
    return value
  case *ast.FunctionExpression:
    return &object.Function{Params:nodeType.Params, Body:nodeType.Block, Env:environment.ExtendNewEnvironment()}
  case *ast.CallFunctionExpression:
    params := evalExpressions(nodeType.Arguments, environment)
    if len(params) == 1 && isError(params[0]) {
      return params[0]
    }
    value := Eval(nodeType.Function, environment)
    if isError(value) {
      return value
    }
    function, ok := value.(*object.Function)
    if !ok {
      return newError(fmt.Sprintf("expected function %s", value.Inspect()))
    }
    return applyArgumentsToFunctionAndCall(function, params)
  case *ast.PrefixExpression:
    rightExpression := Eval(nodeType.RightExpression, environment)
    if isError(rightExpression) {
      return rightExpression
    }
    return evalPrefixExpression(nodeType.Operator, rightExpression)
  case *ast.ReturnStatement:
    value := Eval(nodeType.Value, environment)
    if isError(value) {
      return value
    }
    return &object.ReturnObject{Value:value}
  case *ast.IfExpression:
    condition := Eval(nodeType.Condition, environment)
    if isError(condition) {
      return condition
    }
    if evalTruthValue(condition) {
      return Eval(nodeType.ConditionalBlock, environment)
    } else if nodeType.AlternativeBlock != nil {
      return Eval(nodeType.AlternativeBlock, environment)
    } else {
      return nil
    }
  case *ast.InfixExpression:
    leftExpression := Eval(nodeType.LeftExpression, environment)
    if isError(leftExpression) {
      return leftExpression
    }
    rightExpression := Eval(nodeType.RightExpression, environment)
    if isError(rightExpression) {
      return rightExpression
    }
    return evalInfixExpression(nodeType.Operator, rightExpression, leftExpression)
  }

  return nil
}

func applyArgumentsToFunctionAndCall(function *object.Function, params []object.Object) object.Object {

  if len(params) > len(function.Params) {
    return newError(fmt.Sprintf("this function takes at least %d arguments (%d given)", len(function.Params), len(params)))
  }
  newEnvironment := applyArguments(function, params)
  if len(function.Params) == len(params) {
    returnValue := Eval(function.Body, newEnvironment)
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

func evalExpressions(expressions []ast.Expression, environment *object.Environment) []object.Object {
  expressionsObjects := make([]object.Object, 0)
  for _, expression := range expressions {
    expressionObject := Eval(expression, environment)
    if isError(expressionObject) {
      return []object.Object{expressionObject}
    }
    expressionsObjects = append(expressionsObjects, expressionObject)
  }
  return expressionsObjects
}

func evalProgram(program *ast.Program, environment *object.Environment) object.Object {
  var result object.Object
  for _, statement := range program.Statements {
    result = Eval(statement, environment)
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
  return &object.ErrorObject{Error:error}
}

func evalInfixExpression(operator string, rightValue, leftValue object.Object) object.Object {

  if rightValue.Type() == object.INTEGER_OBJ && leftValue.Type() == object.INTEGER_OBJ {
    return evalIntegerInfixExpression(operator, rightValue, leftValue);
  }
  if (rightValue.Type() == object.STRING_OBJ || leftValue.Type() == object.STRING_OBJ) && operator == "+" {
    return nativeStringToObject(fmt.Sprintf("%s%s", leftValue.Inspect(), rightValue.Inspect()));
  }
  switch operator {
  case "==":
    return nativeToBooleanObject(leftValue == rightValue)
  case "!=":
    return nativeToBooleanObject(leftValue != rightValue)
  default:
    return &object.ErrorObject{Error:fmt.Sprintf("unknow operator for %s %s %s", leftValue.Inspect(), operator, rightValue.Inspect())}
  }

}

func nativeStringToObject(text string) object.Object {
  return &object.String{Value:text}
}

func evalIntegerInfixExpression(operator string, rightValue, leftValue object.Object) object.Object {
  rightIntegerValue := rightValue.(*object.Integer)
  leftIntegerValue := leftValue.(*object.Integer)
  switch operator {
  case "+":
    return &object.Integer{Value:leftIntegerValue.Value + rightIntegerValue.Value }
  case "-":
    return &object.Integer{Value:leftIntegerValue.Value - rightIntegerValue.Value }
  case "/":
    return &object.Integer{Value:leftIntegerValue.Value / rightIntegerValue.Value }
  case "*":
    return &object.Integer{Value:leftIntegerValue.Value * rightIntegerValue.Value }
  case "%":
    return &object.Integer{Value:leftIntegerValue.Value % rightIntegerValue.Value }
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
    return &object.ErrorObject{Error:fmt.Sprintf("unknow operator for -%s", rightValue.Inspect())}
  }
  value := rightValue.(*object.Integer)
  return &object.Integer{Value:-value.Value}
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

func evalStatement(statements []ast.Statement, environment *object.Environment) object.Object {
  var result object.Object
  for _, statement := range statements {
    result = Eval(statement, environment)
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

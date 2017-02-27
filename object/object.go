package object

import (
  "fmt"
  "rootlang/ast"
  "bytes"
  "strings"
)

type ObjectType string

const (
  INTEGER_OBJ  = "INTEGER"
  BOOLEAN_OBJ  = "BOOLEAN"
  NULL_OBJ     = "NULL"
  RETURN_OBJ   = "RETURN"
  ERROR_OBJ    = "ERROR"
  FUNCTION_OBJ = "FUNCTION"
  STRING_OBJ   = "STRING"
)

var (
  TRUE  = &Boolean{Value:true}
  FALSE = &Boolean{Value:false}
  NULL  = &Null{}
)

type Object interface {
  Type() ObjectType
  Inspect() string
}

type Integer struct {
  Value int64
}

func (integer *Integer) Type() ObjectType {
  return INTEGER_OBJ
}

func (integer *Integer) Inspect() string {
  return fmt.Sprintf("%d", integer.Value)
}

type Boolean struct {
  Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type String struct {
  Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnObject struct{ Value Object }

func (r *ReturnObject) Type() ObjectType { return RETURN_OBJ }
func (r *ReturnObject) Inspect() string  { return r.Inspect() }

type ErrorObject struct {
  Error string
}

func (e *ErrorObject) Type() ObjectType { return ERROR_OBJ }
func (e *ErrorObject) Inspect() string  { return e.Error }

type Function struct {
  Params []*ast.Identifier
  Body   *ast.BlockStatement
  Env    *Environment
}

func (f *Function) Clone(newParams []*ast.Identifier, env *Environment) *Function {
  return &Function{Body:f.Body, Env:env, Params:newParams}
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

func (f *Function) Inspect() string {
  buffer := bytes.NewBufferString("(")
  paramsText := make([]string, 0)
  if f.Params != nil {
    for _, param := range f.Params {
      paramsText = append(paramsText, param.String())
    }
  }

  buffer.WriteString(strings.Join(paramsText, ","))
  buffer.WriteString(")=>")
  buffer.WriteString(f.Body.String())
  return buffer.String()
}

type Environment struct {
  vars  map[string]Object
  outer *Environment
}

func (e*Environment) ExtendNewEnvironment() *Environment {
  newEnvironment := NewEnvironment()
  newEnvironment.outer = e
  return newEnvironment
}

func NewEnvironment() *Environment {
  return &Environment{vars:make(map[string]Object), outer:nil}
}

func (e*Environment) GetVar(name string) (Object, bool) {
  value, ok := e.vars[name]
  if ok {
    return value, ok
  }
  if e.outer != nil {
    value, ok = e.outer.GetVar(name)
  }
  return value, ok
}
func (e*Environment) SetVar(name string, value Object) {
  e.vars[name] = value
}

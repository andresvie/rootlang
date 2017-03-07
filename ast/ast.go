package ast

import (
  "rootlang/lexer"
  "bytes"
  "fmt"
  "strings"
)

type Node interface {
  TokenLiteral() string
  String() string
}
type Statement interface {
  Node
  statementNode()
}
type Expression interface {
  Node
  expressionNode()
}

type IntegerLiteral struct {
  Token lexer.Token
  Value int64
}

func (int *IntegerLiteral) expressionNode() {

}

func (int *IntegerLiteral) TokenLiteral() string {
  return int.Token.Literal
}

func (int *IntegerLiteral) String() string {
  return fmt.Sprintf("%d", int.Value)
}


type  ParamsExpression struct{
  Token lexer.Token
  Params []*Identifier
}



func (params *ParamsExpression) expressionNode() {

}

func (params *ParamsExpression) TokenLiteral() string {
  return ""
}

func (params *ParamsExpression) String() string {
  buffer := bytes.NewBufferString("(")
  paramsText := make([]string, 0)
  for _,param := range params.Params{
    paramsText = append(paramsText, param.String())
  }
  buffer.WriteString(strings.Join(paramsText, ","))
  buffer.WriteString(")")
  return buffer.String()
}



type StringExpression struct {
  Token lexer.Token
  Value string
}

func (str *StringExpression) expressionNode() {

}

func (str *StringExpression) TokenLiteral() string {
  return str.Token.Literal
}

func (str *StringExpression) String() string {
  return str.Value
}

type PrefixExpression struct {
  Token           lexer.Token
  Operator        string
  RightExpression Expression
}

func (prefix *PrefixExpression) TokenLiteral() string {
  return prefix.Token.Literal
}

func (prefix *PrefixExpression) expressionNode() {

}

func (prefix *PrefixExpression) String() string {
  var buffer *bytes.Buffer = bytes.NewBufferString(prefix.Operator)
  buffer.WriteString("(")
  buffer.WriteString(prefix.RightExpression.String())
  buffer.WriteString(")")
  return buffer.String()
}

type LetStatement struct {
  Token lexer.Token
  Name  *Identifier
  Value Expression
}

type ImportStatement struct{
  Token lexer.Token
  Path string
  Name  *Identifier
};

func (im *ImportStatement) statementNode() {
}

func (im *ImportStatement) TokenLiteral() string {
  return im.Token.Literal
}

func (im *ImportStatement) String() string {

  return fmt.Sprintf(`import "%s" as %s`, im.Path, im.Name.String())
}

func (let *LetStatement) statementNode() {
}

func (let *LetStatement) TokenLiteral() string {
  return let.Token.Literal
}

func (let *LetStatement) String() string {
  var buffer *bytes.Buffer = bytes.NewBufferString("let ");
  buffer.WriteString(let.Name.String())
  buffer.WriteString(" = ")
  buffer.WriteString(let.Value.String())
  buffer.WriteString(";")
  return buffer.String()
}

type ReturnStatement struct {
  Token lexer.Token
  Value Expression
}

func (ret *ReturnStatement) statementNode() {

}
func (ret *ReturnStatement) String() string {
  var buffer *bytes.Buffer = bytes.NewBufferString("return ");
  buffer.WriteString(ret.Value.String())
  buffer.WriteString(";")
  return buffer.String();
}

func (ret *ReturnStatement) TokenLiteral() string {
  return ret.Token.Literal
}

type ExpressionStatement struct {
  Exp Expression
}

func (exp *ExpressionStatement) statementNode() {}
func (exp *ExpressionStatement) String() string {
  buffer := bytes.NewBufferString("")
  buffer.WriteString(exp.Exp.String())
  buffer.WriteString(";")
  return buffer.String()
}

func (exp *ExpressionStatement) TokenLiteral() string {
  return exp.Exp.TokenLiteral()
}

type BlockStatement struct {
  Token      lexer.Token
  Statements []Statement
}

func (block *BlockStatement) statementNode() {}
func (block *BlockStatement) String() string {
  buffer := bytes.NewBufferString("{")
  for _, statement := range block.Statements {
    buffer.WriteString(statement.String())
  }
  buffer.WriteString("}")
  return buffer.String()
}

func (block *BlockStatement) TokenLiteral() string {
  return block.Token.Literal
}

type IfExpression struct {
  Token            lexer.Token
  Condition        Expression
  ConditionalBlock *BlockStatement
  AlternativeBlock *BlockStatement
}

func (ifExp *IfExpression) expressionNode() {}
func (ifExp *IfExpression) String() string {
  buffer := bytes.NewBufferString("")
  buffer.WriteString(fmt.Sprintf("if(%s)", ifExp.Condition.String()))
  buffer.WriteString(ifExp.ConditionalBlock.String())
  if ifExp.AlternativeBlock != nil {
    buffer.WriteString(fmt.Sprintf("else%s", ifExp.AlternativeBlock.String()))
  }
  return buffer.String()
}

func (ifExp *IfExpression) TokenLiteral() string {
  return ifExp.Token.Literal
}

type InfixExpression struct {
  Token           lexer.Token
  LeftExpression  Expression
  Operator        string
  RightExpression Expression
}

type CallFunctionExpression struct {
  Token     lexer.Token
  Function  Expression
  Arguments []Expression
}

func (funcCall *CallFunctionExpression) expressionNode() {}
func (funcCall *CallFunctionExpression) String() string {
  args := make([]string, 0)
  buffer := bytes.NewBufferString(funcCall.Function.String())
  buffer.WriteString("(");
  for _, arg := range funcCall.Arguments {
    args = append(args, arg.String())
  }
  buffer.WriteString(strings.Join(args, ","))
  buffer.WriteString(")");
  return buffer.String()
}

func (funcCall *CallFunctionExpression) TokenLiteral() string {
  return funcCall.Token.Literal
}

func (infix *InfixExpression) expressionNode() {}

func (infix *InfixExpression) String() string {
  buffer := bytes.NewBufferString("(")
  buffer.WriteString(infix.LeftExpression.String())
  buffer.WriteString(fmt.Sprintf(" %s ", infix.Operator))
  buffer.WriteString(fmt.Sprintf("%s)", infix.RightExpression.String()))
  return buffer.String()
}

func (infix *InfixExpression) TokenLiteral() string {
  return infix.Token.Literal
}

type Identifier struct {
  Token lexer.Token
  Value string
}

func (id *Identifier) expressionNode() {
}

func (id *Identifier) TokenLiteral() string {
  return id.Token.Literal
}

func (id *Identifier) String() string {
  return id.Value
}

type FunctionExpression struct {
  Token  lexer.Token
  Params []*Identifier
  Block  *BlockStatement
}

func (fnExpression *FunctionExpression) String() string {
  buffer := bytes.NewBufferString("(")
  params := make([]string, 0)
  for _, param := range fnExpression.Params {
    params = append(params, param.String())
  }
  buffer.WriteString(strings.Join(params, ","))
  buffer.WriteString(")=>")
  buffer.WriteString(fnExpression.Block.String())
  return buffer.String()
}

func (fnExpression *FunctionExpression) expressionNode() {

}

func (fnExpression *FunctionExpression) TokenLiteral() string {
  return fnExpression.Token.Literal
}

type BoolExpression struct {
  Token lexer.Token
  Value string
}

func (boolExpression *BoolExpression) String() string {
  return boolExpression.Value
}

func (boolExpression *BoolExpression) expressionNode() {

}

func (boolExpression *BoolExpression) TokenLiteral() string {
  return boolExpression.Token.Literal
}

type Program struct {
  Statements []Statement
}

func (p *Program) TokenLiteral() string {
  if len(p.Statements) > 0 {
    return p.Statements[0].TokenLiteral()
  }
  return ""
}

func (p *Program) String() string {
  var buffer *bytes.Buffer = bytes.NewBufferString("");
  for _, statement := range p.Statements {
    buffer.WriteString(statement.String())
  }
  return buffer.String()
}

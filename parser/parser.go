package parser

import (
	"rootlang/lexer"
	"rootlang/ast"
	"strconv"
	"fmt"
	"strings"
)

const (
	_        int = iota
	LOWEST
	EQUALS
	SUM
	PRODUCT
	PREFIX
	CALL
	FUNCTION
)

var precedences = map[lexer.TokenType]int{lexer.EQUAL: EQUALS,
	lexer.NOTEQUAL:                                EQUALS,
	lexer.LESSTHAN:                                EQUALS,
	lexer.MORETHAN:                                EQUALS,
	lexer.PLUS:                                    SUM,
	lexer.MINUS:                                   SUM,
	lexer.MODULE:                                  PRODUCT,
	lexer.MULTIPLY:                                PRODUCT,
	lexer.MOD:                                     PRODUCT,
	lexer.DIV:                                     PRODUCT,
	lexer.LPAREN:                                  CALL,
	lexer.FUNCTION:                                FUNCTION,

}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)
type Parser struct {
	errors          []string
	prefixErrors    []string
	prefixFunctions map[lexer.TokenType]prefixParseFn
	infixFunctions  map[lexer.TokenType]infixParseFn
	l               *lexer.Lexer
	curToken        lexer.Token
	peekToken       lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.prefixFunctions = make(map[lexer.TokenType]prefixParseFn)
	p.infixFunctions = make(map[lexer.TokenType]infixParseFn)
	p.registerPrefixFunction()
	p.registerInfixFunction()
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) GetErrors() []string {
	return p.errors
}

func (p *Parser) registerPrefixFunction() {
	p.prefixFunctions[lexer.INT] = p.parseIntExpression
	p.prefixFunctions[lexer.IDENT] = p.parseIdentifierExpression
	p.prefixFunctions[lexer.STRING] = p.parseStringExpression
	p.prefixFunctions[lexer.MINUS] = p.parsePrefixExpression
	p.prefixFunctions[lexer.NOT] = p.parsePrefixExpression
	p.prefixFunctions[lexer.TRUE] = p.parseBoolExpression
	p.prefixFunctions[lexer.FALSE] = p.parseBoolExpression
	p.prefixFunctions[lexer.LPAREN] = p.parseGroupedExpression
	p.prefixFunctions[lexer.IF] = p.parseIfExpression

}

func (p *Parser) parseStringExpression() ast.Expression {
	return &ast.StringExpression{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExpression := &ast.IfExpression{Token: p.curToken}
	if !p.isNextTokenExpected(lexer.LPAREN) {
		p.errors = append(p.errors, "left parent is expected in if expression")
		return nil
	}
	p.nextToken()
	p.nextToken()
	condition := p.parseExpression(LOWEST)
	if condition == nil {
		p.errors = append(p.errors, "condition is required on if expression")
		return nil
	}
	if !p.isNextTokenExpected(lexer.RPAREN) {
		p.errors = append(p.errors, "right parent is expected in if expression")
		return nil
	}
	p.nextToken()
	if !p.isNextTokenExpected(lexer.LBRACE) {
		p.errors = append(p.errors, "block for if expression is required")
		return nil
	}
	p.nextToken()
	conditionalBlock := p.parserBlockStatement()
	if conditionalBlock == nil {
		p.errors = append(p.errors, "block for if expression is required")
		return nil
	}
	ifExpression.Condition = condition
	ifExpression.ConditionalBlock = conditionalBlock.(*ast.BlockStatement)
	if !p.isNextTokenExpected(lexer.ELSE) {
		return ifExpression
	}
	p.nextToken()
	if !p.isNextTokenExpected(lexer.LBRACE) {
		p.errors = append(p.errors, "begin of block for else expression is expected")
		return nil
	}
	p.nextToken()
	alternativeBlock := p.parserBlockStatement()
	if alternativeBlock == nil {
		p.errors = append(p.errors, "block for else is expected")
		return nil
	}
	ifExpression.AlternativeBlock = alternativeBlock.(*ast.BlockStatement)
	return ifExpression
}

func (p *Parser) parseBoolExpression() ast.Expression {
	return &ast.BoolExpression{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	if p.isTokenExpected(lexer.IDENT) && p.isNextTokenExpected(lexer.COMMA) {
		return p.parseParams()
	} else if p.isTokenExpected(lexer.RPAREN) {
		return &ast.ParamsExpression{Params: make([]*ast.Identifier, 0)}
	} else {
		groupExpression := p.parseExpression(LOWEST)
		if !p.isNextTokenExpected(lexer.RPAREN) {
			return nil
		}
		p.nextToken()
		return groupExpression
	}
}

func (p *Parser) parseParams() *ast.ParamsExpression {
	params := make([]*ast.Identifier, 0)
	for p.isTokenExpected(lexer.IDENT) {
		expression := p.parseIdentifierExpression()
		identifier := expression.(*ast.Identifier)
		params = append(params, identifier)
		p.moveNextTokenExpected(lexer.COMMA)
		p.nextToken()

	}
	return &ast.ParamsExpression{Params: params}
}

func (p *Parser) registerInfixFunction() {
	p.infixFunctions[lexer.PLUS] = p.parseInfixExpression
	p.infixFunctions[lexer.MINUS] = p.parseInfixExpression
	p.infixFunctions[lexer.DIV] = p.parseInfixExpression
	p.infixFunctions[lexer.MULTIPLY] = p.parseInfixExpression
	p.infixFunctions[lexer.MOD] = p.parseInfixExpression
	p.infixFunctions[lexer.EQUAL] = p.parseInfixExpression
	p.infixFunctions[lexer.NOTEQUAL] = p.parseInfixExpression
	p.infixFunctions[lexer.LESSTHAN] = p.parseInfixExpression
	p.infixFunctions[lexer.MORETHAN] = p.parseInfixExpression
	p.infixFunctions[lexer.MODULE] = p.parseInfixExpression
	p.infixFunctions[lexer.LPAREN] = p.parseCallFunctionExpression
	p.infixFunctions[lexer.FUNCTION] = p.parseFunctionExpression

}

func (p *Parser) parseFunctionExpression(params ast.Expression) ast.Expression {
	functionToken := p.curToken
	p.nextToken()
	var paramsExpression []*ast.Identifier = nil
	var blogStatement *ast.BlockStatement = nil
	var ok bool
	switch paramsType := params.(type) {
	case *ast.ParamsExpression:
		paramsExpression = paramsType.Params
	case *ast.Identifier:
		paramsExpression = []*ast.Identifier{paramsType}
	default:
		p.errors = append(p.errors, "params are expected")
		return nil
	}
	if !p.isTokenExpected(lexer.LBRACE) {
		expression := p.parseExpression(LOWEST)
		if expression == nil {
			p.errors = append(p.errors, "expression was expected on lambda function")
			return nil
		}
		token := lexer.Token{Type: lexer.RETURN, Literal: "return"}
		blockStatementToken := lexer.Token{Type: lexer.LBRACE, Literal: "{"}
		returnStatement := &ast.ReturnStatement{Value: expression, Token: token}
		statements := []ast.Statement{returnStatement}
		blogStatement = &ast.BlockStatement{Token: blockStatementToken, Statements: statements}
	} else {
		blogStatement, ok = p.parserBlockStatement().(*ast.BlockStatement)
		if !ok {
			p.errors = append(p.errors, "function definition block is expected")
			return nil
		}
	}

	return &ast.FunctionExpression{functionToken, paramsExpression, blogStatement}

}

func (p *Parser) parseCallFunctionExpression(function ast.Expression) ast.Expression {
	callFunctionExpression := &ast.CallFunctionExpression{Token: p.curToken, Function: function}
	p.nextToken()
	arguments := p.parseArguments()
	if arguments == nil {
		p.errors = append(p.errors, "error parsing function arguments")
		return nil
	}
	callFunctionExpression.Arguments = arguments
	return callFunctionExpression
}
func (p *Parser) parseArguments() []ast.Expression {
	arguments := make([]ast.Expression, 0)
	for !p.isTokenExpected(lexer.RPAREN) {
		expression := p.parseExpression(LOWEST)
		if expression == nil {
			return nil
		}
		arguments = append(arguments, expression)
		if (p.moveNextTokenExpected(lexer.COMMA)) {
			p.nextToken()
		} else {
			break
		}

	}
	if !p.isNextTokenExpected(lexer.SEMICOLON) {
		p.nextToken()
	}

	return arguments
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {

	statements := make([]ast.Statement, 0)

	var statement ast.Statement
	for p.curToken.Type != lexer.EOF {
		statement = p.parseStatement()
		if statement != nil {
			statements = append(statements, statement)
		}
		p.nextToken()
	}
	return &ast.Program{Statements: statements}
}

func (p *Parser) parseStatement() ast.Statement {
	if p.curToken.Type == lexer.LET {
		return p.parseLetStatement()
	} else if p.curToken.Type == lexer.RETURN {
		return p.parseReturnStatement()
	} else if p.curToken.Type == lexer.LBRACE {
		return p.parserBlockStatement()
	} else if p.curToken.Type == lexer.IMPORT {
		return p.parseImportStatement()
	} else {
		return p.parseExpressionStatement()
	}
	return nil
}

func (p *Parser) parseImportStatement() ast.Statement {
	token := p.curToken
	p.nextToken()
	if !p.isTokenExpected(lexer.STRING) {
		p.errors = append(p.errors, "string path is expected")
		return nil
	}
	var identity *ast.Identifier
	path := p.curToken.Literal
	if !p.isNextTokenExpected(lexer.AS) {
		names := strings.Split(path, "/")
		name := names[len(names)-1]
		identity = &ast.Identifier{Token: lexer.Token{Literal: name, Type: lexer.IDENT}, Value: name}
	} else {
		p.nextToken()
		if !p.isNextTokenExpected(lexer.IDENT) {
			p.errors = append(p.errors, "identity is expected")
			return nil
		}
		p.nextToken()
		identity = p.parseIdentifierExpression().(*ast.Identifier)
	}

	return &ast.ImportStatement{Token: token, Path: path, Name: identity}
}

func (p *Parser) parserBlockStatement() ast.Statement {
	statements := make([]ast.Statement, 0)
	blockStatement := &ast.BlockStatement{Token: p.curToken}
	for !(p.isNextTokenExpected(lexer.RBRACE) || p.isNextTokenExpected(lexer.EOF)) {
		p.nextToken()
		statement := p.parseStatement()
		if statement == nil {
			continue
		}
		statements = append(statements, statement)
	}
	blockStatement.Statements = statements
	if !p.isNextTokenExpected(lexer.RBRACE) {
		p.errors = append(p.errors, "right brace is expected")
		return nil
	}
	p.nextToken()
	return blockStatement
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	statement := &ast.ExpressionStatement{}
	exp := p.parseExpression(LOWEST)
	if exp == nil {
		return nil
	}
	if p.isNextTokenExpected(lexer.SEMICOLON) {
		p.nextToken()
	}
	statement.Exp = exp
	return statement
}

func (p *Parser) parseLetStatement() ast.Statement {
	token := p.curToken
	if p.peekToken.Type != lexer.IDENT {
		p.errors = append(p.errors, "ident is expected after let")
		return nil
	}
	p.nextToken()
	if p.peekToken.Type != lexer.ASSIGN {
		p.errors = append(p.errors, "after declaration equal sign is expected")
		return nil
	}
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	p.nextToken()
	value := p.parseExpression(LOWEST)
	if value == nil {
		return nil
	}
	var statement ast.Statement = &ast.LetStatement{Token: token, Name: ident, Value: value}
	p.nextToken()
	return statement

}

func (p *Parser) parseReturnStatement() ast.Statement {
	token := p.curToken
	p.nextToken()
	expression := p.parseExpression(LOWEST)
	if !p.moveNextTokenExpected(lexer.SEMICOLON) {
		p.errors = append(p.errors, "semicolon is expected")
		return nil
	}
	var returnStatement *ast.ReturnStatement = &ast.ReturnStatement{Token: token, Value: expression}
	return returnStatement
}

func (p *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression {
	infixExpression := &ast.InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, LeftExpression: leftExpression}
	precedence := p.getCurrentPrecedenceToken()
	p.nextToken()
	infixExpression.RightExpression = p.parseExpression(precedence)
	return infixExpression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefixExpresion := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	prefixExpresion.RightExpression = p.parseExpression(PREFIX)
	return prefixExpresion
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFunction := p.prefixFunctions[p.curToken.Type]
	if prefixFunction == nil {
		p.registerPrefixParserError()
		return nil
	}
	leftExpression := prefixFunction()
	for !p.isNextTokenExpected(lexer.SEMICOLON) && precedence < p.getPeekPrecedenceToken() {
		infixFunction, ok := p.infixFunctions[p.peekToken.Type]
		if !ok {
			return leftExpression
		}
		p.nextToken()
		leftExpression = infixFunction(leftExpression)
	}
	return leftExpression
}

func (p *Parser) getCurrentPrecedenceToken() int {
	precedence, ok := precedences[p.curToken.Type]
	if !ok {
		return LOWEST
	}
	return precedence
}

func (p *Parser) getPeekPrecedenceToken() int {
	precedence, ok := precedences[p.peekToken.Type]
	if !ok {
		return LOWEST
	}
	return precedence
}

func (p *Parser) registerPrefixParserError() {
	p.prefixErrors = append(p.prefixErrors, fmt.Sprintf("function for %s tokne not found", p.curToken.Type))
}

func (p *Parser) parseIntExpression() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Literal, 10, 0)
	if err != nil {
		p.errors = append(p.errors, "integer is expected")
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: val}
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) moveNextTokenExpected(tokenType lexer.TokenType) bool {
	if !p.isNextTokenExpected(tokenType) {
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) isTokenExpected(tokenType lexer.TokenType) bool {
	return p.curToken.Type == tokenType;
}

func (p *Parser) isNextTokenExpected(tokenType lexer.TokenType) bool {
	return p.peekToken.Type == tokenType;
}

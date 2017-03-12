package parser

import (
	"testing"
	"rootlang/lexer"
	"rootlang/ast"
	"strconv"
)

func TestLetStatementLiteralValue(t *testing.T) {
	input := `
		   let x = 5;
		   let y = 10;
		   let foobar = 838383;
   	`
	l := lexer.New(input)
	p := New(l)
	statementsExpected := []ast.LetStatement{createLiteralLetStatement("x", "5"), createLiteralLetStatement("y", "10"), createLiteralLetStatement("foobar", "838383")}
	program := p.ParseProgram()
	if len(program.Statements) != 3 {
		t.Fatal("program should has 3 statements")
		return
	}
	for i := 0; i < len(program.Statements); i++ {
		var let *ast.LetStatement = program.Statements[i].(*ast.LetStatement)
		if !assertLiteralLetStatement(let, &statementsExpected[i]) {
			t.Errorf("literatel let statement %s should be equal to %s", statementsExpected[i].TokenLiteral(), program.Statements[i].TokenLiteral())
		}
	}
}

func TestStringExpression(t *testing.T) {
	input := `"carlos viera"`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != 1 {
		t.Error("should have 1 statement")
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		return
	}
	expression, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Error("Expression Statements is expected")
		return
	}
	stringExpression, okStringExpression := expression.Exp.(*ast.StringExpression)
	if !okStringExpression {
		t.Error("String Expression is expected")
		return
	}
	if stringExpression.Value != "carlos viera" {
		t.Error(`"carlos viera" is expected`)
		return
	}

}

func TestIntegerExpression(t *testing.T) {
	input := `5`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != 1 {
		t.Error("should statements 1")
		return
	}
	expression, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Error("Expression Statements is expected")
		return
	}
	integerExpression, okIntegerExpression := expression.Exp.(*ast.IntegerLiteral)
	if !okIntegerExpression {
		t.Error("Integer expression is expected")
		return
	}
	if integerExpression.Value != 5 {
		t.Error("5 values is expected")
		return
	}

}

func TestBooleanExpression(t *testing.T) {
	input := `false`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != 1 {
		t.Error("should statements 1")
		return
	}
	expression, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Error("Expression Statements is expected")
		return
	}
	booleanExpression, okBooleanExpression := expression.Exp.(*ast.BoolExpression)
	if !okBooleanExpression {
		t.Error("Boolean expression is expected")
		return
	}
	if booleanExpression.Value != "false" {
		t.Error("false values is expected")
		return
	}

}

func TestGroupedExpression(t *testing.T) {
	input := `
		   let x = (a + b) * c;
		   let y = a + b + (a * b);
   	`
	expectedStatements := []string{
		"let x = ((a + b) * c);", "let y = ((a + b) + (a * b));"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != len(expectedStatements) {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	for index, statement := range program.Statements {
		if expectedStatements[index] != statement.String() {
			showParserErrors(p, t)
			showPrefixParserError(p, t)
			t.Errorf("statement expected is %s and got %s", expectedStatements[index], statement.String())
		}
	}
}

func TestBoolExpression(t *testing.T) {
	input := `
		   let x = false;
		   let y = true;
		   return false == true;
		   return false;
		   return true;
   	`
	expectedStatements := []string{
		"let x = false;", "let y = true;", "return (false == true);", "return false;", "return true;"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != len(expectedStatements) {
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	for index, statement := range program.Statements {
		if expectedStatements[index] != statement.String() {
			showParserErrors(p, t)
			showPrefixParserError(p, t)
			t.Errorf("statement expected is %s and got %s", expectedStatements[index], statement.String())
		}
	}
}

func TestInfixExpression(t *testing.T) {
	/*

			 ,
			 */
	input := `
	let x = a + b;
			 let y = a * b + c;
			 let foobar = a + c * b;
			 return a + b / c;
			 return -a + b - c;
			 return a > b;
			 return a < b;
			 return a == b;
			 return a != b;
		   return net::listen();
   	`
	expectedStatements := []string{"let x = (a + b);", "let y = ((a * b) + c);", "let foobar = (a + (c * b));",
				       "return (a + (b / c));",
				       "return ((-(a) + b) - c);", "return (a > b);", "return (a < b);",
				       "return (a == b);", "return (a != b);",
				       "return (net :: listen());"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != len(expectedStatements) {
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	for index, statement := range program.Statements {
		if expectedStatements[index] != statement.String() {
			showParserErrors(p, t)
			showPrefixParserError(p, t)
			t.Errorf("statement expected is %s and got %s", expectedStatements[index], statement.String())
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `
		  if (x > y)
		  {
		  	return y;
		  }
   	`
	expectedStatements := []string{"if((x > y)){return y;}"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != len(expectedStatements) {
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Error("Expression Statement is expected")
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		return
	}
	ifExpression, ok := expressionStatement.Exp.(*ast.IfExpression)
	if !ok {
		t.Error("if Expression is expected")
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		return
	}
	if (ifExpression.AlternativeBlock != nil) {
		t.Error("else expression is not expected")
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		return
	}
	if (ifExpression.String() != expectedStatements[0]) {
		t.Errorf("if expression expected %s and got %s", expectedStatements[0], ifExpression.String())
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		return
	}

}

func TestIfExpressionWithElse(t *testing.T) {
	input := `
		  if (x > y)
		  {
		  	return y;
		  }else
		  {
		  	return x;
		  }
   	`
	expectedStatements := []string{"if((x > y)){return y;}else{return x;}"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != len(expectedStatements) {
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Error("Expression Statement is expected")
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		return
	}
	ifExpression, ok := expressionStatement.Exp.(*ast.IfExpression)
	if !ok {
		t.Error("if Expression is expected")
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		return
	}
	if (ifExpression.AlternativeBlock == nil) {
		t.Error("else expression is expected")
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		return
	}
	if (ifExpression.String() != expectedStatements[0]) {
		t.Errorf("if expression expected %s and got %s", expectedStatements[0], ifExpression.String())
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		return
	}

}

func TestCallFunctionExpressionWithoutArguments(t *testing.T) {
	input := `add();`
	expectedStatements := []string{"add()"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != len(expectedStatements) {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	expressionStatement, okExpression := program.Statements[0].(*ast.ExpressionStatement)
	if !okExpression {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("expression statement is expected")
		return
	}

	callFunctionExpression, okCallFunctionExpression := expressionStatement.Exp.(*ast.CallFunctionExpression)
	if !okCallFunctionExpression {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("call function expression is expected")
		return
	}
	args := callFunctionExpression.Arguments
	if len(args) != 0 {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("number of argument expected is 0 and got %d", len(args))
		return
	}
	if callFunctionExpression.String() != expectedStatements[0] {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("call function expression should %s and got %s", expectedStatements[0], callFunctionExpression.String())
		return
	}
}

func TestCallFunctionExpression(t *testing.T) {
	input := `add(2, x(2,3))`
	expectedStatements := []string{"add(2,x(2,3))"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(p.GetErrors()) != 0{
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Error("not error was expected\n")
		return
	}
	if len(program.Statements) != len(expectedStatements) {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	expressionStatement, okExpression := program.Statements[0].(*ast.ExpressionStatement)
	if !okExpression {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("expression statement is expected")
		return
	}

	callFunctionExpression, okCallFunctionExpression := expressionStatement.Exp.(*ast.CallFunctionExpression)
	if !okCallFunctionExpression {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("call function expression is expected")
		return
	}
	args := callFunctionExpression.Arguments
	if len(args) != 2 {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("number of argument expected is 2 and got %d", len(args))
		return
	}
	if callFunctionExpression.String() != expectedStatements[0] {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("call function expression should %s and got %s", expectedStatements[0], callFunctionExpression.String())
		return
	}
}

func TestFunctionExpressionWithoutParams(t *testing.T) {
	input := ` () =>
	{
		return 5;
	}
   	`
	expectedStatements := []string{"()=>{return 5;}"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != len(expectedStatements) {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	expressionStatement, okExpression := program.Statements[0].(*ast.ExpressionStatement)
	if !okExpression {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("expression statement is expected")
		return
	}

	functionExpression, okFunctionExpression := expressionStatement.Exp.(*ast.FunctionExpression)
	if !okFunctionExpression {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("function expression is expected")
		return
	}
	if functionExpression.String() != expectedStatements[0] {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("function expression should %s and got %s", expectedStatements[0], functionExpression.String())
		return
	}

}


func TestFunctionLambdaShorcut(t *testing.T) {
	input := ` () => 5;
   	`
	expectedStatements := []string{"()=>{return 5;}"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != len(expectedStatements) {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	expressionStatement, okExpression := program.Statements[0].(*ast.ExpressionStatement)
	if !okExpression {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("expression statement is expected")
		return
	}

	functionExpression, okFunctionExpression := expressionStatement.Exp.(*ast.FunctionExpression)
	if !okFunctionExpression {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("function expression is expected")
		return
	}
	if functionExpression.String() != expectedStatements[0] {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("function expression should %s and got %s", expectedStatements[0], functionExpression.String())
		return
	}

}

func TestFunctionExpression(t *testing.T) {
	input := ` (x,y) =>
	{
		return x;
	}
   	`
	expectedStatements := []string{"(x,y)=>{return x;}"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != len(expectedStatements) {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	expressionStatement, okExpression := program.Statements[0].(*ast.ExpressionStatement)
	if !okExpression {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("expression statement is expected")
		return
	}

	functionExpression, okFunctionExpression := expressionStatement.Exp.(*ast.FunctionExpression)
	if !okFunctionExpression {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("function expression is expected")
		return
	}
	if (functionExpression.String() != expectedStatements[0]) {
		showParserErrors(p, t)
		showPrefixParserError(p, t)
		t.Errorf("function expression should %s and got %s", expectedStatements[0], functionExpression.String())
		return
	}

}

func TestPrefixExpression(t *testing.T) {
	input := `
		   let x = -5;
		   let y = !a;
		   return !a;
		   return -4;
		   return -w;
   	`
	expectedStatements := []string{"let x = -(5);", "let y = !(a);", "return !(a);", "return -(4);", "return -(w);"}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != len(expectedStatements) {
		t.Errorf("number of statement expected is %d and got %d", len(expectedStatements), len(program.Statements))
		return
	}
	for index, statement := range program.Statements {
		if expectedStatements[index] != statement.String() {
			showParserErrors(p, t)
			showPrefixParserError(p, t)
			t.Errorf("statement expected is %s and got %s", expectedStatements[index], statement.String())
		}
	}

}

func showParserErrors(p *Parser, t *testing.T) {
	for _, errorText := range p.errors {
		t.Logf("%s\n", errorText)
	}
}

func showPrefixParserError(p *Parser, t *testing.T) {
	for _, errorText := range p.prefixErrors {
		t.Logf("%s\n", errorText)
	}
}

func TestBlockStatement(t *testing.T) {
	input := `{
		5 + x;
		let x = y != true;
		return x;
	}
   	`
	l := lexer.New(input)
	p := New(l)
	statementsExpected := []string{"{(5 + x);let x = (y != true);return x;}"}
	program := p.ParseProgram()
	if len(program.Statements) != len(statementsExpected) {
		t.Error("program should has 1 statements")
		showPrefixParserError(p, t)
		showParserErrors(p, t)
		return
	}

	statement, ok := program.Statements[0].(*ast.BlockStatement)
	if !ok {
		t.Error("Block Statement is Expected")
		showPrefixParserError(p, t)
		showParserErrors(p, t)
		return
	}
	if len(statement.Statements) != 3 {
		t.Error("Three statement is expected in this block statement")
		showPrefixParserError(p, t)
		showParserErrors(p, t)
		return
	}
	_, okExpression := statement.Statements[0].(*ast.ExpressionStatement)
	if !okExpression {
		t.Error("waiting at index 0 Expression Statement")
		showPrefixParserError(p, t)
		showParserErrors(p, t)
		return
	}
	_, okLet := statement.Statements[1].(*ast.LetStatement)
	if !okLet {
		t.Error("waiting at index 1 Let Statement")
		showPrefixParserError(p, t)
		showParserErrors(p, t)
		return
	}
	_, okRet := statement.Statements[2].(*ast.ReturnStatement)
	if !okRet {
		t.Error("waiting at index 2 Ret Statement")
		showPrefixParserError(p, t)
		showParserErrors(p, t)
		return
	}

	if statementsExpected[0] != statement.String() {
		t.Errorf("Expression expected is %s and got %s", statementsExpected[0], statement.String())
		showPrefixParserError(p, t)
		showParserErrors(p, t)
		return
	}

}

func TestExpressionStatement(t *testing.T) {
	input := `
		   5 + x;
		   x != true

   	`
	l := lexer.New(input)
	p := New(l)
	statementsExpected := []string{"(5 + x);", "(x != true);"}
	program := p.ParseProgram()
	if len(program.Statements) != len(statementsExpected) {
		t.Error("program should has 2 statements")
		showPrefixParserError(p, t)
		showParserErrors(p, t)
		return
	}
	for i := 0; i < len(program.Statements); i++ {
		statement, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Error("Expression Statement is Expected")
			showPrefixParserError(p, t)
			showParserErrors(p, t)
			return
		}
		if statementsExpected[i] != statement.String() {
			t.Errorf("Expression expected is %s and got %s", statementsExpected[i], statement.String())
			showPrefixParserError(p, t)
			showParserErrors(p, t)
			return
		}

	}
}

func TestImportStatement(t *testing.T) {
	tests := []struct {
		input           string
		importStatement *ast.ImportStatement
	}{
		{`import "net"`, &ast.ImportStatement{Path: "net", Token: createToken(lexer.IMPORT, "import"), Name: &ast.Identifier{Value: "net"}}},
		{`import "tmp/carlos" as test`, &ast.ImportStatement{Path: "tmp/carlos", Token: createToken(lexer.IMPORT, "import"), Name: &ast.Identifier{Value: "test"}}},
		{`import "multiprocessing/threads/green"`, &ast.ImportStatement{Path: "multiprocessing/threads/green", Token: createToken(lexer.IMPORT, "import"), Name: &ast.Identifier{Value: "green"}}},
	}
	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()
		if len(program.Statements) != 1 {
			t.Fatal("program should has 1 statements")
			return
		}
		importStatement, ok := program.Statements[0].(*ast.ImportStatement)
		if !ok {
			t.Fatal("expected import statement")
			return
		}
		if importStatement.Path != test.importStatement.Path {
			t.Errorf("import expected path %s and got %s", importStatement.Path, test.importStatement.Path)
			return
		}
		if importStatement.Name.Value != test.importStatement.Name.Value {
			t.Errorf("import expected name %s and got %s", importStatement.Name.Value, test.importStatement.Name.Value)
			return
		}
	}
}

func TestReturnStatementLiteralValue(t *testing.T) {
	input := `
		   return 5;
		   return x;
		   return 838383;
   	`
	l := lexer.New(input)
	p := New(l)
	statementsExpected := []ast.ReturnStatement{createReturnValue("5"), createReturnIdentifier("x"), createReturnValue("838383")}
	program := p.ParseProgram()
	if len(program.Statements) != 3 {
		t.Fatal("program should has 3 statements")
		return
	}
	for i := 0; i < len(program.Statements); i++ {
		var ret *ast.ReturnStatement = program.Statements[i].(*ast.ReturnStatement)
		if !assertReturnStatement(ret, &statementsExpected[i]) {
			t.Errorf("return literal let statement %s should be equal to %s", statementsExpected[i].TokenLiteral(), program.Statements[i].TokenLiteral())
		}
	}
}

func assertReturnStatement(ret, retExpected *ast.ReturnStatement) bool {
	isEqual := isTokenEqual(ret.Token, retExpected.Token)
	return isEqual && assertExpression(ret.Value, retExpected.Value)

}

func assertExpression(exp, expExpected ast.Expression) bool {

	if expectedIntExpression(exp, expExpected) {
		ex1 := exp.(*ast.IntegerLiteral)
		ex2 := expExpected.(*ast.IntegerLiteral)
		return ex1.Value == ex2.Value && isTokenEqual(ex1.Token, ex2.Token)
	}
	if expectedIdentifierExpression(exp, expExpected) {
		ex1 := exp.(*ast.Identifier)
		ex2 := expExpected.(*ast.Identifier)
		return ex1.Value == ex2.Value && isTokenEqual(ex1.Token, ex2.Token)
	}
	return false

}

func expectedIntExpression(exp, expExpected ast.Expression) bool {
	return isIntExpression(exp) && isIntExpression(expExpected)
}

func expectedIdentifierExpression(exp, expExpected ast.Expression) bool {
	return isIdentifierExpression(exp) && isIdentifierExpression(expExpected)
}

func isIntExpression(exp ast.Expression) bool {
	_, ok := exp.(*ast.IntegerLiteral)
	return ok
}

func isIdentifierExpression(exp ast.Expression) bool {
	_, ok := exp.(*ast.Identifier)
	return ok
}

func createReturnIdentifier(name string) ast.ReturnStatement {
	id := &ast.Identifier{Token: createToken(lexer.IDENT, name), Value: name}
	return ast.ReturnStatement{Token: createToken(lexer.RETURN, "return"), Value: id}
}

func createReturnValue(value string) ast.ReturnStatement {
	val, _ := strconv.ParseInt(value, 10, 0)
	id := &ast.IntegerLiteral{Token: createToken(lexer.INT, value), Value: val}
	return ast.ReturnStatement{Token: createToken(lexer.RETURN, "return"), Value: id}

}

func assertLiteralLetStatement(l, lExpected *ast.LetStatement) bool {
	isEqual := isTokenEqual(l.Token, lExpected.Token)
	isEqual = isEqual && isTokenEqual(l.Name.Token, lExpected.Name.Token)
	isEqual = isEqual && l.Name.Value == lExpected.Name.Value
	isEqual = isEqual && assertIntegerExpression(l.Value, lExpected.Value)
	return isEqual
}

func assertIntegerExpression(intExpression, intExpressionExpected ast.Expression) bool {
	var value *ast.IntegerLiteral = intExpression.(*ast.IntegerLiteral)
	var valueExpected *ast.IntegerLiteral = intExpressionExpected.(*ast.IntegerLiteral)
	isEqual := isTokenEqual(value.Token, valueExpected.Token)
	return isEqual && value.Value == valueExpected.Value
}

func isTokenEqual(token, tokenExpected lexer.Token) bool {
	isEqual := token.Literal == tokenExpected.Literal
	return isEqual && token.Type == tokenExpected.Type
}

func createLiteralLetStatement(name, value string) ast.LetStatement {
	id := &ast.Identifier{Token: createToken(lexer.IDENT, name), Value: name}
	val, _ := strconv.ParseInt(value, 10, 0)
	integerLiteral := ast.IntegerLiteral{Token: createToken(lexer.INT, value), Value: val}
	return ast.LetStatement{Token: createToken(lexer.LET, "let"), Name: id, Value: &integerLiteral}
}

func createToken(typeToken lexer.TokenType, value string) lexer.Token {
	return lexer.Token{Type: typeToken, Literal: value}
}

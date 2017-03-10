package evaluator

import (
	"testing"
	"rootlang/lexer"
	"rootlang/parser"
	"rootlang/object"
	"rootlang/builtin"
	"io/ioutil"
	"bytes"
	"os"
)

func TestIntegerEvaluator(t *testing.T) {
	input := `5`
	l := lexer.New(input)
	programParser := parser.New(l)
	program := programParser.ParseProgram()
	returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
	returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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

func TestBuiltFunctionExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`let a=len("Carlos");a;`, 6},
	}
	for _, test := range tests {
		l := lexer.New(test.input)
		programParser := parser.New(l)
		program := programParser.ParseProgram()
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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
		returnValue := Eval(program, object.NewEnvironment(), builtin.New())
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

func TestImportModule(t *testing.T) {
	moduleContent := "let test = (x,y)=>{return x+y;};let test1 = test(1);let result = test1(2);"
	modulePath := "/tmp/testImport.rl"
	createModule(moduleContent, modulePath)
	input := `import "testImport" as test;`
	l := lexer.New(input)
	programParser := parser.New(l)
	program := programParser.ParseProgram()
	builtinSymbols := builtin.New()
	builtinSymbols.RegisterPath("/tmp/")
	environment := object.NewEnvironment()
	returnValue := Eval(program, environment, builtinSymbols)
	if returnValue != nil && returnValue.Type() == object.ERROR_OBJ {
		t.Errorf("fail eval proram %s -> %s", input, returnValue.Inspect())
		return
	}
	value, ok := environment.GetVar("test")
	if !ok {
		t.Error("should has a entry for test module and got nothing")
		return
	}
	if value.Type() != object.MODULE_OBJ {
		t.Errorf("should has a module object and got %s", value.Type())
		return
	}
	module := value.(*object.Module)

	if module.Name != "test" {
		t.Errorf("should has a module object name and got %s", module.Name)
		return
	}
	testFunction, ok := module.Env.GetVar("test")
	if !ok {
		t.Error("test function is expected on module test")
		return
	}
	if testFunction.Type() != object.FUNCTION_OBJ {
		t.Errorf("inside module test function test is expected and got %s", testFunction.Type())
		return
	}

	test1Function, ok := module.Env.GetVar("test1")
	if !ok {
		t.Error("test1 function is expected on module test")
		return
	}
	if test1Function.Type() != object.FUNCTION_OBJ {
		t.Errorf("inside module test function test1 is expected and got %s", test1Function.Type())
		return
	}
	result, ok := module.Env.GetVar("result")
	if !ok {
		t.Error("result integer is expected on module test")
		return
	}
	if result.Type() != object.INTEGER_OBJ {
		t.Errorf("inside module test  result integer is expected and got %s", test1Function.Type())
		return
	}
	integerResult := result.(*object.Integer)
	if integerResult.Value != 3 {
		t.Errorf("inside de module the result variable should be has value of 3 and %d", integerResult.Value)
		return
	}
	os.Remove(modulePath)
}

func TestModuleCall(t *testing.T) {
	moduleContent := "let y = 5; let addToX = x=>{return x+y;};"
	modulePath := "/tmp/testImport.rl"
	createModule(moduleContent, modulePath)
	expectedValue := int64(15)
	input := `import "testImport" as test;
	let y = "carlos test"; let x = test::addToX(10);x`
	l := lexer.New(input)
	programParser := parser.New(l)
	program := programParser.ParseProgram()
	builtinSymbols := builtin.New()
	builtinSymbols.RegisterPath("/tmp/")
	environment := object.NewEnvironment()
	returnValue := Eval(program, environment, builtinSymbols)
	if returnValue == nil || returnValue.Type() != object.INTEGER_OBJ {
		t.Error("expected integer")
		return
	}
	integerValue := returnValue.(*object.Integer)
	if integerValue.Value != expectedValue {
		t.Errorf("expected integer with value %d and got %d", expectedValue, integerValue.Value)
		return
	}

	os.Remove(modulePath)
}

func TestModuleCall2(t *testing.T) {
	moduleContent := "let y = 5; let addToX = x=>{return x+y;};"
	modulePath := "/tmp/testImport.rl"
	createModule(moduleContent, modulePath)
	expectedValue := int64(15)
	input := `import "testImport" as test;
	let y = test::addToX; let x = y(10);x`
	l := lexer.New(input)
	programParser := parser.New(l)
	program := programParser.ParseProgram()
	builtinSymbols := builtin.New()
	builtinSymbols.RegisterPath("/tmp/")
	environment := object.NewEnvironment()
	returnValue := Eval(program, environment, builtinSymbols)
	if returnValue == nil || returnValue.Type() != object.INTEGER_OBJ {
		t.Error("expected integer")
		return
	}
	integerValue := returnValue.(*object.Integer)
	if integerValue.Value != expectedValue {
		t.Errorf("expected integer with value %d and got %d", expectedValue, integerValue.Value)
		return
	}

	os.Remove(modulePath)
}

func TestModuleCall3(t *testing.T) {
	moduleContent := "let y = 5;"
	modulePath := "/tmp/testImport.rl"
	createModule(moduleContent, modulePath)
	expectedValue := int64(15)
	input := `import "testImport" as test;
	let y = 20; let x = test::(a => { return a + y;})(10);x`
	l := lexer.New(input)
	programParser := parser.New(l)
	program := programParser.ParseProgram()
	builtinSymbols := builtin.New()
	builtinSymbols.RegisterPath("/tmp/")
	environment := object.NewEnvironment()
	returnValue := Eval(program, environment, builtinSymbols)
	if returnValue == nil || returnValue.Type() != object.INTEGER_OBJ {
		t.Error("expected integer")
		return
	}
	integerValue := returnValue.(*object.Integer)
	if integerValue.Value != expectedValue {
		t.Errorf("expected integer with value %d and got %d", expectedValue, integerValue.Value)
		return
	}

	os.Remove(modulePath)
}

func TestMainModuleModule(t *testing.T) {

	input := `import "/net";
import "/bytes";
let main = () => {
        let get_clients_to_write = (server, client)=> {
                let is_current_client = client_to_write => {
                        let client_id = net::get_client_id(client);
                        let client_write_id = net::get_client_id(client_to_write);
                        return client_id != client_write_id;
                };
                let clients = net::get_clients(server);
                return filter(is_current_client, clients);
        };
        let on_client_connect = (server, client) => {
                print("new client arrive --> ", client);
                let clients = get_clients_to_write(server,client);
                let client_id = net::get_client_id(client);
                let message_to_send = bytes::create_writer("new-client :) ",client_id);
                let clients_write = map(client_to_write => { return net::write_to_client(client_to_write, message_to_send);}, clients);
                return clients_write;

        };
        let on_client_write = (server,client, message)=> {
                let message_text = bytes::read_string(message);
                print(message_text);
                let clients = get_clients_to_write(server, client);
                print("testing 1");
                let client_id =  net::get_client_id(client);
                print(clients, client_id);
                let message_to_send = bytes::create_writer(client_id, ": ", message_text);
                print(message_to_send);
                let clients_write = map(client_to_write => { return net::write_to_client(client_to_write, message_to_send);}, clients);
                print(clients_write);
                return clients_write;
        };
        print("hola mundo");
        net::listen(3000,on_client_connect, on_client_write);
        return 0;
};`
	l := lexer.New(input)
	programParser := parser.New(l)
	program := programParser.ParseProgram()
	env := object.NewEnvironment()
	builtinSymbols := builtin.New()
	returnValue := Eval(program, env, builtinSymbols)
	if len(programParser.GetErrors()) != 0 {
		t.Error(programParser.GetErrors()[0])
		return
	}
	if returnValue != nil && returnValue.Type() == object.ERROR_OBJ {
		t.Error("Error Evaluation main module")
		return
	}
	mainFunction, hasMain := env.GetVar("main")
	if !hasMain {
		t.Error("Module Has No Main Function")
		return
	}
	if mainFunction.Type() != object.FUNCTION_OBJ {
		t.Error("Main is not a function")
		return
	}
	CallMainFunction(mainFunction.(*object.Function), builtinSymbols)

}

func createModule(module, path string) {
	buffer := bytes.NewBufferString(module)
	ioutil.WriteFile(path, buffer.Bytes(), 0644)
}

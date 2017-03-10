package main

import (
	"rootlang/parser"
	"io"
	"bufio"
	"fmt"
	"rootlang/lexer"
	"rootlang/evaluator"
	"rootlang/object"
	"rootlang/builtin"
	"os"
)

var PROMPT string = "rootlang>"

func main() {
	if len(os.Args) == 1 {
		start(os.Stdin, os.Stdout)
	} else {
		modulePath := os.Args[1]
		builtinSymbols := builtin.New()
		env, err := evaluator.ReadPrincipalModule(modulePath)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error On Module %s  --> %s\n", err.Error(), modulePath))
			return
		}
		mainFunction, hasMain := env.GetVar("main")
		if !hasMain {
			os.Stderr.WriteString("Module Has No Main Function\n")
			return
		}
		if mainFunction.Type() != object.FUNCTION_OBJ {
			os.Stderr.WriteString("Main is not a function\n")
			return
		}
		returnValue := evaluator.CallMainFunction(mainFunction.(*object.Function), builtinSymbols)
		if returnValue.Type() == object.ERROR_OBJ {
			os.Stderr.WriteString(fmt.Sprintf("%s\n", returnValue.Inspect()))
			os.Exit(-1)
		}
	}

}

func start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	environment := object.NewEnvironment()
	builtinSymbols := builtin.New()
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.GetErrors()) != 0 {
			printParserErrors(out, p.GetErrors())
			continue
		}
		evaluated := evaluator.Eval(program, environment, builtinSymbols)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, error := range errors {
		io.WriteString(out, fmt.Sprintf("%s\n", error))
	}
}

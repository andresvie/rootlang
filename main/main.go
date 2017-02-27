package main

import (
  "rootlang/parser"
  "io"
  "bufio"
  "fmt"
  "rootlang/lexer"
  "rootlang/evaluator"
  "rootlang/object"
  "os"
)

var PROMPT string = "rootlang>"

func main() {
  start(os.Stdin, os.Stdout)
}

func start(in io.Reader, out io.Writer) {
  scanner := bufio.NewScanner(in)
  environment := object.NewEnvironment()
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
    evaluated := evaluator.Eval(program, environment)
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

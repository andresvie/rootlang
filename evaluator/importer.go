package evaluator

import (
	"rootlang/object"
	"rootlang/ast"
	"rootlang/builtin"
	"os"
	"path"
	"fmt"
	"io/ioutil"
	"rootlang/lexer"
	"rootlang/parser"
	"strings"
)

func importModule(importStatement *ast.ImportStatement, builtinSymbols *builtin.Builtin) object.Object {
	modulePath := getModulePath(importStatement.Path, builtinSymbols.GetPaths())

	if modulePath == "" {
		return &object.ErrorObject{Error: fmt.Sprintf("not module %s found", importStatement.Path)}
	}
	moduleContent, err := readModuleFile(modulePath)
	if err != nil {
		return &object.ErrorObject{Error: fmt.Sprintf("the module %s can not be read", importStatement.Path)}
	}
	newEnvironment := object.NewEnvironment()
	l := lexer.New(moduleContent)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.GetErrors()) != 0 {
		return &object.ErrorObject{Error: fmt.Sprintf("error parsing the module %s %s", importStatement.Path, strings.Join(p.GetErrors(), "\n"))}
	}
	evalResult := Eval(program, newEnvironment, builtinSymbols)
	if evalResult != nil && evalResult.Type() == object.ERROR_OBJ {
		return evalResult
	}
	return &object.Module{Path: importStatement.Path, Name: importStatement.Name.Value, Env: newEnvironment}

}

func readModuleFile(modulePath string) (string, error) {
	content, err := ioutil.ReadFile(modulePath)
	if err != nil {
		return "", err
	}
	return string(content), err
}

func getModulePath(module string, paths []string) string {
	modulePathWithExtension := fmt.Sprintf("%s.rl", module)
	for _, modulePath := range paths {
		if existPath(modulePath, modulePathWithExtension) {
			return path.Join(modulePath, modulePathWithExtension)
		}
	}
	return ""
}

func existPath(module, modulePath string) bool {
	fullPath := path.Join(module, modulePath)
	_, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

package builtin

import (
	"io"
	"fmt"
	"crypto/rand"
	"rootlang/object"
	"rootlang/ast"
)

const (
	LEN    = "len"
	LIST   = "list"
	APPEND = "append"
	MAP    = "map"
	FILTER = "filter"
	ZIP    = "zip"
	REDUCE = "reduce"
	PRINT  = "print"
	NET    = "net"
	BYTES = "bytes"
)

type function func(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object

type BuiltinFunction struct {
	Name     string
	Function function
	Params   []object.Object
}

func (builtinFunction *BuiltinFunction) Type() object.ObjectType {
	return object.BUILTIN_FUNCTION_OBJ
}

func (builtinFunction *BuiltinFunction) Inspect() string {
	return builtinFunction.Name
}

type Builtin struct {
	symbols map[string]object.Object
	modules map[string]*object.Module
	paths   []string
}

func New() *Builtin {
	symbols := registerSymbols()
	paths := make([]string, 0)
	return &Builtin{symbols: symbols, paths: paths}
}

func registerSymbols() map[string]object.Object {
	symbols := make(map[string]object.Object)
	symbols[LEN] = getBuiltinFunction(_len, LEN)
	symbols[LIST] = getBuiltinFunction(_list, LIST)
	symbols[APPEND] = getBuiltinFunction(_append, APPEND)
	symbols[MAP] = getBuiltinFunction(_map, MAP)
	symbols[FILTER] = getBuiltinFunction(_filter, FILTER)
	symbols[ZIP] = getBuiltinFunction(_zip, ZIP)
	symbols[REDUCE] = getBuiltinFunction(_reduce, REDUCE)
	symbols[PRINT] = getBuiltinFunction(_print, PRINT)
	symbols[NET] = buildNetModule()
	symbols[BYTES] = buildBytesModule()
	return symbols
}

func (b *Builtin) RegisterPath(path string) {
	b.paths = append(b.paths, path)
}

func (b *Builtin) GetPaths() []string {
	return b.paths
}

func (b *Builtin) GetObject(name string) (object.Object, bool) {
	value, ok := b.symbols[name]
	return value, ok
}

func getBuiltinFunction(f function, symbol string) *BuiltinFunction {
	return &BuiltinFunction{Name: symbol, Function: f, Params: make([]object.Object, 0)}
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

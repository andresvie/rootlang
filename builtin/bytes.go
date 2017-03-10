package builtin

import "rootlang/ast"
import (
	"bytes"
	"rootlang/object"
	"fmt"
)

const (
	BYTE_OBJ      = "byte"
	READER_BUFFER = "reader_buffer"
	WRITER_BUFFER = "writer_buffer"
)

type ReaderBuffer interface {
	readString(object.Object) object.Object
}

type WriterBuffer interface {
	writeString(object.Object) object.Object
}

type WriterBufferObject struct {
	data *bytes.Buffer
}

func (writer *WriterBufferObject) writeString(obj object.Object) object.Object {
	writer.data.WriteString(obj.Inspect())
	return writer
}

func (writer *WriterBufferObject) Type() object.ObjectType {
	return WRITER_BUFFER
}

func (writer *WriterBufferObject) Inspect() string {
	return fmt.Sprintf("%d", writer.data.Len())
}

type ReaderBufferObject struct {
	data *bytes.Buffer
}

func (reader *ReaderBufferObject) Type() object.ObjectType {
	return READER_BUFFER
}

func (reader *ReaderBufferObject) Inspect() string {
	return fmt.Sprintf("%d", reader.data.Len())
}

func (reader *ReaderBufferObject) readString() object.Object {
	text, err := reader.data.ReadString(0)
	if err != nil && err.Error() != "EOF" {
		return &object.ErrorObject{Error: err.Error()}
	}
	return &object.String{Value: text}

}

func createReaderBufferFromString(text string) *ReaderBufferObject {
	return &ReaderBufferObject{bytes.NewBufferString(text)}
}

func buildBytesModule() *object.Module {
	env := object.NewEnvironment()
	env.SetVar("create_writer", getBuiltinFunction(_create_writer, "create_writer"))
	env.SetVar("read_string", getBuiltinFunction(_read_string, "read_string"))

	return &object.Module{Env: env, Name: "bytes", Path: "/bytes"}

}
func _read_string(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	if len(params) != 1 || params[0].Type() != READER_BUFFER {
		return &object.ErrorObject{Error: "expected reader buffer"}
	}
	reader := params[0].(*ReaderBufferObject)
	return reader.readString()
}
func _create_writer(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	buffer := bytes.NewBufferString("")
	for _, value := range params {
		switch valueType := value.(type) {
		case *object.String:
			buffer.WriteString(valueType.Value)
		case *object.Integer:
			buffer.WriteString(valueType.Inspect())
		default:
			return &object.ErrorObject{Error: fmt.Sprintf("can not writer to buffer type %s", value.Type())}
		}
	}
	return &WriterBufferObject{data: buffer}
}

package builtin

import "rootlang/object"
import "rootlang/ast"
import "net"
import (
	"fmt"
	"bufio"
)

const (
	SERVER_OBJ = "SERVER"
	CLIENT_OBJ = "CLIENT"
)

type Client struct {
	id  string
	con net.Conn
}

func (client *Client) Type() object.ObjectType {
	return CLIENT_OBJ
}

func (client *Client) Inspect() string {
	return fmt.Sprintf("%s", client.id)
}

type Server struct {
	listener net.Listener
	port     int64
	clients  map[string]*Client
}

func (server *Server) Type() object.ObjectType {
	return SERVER_OBJ
}

func (server *Server) Inspect() string {
	return fmt.Sprintf("tcp::%d", server.port)
}

func buildNetModule() *object.Module {
	env := object.NewEnvironment()
	env.SetVar("listen", getBuiltinFunction(_listen, "listen"))
	env.SetVar("get_client_id", getBuiltinFunction(_get_client_id, "get_client_id"))
	env.SetVar("get_clients", getBuiltinFunction(_get_clients, "get_clients"))
	env.SetVar("write_to_client", getBuiltinFunction(_write_to_client, "write_to_client"))

	return &object.Module{Env: env, Name: "net", Path: "/net"}

}

func _write_to_client(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	if len(params) != 2 {
		return &object.ErrorObject{Error: "expected 2 params"}
	}
	if params[0].Type() != CLIENT_OBJ || params[1].Type() != WRITER_BUFFER {
		return &object.ErrorObject{Error: "expected params with type client and writer_buffer"}
	}
	client := params[0].(*Client)
	writer := params[1].(*WriterBufferObject)
	numberOfBytes, err := client.con.Write(writer.data.Bytes())
	if err != nil {
		return &object.ErrorObject{Error: err.Error()}
	}
	return &object.Integer{Value: int64(numberOfBytes)}

}

func _get_client_id(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	if len(params) != 1 || params[0].Type() != CLIENT_OBJ {
		return &object.ErrorObject{Error: "expected client object"}
	}
	return &object.String{Value: params[0].Inspect()}
}

func _get_clients(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	if len(params) != 1 || params[0].Type() != SERVER_OBJ {
		return &object.ErrorObject{Error: "expected server object"}
	}
	server := params[0].(*Server)
	values := make([]object.Object, 0)
	for _, value := range server.clients {
		values = append(values, value)
	}
	return &object.List{Elements: values}
}

func _listen(env *object.Environment, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object, params ...object.Object) object.Object {
	if len(params) != 3 {
		return &object.ErrorObject{Error: fmt.Sprintf("net::listen expected 3 params and got %s", len(params))}
	}
	if params[0].Type() != object.INTEGER_OBJ || params[1].Type() != object.FUNCTION_OBJ || params[2].Type() != object.FUNCTION_OBJ {
		return &object.ErrorObject{Error: "the signature expected is net::listen(port, (server, new-client) => {}, (client, data) => {});"}
	}
	port := params[0].(*object.Integer).Value
	onClientConnect := params[1].(*object.Function)
	onClientWrite := params[2].(*object.Function)
	serverConnection, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return &object.ErrorObject{Error: err.Error() }
	}
	clients := make(map[string]*Client)
	server := &Server{listener: serverConnection, clients: clients, port: port}
	for {
		conn, err := serverConnection.Accept()
		if err != nil {
			return &object.ErrorObject{Error: err.Error() }
		}
		client := createClient(conn)
		params := []object.Object{server, client}
		server.clients[client.id] = client
		returnValue := applyArgumentsToFunctionAndCall(onClientConnect, params, b, eval)
		if returnValue != nil && isErrorObject(returnValue) {
			return returnValue
		}
		go handleClient(server, client, onClientWrite, b, eval)
	}

}

func handleClient(server *Server, client *Client, onClientWrite *object.Function, b *Builtin, eval func(node ast.Node, environment *object.Environment, builtinSymbols *Builtin) object.Object) {
	for {

		message, _ := bufio.NewReader(client.con).ReadString('\n')
		buffer := createReaderBufferFromString(message)
		applyArgumentsToFunctionAndCall(onClientWrite, []object.Object{server, client, buffer}, b, eval)
	}
}

func createClient(con net.Conn) *Client {
	id, _ := newUUID()
	return &Client{id: id, con: con}
}

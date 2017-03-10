import "/net";
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
		let clients = get_clients_to_write(server, client);
		let client_id =  net::get_client_id(client);
		let message_to_send = bytes::create_writer(client_id, ": ", message_text);
		let clients_write = map(client_to_write => { return net::write_to_client(client_to_write, message_to_send);}, clients);
		return clients_write;	
	};
	print("server listen on port 3000");
	net::listen(3000,on_client_connect, on_client_write);
	return 0;
};

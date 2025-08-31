import socket
import logging

class Socket:

    def __init__(self, sock=None):
        # Initialize socket
        if sock == None:
        	self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        else:
        	self.socket = sock
        
    def bind_and_listen(self, port, listen_backlog):
        self.socket.bind(('', port))
        self.socket.listen(listen_backlog)
        
    def accept(self):
    	client_socket, addr = self.socket.accept()
    	return Socket(client_socket), addr
    	
    #def recv(self, length):
    	#return self.socket.recv(length).rstrip().decode('utf-8')
    	
    def getpeername(self):
    	return self.socket.getpeername()
    
    def send(self, msg):
    	self.socket.sendall("{}\n".format(msg).encode('utf-8'))
    	
    def close(self):
    	self.socket.close()
    	
    def recv_length(self):
        LENGTH = 4
        data = b''
        while len(data) < LENGTH:
            chunk = self.socket.recv(LENGTH - len(data))
            if not chunk:
                raise ConnectionError("Connection closed unexpectedly")
            data += chunk

        length = int.from_bytes(data, 'big')
        return length
	    
    def recv_msg(self, length):
        data = b''
        while len(data) < length:

            chunk = self.socket.recv(length - len(data))
            if not chunk:
                raise ConnectionError("Connection closed unexpectedly")
            data += chunk
        msg = data.decode('utf-8')
        return msg

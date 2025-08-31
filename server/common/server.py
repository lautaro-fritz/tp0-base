import logging
import signal

from .my_socket import Socket
from .utils import Bet, store_bets

class Server: 

    def _graceful_exit_handler(self, signum, frame):
        logging.info(f"Signal {signum} received, shutting down gracefully.")
        self._is_running = False
        try:
            self._server_socket.close()
            logging.info(f'action: server_socket_close | result: success')
        except Exception as e:
            logging.error(f"Error closing server socket: {e}")

    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = Socket()
        self._server_socket.bind_and_listen(port, listen_backlog)
        signal.signal(signal.SIGTERM, self._graceful_exit_handler)
        self._is_running = True

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while self._is_running:
            client_sock = self.__accept_new_connection()
            if client_sock:
                self.__handle_client_connection(client_sock)
        
        logging.info(f"Server shutdown complete.")

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # primero leo la longitud exacta del mensaje del cliente
            length = client_sock.recv_length()
            addr = client_sock.getpeername()
            logging.info(f'action: receive_length | result: success | ip: {addr[0]} | msg: {length}')
            
            # luego leo el mensaje con la longitud previamente leida
            msg = client_sock.recv_msg(length)
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            
            # tendria que mover esto a una clase BetManager o algo asi
            values = msg.split('/')
            bet = Bet(values[0], values[1], values[2], values[3], values[4], values[5])
            store_bets([bet])
            logging.info(f'action: apuesta_almacenada | result: success | dni: {values[3]} | numero: {values[5]}.')
            # TODO: Modify the send to avoid short-writes
            
            confirmation_msg = "OK"
            client_sock.send(confirmation_msg)
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        try:
            client_socket, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return client_socket
        except OSError as e:
            logging.debug(f'Server socket closed or error accepting connection: {e}')
            return None

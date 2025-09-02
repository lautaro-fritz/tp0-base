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
            length = client_sock.recv_length()
            addr = client_sock.getpeername()
            logging.info(f'action: receive_length | result: success | ip: {addr[0]} | msg: {length}')
            
            msg = client_sock.recv_msg(length)
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            
            parts = msg.split('#')
            
            if len(parts) < 2:
                logging.error(f'action: apuesta_recibida | result: fail | cantidad: 0 | reason: no_bets_received')
                client_sock.send("ERROR")
                return
            
            client_id = parts[0]
            bets_raw = parts[1:]
            
            bets = []
            for bet_str in bets_raw:
                fields = bet_str.split('/')
                if len(fields) != 5:
                    logging.error(f'action: apuesta_recibida | result: fail | cantidad: {len(bets_raw)} | reason: malformed_bet')
                    client_sock.send("ERROR")
                    return
                
                bet = Bet(client_id, fields[0], fields[1], fields[2], fields[3], fields[4])
                bets.append(bet)
            
            try:
                store_bets(bets)
            except Exception as e:
                logging.error(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)} | error: {e}')
                client_sock.send("ERROR")
                return
            
            logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
            client_sock.send("OK")
        
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
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

import logging
import signal
import threading

from .my_socket import Socket
from .utils import Bet, store_bets, load_bets, has_won


class Server:

    def _graceful_exit_handler(self, signum, frame):
        logging.info(f"Signal {signum} received, shutting down gracefully.")
        self._is_running = False
        try:
            self._server_socket.close()
            logging.info('action: server_socket_close | result: success')
        except Exception as e:
            logging.error(f"Error closing server socket: {e}")

        # Wait for all threads to finish
        with self._threads_lock:
            for thread in self._threads:
                thread.join()
        logging.info("All client threads have been joined. Server shutdown complete.")

    def __init__(self, port, listen_backlog, clients_amount):
        self._server_socket = Socket()
        self._server_socket.bind_and_listen(port, listen_backlog)

        signal.signal(signal.SIGTERM, self._graceful_exit_handler)
        self._is_running = True

        self.registered_agencies = [False] * clients_amount
        self.winners = []

        # Concurrency management
        self._threads = []
        self._threads_lock = threading.Lock()
        self._agencies_lock = threading.Lock()
        self._bets_lock = threading.Lock()
        self._winners_lock = threading.Lock()
        self._winner_selected = False
        self._winner_lock = threading.Lock()

    def run(self):
        while self._is_running:
            client_sock = self.__accept_new_connection()
            if client_sock:
                thread = threading.Thread(
                    target=self.__handle_client_connection,
                    args=(client_sock,)
                )
                thread.start()
                with self._threads_lock:
                    self._threads.append(thread)

    def __handle_client_connection(self, client_sock):
        try:
            length = client_sock.recv_length()
            addr = client_sock.getpeername()
            logging.info(f'action: receive_length | result: success | ip: {addr[0]} | msg: {length}')

            msg = client_sock.recv_msg(length)
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')

            parts = msg.split('#')
            if len(parts) < 2:
                logging.error(f'action: mensaje_recibido | result: fail | reason: no op code')
                client_sock.send("ERROR")
                return

            if len(parts) > 2:
                client_id = parts[0]
                bets_raw = parts[2:]
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
                    with self._bets_lock:
                        store_bets(bets)
                except Exception as e:
                    logging.error(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)} | error: {e}')
                    client_sock.send("ERROR")
                    return

                logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
                client_sock.send("OK")

            elif len(parts) == 2:
                client_id = int(parts[0])

                if parts[1] == "D":
                    logging.info(f'action: done_recibido | result: success')
                    with self._agencies_lock:
                        self.registered_agencies[client_id - 1] = True
                    client_sock.send("OK")

                elif parts[1] == "W":
                    logging.info(f'action: get_winners_recibido | result: success')
                    with self._agencies_lock:
                        if self.registered_agencies.count(False) >= 1:
                            client_sock.send("ERROR")
                            return

                    with self._winner_lock:
                        if not self._winner_selected:
                            with self._bets_lock:
                                all_bets = load_bets()
                            for bet in all_bets:
                                if has_won(bet):
                                    with self._winners_lock:
                                        self.winners.append(bet)
                            self._winner_selected = True
                            logging.info("action: sorteo | result: success")

                    # Send winners
                    with self._winners_lock:
                        agency_winners = [w for w in self.winners if w.agency == client_id]
                        documents_str = "#".join(w.document for w in agency_winners)
                        response = f"W#{documents_str}"
                        
                    client_sock.send(response)
                    return

        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        logging.info('action: accept_connections | result: in_progress')
        try:
            client_socket, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return client_socket
        except OSError as e:
            logging.debug(f'Server socket closed or error accepting connection: {e}')
            return None


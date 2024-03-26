import socket
import logging
import signal
import sys

from .utils import receive_bet, store_bets


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        
        # Register signal handlers
        signal.signal(signal.SIGTERM, self.handle_sigterm)
        
        # Initialize client list 
        self.clients = []

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        self.running = True
        
        while self.running:
            client_sock = self.__accept_new_connection()
            self.clients.append(client_sock)
            self.__handle_client_connection(client_sock)
        
        self._server_socket.close()

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            while True:
                bets = receive_bet(client_sock)
                if len(bets) == 0:
                    logging.info("action: receive_bets | result: success")
                    break
                store_bets(bets)
                
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
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def handle_sigterm(self, signum, frame):
        """
        Handle SIGTERM signal

        Function that is called when SIGTERM signal is received.
        It closes all the open connections and sets the running flag to false
        """
        
        logging.info('action: handle_sigterm | result: in_progress')
        self.running = False
        for client in self.clients:
            client.close()
        logging.info(f'action: handle_sigterm | result: success')
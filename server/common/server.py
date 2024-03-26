import socket
import logging
import signal
import sys

from .utils import get_all_winners, get_winners, receive_bet, store_bets


class Server:
    def __init__(self, port, listen_backlog, clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.all_clients_ready = False
        self.finished_clients = 0
        self.clients_number = clients
        
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

    def receive_action(client_sock):
        """
        Receive an action from a client socket.

        Args:
            client_sock (socket.socket): The socket from which to receive the action.

        Returns:
            a tuple containing the id and the action
        """
        try : 
            message_length = int.from_bytes(client_sock.recv(4), byteorder='big')
        except ValueError as e:
            return None, None
        # Luego, leer el mensaje completo
        msg = client_sock.recv(message_length).decode('utf-8').strip()
        return msg.split(',')

    def send_winners(self, client_sock: socket.socket, id: int):
        """
        Send the winners to the given client socket.

        Args:
            client_sock (socket.socket): The socket to send the winners to.
            id (int): The id of the game.

        Returns:
            None

        Raises:
            ValueError: If the server is not ready to send the winners.
        """
        if self.all_clients_ready:
            message = get_winners(id)
            logging.info(f"action: command_winners | result: success | id: {id}")
            msg_length_bytes = len(message).to_bytes(4, byteorder="big")
            client_sock.sendall(msg_length_bytes + message.encode())
        else:
            message = "NOT_WINNERS_YET"
            msg_length_bytes = len(message).to_bytes(4, byteorder="big")
            client_sock.sendall(msg_length_bytes + message.encode())
            
    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            while self.running:
                id, action = Server.receive_action(client_sock)
            
                if action != 'BET':
                    logging.info(f'action: command_receive | result: in_progress | id: {id} | action: {action}')
                if action == 'CLOSE_CONNECTION':
                    break                
                elif action == 'BET':
                    store_bets(receive_bet(client_sock))
                elif action == 'FINISH_BET':
                    logging.info("action: storing_bets | result: success")
                    self.finished_clients += 1
                    if self.finished_clients == self.clients_number:
                        self.all_clients_ready = True
                        logging.info(f'action: sorteo | result: success | Cant : {get_all_winners()} ')
                elif action == 'WINNERS':
                    self.send_winners(client_sock, id)
                else:
                    logging.error("action: command_receive | result: fail | error: {e}")
                
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
            logging.info(f'action: close_client | result: success')
        logging.info(f'action: handle_sigterm | result: success')
import logging
import os
import subprocess
from configparser import ConfigParser

#running docker run --rm -it --network=tp0_testing_net health_test

def send_message_to_server(message,server, port, logger = None):
    process = subprocess.Popen(["nc", server, port], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    process.stdin.write(message)
    process.stdin.close()
    response = process.stdout.readline().strip()
    if logger is not None:
        if response == message:
            logger.info("Received message from server successfully")
        else:
            logger.error("Error sending or receiving message from server")
    
    
if __name__ == "__main__":
    logging.basicConfig(
        format='%(asctime)s %(levelname)-8s %(message)s',
        level= 'INFO',
        datefmt='%Y-%m-%d %H:%M:%S',
    )
    config = ConfigParser(os.environ)
    config.read("config.ini")
    logger = logging.getLogger('logger')
    try : 
        server = str(os.getenv('SERVER_IP', config["DEFAULT"]["SERVER_IP"]))
        port = str(os.getenv('SERVER_PORT', config["DEFAULT"]["SERVER_PORT"]))
        send_message_to_server(b"hello server", server, port, logger)
    except Exception as e:
        logger.error(e)

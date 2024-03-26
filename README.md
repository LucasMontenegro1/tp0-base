# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un ejemplo de cliente-servidor el cual corre en containers con la ayuda de [docker-compose](https://docs.docker.com/compose/). El mismo es un ejemplo práctico brindado por la cátedra para que los alumnos tengan un esqueleto básico de cómo armar un proyecto de cero en donde todas las dependencias del mismo se encuentren encapsuladas en containers. El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers.

Por otro lado, se presenta una guía de ejercicios que los alumnos deberán resolver teniendo en cuenta las consideraciones generales descriptas al pie de este archivo.

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que posee encapsulado diferentes comandos utilizados recurrentemente en el proyecto en forma de targets. Los targets se ejecutan mediante la invocación de:

* **make \<target\>**:
Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de _debugging_ y _troubleshooting_.

Los targets disponibles son:
* **docker-compose-up**: Inicializa el ambiente de desarrollo (buildear docker images del servidor y cliente, inicializar la red a utilizar por docker, etc.) y arranca los containers de las aplicaciones que componen el proyecto.
* **docker-compose-down**: Realiza un `docker-compose stop` para detener los containers asociados al compose y luego realiza un `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene.
* **docker-compose-logs**: Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose.
* **docker-image**: Buildea las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para testear nuevos cambios en las imágenes antes de arrancar el proyecto.
* **build**: Compila la aplicación cliente para ejecución en el _host_ en lugar de en docker. La compilación de esta forma es mucho más rápida pero requiere tener el entorno de Golang instalado en la máquina _host_.

### Servidor
El servidor del presente ejemplo es un EchoServer: los mensajes recibidos por el cliente son devueltos inmediatamente. El servidor actual funciona de la siguiente forma:
1. Servidor acepta una nueva conexión.
2. Servidor recibe mensaje del cliente y procede a responder el mismo.
3. Servidor desconecta al cliente.
4. Servidor procede a recibir una conexión nuevamente.

### Cliente
El cliente del presente ejemplo se conecta reiteradas veces al servidor y envía mensajes de la siguiente forma.
1. Cliente se conecta al servidor.
2. Cliente genera mensaje incremental.
recibe mensaje del cliente y procede a responder el mismo.
3. Cliente envía mensaje al servidor y espera mensaje de respuesta.
Servidor desconecta al cliente.
4. Cliente vuelve al paso 2.

Al ejecutar el comando `make docker-compose-up` para comenzar la ejecución del ejemplo y luego el comando `make docker-compose-logs`, se observan los siguientes logs:

```
$ make docker-compose-logs
docker compose -f docker-compose-dev.yaml logs -f
client1  | time="2023-03-17 04:36:59" level=info msg="action: config | result: success | client_id: 1 | server_address: server:12345 | loop_lapse: 20s | loop_period: 5s | log_level: DEBUG"
client1  | time="2023-03-17 04:36:59" level=info msg="action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1\n"
server   | 2023-03-17 04:36:59 DEBUG    action: config | result: success | port: 12345 | listen_backlog: 5 | logging_level: DEBUG
server   | 2023-03-17 04:36:59 INFO     action: accept_connections | result: in_progress
server   | 2023-03-17 04:36:59 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2023-03-17 04:36:59 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2023-03-17 04:36:59 INFO     action: accept_connections | result: in_progress
server   | 2023-03-17 04:37:04 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2023-03-17 04:37:04 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2023-03-17 04:37:04 INFO     action: accept_connections | result: in_progress
client1  | time="2023-03-17 04:37:04" level=info msg="action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2\n"
server   | 2023-03-17 04:37:09 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2023-03-17 04:37:09 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
server   | 2023-03-17 04:37:09 INFO     action: accept_connections | result: in_progress
client1  | time="2023-03-17 04:37:09" level=info msg="action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3\n"
server   | 2023-03-17 04:37:14 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2023-03-17 04:37:14 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
client1  | time="2023-03-17 04:37:14" level=info msg="action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4\n"
server   | 2023-03-17 04:37:14 INFO     action: accept_connections | result: in_progress
client1  | time="2023-03-17 04:37:19" level=info msg="action: timeout_detected | result: success | client_id: 1"
client1  | time="2023-03-17 04:37:19" level=info msg="action: loop_finished | result: success | client_id: 1"
client1 exited with code 0
```

## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Modificar la definición del DockerCompose para agregar un nuevo cliente al proyecto.

#### Resolución 

Para este caso se agrego un cliente nuevo modigicando el archivo **docker-compose-dev.yaml**

```yaml
    client2:
    container_name: client2
    image: client:latest
    entrypoint: /client
    environment:
    - CLI_ID=2
    - CLI_LOG_LEVEL=DEBUG
    networks:
    - testing_net
    depends_on:
    - server
```

Como se puede ver se agrega el container name como client2 y se modifica su ID.


### Ejercicio N°1.1:
Definir un script (en el lenguaje deseado) que permita crear una definición de DockerCompose con una cantidad configurable de clientes.

#### Resolución

Para el ejercicio 1.1. se creo un script llamado **gemerate_compose.py**. Para su utilización:

```bash
    python generate_compose.py n 
```
siendo n la cantidad de clientes requeridos en este caso

para mas información

```bash
    python generate_compose.py --help
```

### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: `docker volumes`).

#### Resolución 

Para este caso se modificó el archivo **docker-compose-dev.yaml** utilizando docker volumes 

```yaml
    client2:
        container_name: client2
        image: client:latest
        entrypoint: /client
        environment:
        - CLI_ID=2
        - CLI_LOG_LEVEL=DEBUG
        networks:
        - testing_net
        depends_on:
        - server
        volumes:
        - ./client/config.yaml:/config.yaml
```

Como se puede ver se agrega la linea 

```yaml
        volumes:
        - ./client/config.yaml:/config.yaml
```

A su vez se modifica el archivo **generate_compose.py**

### Ejercicio N°3:
Crear un script que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un EchoServer, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado. Netcat no debe ser instalado en la máquina _host_ y no se puede exponer puertos del servidor para realizar la comunicación (hint: `docker network`).


#### Resolución

Para este caso se creo la carpeta server_test, dentro, se encuentran 3 archivos: 
 1. **config.ini** : con la configuración del test
 2. **Dockerfile** : el dockerfile para correr el test
 3. **server_checker.py** : el script para probar el servidor

para correr el script (mientras se corre el servidor):

1. buildear:
```bash
    docker build -t tests .
```

2. correr:
```bash
     docker run --rm -it --network=tp0_testing_net tests
```
Como se puede ver se utiliza el comando network

por ultimo se recibira un mensaje:

```bash
2024-03-26 01:30:24 INFO     Received message from server successfully
```

A su vez se deberá visualizar en los logs del servidor:

```bash
server   | 2024-03-26 01:30:20 INFO     action: accept_connections | result: in_progress
server   | 2024-03-26 01:30:24 INFO     action: accept_connections | result: success | ip: 172.25.125.5
server   | 2024-03-26 01:30:24 INFO     action: receive_message | result: success | ip: 172.25.125.5 | msg: hello server
```


### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

### Resolución

- Servidor: En el caso del servidor se utilizo **signal** para detectar el SIGTERM.

Al producirse un sigterm se ejecuta la siguiente fincion, estableciendo la flag running en false, y cerrando la comunicación con los clientes.
```python
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
```

- Cliente : Para el caso del cliente se utilizaron channels para handelear el sigterm:

al loop del cliente se agrega 

```golang
    case <-c.channel:
        log.Infof("action: sigterm_detected | result: success | client_id: %v",
            c.config.ID,
        )
        break loop
```

## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base al código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).


#### Resolución 

Para la implemación del ejercicio se cambio tanto la estructura del cliente como la estructura del servidor.
-Cliente:
```golang
    func (c *Client) StartClientLoop() {
        c.createClientSocket()
        err := c.bet.sendBet(c.conn, c.config.ID)

        if err != nil {
            log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
                c.config.ID,
                err,
            )
        } else {
            res, err := getResponse(c.conn)
            if err != nil {
                log.Errorf("action: get_response | result: fail | client_id: %v | error: %v",
                    c.config.ID,
                    err,
                )
            }
            log.Infof("action: respuesta_servidor | result: success | message : %v ", res)
            log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", c.bet.ID, c.bet.Number)
        }

        c.conn.Close()
    }

```
En este caso el cliente envia la apuesta a traves de sendBet, esperando luego una respuesta del servidor

```golang
    func (bet *Bet) sendBet(conn net.Conn, ID string) error {
        formatted_bet := bet.GetFormatedBet()
        message := fmt.Sprintf("%s,%s", ID, formatted_bet)
        message_length := len(message)
        buf := make([]byte, 4+message_length)
        binary.BigEndian.PutUint32(buf[:4], uint32(message_length))
        copy(buf[4:], message)
        _, err := conn.Write(buf)
        if err != nil {
            return fmt.Errorf("error al enviar la apuesta al servidor: %v", err)
        }
        return err
}
```

Se obtiene el formato de la apuesta y el tamaño del mensaje y luego se lo envia al servidor

-Servidor: 

```python
        def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bet = receive_bet(client_sock)
            store_bets([bet])
            logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
```
en este caso se utiliza la funcion receive bets para recibir las distintas apuestas de los clientes 

```python
def receive_bet(client_sock):
    # Leer primero la longitud del mensaje
    message_length = int.from_bytes(client_sock.recv(4), byteorder='big')
    # Luego, leer el mensaje completo
    msg = client_sock.recv(message_length).decode('utf-8') 
    
    #envio la respuesta al cliente
    message = 'OK'
    msg_length_bytes = len(message).to_bytes(4, byteorder='big')

    # Enviar los 4 bytes de longitud seguidos del mensaje
    client_sock.sendall(msg_length_bytes + message.encode())
    
    return Bet(*msg.split(","))
```

en primer lugar recibe el tamaño del mensaje para luego leerlo y formatear segun corresponda, luego envia una respuesta 'OK' al cliente para indicar que recibe correctamente el mensaje.

Se evitan short writes y short reads enviando el tamaño y esperando el mismo.



### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento. La cantidad de apuestas dentro de cada _batch_ debe ser configurable. Realizar una implementación genérica, pero elegir un valor por defecto de modo tal que los paquetes no excedan los 8kB. El servidor, por otro lado, deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.


#### Resolución 

Como primera medida se modificaron los archivos de configuracion de los clientes de manera de establecer el tamaño del __batch__. Por otro lado se cambio el archivo **docker-compose-dev.yaml** para tener como volumen el csv perteneciente a cada uno de los clientes.

```yaml
server:
  address: "server:12345"
loop:
  lapse: "0m20s"
  period: "5s"
log:
  level: "info"
file:
  csv: "./agency.csv"
batch_size: 5
```
```yaml
client1:
    container_name: client1
    image: client:latest
    entrypoint: /client
    environment:
    - CLI_ID=1
    - CLI_LOG_LEVEL=DEBUG
    networks:
    - testing_net
    depends_on:
    - server
    volumes:
    - ./client/config.yaml:/config.yaml
    - ./.data/dataset/agency-1.csv:/agency.csv
```

Por otro lado se modifico el cliente de manera de aceptar varias apuestas a la vez

```golang
func (c *Client) StartClientLoop() {
	c.createClientSocket()
	data := csv.NewReader(c.file)
loop:
	for {
		select {
		case <-c.channel:
			break loop
		default:
		}
		bets := GetBetBatch(data, c.config.BatchSize)
		if len(bets) == 0 {
			log.Infof("action: apuestas_enviadas | result: success")
			break loop
		}
		err := SendBets(c.conn, bets, c.config.ID)
		if err != nil {
			log.Infof("action: apuestas_enviadas | result: fail | %v", err.Error())
			break loop
		}

		_, err2 := getResponse(c.conn)

		if err2 != nil {
			log.Errorf("action: get_response | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
		}
	}
	sendCloseMessage(c.conn)
	c.file.Close()
	c.conn.Close()

}

```

Como se puede ver, se creo la función GetBetBatch de manera de obtener el batch de bets a enviar

```golang
func GetBetBatch(reader *csv.Reader, lines int) []*Bet {
	bets := make([]*Bet, 0, lines)
	for i := 0; i < lines; i++ {
		record, err := reader.Read()
		if err != nil {
			return bets[:i]
		}
		bets = append(bets, NewBet(record[0], record[1], record[2], record[3], record[4]))
	}
	return bets
}
```

Nuevamente, para el envio de los mismos se establece el tamaño del mensaje primero

```golang
func SendBets(conn net.Conn, bets []*Bet, ID string) error {
	var message string
	for _, bet := range bets {
		message += bet.GetFormatedBet(ID)
	}
	message_length := len(message)
	buf := make([]byte, 4+message_length)
	binary.BigEndian.PutUint32(buf[:4], uint32(message_length))
	copy(buf[4:], message)
	_, err := conn.Write(buf)
	if err != nil {
		return fmt.Errorf("error al enviar la apuesta al servidor: %v", err)
	}
	return err
}
```

### Ejercicio N°7:
Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo, no podrá responder consultas por la lista de ganadores.
Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

#### Resolución

Para el ejercicio 7 se cambio principalemente el protocolo de mensajeria:

 - `'BET':` indica el envio de una apuesta
 - `'FINISH_BET'` : indica la fincalización del envio de las apuestas
 - `'CLOSE_CONNECTION'` : indica el cierre de la conexion 
 - `'WINNERS' ` : se envia para obtener los ganadores, el servidor podra responder con los mismos o con `NOT_WINNERS_YET`. en el caso de no existir ganadores se retorna `'NOT_A_WINNER'`


__Servidor__:

```python
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
                        logging.info(f'action: sorteo | result: success')
                elif action == 'WINNERS':
                    self.send_winners(client_sock, id)
                else:
                    logging.error("action: command_receive | result: fail | error: {e}")
                
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

```

__Cliente__ :

```golang
func (c *Client) GetWinners() {
loop:
	for {
		select {
		case <-c.channel:
			log.Info("action : handle_sigterm | result : success")
			break loop
		default:
		}
		log.Infof("Asking server for winners: %v", c.config.ID)
		// Asks the server for the winners
		AskForWinners(c.conn, c.config.ID)
		res, err := GetWinnersFromServer(c.conn)
		if err != nil {
			log.Errorf("action: get_winners | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			break loop
		} else {
			if res != "NOT_WINNERS_YET" {
				// Received winners
				log.Infof("action: get_winners | result: success | client_id: %v | winners: %v",
					c.config.ID,
					res,
				)
				break loop
			}

		}

	}
	// Close the connection
	SendCloseMessage(c.conn, c.config.ID)
	c.conn.Close()

}
```

## Parte 3: Repaso de Concurrencia

### Ejercicio N°8:
Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo.
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

En caso de que el alumno implemente el servidor Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

#### Resolución

Para la resolución del ejercicio propuesto se utilizó la biblioteca `multiprocessing` de python. Para el procesamiento de consultas en paralelo se utilizaron locks para garantizar la exclusion mutua.

```python
        self.locks = {
            'bets': manager.Lock(),
            'clients': manager.Lock()
        }
        
        self.data = manager.dict({
            'finished_clients' : 0,
            'all_clients_ready' : False
        })
```

En cuanto al loop aceptador

```python
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
            client = multiprocessing.Process(target=self.__handle_client_connection, args=(client_sock,))
            self.clients.append(client)
            client.start()
        
        self._server_socket.close()

```
Como se puede ver por cada cliente nuevo se lanza un proceso

El loop principal del cliente no sufrió de mayores cambios a excepcion del uso de locks, ejemplo:


```python
    elif action == 'BET':
        bet = receive_bet(client_sock)
        self.locks['bets'].acquire()
        store_bets(bet)
        self.locks['bets'].release()

```


## Consideraciones Generales
Se espera que los alumnos realicen un _fork_ del presente repositorio para el desarrollo de los ejercicios.
El _fork_ deberá contar con una sección de README que indique como ejecutar cada ejercicio.
La Parte 2 requiere una sección donde se explique el protocolo de comunicación implementado.
La Parte 3 requiere una sección que expliquen los mecanismos de sincronización utilizados.

Finalmente, se pide a los alumnos leer atentamente y **tener en cuenta** los criterios de corrección provistos [en el campus](https://campusgrado.fi.uba.ar/mod/page/view.php?id=73393).

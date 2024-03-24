package common

import (
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config  ClientConfig
	conn    net.Conn
	channel chan os.Signal
	bet     *Bet
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, channel chan os.Signal, bet *Bet) *Client {
	client := &Client{
		config:  config,
		channel: channel,
		bet:     bet,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
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

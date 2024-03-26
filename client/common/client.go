package common

import (
	"encoding/csv"
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
	BatchSize     int
}

// Client Entity that encapsulates how
type Client struct {
	config  ClientConfig
	conn    net.Conn
	channel chan os.Signal
	file    *os.File
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, channel chan os.Signal, file *os.File) *Client {
	client := &Client{
		config:  config,
		channel: channel,
		file:    file,
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

func (c *Client) SendBets() {
	data := csv.NewReader(c.file)
loop:
	for {
		select {
		case <-c.channel:
			break loop
		default:
		}

		// Gets the bets to send
		bets := GetBetBatch(data, c.config.BatchSize)
		if len(bets) == 0 {
			log.Infof("action: apuestas_enviadas | result: success")
			break loop
		}
		// Sends the bets if the len is greater than 0
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
	// Closes the connection
	SendEndOfBets(c.conn, c.config.ID)
	SendCloseMessage(c.conn, c.config.ID)
	c.file.Close()
	c.conn.Close()

}

// GetWinners sends a request to the server to get the winners and handles the response
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

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	c.createClientSocket()
	c.SendBets()
	c.createClientSocket()
	c.GetWinners()
}

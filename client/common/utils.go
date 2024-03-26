package common

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"net"
)

type action string

const (
	action_bet     action = "BET"
	action_close   action = "CLOSE_CONNECTION"
	action_winners action = "WINNERS"
	action_finish  action = "FINISH_BET"
)

type Bet struct {
	Name      string
	LastName  string
	ID        string
	BirthDate string
	Number    string
}

// NewBet returns a new instance of a Bet.
func NewBet(name string, lastname string, id string, birthdate string, number string) *Bet {
	return &Bet{
		Name:      name,
		LastName:  lastname,
		ID:        id,
		BirthDate: birthdate,
		Number:    number,
	}
}

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

// GetFormatedBet returns a formatted string of the Bet.
func (bet *Bet) GetFormatedBet(id string) string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s;", id, bet.Name, bet.LastName, bet.ID, bet.BirthDate, bet.Number)
}

// Sends a message to the server with the given action and id
func SendAction(conn net.Conn, value action, ID string) error {
	actionWithID := fmt.Sprintf("%s,%s", ID, value)
	actionBytes := []byte(actionWithID)
	messageLength := len(actionBytes)
	messageBuffer := make([]byte, 4+messageLength)
	binary.BigEndian.PutUint32(messageBuffer[:4], uint32(messageLength))
	copy(messageBuffer[4:], actionBytes)
	_, err := conn.Write(messageBuffer)
	if err != nil {
		return fmt.Errorf("error al enviar la acción al servidor: %v", err)
	}
	return nil

}

// SendBets sends a batch of bets to the server over the given connection.
func SendBets(conn net.Conn, bets []*Bet, ID string) error {
	SendAction(conn, action_bet, ID)
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
		return fmt.Errorf("error sending to server: %v", err)
	}
	return err
}

// sendCloseMessage sends a close message to the given connection.
func SendCloseMessage(conn net.Conn, ID string) error {
	return SendAction(conn, action_close, ID)
}

// SendEndOfBets sends the end of bets message to the server.
func SendEndOfBets(conn net.Conn, ID string) error {
	// Send the end of bets message to the server.
	return SendAction(conn, action_finish, ID)
}

// AskForWinners sends the "winners" action to the server, indicating that the client is ready to receive the winners list.
func AskForWinners(conn net.Conn, ID string) error {
	return SendAction(conn, action_winners, ID)
}

// GetWinnersFromServer receives the winners list from the server.
func GetWinnersFromServer(conn net.Conn) (string, error) {
	return getResponse(conn)
}

// getResponse reads a response from the given connection.
func getResponse(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	lengthBytes := make([]byte, 4)
	if _, err := io.ReadFull(reader, lengthBytes); err != nil {
		return "", err
	}
	messageLength := binary.BigEndian.Uint32(lengthBytes)
	messageBytes := make([]byte, messageLength)
	if _, err := io.ReadFull(reader, messageBytes); err != nil {
		return "", err
	}
	return string(messageBytes), nil
}

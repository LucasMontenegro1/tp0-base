package common

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"net"
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

// SendBets sends a batch of bets to the server over the given connection.
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

// sendCloseMessage sends a close message to the given connection.
func sendCloseMessage(conn net.Conn) {
	closeMessage := "CLOSE_CONNECTION"
	messageLength := uint32(len(closeMessage))

	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, messageLength)
	conn.Write(lenBuf)

	conn.Write([]byte(closeMessage))
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

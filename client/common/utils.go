package common

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

type Bet struct {
	Name      string
	LastName  string
	ID        string
	BirthDate string
	Number    string
}

// NewBet returns a new instance of a Bet.
func NewBet() *Bet {
	return &Bet{
		Name:      os.Getenv("NOMBRE"),
		LastName:  os.Getenv("APELLIDO"),
		ID:        os.Getenv("DOCUMENTO"),
		BirthDate: os.Getenv("NACIMIENTO"),
		Number:    os.Getenv("NUMERO"),
	}
}

// GetFormatedBet returns a formatted string of the Bet.
func (bet *Bet) GetFormatedBet() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s\n", bet.Name, bet.LastName, bet.ID, bet.BirthDate, bet.Number)
}

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

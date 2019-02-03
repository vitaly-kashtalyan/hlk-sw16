package hlk_sw16

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strings"
)

const (
	prefix        = "\xaa"
	suffix        = "\xbb"
	statusReading = "\x1e"
	switchOnAll   = "\x01"
	switchOffAll  = "\x02"
	action        = "\x0f"
	actionOnAll   = "\x0a"
	actionOffAll  = "\x0b"
	relayOn       = "\x01"
	relayOff      = "\x02"
	defaultByte   = "\x01"
)

var (
	relays = []string{"\x00", "\x01", "\x02", "\x03", "\x04", "\x05", "\x06", "\x07", "\x08", "\x09", "\x0a", "\x0b", "\x0c", "\x0d", "\x0e", "\x0f"}
)

type Connection struct {
	Conn net.Conn
	Err  error
}

func New(host string, port string) (c *Connection) {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Println("Error establishing connection:", err.Error())
	}

	return &Connection{
		Conn: conn,
		Err:  err,
	}
}

func (c *Connection) Close() (err error) {
	if c.Conn != nil {
		err = c.Conn.Close()
	}
	return
}

func (c *Connection) SwitchAllOff() (err error) {
	return c.relaySwitch(actionOffAll, switchOffAll, relayOff)
}

func (c *Connection) SwitchAllOn() (err error) {
	return c.relaySwitch(actionOnAll, switchOnAll, relayOn)
}

func (c *Connection) RelayOn(id int) (err error) {
	if id < 0 || id > 15 {
		err = errors.New("argument id has an invalid value: please use only 0-15")
	} else {
		err = c.relaySwitch(action, relays[id], relayOn)
	}
	return
}

func (c *Connection) RelayOff(id int) (err error) {
	if id < 0 || id > 15 {
		err = errors.New("argument id has an invalid value: please use only 0-15")
	} else {
		err = c.relaySwitch(action, relays[id], relayOff)
	}
	return
}

func (c *Connection) relaySwitch(action string, relay string, state string) (err error) {
	return c.WriteMessage(prefix + action + relay + state + strings.Repeat(defaultByte, 15) + suffix)
}

func (c *Connection) StatusRelays() (err error) {
	return c.WriteMessage(prefix + statusReading + strings.Repeat(defaultByte, 17) + suffix)
}

func (c *Connection) WriteMessage(msg string) (err error) {
	if _, err := c.Conn.Write([]byte(msg)); err != nil {
		log.Println("Write to server failed:", err.Error())
	}
	return
}

func (c *Connection) ReadMessage() (msg []byte, err error) {
	reader := bufio.NewReader(c.Conn)
	b, err := reader.ReadByte()
	if err != nil {
		log.Println("Error reading:", err.Error())
		return
	}

	if reader.Buffered() > 0 {
		msg = append(msg, b)
		for reader.Buffered() > 0 {
			// read byte by byte until the buffered data is not empty
			b, err := reader.ReadByte()
			if err == nil {
				msg = append(msg, b)
			} else {
				log.Println("Unreadable character:", b)
			}
		}
	}
	return
}

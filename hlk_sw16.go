package hlk_sw16

import (
	"bufio"
	"log"
	"net"
	"strings"
)

const (
	Prefix        = "\xaa"
	Suffix        = "\xbb"
	StatusReading = "\x1e"
	SwitchOnAll   = "\x01"
	SwitchOffAll  = "\x02"
	Action        = "\x0f"
	ActionOnAll   = "\x0a"
	ActionOffAll  = "\x0b"
	RelayOn       = "\x01"
	RelayOff      = "\x02"
	Default       = "\x01"
)

var (
	Relays = []string{"\x00", "\x01", "\x02", "\x03", "\x04", "\x05", "\x06", "\x07", "\x08", "\x09", "\x0a", "\x0b", "\x0c", "\x0d", "\x0e", "\x0f"}
)

type Connection struct {
	conn net.Conn
	err  error
}

func new(host string, port string) (c *Connection) {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Println("Error establishing connection:", err.Error())
	}

	return &Connection{
		conn: conn,
		err:  err,
	}
}

func (c *Connection) close() (err error) {
	if c.conn != nil {
		err = c.conn.Close()
	}
	return
}

func (c *Connection) switchAllOff() (err error) {
	return c.relaySwitch(ActionOffAll, SwitchOffAll, RelayOff)
}

func (c *Connection) switchAllOn() (err error) {
	return c.relaySwitch(ActionOnAll, SwitchOnAll, RelayOn)
}

func (c *Connection) relayOn(id int) (err error) {
	return c.relaySwitch(Action, Relays[id], RelayOn)
}

func (c *Connection) relayOff(id int) (err error) {
	return c.relaySwitch(Action, Relays[id], RelayOff)
}

func (c *Connection) relaySwitch(action string, relay string, state string) (err error) {
	return c.writeMessage(Prefix + action + relay + state + strings.Repeat(Default, 15) + Suffix)
}

func (c *Connection) readStatusRelays() (err error) {
	return c.writeMessage(Prefix + StatusReading + strings.Repeat(Default, 17) + Suffix)
}

func (c *Connection) writeMessage(msg string) (err error) {
	if _, err := c.conn.Write([]byte(msg)); err != nil {
		log.Println("Write to server failed:", err.Error())
	}
	return
}

func (c *Connection) readMessage() (msg []byte, err error) {
	reader := bufio.NewReader(c.conn)
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
				return
			}
		}
	}
	return
}

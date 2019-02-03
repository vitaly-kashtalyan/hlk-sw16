package hlk_sw16

import (
	"bufio"
	"log"
	"net"
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

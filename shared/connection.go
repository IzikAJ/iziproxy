package shared

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
)

const (
	breakPoint = byte(0)
)

// Connection - net conection with mutex
type Connection struct {
	// mux    sync.RWMutex
	reader *bufio.Reader
	net.Conn
}

func md5Hash(data []byte) string {
	cntMD5 := md5.New()
	cntMD5.Write(data)
	return hex.EncodeToString(cntMD5.Sum(nil))
}

func info(prefix string, data []byte) {
	max := 150
	dlen := len(data)
	if dlen > max {
		fmt.Printf("[%s (%d)\n%s\n...\n%s\n]\n", prefix, dlen, data[:(max/2)], data[dlen-(max/2):])
	} else if dlen > 0 {
		fmt.Printf("[%s (%d)\n%s\n]\n", prefix, dlen, data)
	}
}

// Init - initialize connection instance
func (conn *Connection) Init() {
	fmt.Println("INITIALIZE Connection", conn.RemoteAddr().String())
	(*conn).reader = bufio.NewReader(conn)
}

// ReadRaw will read data from stream
func (conn *Connection) ReadRaw() (data []byte, err error) {
	data, err = conn.reader.ReadBytes(breakPoint)
	if err != nil {
		fmt.Println("CONNECTION_READ ERROR", err)
	}
	dlen := len(data)
	if dlen > 0 && data[dlen-1] == breakPoint {
		data = data[:dlen-1]
	}
	// info("R", data)
	return
}

// WriteRaw will write data to stream
func (conn *Connection) WriteRaw(data []byte) (err error) {
	_, err = (*conn).Write(append(data, breakPoint))
	if err != nil {
		fmt.Println("CONNECTION_WRITE ERROR", err)
		return
	}
	// info("W", data)
	return
}

// SendMessage just encode and send message to connection
func (conn *Connection) SendMessage(msg Message) (err error) {
	data, err := msg.dump()
	if err != nil {
		fmt.Println("Cant dump message")
		return err
	}
	err = conn.WriteRaw(data)
	if err != nil {
		fmt.Println("CLIENT DISCONNECTED", (*conn).RemoteAddr())
		return err
	}
	return nil
}

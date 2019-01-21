package requestsender

import (
	"fmt"
	"net"
	"strconv"
)

type connectionInfo struct {
	protocol string
	address  string
	port     int
}

func (c *connectionInfo) SetConnectionInfo(protocol string, address string, port int) {
	c.protocol = protocol
	c.address = address
	c.port = port
}

func (c *connectionInfo) connect() net.Conn {
	conn, err := net.Dial(c.protocol, c.address+":"+strconv.Itoa(c.port))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connect to", c.address, ":", c.port, "by", c.protocol)
	return conn
}

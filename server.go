package main

import (
	"./structs"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":2000")

	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()
	fmt.Println("Server is listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		t := structs.Hid_t{}
		buff := make([]byte, (1024 * 4))
		n, err := conn.Read(buff)
		if n == 0 || err != nil {
			fmt.Println("Read error: ", err)
			break
		}

		fmt.Println(buff[:n])

		t.Deserialization(buff, n)
		fmt.Println(t)
	}

}

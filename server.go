package main

import (
	"log"
	"net"

	"./structs"
)

func main() {
	listener, err := net.Listen("tcp", ":2000")

	if err != nil {
		log.Println(err)
	}
	defer listener.Close()
	log.Println("Server is listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
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
			log.Println("Read error: ", err)
			break
		}

		log.Println(buff[:n])

		t.Deserialization(buff, n)
		log.Println(t)
	}

}

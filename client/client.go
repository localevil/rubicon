package main

import (
	"fmt"
	"net"

	"../structs"
)

var transactionCount = uint16(0)

func main() {
	conn, err := net.Dial("tcp", "10.0.40.199:2000")
	if err != nil {
		panic(err)
	}
	send(conn)
}

func send(conn net.Conn) {
	defer conn.Close()
	for {
		message := structs.NewPackage_t(true,
			structs.Hid_t{Type_t: 0x22, Serial: 1975},
			&structs.FirmwareVersionData{},
			&transactionCount)

		fmt.Println(message.ToByteSlice())
		n, err := conn.Write(message.ToByteSlice())
		if err != nil {
			fmt.Println(err)
			continue
		}
		buf := make([]byte, 265)
		n, err = conn.Read(buf)
		if n == 0 || err != nil {
			fmt.Println(err)
			fmt.Scanln()
			continue
		}

		fmt.Println(buf[:n])
		fmt.Scanln()
	}
}

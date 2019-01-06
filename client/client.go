package main

import (
	"../structs"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:2000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for {
		message := structs.Hid_t{}
		fmt.Println("Enter type: ")
		_, err := fmt.Scanln(&message.Type_t)
		if err != nil {
			panic(err)
		}

		fmt.Println("Enter serial: ")
		_, err = fmt.Scanln(&message.Serial)
		if err != nil {
			panic(err)
		}
		fmt.Println(message.Serialization())
		conn.Write(message.Serialization())
	}
}

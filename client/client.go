package main

import (
	"../structs"
	"fmt"
	"net"
)

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
		str := "Hello"
		message := structs.Package_t{}
		message.BeginSequence = [2]byte{0xB6, 0x49}
		message.ReciverAddres = structs.Hid_t{0x22, 1975}
		message.InfoPartLen = uint8(len(str))
		copy(message.InfoPart[:message.InfoPartLen], str[:])
		//buf := message.ToByteSlice()
		//size := uint32(len(buf))
		fmt.Println(str)
		//message.CrcValue = structs.CalcCrcCcitt(buf, size, 0xFFFE)

		//fmt.Println(message.Serialization())
		//conn.Write(message.Serialization())
		buffer := make([]byte, (1024 * 4))
		length, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(buffer[:length])
	}
}

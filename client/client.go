package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"unsafe"

	"../structs"
)

var transactionCount = uint16(0)

func main() {
	drv := driver{}
	drv.connect("tcp", "10.0.40.199", 2000)
}

type driver struct {
	hit  structs.HidT
	conn net.Conn
}

func (d *driver) connect(protocol string, address string, port int) {
	conn, err := net.Dial("tcp", "10.0.40.199:"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	d.conn = conn
	d.start()
}

func (d *driver) start() {
	message := structs.NewRequestPackage(true,
		structs.HidT{TypeT: 0xff, Serial: 0xffff},
		&structs.HandShake{},
		&transactionCount)
	_, err := d.conn.Write(message.ToByteSlice())
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 265)
	n, err := d.conn.Read(buf)
	if n == 0 || err != nil {
		panic(err)
	}
	response := structs.NewResponsePackage(buf[:n])
	newHandShakeResponce := structs.NewHandShakeResponse(response.InfoPart)
	fmt.Printf("%x \n", *newHandShakeResponce)
	d.hit = response.ReceiverAddress

	d.checkStatus()
	takeAlive()
}

func takeAlive() {
	var quit string
	for quit != "q" {
		fmt.Scanln(quit)
	}
}

func (d *driver) checkStatus() {
	message := structs.NewRequestPackage(true,
		d.hit,
		&structs.StatusWord{},
		&transactionCount)
	_, err := d.conn.Write(message.ToByteSlice())
	if err != nil {
		panic(err)
	}
	d.handleChackStatus()
}

func (d *driver) handleChackStatus() {
	buf := make([]byte, 265)
	n, err := d.conn.Read(buf)
	if n == 0 || err != nil {
		panic(err)
	}
	response := structs.NewResponsePackage(buf[:n])
	//structs.NewStatusWordResponse(response.InfoPart)
	statusWordResponse := (*(*structs.StatusWordResponse)(unsafe.Pointer(&response.InfoPart[0])))
	fmt.Println(statusWordResponse)
	time.Sleep(2 * time.Second)
	go d.checkStatus()
}

func (d *driver) TakeVersion() {
	message := structs.NewRequestPackage(true,
		d.hit,
		&structs.FirmwareVersion{},
		&transactionCount)
	_, err := d.conn.Write(message.ToByteSlice())
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 265)
	n, err := d.conn.Read(buf)
	if n == 0 || err != nil {
		panic(err)
	}
	response := structs.NewResponsePackage(buf[:n])
	firmwareVersionResponse := structs.NewFirmwareVersionResponse(response.InfoPart)
	fmt.Println(*firmwareVersionResponse)
	fmt.Printf("%x \n", *firmwareVersionResponse)
}

func handleResponse(conn net.Conn) {
}

func timer(conn net.Conn) {
	time.Sleep(5 * time.Second)
	go handleResponse(conn)
}

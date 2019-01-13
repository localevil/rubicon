package main

import (
	"sync"

	"./requestsender"

	"../structs"
)

//device struct
type device struct {
	hid    structs.Hid
	sender requestsender.Sender
	wg     sync.WaitGroup
}

func (d *device) connect(protocol string, address string, port int) {
	d.sender.SetWaitGroup(&d.wg)
	d.sender.SetConnectionInfo(protocol, address, port)
	d.Start()
}

func (d *device) Start() {
	go d.Tick()
}

func (d *device) Stop() {
	d.wg.Wait()
}

func (d *device) Tick() {
	for {
		d.sender.AddCommand(&structs.StatusWord{})
	}
}

func (d *device) SendHandShake() {
	emptyHid := structs.Hid{TypeT: 0xff, Serial: 0xffff}
	d.sender.AddCommand(&structs.HandShake{}, emptyHid)
}

//TakeVersion send command to take version
func (d *device) TakeVersion() {
	d.sender.AddCommand(&structs.FirmwareVersion{})
}

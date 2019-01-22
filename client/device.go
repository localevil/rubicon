package main

import (
	"encoding/binary"
	"log"
	"os"
	"sync"
	"time"

	"./requestsender"
)

//device struct
type device struct {
	sender           requestsender.Sender
	wg               sync.WaitGroup
	hid              hid
	transactionCount uint16
}

func (d *device) connect(protocol string, address string, port int) {
	d.sender.SetWaitGroup(&d.wg)
	d.sender.SetConnectionInfo(protocol, address, port)
}

func (d *device) Start() {
	d.sender.Start()
	d.SendHandShake()
}

func (d *device) Stop() {
	d.wg.Wait()
	d.sender.Stop()
}

func readIndex() uint32 {
	file, err := os.Open("index.bin")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	data := make([]byte, 4)
	var n int
	n, err = file.Read(data)
	return binary.LittleEndian.Uint32(data[:n])
}

func writeIndex(index uint32) {
	file, err := os.Create("index.bin")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, index)
	file.Write(data)
}

func isAddressLoop(id uint16) bool {
	return id >= 0xA000 || id <= 0xBFFF
}

func printSnid(s snid) {
	log.Println("Serial nummber:", s.typeOfSind(), s.serialNumber())
	if isAddressLoop(s.id) {
		log.Println("ID:", s.addressLoop(), s.shortAddress(), s.logicSubDevNum())
	}
}

func (d *device) Tick() {
	sequence := uint8(0)
	index := readIndex()
	req := newRequestPackage(true, d.hid, newStatusWord(), &d.transactionCount)
	request := requestsender.Request{
		Buffer: req.ToByteSlice(),
		Handler: func(buf []byte) {
			response := newResponsePackage(buf)
			command := newStatusWordResponse(response.InfoPart)
			if command.isNewEvent() {
				go d.TakeEvent(index)
				index = readIndex()
			}
			if command.isNewStatuses() {
				go d.TakeNewStatus(sequence)
				sequence++
			}
		}}
	for {
		d.sender.AddCommand(request)
		time.Sleep(2 * time.Second)
	}
}

func (d *device) SendHandShake() {
	emptyHid := hid{TypeT: 0xff, Serial: 0xffff}
	req := newRequestPackage(true, emptyHid, newHandShake(), &d.transactionCount)
	request := requestsender.Request{
		Buffer: req.ToByteSlice(),
		Handler: func(buf []byte) {
			responce := newResponsePackage(buf)
			d.hid = responce.ReceiverAddress
			command := handShakeResponse{}
			command.deserialization(responce.InfoPart)
		},
		Name: "SendHandShake"}
	d.sender.AddCommand(request)
	time.Sleep(2 * time.Second)
}

//TakeVersion send command to take version
func (d *device) TakeVersion() {
	req := newRequestPackage(true, d.hid, newFirmwareVersion(), &d.transactionCount)
	request := requestsender.Request{
		Buffer: req.ToByteSlice(),
		Handler: func(buf []byte) {
			response := newResponsePackage(buf)
			newFirmwareVersionResponse(response.InfoPart)
		},
		Name: "TakeVersion"}
	d.sender.AddCommand(request)
}

var sizeOfFile uint8

func (d *device) TakeFileList() {
	req := newRequestPackage(true, d.hid, newFileList(), &d.transactionCount)
	request := requestsender.Request{
		Buffer: req.ToByteSlice(),
		Handler: func(buf []byte) {
			response := newResponsePackage(buf)
			command := fileListResponse{}
			command.deserialization(response.InfoPart)
			for i := uint16(0); i < command.num; i++ {
				id := command.files[i].ID
				log.Printf("% x", id)
				if id == 0x1201 {
					size := command.files[i].size
					sizeOfFile = uint8(size)
				}
			}
		},
		Name: "TakeFileList"}
	d.sender.AddCommand(request)
}

func (d *device) SendCommandToUnit(snid snid, evcode uint16, userNum uint16, args []byte) {
	req := newRequestPackage(true, d.hid, newCommandToUnit(snid, evcode, userNum, args), &d.transactionCount)
	request := requestsender.Request{
		Buffer: req.ToByteSlice(),
		Handler: func(buf []byte) {
			response := newResponsePackage(buf)
			infoPart := newCommandToUnitResponse(response.InfoPart)
			printSnid(infoPart.snid)
		},
		Name: "CommandToUnit"}
	d.sender.AddCommand(request)
}

func (d *device) TakeFindSnOnAS() {
	req := newRequestPackage(true, d.hid, newTakeFindSnOnAS(), &d.transactionCount)
	request := requestsender.Request{
		Buffer: req.ToByteSlice(),
		Handler: func(buf []byte) {
			response := newResponsePackage(buf)
			command := takeFindSnOnASResponse{}
			command.deserialization(response.InfoPart)
			for i := uint8(0); i < command.num; i++ {
				printSnid(command.serials[i])
			}
		},
		Name: "TakeFindSnOnAS"}
	d.sender.AddCommand(request)
}

func (d *device) ReadFile(id uint16, size uint8, addr uint32) {
	req := newRequestPackage(true, d.hid, newReadFile(id, size, addr), &d.transactionCount)
	request := requestsender.Request{
		Buffer: req.ToByteSlice(),
		Handler: func(buf []byte) {
			response := newResponsePackage(buf)
			command := readFileResponse{}
			command.deserialization(response.InfoPart)
			log.Println(command.data)
		},
		Name: "ReadFile"}
	d.sender.AddCommand(request)
}

func uint16ToBoll(v uint16) bool {
	if v == 0 {
		return false
	}
	return true
}

func (d *device) TakeNewStatus(sequence uint8) {
	req := newRequestPackage(true, d.hid, newNewStatus(sequence), &d.transactionCount)
	request := requestsender.Request{
		Buffer: req.ToByteSlice(),
		Handler: func(buf []byte) {
			response := newResponsePackage(buf)
			if returnCodeMap[binary.LittleEndian.Uint16(buf[8:10])] == returnCodeMap[30] {
				log.Println(returnCodeMap[30])
				return
			}
			command := newNewStatusResponse(response.InfoPart)
			for i := uint8(0); i < command.num; i++ {
				stat := command.sindInfos[i].statusT
				if !stat.isNormal() {
					log.Printf("Fire: %t Alarm: %t Fault: %t AP: %t Bypass: %t Wait: %t Not ready: %t TS: %t Armed: %t On: %t Error: %t",
						uint16ToBoll(stat.isFire()),
						uint16ToBoll(stat.isAlarm()),
						uint16ToBoll(stat.isFault()),
						uint16ToBoll(stat.isAp()),
						uint16ToBoll(stat.isBypass()),
						uint16ToBoll(stat.isWait()),
						uint16ToBoll(stat.isNotReady()),
						uint16ToBoll(stat.isTechSig()),
						uint16ToBoll(stat.isArmed()),
						uint16ToBoll(stat.isOn()),
						uint16ToBoll(stat.isError()))
				}

				printSnid(command.sindInfos[i].area)
				printSnid(command.sindInfos[i].snidT)

				log.Println("Dop status code:", statusCodeMap[command.sindInfos[i].code])
			}
		},
		Name: "TakeNewStatus"}
	d.sender.AddCommand(request)
	time.Sleep(time.Second)
}

func (d *device) TakeEvent(index uint32) {
	req := newRequestPackage(true, d.hid, newTakeEvent(index), &d.transactionCount)
	request := requestsender.Request{
		Buffer: req.ToByteSlice(),
		Handler: func(buf []byte) {
			response := newResponsePackage(buf)
			command := newTakeEventResponse(response.InfoPart)

			printSnid(command.event.dst)
			printSnid(command.event.src)
			log.Println(eventCodeMap[command.event.evcode])
		},
		Name: "TakeEvent"}
	d.sender.AddCommand(request)
}

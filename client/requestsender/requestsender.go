package requestsender

import (
	"net"
	"sync"

	"../../structs"
)

//Sender struct
type Sender struct {
	connectionInfo
	conn             net.Conn
	query            chan structs.RequestPackage
	wg               *sync.WaitGroup
	transactionCount uint16
	hid              structs.Hid
}

//Start sending command from query
func (s *Sender) Start() {
	s.conn = s.connect()
	go s.sendCommands()
	s.wg.Add(1)
}

//Stop stop
func (s *Sender) Stop() {
	s.conn.Close()
}

//SetWaitGroup setting wait group
func (s *Sender) SetWaitGroup(wg *sync.WaitGroup) {
	s.wg = wg
}

//SetHid setting HID
func (s *Sender) SetHid(hid structs.Hid) {
	s.hid = hid
}

func (s *Sender) handleCommand(command structs.Command) structs.RecivedDate {
	buffer := make([]byte, 265)
	size, err := s.conn.Read(buffer)
	if size == 0 || err != nil {
		panic(err)
	}
	response := structs.NewResponsePackage(buffer[:size])
	defer s.wg.Done()
	return command.GetResponse(response.InfoPart)
}

func (s *Sender) sendCommands() {
	for sendPackage := range s.query {
		_, err := s.conn.Write(sendPackage.ToByteSlice())
		if err != nil {
			panic(err)
		}
		s.handleCommand(sendPackage.InfoPart)
	}
}

//AddCommand put conmmand to query in sender
func (s *Sender) AddCommand(command structs.Command, hid ...structs.Hid) {
	if len(hid) == 0 {
		s.query <- structs.NewRequestPackage(true, s.hid, command, &s.transactionCount)
	} else {
		for i := 0; i < len(hid); i++ {
			s.query <- structs.NewRequestPackage(true, hid[i], command, &s.transactionCount)
		}
	}

}

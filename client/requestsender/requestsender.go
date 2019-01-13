package requestsender

import (
	"net"

	"../../structs"
)

//Sender struct
type Sender struct {
	conn   net.Conn
	query  []structs.DataT
	active bool
}

//Start sending command frome query
func (s *Sender) Start() {
	go s.sendCommands()

}
func (s *Sender) sendCommands() {
	for i := 0; i <= len(s.query); i++ {
		_, err := s.conn.Write(s.query[i].ToByteSlice())
		if err != nil {
			panic(err)
		}
		if
	}
}

//AddCommand put conmmand to query in sender
func (s *Sender) AddCommand(command structs.DataT) {
	s.query = append(s.query, command)
	if !s.active {
		go s.Start()
	}
}

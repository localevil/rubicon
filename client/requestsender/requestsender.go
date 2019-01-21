package requestsender

import (
	"fmt"
	"net"
	"sync"
)

//ResponseHandlerT type
type ResponseHandlerT func([]byte)

//Request 123
type Request struct {
	Buffer  []byte
	Handler ResponseHandlerT
	Name    string
}

//Sender struct
type Sender struct {
	connectionInfo
	conn  net.Conn
	queue chan Request
	wg    *sync.WaitGroup
}

//Start sending command from queue
func (s *Sender) Start() {
	s.conn = s.connect()
	go s.process()
}

//Stop stop
func (s *Sender) Stop() {
	s.conn.Close()
	close(s.queue)
}

//SetWaitGroup setting wait group
func (s *Sender) SetWaitGroup(wg *sync.WaitGroup) {
	s.queue = make(chan Request)
	s.wg = wg
}

func (s *Sender) tryReconnect() {
	s.conn = s.connect()
}

func (s *Sender) handleCommand(request Request) {
	buffer := make([]byte, 265)
	size, err := s.conn.Read(buffer)
	if size == 0 || err != nil {
		fmt.Println(err)
		return
	}
	if request.Handler != nil {
		request.Handler(buffer[:size])
	}
	fmt.Println("Finish sending: ", request.Name)
}

func (s *Sender) process() {
	s.wg.Add(1)
	defer s.wg.Done()
	fmt.Println("Start sending commands")
	for request := range s.queue {
		fmt.Println("Start sending: ", request.Name)
		fmt.Printf("> % x\n", request.Buffer)
		_, err := s.conn.Write(request.Buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		s.handleCommand(request)
	}
	fmt.Println("Stop sending commands")
}

//AddCommand put conmmand to queue in sender
func (s *Sender) AddCommand(request Request) {
	s.queue <- request
}

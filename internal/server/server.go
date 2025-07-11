package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener    net.Listener
	Close_state atomic.Bool
}

func Serve(port int) (*Server, error) {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}
	server := Server{
		Listener: listner,
	}
	server.Close_state.Store(false)
	go server.listen()
	return &server, nil
}

func (s *Server) Close() error {
	s.Close_state.Store(true)
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) listen() {
	for {
		if s.Close_state.Load() == true {
			break
		}
		connection, err := s.Listener.Accept()
		if err != nil {
			if s.Close_state.Load() == true {
				break
			}
			log.Println("error accepting connection: ", err)
		}
		go s.handle(connection)
	}
}

func (s *Server) handle(conn net.Conn) {
	response_string := "" +
		"HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n\r\n" +
		"Hello World!\r\n"
	_, err := conn.Write([]byte(response_string))
	if err != nil {
		log.Println("error writing to connection: ", err)
	}
	conn.Close()
}

package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/Drnel/btdv_http_from_tcp/internal/response"
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
	err := response.WriteStatusLine(conn, response.Ok)
	if err != nil {
		log.Println("error writing Status line to connection: ", err)
	}
	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		log.Println("error writing headers to connection: ", err)
	}
	conn.Close()
}

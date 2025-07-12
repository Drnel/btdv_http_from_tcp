package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/Drnel/btdv_http_from_tcp/internal/request"
	"github.com/Drnel/btdv_http_from_tcp/internal/response"
)

type Server struct {
	listener    net.Listener
	handler     Handler
	close_state atomic.Bool
}

type HandlerError struct {
	Status_code int
	Message     string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (h *HandlerError) Write(w io.Writer) error {
	default_headers := response.GetDefaultHeaders(len(h.Message))
	err := response.WriteStatusLine(w, response.StatusCode(h.Status_code))
	if err != nil {
		return err
	}
	err = response.WriteHeaders(w, default_headers)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(h.Message))
	if err != nil {
		return err
	}
	return nil
}

func Serve(port int, handler Handler) (*Server, error) {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}
	server := Server{
		listener: listner,
		handler:  handler,
	}
	server.close_state.Store(false)
	go server.listen()
	return &server, nil
}

func (s *Server) Close() error {
	s.close_state.Store(true)
	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) listen() {
	for {
		if s.close_state.Load() == true {
			break
		}
		connection, err := s.listener.Accept()
		if err != nil {
			if s.close_state.Load() == true {
				break
			}
			log.Println("error accepting connection: ", err)
		}
		go s.handle(connection)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	request, err := request.RequestFromReader(conn)
	if err != nil {
		log.Println("error while parsing request:", err)
	}
	var buf bytes.Buffer
	handler_err := s.handler(&buf, request)
	if handler_err != nil {
		err := handler_err.Write(conn)
		if err != nil {
			log.Println("error while writing handler error to connection:", err)
		}
		return
	}
	default_response_headers := response.GetDefaultHeaders(buf.Len())
	err = response.WriteStatusLine(conn, response.Ok)
	if err != nil {
		log.Println("error writing status line:", err)
		return
	}
	err = response.WriteHeaders(conn, default_response_headers)
	if err != nil {
		log.Println("error writing headers:", err)
		return
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Println("error writing body:", err)
		return
	}
}

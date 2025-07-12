package server

import (
	"fmt"
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

type Handler func(w *response.Writer, req *request.Request) *HandlerError

func (h *HandlerError) Write(w *response.Writer) error {
	default_headers := response.GetDefaultHeaders(len(h.Message))
	default_headers["Content-Type"] = "text/html"
	err := w.WriteStatusLine(response.StatusCode(h.Status_code))
	if err != nil {
		return err
	}
	err = w.WriteHeaders(default_headers)
	if err != nil {
		return err
	}
	_, err = w.WriteBody([]byte(h.Message))
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
	response_writer := response.Writer{
		WriterState: response.WriteStateStatusLine,
		Writer:      conn,
	}
	handler_err := s.handler(&response_writer, request)
	if handler_err != nil {
		err := handler_err.Write(&response_writer)
		if err != nil {
			log.Println("error while writing handler error to connection:", err)
		}
		return
	}
}

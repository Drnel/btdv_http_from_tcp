package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Drnel/btdv_http_from_tcp/internal/request"
	"github.com/Drnel/btdv_http_from_tcp/internal/response"
	"github.com/Drnel/btdv_http_from_tcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")

}

func handler(w *response.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			Status_code: 400,
			Message: `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`,
		}
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			Status_code: 500,
			Message: `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`,
		}
	}
	err := w.WriteStatusLine(response.Ok)
	if err != nil {
		log.Println("error writing status line:", err)
	}
	message := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	default_headers := response.GetDefaultHeaders(len(message))
	default_headers["Content-Type"] = "text/html"
	err = w.WriteHeaders(default_headers)
	if err != nil {
		log.Println("error writing headers:", err)
	}
	_, err = w.WriteBody([]byte(message))
	if err != nil {
		log.Println("error writing body:", err)
	}
	return nil
}

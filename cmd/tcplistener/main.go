package main

import (
	"fmt"
	"github.com/Drnel/btdv_http_from_tcp/internal/request"
	"net"
)

func main() {
	listner, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Error creating tcp listner: %s", err)
		return
	}
	defer listner.Close()
	for {
		connection, err := listner.Accept()
		if err != nil {
			fmt.Printf("Error while waiting for connection: %s", err)
			return
		}
		fmt.Printf("A connection has been accepted\n")
		request, err := request.RequestFromReader(connection)
		if err != nil {
			fmt.Printf("Error parsing request from reader: %s", err)
			return
		}
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for field_name, field_value := range request.Headers {
			fmt.Printf("- %s: %s\n", field_name, field_value)
		}

		connection.Close()
		fmt.Printf("The connection has been closed\n")
	}
}

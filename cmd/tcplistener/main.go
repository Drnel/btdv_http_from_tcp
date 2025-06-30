package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
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
		lines := getLinesChannel(connection)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}
		connection.Close()
		fmt.Printf("The connection has been closed\n")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)
	go func() {
		defer close(channel)
		defer f.Close()
		file_read_buffer := make([]byte, 8)
		current_line := ""
		for {
			n, err := f.Read(file_read_buffer)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Printf("File read error: %s", err)
					os.Exit(1)
				}
			}
			current_line = current_line + string(file_read_buffer[:n])
			parts := strings.Split(current_line, "\n")
			if len(parts) > 1 {
				channel <- parts[0]
				current_line = "" + parts[1]
			}
		}
		if current_line != "" {
			channel <- current_line
		}

	}()
	return channel
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udp_address, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("error resolving udp address: %s", err)
		return
	}
	udp_connection, err := net.DialUDP("udp", nil, udp_address)
	if err != nil {
		fmt.Printf("Error creating udp connection: %s", err)
		return
	}
	defer udp_connection.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading line: %s", err)
		}
		_, err = udp_connection.Write([]byte(line))
		if err != nil {
			fmt.Printf("Error writing to connection: %s", err)
		}

	}
}

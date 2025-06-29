package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("File open error: %s", err)
		os.Exit(1)
	}
	defer file.Close()
	file_read_buffer := make([]byte, 8)
	for {
		_, err := file.Read(file_read_buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("File read error: %s", err)
				os.Exit(1)
			}
		}
		fmt.Printf("read: %s\n", file_read_buffer)
	}
}

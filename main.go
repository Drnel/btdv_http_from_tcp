package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("File open error: %s", err)
		os.Exit(1)
	}
	defer file.Close()
	file_read_buffer := make([]byte, 8)
	current_line := ""
	for {
		n, err := file.Read(file_read_buffer)
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
			fmt.Printf("read: %s\n", parts[0])
			current_line = "" + parts[1]
		}
	}
	if current_line != "" {
		fmt.Printf("read: %s\n", current_line)
	}
}

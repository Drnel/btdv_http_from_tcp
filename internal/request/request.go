package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request_string, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	request_lines := strings.Split(string(request_string), "\r\n")
	request_line, err := parseRequest(request_lines[0])
	if err != nil {
		return nil, err
	}
	request := Request{request_line}
	return &request, nil
}

func parseRequest(request_line string) (RequestLine, error) {
	parts := strings.Split(request_line, " ")
	if len(parts) != 3 {
		return RequestLine{}, errors.New("Error parsing request line")
	}
	rl := RequestLine{}
	rl.Method = parts[0]
	for _, char := range rl.Method {
		if char < 'A' || char > 'Z' {
			return RequestLine{}, errors.New("Got invalid method")
		}
	}
	rl.RequestTarget = parts[1]
	if parts[2] != "HTTP/1.1" {
		return RequestLine{}, errors.New("Invalid http version")
	}
	rl.HttpVersion = "1.1"
	return rl, nil
}

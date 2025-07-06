package request

import (
	"errors"
	"io"
	"strings"
)

const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	parserState int // 0 initialized | 1 done
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	request := Request{}
	request.parserState = 0
	for request.parserState != 1 {
		if readToIndex == len(buf) {
			new_buf := make([]byte, len(buf)*2, len(buf)*2)
			copy(new_buf, buf)
			buf = new_buf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.parserState = 1
				break
			}
			return &request, err
		}
		readToIndex += n
		bytes_read, err := request.parse(buf[:readToIndex])
		if err != nil {
			return &request, err
		}
		copy(buf, buf[bytes_read:])
		readToIndex = readToIndex - bytes_read
	}
	return &request, nil
}

func parseRequestLine(request_string string) (RequestLine, int, error) {
	request_lines := strings.Split(request_string, "\r\n")
	if len(request_lines) == 1 {
		return RequestLine{}, 0, nil
	}
	parts := strings.Split(request_lines[0], " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("Error parsing request line")
	}
	rl := RequestLine{}
	rl.Method = parts[0]
	for _, char := range rl.Method {
		if char < 'A' || char > 'Z' {
			return RequestLine{}, 0, errors.New("Got invalid method")
		}
	}
	rl.RequestTarget = parts[1]
	if parts[2] != "HTTP/1.1" {
		return RequestLine{}, 0, errors.New("Invalid http version")
	}
	rl.HttpVersion = "1.1"
	return rl, len(request_lines[0]), nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.parserState == 0 {
		RequestLine, bytes_read, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if bytes_read == 0 {
			return 0, nil
		}
		r.RequestLine = RequestLine
		r.parserState = 1
		return bytes_read, nil
	}
	if r.parserState == 1 {
		return 0, errors.New("error: trying to read data in a done state")
	}
	return 0, errors.New("error: unknown state")
}

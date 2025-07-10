package request

import (
	"errors"
	"github.com/Drnel/btdv_http_from_tcp/internal/headers"
	"io"
	"strings"
)

const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	ParserState ParserState // 0 initialized | 1 done
}

type ParserState int

const (
	initialized = iota
	requestStateParseHeaders
	requestStateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	request := Request{}
	request.Headers = make(headers.Headers)
	request.ParserState = initialized
	for request.ParserState != requestStateDone {
		if readToIndex == len(buf) {
			new_buf := make([]byte, len(buf)*2, len(buf)*2)
			copy(new_buf, buf)
			buf = new_buf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if request.ParserState != requestStateDone {
					return nil, errors.New("unexpected eof")
				}
				break
			}
			return &request, err
		}
		readToIndex += n
		bytes_parsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return &request, err
		}
		copy(buf, buf[bytes_parsed:])
		readToIndex = readToIndex - bytes_parsed
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
	return rl, len(request_lines[0]) + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.ParserState == initialized {
		RequestLine, bytes_parsed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if bytes_parsed == 0 {
			return 0, nil
		}
		r.RequestLine = RequestLine
		r.ParserState = requestStateParseHeaders
		return bytes_parsed, nil
	}
	if r.ParserState == requestStateParseHeaders {
		totalBytesParsed := 0
		for r.ParserState != requestStateDone {
			bytes_parsed, done, err := r.Headers.Parse(data[totalBytesParsed:])
			if err != nil {
				return totalBytesParsed, err
			}
			if bytes_parsed == 0 {
				return totalBytesParsed, nil
			}
			totalBytesParsed += bytes_parsed
			if done {
				r.ParserState = requestStateDone
			}
		}
		return totalBytesParsed, nil
	}
	if r.ParserState == requestStateDone {
		return 0, errors.New("error: trying to read data in a done state")
	}
	return 0, errors.New("error: unknown state")
}

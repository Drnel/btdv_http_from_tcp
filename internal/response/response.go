package response

import (
	"fmt"
	"io"

	"github.com/Drnel/btdv_http_from_tcp/internal/headers"
)

type StatusCode int

const (
	Ok                    = 200
	Bad_request           = 400
	Internal_server_error = 500
)

type Writer struct {
	WriterState WriterState
	Writer      io.Writer
}

type WriterState int

const (
	WriteStateStatusLine = iota
	writeHeaders
	writeBody
	done
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState == WriteStateStatusLine {
		err := WriteStatusLine(w.Writer, statusCode)
		w.WriterState = writeHeaders
		return err
	} else {
		return fmt.Errorf("Wrong state of writer: %v", w.WriterState)
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterState == writeHeaders {
		err := WriteHeaders(w.Writer, headers)
		w.WriterState = writeBody
		return err
	} else {
		return fmt.Errorf("Wrong state of writer: %v", w.WriterState)
	}
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.WriterState == writeBody {
		n, err := w.Writer.Write(p) //(w.writer, statusCode)
		w.WriterState = done
		return n, err
	} else {
		return 0, fmt.Errorf("Wrong state of writer: %v", w.WriterState)
	}
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reason_phrase := ""
	switch statusCode {
	case Ok:
		reason_phrase = "OK"
	case Bad_request:
		reason_phrase = "Bad Request"
	case Internal_server_error:
		reason_phrase = "Internal Server Error"
	default:
		reason_phrase = ""
	}
	start_line := fmt.Sprintf("HTTP/1.1 %v ", statusCode) + reason_phrase + "\r\n"
	_, err := w.Write([]byte(start_line))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	default_headers := make(headers.Headers)
	default_headers["Content-Length"] = fmt.Sprintf("%v", contentLen)
	default_headers["Connection"] = "close"
	default_headers["Content-Type"] = "text/plain"
	return default_headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}

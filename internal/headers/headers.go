package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	lines := strings.Split(string(data), "\r\n")
	if len(lines) == 1 {
		return 0, false, nil
	}
	if lines[0] == "" {
		return 2, true, nil
	}
	header := strings.TrimSpace(lines[0])
	header_parts := strings.SplitN(header, ":", 2)
	if header_parts[0] != strings.TrimSpace(header_parts[0]) {
		return 0, false, errors.New("Space between key and ':'")
	}
	h[header_parts[0]] = strings.TrimSpace(header_parts[1])
	return len(lines[0]) + 2, false, nil
}

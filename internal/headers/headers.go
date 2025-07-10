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
	if len(header_parts) != 2 {
		return 0, false, errors.New("couldnt find ':' in the field")
	}
	if header_parts[0] != strings.TrimSpace(header_parts[0]) {
		return 0, false, errors.New("Space between key and ':'")
	}
	if !valid_field_name(header_parts[0]) {
		return 0, false, errors.New("Invalid character in field_name")
	}
	value, ok := h[strings.ToLower(header_parts[0])]
	if ok {
		h[strings.ToLower(header_parts[0])] = value + ", " + strings.TrimSpace(header_parts[1])
	} else {
		h[strings.ToLower(header_parts[0])] = strings.TrimSpace(header_parts[1])
	}
	return len(lines[0]) + 2, false, nil
}

func valid_field_name(field_name string) bool {
	allowedRunes := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"
	allowedRunesMap := make(map[rune]bool)
	for _, r := range allowedRunes {
		allowedRunesMap[r] = true
	}
	for _, r := range field_name {
		if !allowedRunesMap[r] {
			return false
		}
	}
	return true
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

package jsonjoy

import (
	"errors"
	"strings"
)

// ParseJSONPointer parses JSON Pointer from canonical string form into a Go slice.
func ParseJSONPointer(str string) ([]string, error) {
	if len(str) == 0 {
		return []string{}, nil
	}
	if !strings.HasPrefix(str, "/") {
		return nil, errors.New("Invalid pointer")
	}
	return strings.Split(str[1:], "/"), nil
}

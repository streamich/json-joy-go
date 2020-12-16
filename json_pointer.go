package jsonjoy

import (
	"errors"
	"strings"
)

// DecodeReferenceToken decodes a single JSON Pointer reference token.
func DecodeReferenceToken(token string) string {
	token = strings.Replace(token, `~1`, `/`, -1)
	token = strings.Replace(token, `~0`, `~`, -1)
	return token
}

// EncodeReferenceToken encodes a single JSON Pointer reference token.
func EncodeReferenceToken(token string) string {
	token = strings.Replace(token, `~`, `~0`, -1)
	token = strings.Replace(token, `/`, `~1`, -1)
	return token
}

// ParseJSONPointer parses JSON Pointer from canonical string form into a Go slice.
func ParseJSONPointer(str string) ([]string, error) {
	if len(str) == 0 {
		return []string{}, nil
	}
	if str[0] != '/' {
		return nil, errors.New("Invalid pointer")
	}
	return strings.Split(str[1:], "/"), nil
}

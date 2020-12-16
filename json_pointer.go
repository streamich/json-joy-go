package jsonjoy

import (
	"errors"
	"strings"
)

const (
	rootPointer      = ""
	pointerSeparator = "/"
)

// JSONPointer a list of decoded JSON Pointer reference tokens.
type JSONPointer = []string

// DecodeReferenceToken decodes a single JSON Pointer reference token.
func DecodeReferenceToken(token string) string {
	token = strings.Replace(token, `~1`, pointerSeparator, -1)
	token = strings.Replace(token, `~0`, `~`, -1)
	return token
}

// EncodeReferenceToken encodes a single JSON Pointer reference token.
func EncodeReferenceToken(token string) string {
	token = strings.Replace(token, `~`, `~0`, -1)
	token = strings.Replace(token, pointerSeparator, `~1`, -1)
	return token
}

// ParseJSONPointer parses JSON Pointer from canonical string form into a Go
// slice of decoded tokens.
func ParseJSONPointer(str string) (JSONPointer, error) {
	if len(str) == 0 {
		return []string{}, nil
	}
	if str[0] != '/' {
		return nil, errors.New("Invalid pointer")
	}
	tokens := strings.Split(str[1:], pointerSeparator)
	for index, token := range tokens {
		tokens[index] = DecodeReferenceToken(token)
	}
	return tokens, nil
}

// FormatJSONPointer formats JSON Pointer tokens into the canonical string form.
func FormatJSONPointer(tokens JSONPointer) string {
	if len(tokens) == 0 {
		return rootPointer
	}
	encoded := make([]string, len(tokens))
	for index, token := range tokens {
		encoded[index] = EncodeReferenceToken(token)
	}
	return pointerSeparator + strings.Join(encoded, pointerSeparator)
}

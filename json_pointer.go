package jsonjoy

import (
	"errors"
	"strings"
)

// JSONValue represents a type of JSON value.
type JSONValue int

// JSON value types.
const (
	JSONObject JSONValue = iota
	JSONArray
	JSONPrimitive
)

const (
	rootPointer    = ""
	tokenSeparator = "/"
	escapeChar     = "~"
)

// JSONPointer a list of decoded JSON Pointer reference tokens.
type JSONPointer []string

// ErrNotFound is returned when JSONPointer.Find cannot locate a value.
var ErrNotFound = errors.New("not found")

// DecodeReferenceToken decodes a single JSON Pointer reference token.
func DecodeReferenceToken(token string) string {
	token = strings.Replace(token, `~1`, tokenSeparator, -1)
	token = strings.Replace(token, `~0`, escapeChar, -1)
	return token
}

// EncodeReferenceToken encodes a single JSON Pointer reference token.
func EncodeReferenceToken(token string) string {
	token = strings.Replace(token, escapeChar, `~0`, -1)
	token = strings.Replace(token, tokenSeparator, `~1`, -1)
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
	tokens := strings.Split(str[1:], tokenSeparator)
	for index, token := range tokens {
		tokens[index] = DecodeReferenceToken(token)
	}
	return tokens, nil
}

// Format formats JSON Pointer tokens into the canonical string form.
func (tokens JSONPointer) Format() string {
	if len(tokens) == 0 {
		return rootPointer
	}
	encoded := make([]string, len(tokens))
	for index, token := range tokens {
		encoded[index] = EncodeReferenceToken(token)
	}
	return tokenSeparator + strings.Join(encoded, tokenSeparator)
}

// Find a value in parsed JSON document.
func (tokens JSONPointer) Find(doc interface{}) (interface{}, error) {
	if len(tokens) == 0 {
		return doc, nil
	}
	parent := doc
	child := doc
	for _, token := range tokens {
		parent = child
		switch typedParent := parent.(type) {
		case map[string]interface{}:
			if val, ok := typedParent[token]; ok {
				child = val
				continue
			}
			return nil, ErrNotFound
		case []interface{}:
			return nil, ErrNotFound
		default:
			return nil, ErrNotFound
		}
	}
	return child, nil
}

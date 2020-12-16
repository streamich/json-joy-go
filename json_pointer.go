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
	rootPointer           = ""
	tokenSeparator        = "/"
	escapeChar            = "~"
	tokenSeparatorEncoded = "~1"
	escapeCharEncoded     = "~0"
)

// JSONPointer a list of decoded JSON Pointer reference tokens.
type JSONPointer []string

// ErrNotFound is returned when JSONPointer.Find cannot locate a value.
var ErrNotFound = errors.New("not found")

// DecodeReferenceToken decodes a single JSON Pointer reference token.
func DecodeReferenceToken(token string) string {
	token = strings.Replace(token, tokenSeparatorEncoded, tokenSeparator, -1)
	token = strings.Replace(token, escapeCharEncoded, escapeChar, -1)
	return token
}

// EncodeReferenceToken encodes a single JSON Pointer reference token.
func EncodeReferenceToken(token string) string {
	token = strings.Replace(token, escapeChar, escapeCharEncoded, -1)
	token = strings.Replace(token, tokenSeparator, tokenSeparatorEncoded, -1)
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

// Get a specific value from JSON document identified by a JSON Pointer.
func (tokens JSONPointer) Get(doc interface{}) (interface{}, error) {
	if len(tokens) == 0 {
		return doc, nil
	}
	var key string
	for _, token := range tokens {
		key = token
		switch typedParent := doc.(type) {
		case map[string]interface{}:
			if child, ok := typedParent[key]; ok {
				doc = child
				continue
			}
			return nil, ErrNotFound
		case []interface{}:
			return nil, ErrNotFound
		default:
			return nil, ErrNotFound
		}
	}
	return doc, nil
}

// Resolve all values of a JSON document on the path of a JSON Pointer. Each
// entry in the return list is a value that corresponds to a JSON Pointer token
// with the same index. Returns nil of the JSON Pointer points to root.
func (tokens JSONPointer) Resolve(doc interface{}) ([]interface{}, error) {
	if len(tokens) == 0 {
		return nil, nil
	}
	values := make([]interface{}, len(tokens))
	val := doc
	var key string
	for index, token := range tokens {
		key = token
		switch typedParent := val.(type) {
		case map[string]interface{}:
			if child, ok := typedParent[key]; ok {
				val = child
				values[index] = val
				continue
			}
			return nil, ErrNotFound
		case []interface{}:
			return nil, ErrNotFound
		default:
			return nil, ErrNotFound
		}
	}
	return values, nil
}

package jsonjoy

import (
	"errors"
	"strconv"
	"strings"
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

// ErrInvalidIndex is returned when JSON Pointer array index is not valid.
var ErrInvalidIndex = errors.New("invalid index")

// UnescapeReferenceToken decodes a single JSON Pointer reference token.
func UnescapeReferenceToken(token string) string {
	token = strings.Replace(token, tokenSeparatorEncoded, tokenSeparator, -1)
	token = strings.Replace(token, escapeCharEncoded, escapeChar, -1)
	return token
}

// EscapeReferenceToken encodes a single JSON Pointer reference token.
func EscapeReferenceToken(token string) string {
	token = strings.Replace(token, escapeChar, escapeCharEncoded, -1)
	token = strings.Replace(token, tokenSeparator, tokenSeparatorEncoded, -1)
	return token
}

// NewJSONPointer parses JSON Pointer from canonical string form into a Go
// slice of decoded tokens.
func NewJSONPointer(str string) (JSONPointer, error) {
	if len(str) == 0 {
		return []string{}, nil
	}
	if str[0] != '/' {
		return nil, errors.New("Invalid pointer")
	}
	tokens := strings.Split(str[1:], tokenSeparator)
	for index, token := range tokens {
		tokens[index] = UnescapeReferenceToken(token)
	}
	return tokens, nil
}

// ParseTokenAsArrayIndex parses JSON Pointer reference token to an integer,
// which can be used as array index.
func ParseTokenAsArrayIndex(token string, arr []interface{}) (int, error) {
	index, err := strconv.Atoi(token)
	if err != nil {
		return 0, ErrInvalidIndex
	}
	if index < 0 {
		return 0, ErrInvalidIndex
	}
	if arr != nil {
		if index >= len(arr) {
			return 0, ErrInvalidIndex
		}
	}
	return index, nil
}

// IsRoot returns true if JSON Pointer points to the root of a document.
func (tokens JSONPointer) IsRoot() bool {
	return len(tokens) == 0
}

// Format formats JSON Pointer tokens into the canonical string form.
func (tokens JSONPointer) Format() string {
	if tokens.IsRoot() {
		return rootPointer
	}
	encoded := make([]string, len(tokens))
	for index, token := range tokens {
		encoded[index] = EscapeReferenceToken(token)
	}
	return tokenSeparator + strings.Join(encoded, tokenSeparator)
}

// Get a specific value from JSON document identified by a JSON Pointer.
func (tokens JSONPointer) Get(doc JSON) (JSON, error) {
	if tokens.IsRoot() {
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
			tokenIndex, err := ParseTokenAsArrayIndex(token, typedParent)
			if err != nil {
				return nil, err
			}
			doc = typedParent[tokenIndex]
		default:
			return nil, ErrNotFound
		}
	}
	return doc, nil
}

// Resolve all values of a JSON document on the path of a JSON Pointer. Each
// entry in the return list corresponds to a JSON Pointer token reference
// with the same index. Returns nil if the JSON Pointer points to root.
func (tokens JSONPointer) Resolve(doc JSON) ([]JSON, error) {
	if tokens.IsRoot() {
		return nil, nil
	}
	values := make([]JSON, len(tokens))
	var key string
	for index, token := range tokens {
		key = token
		switch typedParent := doc.(type) {
		case map[string]interface{}:
			if child, ok := typedParent[key]; ok {
				doc = child
				values[index] = doc
				continue
			}
			return nil, ErrNotFound
		case []interface{}:
			tokenIndex, err := ParseTokenAsArrayIndex(token, typedParent)
			if err != nil {
				return nil, err
			}
			doc = typedParent[tokenIndex]
			values[index] = doc
		default:
			return nil, ErrNotFound
		}
	}
	return values, nil
}

// Locate finds an object or an array which contains value located by JSON
// pointer and returns that object as well as the last reference token.
func (tokens JSONPointer) Locate(doc JSON) (JSON, *string, error) {
	if tokens.IsRoot() {
		return nil, nil, nil
	}
	obj := doc
	var key string
	for _, token := range tokens {
		key = token
		switch typedParent := doc.(type) {
		case map[string]interface{}:
			if child, ok := typedParent[key]; ok {
				obj = doc
				doc = child
				continue
			}
			return nil, nil, ErrNotFound
		case []interface{}:
			tokenIndex, err := ParseTokenAsArrayIndex(token, typedParent)
			if err != nil {
				return nil, nil, err
			}
			obj = doc
			doc = typedParent[tokenIndex]
		default:
			return nil, nil, ErrNotFound
		}
	}
	return obj, &key, nil
}

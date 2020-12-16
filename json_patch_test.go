package jsonjoy

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Add_CanModifyKeyValue(t *testing.T) {
	b := []byte(`{"foo":"bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	doc2, _ := Add(doc, JSONPointer{"foo"}, "baz")
	assert.Equal(t, "map[foo:baz]", fmt.Sprint(doc))
	assert.Equal(t, "map[foo:baz]", fmt.Sprint(doc2))
}

func Test_Add_CanModifyKeyValueWithoutModifyingOriginal(t *testing.T) {
	b := []byte(`{"foo":"bar"}`)
	var doc JSON
	json.Unmarshal(b, &doc)
	doc1 := Copy(doc)
	doc2, _ := Add(doc1, JSONPointer{"foo"}, "baz")
	assert.Equal(t, "map[foo:bar]", fmt.Sprint(doc))
	assert.Equal(t, "map[foo:baz]", fmt.Sprint(doc1))
	assert.Equal(t, "map[foo:baz]", fmt.Sprint(doc2))
}

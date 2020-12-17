package jsonjoy

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JsonPatch_insert_CanInsertIntoMiddleOfSlice(t *testing.T) {
	b := []byte(`[1, 2]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	s, _ := doc.([]interface{})
	newSlice := insert(s, 1, 3)
	assert.Equal(t, 3, len(newSlice))
	assert.Equal(t, "[1 3 2]", fmt.Sprint(newSlice))
}

func Test_JsonPatch_insert_CanInsertInFrontOfSlice(t *testing.T) {
	b := []byte(`[1, 2]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	s, _ := doc.([]interface{})
	newSlice := insert(s, 0, 3)
	assert.Equal(t, 3, len(newSlice))
	assert.Equal(t, "[3 1 2]", fmt.Sprint(newSlice))
}

func Test_JsonPatch_insert_CanInsertAtTheBackOfSlice(t *testing.T) {
	b := []byte(`[1, 2]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	s, _ := doc.([]interface{})
	newSlice := insert(s, 2, 3)
	assert.Equal(t, 3, len(newSlice))
	assert.Equal(t, "[1 2 3]", fmt.Sprint(newSlice))
}

func Test_JsonPatch_insert_CanInsertStringAtTheBackOfSlice(t *testing.T) {
	b := []byte(`[1, 2]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	s, _ := doc.([]interface{})
	var str interface{} = "abc"
	newSlice := insert(s, 2, str)
	assert.Equal(t, 3, len(newSlice))
	assert.Equal(t, "[1 2 abc]", fmt.Sprint(newSlice))
}

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

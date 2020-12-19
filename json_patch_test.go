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

func Test_JsonPatch_insert_CanInsertIntoEmptyArray(t *testing.T) {
	b := []byte(`[]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	s, _ := doc.([]interface{})
	var str interface{} = "abcd"
	newSlice := insert(s, 0, str)
	assert.Equal(t, 1, len(newSlice))
	assert.Equal(t, "[abcd]", fmt.Sprint(newSlice))
}

func Test_JsonPatch_Add_CanModifyKeyValueWithoutModifyingOriginal(t *testing.T) {
	b := []byte(`{"foo":"bar"}`)
	var doc JSON
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	Add(&doc2, JSONPointer{"foo"}, "baz")
	assert.Equal(t, "map[foo:bar]", fmt.Sprint(doc))
	assert.Equal(t, "map[foo:baz]", fmt.Sprint(doc2))
}

func Test_JsonPatch_Add_CanInsertItemIntoAnArray(t *testing.T) {
	b := []byte(`[1, 2]`)
	var doc JSON
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	Add(&doc2, JSONPointer{"1"}, 3)
	assert.Equal(t, "[1 2]", fmt.Sprint(doc))
	assert.Equal(t, "[1 3 2]", fmt.Sprint(doc2))
}

func Test_JsonPatch_Add_CanInsertItemIntoAnArrayAtSecondLevel(t *testing.T) {
	b := []byte(`{"a": [1, 2]}`)
	var doc JSON
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	Add(&doc2, JSONPointer{"a", "1"}, 3)
	assert.Equal(t, "map[a:[1 2]]", fmt.Sprint(doc))
	assert.Equal(t, "map[a:[1 3 2]]", fmt.Sprint(doc2))
}

func Test_JsonPatch_Add_CanInsertItemIntoAnEmptyArray(t *testing.T) {
	b := []byte(`{"a": []}`)
	var doc JSON
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	Add(&doc2, JSONPointer{"a", "0"}, "asdf")
	assert.Equal(t, "map[a:[]]", fmt.Sprint(doc))
	assert.Equal(t, "map[a:[asdf]]", fmt.Sprint(doc2))
}

func Test_JsonPatch_Add_CanInsertItemsAtTheEndOfArray(t *testing.T) {
	b := []byte(`{"a": []}`)
	var doc JSON
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	Add(&doc2, JSONPointer{"a", "-"}, 3)
	Add(&doc2, JSONPointer{"a", "-"}, 4)
	Add(&doc2, JSONPointer{"a", "-"}, 5)
	assert.Equal(t, "map[a:[]]", fmt.Sprint(doc))
	assert.Equal(t, "map[a:[3 4 5]]", fmt.Sprint(doc2))
}

func Test_JsonPatch_Add_CanInsertItemsAtTheEndOfArrayWhenArrayIsRoot(t *testing.T) {
	b := []byte(`[]`)
	var doc JSON
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	Add(&doc2, JSONPointer{"-"}, 1)
	Add(&doc2, JSONPointer{"-"}, 1)
	Add(&doc2, JSONPointer{"-"}, 2)
	assert.Equal(t, "[]", fmt.Sprint(doc))
	assert.Equal(t, "[1 1 2]", fmt.Sprint(doc2))
}

func Test_JsonPatch_Add_ReturnsErrorWhenInsertingPastArrayLength(t *testing.T) {
	b := []byte(`[]`)
	var doc JSON
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	err := Add(&doc2, JSONPointer{"1"}, nil)
	assert.Equal(t, "[]", fmt.Sprint(doc))
	assert.Equal(t, "INVALID_INDEX", fmt.Sprint(err))
}

func Test_JsonPatch_ApplyOps_AppliesOperations(t *testing.T) {
	b1 := []byte(`{
		"foo": "bar"
	}`)
	b2 := []byte(`[
		{"op": "replace", "path": "/foo", "value": "baz"},
		{"op": "add", "path": "/gg", "value": [123]}
		]`)
	var doc interface{}
	var patch interface{}
	json.Unmarshal(b1, &doc)
	json.Unmarshal(b2, &patch)
	ops, _, _ := CreateOps(patch)
	ApplyOps(&doc, ops)
	assert.Equal(t, "map[foo:baz gg:[123]]", fmt.Sprint(doc))
}

func Test_JsonPatch_ApplyOps_AppliesRemoveOperation(t *testing.T) {
	b1 := []byte(`{
		"foo": "bar",
		"baz": "qux"
	}`)
	b2 := []byte(`[
		{"op": "remove", "path": "/foo"}
	]`)
	var doc interface{}
	var patch interface{}
	json.Unmarshal(b1, &doc)
	json.Unmarshal(b2, &patch)
	ops, _, _ := CreateOps(patch)
	ApplyOps(&doc, ops)
	assert.Equal(t, "map[baz:qux]", fmt.Sprint(doc))
}

func Test_JsonPatch_ApplyOps_AppliesMoveOperation(t *testing.T) {
	b1 := []byte(`{
		"foo": "bar"
	}`)
	b2 := []byte(`[
		{"op": "move", "path": "/a", "from": "/foo"}
	]`)
	var doc interface{}
	var patch interface{}
	json.Unmarshal(b1, &doc)
	json.Unmarshal(b2, &patch)
	ops, _, _ := CreateOps(patch)
	ApplyOps(&doc, ops)
	assert.Equal(t, "map[a:bar]", fmt.Sprint(doc))
}

func Test_JsonPatch_ApplyOps_AppliesCopyOperation(t *testing.T) {
	b1 := []byte(`{
		"foo": "bar"
	}`)
	b2 := []byte(`[
		{"op": "copy", "path": "/baz", "from": "/foo"}
	]`)
	var doc interface{}
	var patch interface{}
	json.Unmarshal(b1, &doc)
	json.Unmarshal(b2, &patch)
	ops, _, _ := CreateOps(patch)
	ApplyOps(&doc, ops)
	m := doc.(map[string]interface{})
	assert.Equal(t, "bar", m["foo"])
	assert.Equal(t, "bar", m["baz"])
}

func Test_JsonPatch_ApplyOps_TwoTestOperations(t *testing.T) {
	b1 := []byte(`{
		"baz": "qux",
	  	"foo": ["a", 2, "c"]
	}`)
	b2 := []byte(`[
		{"op": "test", "path": "/baz", "value": "qux"},
	  	{"op": "test", "path": "/foo/1", value: 2}
	]`)
	var doc interface{}
	var patch interface{}
	json.Unmarshal(b1, &doc)
	json.Unmarshal(b2, &patch)
	ops, _, _ := CreateOps(patch)
	err := ApplyOps(&doc, ops)
	m := doc.(map[string]interface{})
	assert.Nil(t, err)
	assert.Equal(t, "qux", m["baz"])
	assert.Equal(t, "[a 2 c]", fmt.Sprint(m["foo"]))
}

func Test_JsonPatch_ApplyOps_TestOperation(t *testing.T) {
	b1 := []byte(`{
		"foo": "bar"
	}`)
	b2 := []byte(`[
		{"op": "test", "path": "/foo", "value": "bar"}
	]`)
	var doc interface{}
	var patch interface{}
	json.Unmarshal(b1, &doc)
	json.Unmarshal(b2, &patch)
	ops, _, _ := CreateOps(patch)
	err := ApplyOps(&doc, ops)
	m := doc.(map[string]interface{})
	assert.Nil(t, err)
	assert.Equal(t, "bar", m["foo"])
}

func Test_JsonPatch_ApplyOps_NotFoundPath(t *testing.T) {
	b1 := []byte(`{
		"q": {"bar": 2}
	}`)
	b2 := []byte(`[
		{"op": "add", "path": "/a/b", "value": 1}
	]`)
	var doc interface{}
	var patch interface{}
	json.Unmarshal(b1, &doc)
	json.Unmarshal(b2, &patch)
	ops, _, _ := CreateOps(patch)
	err := ApplyOps(&doc, ops)
	assert.Equal(t, "NOT_FOUND", fmt.Sprint(err))
}

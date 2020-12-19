package jsonjoy

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JsonPatchOperations_CreateOps_ReturnsErrorOnInvalidPatch(t *testing.T) {
	b := []byte(`{}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	_, index, err := CreateOps(doc)
	assert.Equal(t, -1, index)
	assert.Equal(t, ErrPatchInvalid, err)
}

func Test_JsonPatchOperations_CreateOps_EmptyPatchIsNotAnError(t *testing.T) {
	b := []byte(`[]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, index, err := CreateOps(patch)
	assert.Equal(t, -1, index)
	assert.Nil(t, err)
}

func Test_JsonPatchOperations_CreateOps_ReturnsErrorOnInvalidOperation(t *testing.T) {
	b := []byte(`[123]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, index, err := CreateOps(patch)
	assert.Equal(t, 0, index)
	assert.Equal(t, ErrOperationInvalid, err)
}

func Test_JsonPatchOperations_CreateOps_ReturnsErrorOnMissingOpField(t *testing.T) {
	b := []byte(`[{}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, index, err := CreateOps(patch)
	assert.Equal(t, 0, index)
	assert.Equal(t, ErrOperationInvalid, err)
}

func Test_JsonPatchOperations_CreateOps_ReturnsErrorOnInvalidOpField(t *testing.T) {
	b := []byte(`[{"op": 123}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, index, err := CreateOps(patch)
	assert.Equal(t, 0, index)
	assert.Equal(t, ErrOperationInvalid, err)
}

func Test_JsonPatchOperations_CreateOps_ReturnsErrorOnUnknownOperation(t *testing.T) {
	b := []byte(`[{"op": "unknown_test_op"}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, index, err := CreateOps(patch)
	assert.Equal(t, 0, index)
	assert.Equal(t, ErrOperationUnknown, err)
}

func Test_JsonPatchOperations_CreateOps_ReturnsErrorOnMissingPathInAddOperation(t *testing.T) {
	b := []byte(`[{"op": "add"}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, index, err := CreateOps(patch)
	assert.Equal(t, 0, index)
	assert.Equal(t, ErrOperationMissingPath, err)
}

func Test_JsonPatchOperations_CreateOps_ReturnsErrorOnInvalidAddOperationPath(t *testing.T) {
	b := []byte(`[{"op": "add", "path": 123}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, index, err := CreateOps(patch)
	assert.Equal(t, 0, index)
	assert.Equal(t, ErrOperationInvalidPath, err)
}

func Test_JsonPatchOperations_CreateOps_ReturnsErrorOnInvalidAddOperationPathPointer(t *testing.T) {
	b := []byte(`[{"op": "add", "path": "asdf/adsf"}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, index, err := CreateOps(patch)
	assert.Equal(t, 0, index)
	assert.Equal(t, ErrPointerInvalid, err)
}

func Test_JsonPatchOperations_CreateOps_ReturnsErrorOnAddOperationMissingValueField(t *testing.T) {
	b := []byte(`[{"op": "add", "path": "/adsf"}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, index, err := CreateOps(patch)
	assert.Equal(t, 0, index)
	assert.Equal(t, ErrOperationMissingValue, err)
}

func Test_JsonPatchOperations_CreateOps_ReturnsAddOpOnSuccess(t *testing.T) {
	b := []byte(`[{"op": "add", "path": "/foo/bar/baz", "value": {"a": "b"}}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	ops, index, err := CreateOps(patch)
	assert.Nil(t, err)
	assert.Equal(t, -1, index)
	assert.NotNil(t, ops)
	assert.Equal(t, 1, len(ops))
	op, ok := ops[0].(*OpAdd)
	assert.Equal(t, true, ok)
	assert.Equal(t, 3, len(op.path))
	assert.Equal(t, "foo", op.path[0])
	assert.Equal(t, "bar", op.path[1])
	assert.Equal(t, "baz", op.path[2])
	assert.Equal(t, "map[a:b]", fmt.Sprint(op.value))
}

func Test_JsonPatchOperations_CreateOps_CanCreateMultipleOperations(t *testing.T) {
	b := []byte(`[
		{"op": "add", "path": "/foo/bar/baz", "value": {"a": "b"}},
		{"op": "replace", "path": "/a/1", "value": {"a": "c"}},
		{"op": "test", "path": "", "value": 123},
		{"op": "remove", "path": "/asdf"},
		{"op": "move", "path": "/a/b", "from": "/c/d"},
		{"op": "copy", "path": "/1/2", "from": "/3/4"}
	]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	ops, index, err := CreateOps(patch)
	assert.Nil(t, err)
	assert.Equal(t, -1, index)
	assert.NotNil(t, ops)
	assert.Equal(t, 6, len(ops))

	op1, ok1 := ops[0].(*OpAdd)
	assert.Equal(t, true, ok1)
	assert.Equal(t, 3, len(op1.path))
	assert.Equal(t, "foo", op1.path[0])
	assert.Equal(t, "bar", op1.path[1])
	assert.Equal(t, "baz", op1.path[2])
	assert.Equal(t, "map[a:b]", fmt.Sprint(op1.value))

	op2, ok2 := ops[1].(*OpReplace)
	assert.Equal(t, true, ok2)
	assert.Equal(t, 2, len(op2.path))
	assert.Equal(t, "a", op2.path[0])
	assert.Equal(t, "1", op2.path[1])
	assert.Equal(t, "map[a:c]", fmt.Sprint(op2.value))

	op3, ok3 := ops[2].(*OpTest)
	assert.Equal(t, true, ok3)
	assert.Equal(t, 0, len(op3.path))
	assert.Equal(t, 123.0, op3.value)

	op4, ok4 := ops[3].(*OpRemove)
	assert.Equal(t, true, ok4)
	assert.Equal(t, 1, len(op4.path))
	assert.Equal(t, "asdf", op4.path[0])

	op5, ok5 := ops[4].(*OpMove)
	assert.Equal(t, true, ok5)
	assert.Equal(t, 2, len(op5.path))
	assert.Equal(t, "a", op5.path[0])
	assert.Equal(t, "b", op5.path[1])
	assert.Equal(t, 2, len(op5.from))
	assert.Equal(t, "c", op5.from[0])
	assert.Equal(t, "d", op5.from[1])

	op6, ok6 := ops[5].(*OpCopy)
	assert.Equal(t, true, ok6)
	assert.Equal(t, 2, len(op6.path))
	assert.Equal(t, "1", op6.path[0])
	assert.Equal(t, "2", op6.path[1])
	assert.Equal(t, 2, len(op6.from))
	assert.Equal(t, "3", op6.from[0])
	assert.Equal(t, "4", op6.from[1])
}

package jsonjoy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JsonPatchOperations_ValidateOperations_ReturnsErrorOnInvalidPatch(t *testing.T) {
	b := []byte(`{}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	_, err := ValidateOperations(doc)
	assert.Equal(t, ErrPatchInvalid, err)
}

func Test_JsonPatchOperations_ValidateOperations_ReturnsErrorOnEmptyPatch(t *testing.T) {
	b := []byte(`[]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, err := ValidateOperations(patch)
	assert.Equal(t, ErrPatchEmpty, err)
}

func Test_JsonPatchOperations_ValidateOperations_ReturnsErrorOnInvalidOperation(t *testing.T) {
	b := []byte(`[123]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, err := ValidateOperations(patch)
	assert.Equal(t, ErrOperationInvalid, err)
}

func Test_JsonPatchOperations_ValidateOperations_ReturnsErrorOnMissingOpField(t *testing.T) {
	b := []byte(`[{}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, err := ValidateOperations(patch)
	assert.Equal(t, ErrOperationInvalid, err)
}

func Test_JsonPatchOperations_ValidateOperations_ReturnsErrorOnInvalidOpField(t *testing.T) {
	b := []byte(`[{"op": 123}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, err := ValidateOperations(patch)
	assert.Equal(t, ErrOperationInvalid, err)
}

func Test_JsonPatchOperations_ValidateOperations_ReturnsErrorOnUnknownOperation(t *testing.T) {
	b := []byte(`[{"op": "unknown_test_op"}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, err := ValidateOperations(patch)
	assert.Equal(t, ErrOperationUnknown, err)
}

func Test_JsonPatchOperations_ValidateOperations_ReturnsErrorOnMissingPathInAddOperation(t *testing.T) {
	b := []byte(`[{"op": "add"}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, err := ValidateOperations(patch)
	assert.Equal(t, ErrOperationMissingPath, err)
}

func Test_JsonPatchOperations_ValidateOperations_ReturnsErrorOnInvalidAddOperationPath(t *testing.T) {
	b := []byte(`[{"op": "add", "path": 123}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, err := ValidateOperations(patch)
	assert.Equal(t, ErrOperationInvalidPath, err)
}

func Test_JsonPatchOperations_ValidateOperations_ReturnsErrorOnInvalidAddOperationPathPointer(t *testing.T) {
	b := []byte(`[{"op": "add", "path": "asdf/adsf"}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, err := ValidateOperations(patch)
	assert.Equal(t, ErrPointerInvalid, err)
}

func Test_JsonPatchOperations_ValidateOperations_ReturnsErrorOnAddOperationMissingValueField(t *testing.T) {
	b := []byte(`[{"op": "add", "path": "/adsf"}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	patch, _ := doc.([]JSON)
	_, err := ValidateOperations(patch)
	assert.Equal(t, ErrOperationMissingValue, err)
}

package jsonjoy

import (
	"errors"
)

// OpAdd JSON Patch "add" operation.
type OpAdd struct {
	path  JSONPointer
	value JSON
}

// OpRemove JSON Patch "remove" operation.
type OpRemove struct {
	path JSONPointer
}

// OpReplace JSON Patch "replace" operation.
type OpReplace struct {
	path  JSONPointer
	value JSON
}

// OpMove JSON Patch "move" operation.
type OpMove struct {
	path JSONPointer
	from JSONPointer
}

// OpCopy JSON Patch "copy" operation.
type OpCopy struct {
	path JSONPointer
	from JSONPointer
}

// OpTest JSON Patch "test" operation.
type OpTest struct {
	path  JSONPointer
	value JSON
	not   bool
}

// ErrPatchInvalid returned when JSON Patch is invalid.
var ErrPatchInvalid = errors.New("patch invalid")

// ErrPatchEmpty returned when JSON Patch array is empty.
var ErrPatchEmpty = errors.New("patch empty")

// CreateOps validates a list of JSON Patch operations.
func CreateOps(patch JSON) ([]interface{}, int, error) {
	arr, ok := patch.([]JSON)
	if !ok {
		return nil, -1, ErrPatchInvalid
	}
	length := len(arr)
	if length == 0 {
		return nil, -1, ErrPatchEmpty
	}
	ops := make([]interface{}, length)
	for index, operation := range arr {
		op, err := CreateOp(operation)
		if err != nil {
			return nil, index, err
		}
		ops[index] = op
	}
	return ops, -1, nil
}

// ErrOperationInvalid returned when JSON Patch operation is invalid.
var ErrOperationInvalid = errors.New("operation invalid")

// ErrOperationUnknown returned when JSON Patch operation opcode is not recognized.
var ErrOperationUnknown = errors.New("operation unknown")

// CreateOp validates a single JSON Patch operation.
func CreateOp(operation JSON) (interface{}, error) {
	obj, ok := operation.(map[string]JSON)
	if !ok {
		return nil, ErrOperationInvalid
	}
	opInterface, ok := obj["op"]
	if !ok {
		return nil, ErrOperationInvalid
	}
	op, ok := opInterface.(string)
	if !ok {
		return nil, ErrOperationInvalid
	}
	switch op {
	case "add":
		return createAddOp(obj)
	default:
		return nil, ErrOperationUnknown
	}
}

// ErrOperationMissingPath returned when JSON Patch operation is missing the "path" field.
var ErrOperationMissingPath = errors.New("op_missing_path")

// ErrOperationInvalidPath returned when operation "path" field is invalid.
var ErrOperationInvalidPath = errors.New("op_invalid_path")

// ErrOperationMissingValue returned when operation is missing "value" field.
var ErrOperationMissingValue = errors.New("op_missing_value")

func createAddOp(operation map[string]JSON) (interface{}, error) {
	pathInterface, ok := operation["path"]
	if !ok {
		return nil, ErrOperationMissingPath
	}
	pathString, ok := pathInterface.(string)
	if !ok {
		return nil, ErrOperationInvalidPath
	}
	path, err := NewJSONPointer(pathString)
	if err != nil {
		return nil, err
	}
	value, ok := operation["value"]
	if !ok {
		return nil, ErrOperationMissingValue
	}
	return OpAdd{path: path, value: value}, nil
}

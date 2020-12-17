package jsonjoy

import (
	"errors"
)

// OpAdd JSON Patch "add" operation.
type OpAdd struct {
	operation *map[string]JSON
	path      JSONPointer
	value     JSON
}

// OpRemove JSON Patch "remove" operation.
type OpRemove struct {
	operation *map[string]JSON
	path      JSONPointer
}

// OpReplace JSON Patch "replace" operation.
type OpReplace struct {
	operation *map[string]JSON
	path      JSONPointer
	value     JSON
}

// OpMove JSON Patch "move" operation.
type OpMove struct {
	operation *map[string]JSON
	path      JSONPointer
	from      JSONPointer
}

// OpCopy JSON Patch "copy" operation.
type OpCopy struct {
	operation *map[string]JSON
	path      JSONPointer
	from      JSONPointer
}

// OpTest JSON Patch "test" operation.
type OpTest struct {
	operation *map[string]JSON
	path      JSONPointer
	value     JSON
	not       bool
}

// ErrPatchInvalid returned when JSON Patch is invalid.
var ErrPatchInvalid = errors.New("patch invalid")

// ErrPatchEmpty returned when JSON Patch array is empty.
var ErrPatchEmpty = errors.New("patch empty")

// CreateOps validates a list of JSON Patch operations and returns a list of
// Op* structs. Second return argument integer represents operation in which
// error happened, or is set to -1 if validation error did not happen in an operation.
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
	case "replace":
		return createReplaceOp(obj)
	case "remove":
		return createRemoveOp(obj)
	case "move":
		return createMoveOp(obj)
	case "copy":
		return createCopyOp(obj)
	case "test":
		return createTestOp(obj)
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

func createAddOp(operation map[string]JSON) (*OpAdd, error) {
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
	op := OpAdd{operation: &operation, path: path, value: value}
	return &op, nil
}

func createReplaceOp(operation map[string]JSON) (*OpReplace, error) {
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
	op := OpReplace{operation: &operation, path: path, value: value}
	return &op, nil
}

func createTestOp(operation map[string]JSON) (*OpTest, error) {
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
	op := OpTest{operation: &operation, path: path, value: value}
	return &op, nil
}

func createRemoveOp(operation map[string]JSON) (*OpRemove, error) {
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
	op := OpRemove{operation: &operation, path: path}
	return &op, nil
}

func createMoveOp(operation map[string]JSON) (*OpMove, error) {
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
	fromInterface, ok := operation["from"]
	if !ok {
		return nil, ErrOperationMissingPath
	}
	fromString, ok := fromInterface.(string)
	if !ok {
		return nil, ErrOperationInvalidPath
	}
	from, err := NewJSONPointer(fromString)
	if err != nil {
		return nil, err
	}
	op := OpMove{operation: &operation, path: path, from: from}
	return &op, nil
}

func createCopyOp(operation map[string]JSON) (*OpCopy, error) {
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
	fromInterface, ok := operation["from"]
	if !ok {
		return nil, ErrOperationMissingPath
	}
	fromString, ok := fromInterface.(string)
	if !ok {
		return nil, ErrOperationInvalidPath
	}
	from, err := NewJSONPointer(fromString)
	if err != nil {
		return nil, err
	}
	op := OpCopy{operation: &operation, path: path, from: from}
	return &op, nil
}

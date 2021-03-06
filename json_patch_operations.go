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

// OpStrIns JSON Patch+ "str_ins" operation.
type OpStrIns struct {
	operation *map[string]JSON
	path      JSONPointer
	pos       int
	str       string
}

// OpStrDel JSON Patch+ "str_del" operation.
type OpStrDel struct {
	operation *map[string]JSON
	path      JSONPointer
	pos       int
	len       int
	str       string
}

// OpFlip JSON Patch+ "flip" operation.
type OpFlip struct {
	operation *map[string]JSON
	path      JSONPointer
}

// OpInc JSON Patch+ "inc" operation.
type OpInc struct {
	operation *map[string]JSON
	path      JSONPointer
	inc       float64
}

// ErrPatchInvalid returned when JSON Patch is invalid.
var ErrPatchInvalid = errors.New("PATCH_INVALID")

// ErrPatchEmpty returned when JSON Patch array is empty.
var ErrPatchEmpty = errors.New("PATCH_EMPTY")

// CreateOps validates a list of JSON Patch operations and returns a list of
// Op* structs. Second return argument integer represents operation in which
// error happened, or is set to -1 if validation error did not happen in an operation.
func CreateOps(patch JSON) ([]interface{}, int, error) {
	arr, ok := patch.([]JSON)
	if !ok {
		return nil, -1, ErrPatchInvalid
	}
	length := len(arr)
	// if length == 0 {
	// 	return nil, -1, ErrPatchEmpty
	// }
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
var ErrOperationInvalid = errors.New("OP_INVALID")

// ErrOperationUnknown returned when JSON Patch operation opcode is not recognized.
var ErrOperationUnknown = errors.New("OP_UNKNOWN")

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
	case "str_ins":
		return createStrInsOp(obj)
	case "str_del":
		return createStrDelOp(obj)
	case "flip":
		return createFlipOp(obj)
	case "inc":
		return createIncOp(obj)
	default:
		return nil, ErrOperationUnknown
	}
}

// ErrOperationInvalidPath returned when operation "path" field is invalid.
var ErrOperationInvalidPath = errors.New("OP_PATH_INVALID")

// ErrOperationMissingValue returned when operation is missing "value" field.
var ErrOperationMissingValue = errors.New("OP_VALUE_MISSING")

func getPath(operation map[string]JSON) (JSONPointer, error) {
	pathInterface, ok := operation["path"]
	if !ok {
		return nil, ErrOperationInvalidPath
	}
	path, ok := pathInterface.(string)
	if !ok {
		return nil, ErrOperationInvalidPath
	}
	pointer, err := NewJSONPointer(path)
	if err != nil {
		return nil, err
	}
	return pointer, nil
}

func createAddOp(operation map[string]JSON) (*OpAdd, error) {
	pathInterface, ok := operation["path"]
	if !ok {
		return nil, ErrOperationInvalidPath
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
		return nil, ErrOperationInvalidPath
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
		return nil, ErrOperationInvalidPath
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
		return nil, ErrOperationInvalidPath
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

// ErrOperationInvalidFrom returned when operation "path" field is invalid.
var ErrOperationInvalidFrom = errors.New("OP_FROM_INVALID")

func createMoveOp(operation map[string]JSON) (*OpMove, error) {
	pathInterface, ok := operation["path"]
	if !ok {
		return nil, ErrOperationInvalidPath
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
		return nil, ErrOperationInvalidFrom
	}
	fromString, ok := fromInterface.(string)
	if !ok {
		return nil, ErrOperationInvalidFrom
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
		return nil, ErrOperationInvalidPath
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
		return nil, ErrOperationInvalidFrom
	}
	fromString, ok := fromInterface.(string)
	if !ok {
		return nil, ErrOperationInvalidFrom
	}
	from, err := NewJSONPointer(fromString)
	if err != nil {
		return nil, err
	}
	op := OpCopy{operation: &operation, path: path, from: from}
	return &op, nil
}

func createStrInsOp(operation map[string]JSON) (*OpStrIns, error) {
	pathInterface, ok := operation["path"]
	if !ok {
		return nil, ErrOperationInvalidPath
	}
	pathString, ok := pathInterface.(string)
	if !ok {
		return nil, ErrOperationInvalidPath
	}
	path, err := NewJSONPointer(pathString)
	if err != nil {
		return nil, err
	}
	posInterface, ok := operation["pos"]
	if !ok {
		return nil, ErrOperationInvalid
	}
	posFloat, ok := posInterface.(float64)
	if !ok {
		return nil, ErrOperationInvalid
	}
	pos := int(posFloat)
	strInterface, ok := operation["str"]
	if !ok {
		return nil, ErrOperationInvalid
	}
	str, ok := strInterface.(string)
	if !ok {
		return nil, ErrOperationInvalid
	}
	op := OpStrIns{operation: &operation, path: path, pos: pos, str: str}
	return &op, nil
}

func createStrDelOp(operation map[string]JSON) (*OpStrDel, error) {
	pathInterface, ok := operation["path"]
	if !ok {
		return nil, ErrOperationInvalidPath
	}
	pathString, ok := pathInterface.(string)
	if !ok {
		return nil, ErrOperationInvalidPath
	}
	path, err := NewJSONPointer(pathString)
	if err != nil {
		return nil, err
	}
	posInterface, ok := operation["pos"]
	if !ok {
		return nil, ErrOperationInvalid
	}
	posFloat, ok := posInterface.(float64)
	if !ok {
		return nil, ErrOperationInvalid
	}
	pos := int(posFloat)
	var str string = ""
	var deletionLength int = -1
	if lenInterface, ok := operation["len"]; ok {
		lenFloat, ok := lenInterface.(float64)
		if !ok {
			return nil, ErrOperationInvalid
		}
		deletionLength = int(lenFloat)
		if deletionLength < 0 {
			return nil, ErrOperationInvalid
		}
		if _, ok := operation["str"]; ok {
			return nil, ErrOperationInvalid
		}
	}
	if strInterface, ok := operation["str"]; ok {
		strValue, ok := strInterface.(string)
		if !ok {
			return nil, ErrOperationInvalid
		}
		str = strValue
		deletionLength = len(str)
	}
	if deletionLength < 0 {
		return nil, ErrOperationInvalid
	}
	op := OpStrDel{operation: &operation, path: path, pos: pos, len: deletionLength}
	return &op, nil
}

func createFlipOp(operation map[string]JSON) (*OpFlip, error) {
	path, err := getPath(operation)
	if err != nil {
		return nil, err
	}
	op := OpFlip{operation: &operation, path: path}
	return &op, nil
}

func createIncOp(operation map[string]JSON) (*OpInc, error) {
	path, err := getPath(operation)
	if err != nil {
		return nil, err
	}
	incInterface, ok := operation["inc"]
	if !ok {
		return nil, ErrOperationInvalid
	}
	inc, ok := incInterface.(float64)
	if !ok {
		return nil, ErrOperationInvalid
	}
	op := OpInc{operation: &operation, path: path, inc: inc}
	return &op, nil
}

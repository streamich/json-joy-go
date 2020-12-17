package jsonjoy

import (
	"errors"
)

// ErrPatchInvalid returned when JSON Patch is invalid.
var ErrPatchInvalid = errors.New("patch invalid")

// ErrPatchEmpty returned when JSON Patch array is empty.
var ErrPatchEmpty = errors.New("patch empty")

// ValidateOperations validates a list of JSON Patch operations.
func ValidateOperations(patch JSON) (int, error) {
	arr, ok := patch.([]JSON)
	if !ok {
		return -1, ErrPatchInvalid
	}
	if len(arr) == 0 {
		return -1, ErrPatchEmpty
	}
	for index, operation := range arr {
		err := ValidateOperation(operation)
		if err != nil {
			return index, err
		}
	}
	return -1, nil
}

// ErrOperationInvalid returned when JSON Patch operation is invalid.
var ErrOperationInvalid = errors.New("operation invalid")

// ErrOperationUnknown returned when JSON Patch operation opcode is not recognized.
var ErrOperationUnknown = errors.New("operation unknown")

// ValidateOperation validates a single JSON Patch operation.
func ValidateOperation(operation JSON) error {
	obj, ok := operation.(map[string]JSON)
	if !ok {
		return ErrOperationInvalid
	}
	opInterface, ok := obj["op"]
	if !ok {
		return ErrOperationInvalid
	}
	op, ok := opInterface.(string)
	if !ok {
		return ErrOperationInvalid
	}
	switch op {
	case "add":
		return validateOperationWithPathAndValue(obj)
	default:
		return ErrOperationUnknown
	}
}

// ErrOperationMissingPath returned when JSON Patch operation is missing the "path" field.
var ErrOperationMissingPath = errors.New("op_missing_path")

// ErrOperationInvalidPath returned when operation "path" field is invalid.
var ErrOperationInvalidPath = errors.New("op_invalid_path")

// ErrOperationMissingValue returned when operation is missing "value" field.
var ErrOperationMissingValue = errors.New("op_missing_value")

func validateOperationWithPathAndValue(operation map[string]JSON) error {
	pathInterface, ok := operation["path"]
	if !ok {
		return ErrOperationMissingPath
	}
	path, ok := pathInterface.(string)
	if !ok {
		return ErrOperationInvalidPath
	}
	if err := ValidateJSONPointer(path); err != nil {
		return err
	}
	if _, ok := operation["value"]; !ok {
		return ErrOperationMissingValue
	}
	return nil
}

// ErrPointerInvalid returned when JSON Pointer is invalid.
var ErrPointerInvalid = errors.New("pointer_invalid")

// ValidateJSONPointer returns error if JSON Pointer in string form is invalid.
func ValidateJSONPointer(pointer string) error {
	if len(pointer) == 0 {
		return nil
	}
	if pointer[0] != '/' {
		return ErrPointerInvalid
	}
	return nil
}

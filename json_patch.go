package jsonjoy

import (
	"errors"
)

// ErrTest is returned when JSON Patch "error" operations was not passed.
var ErrTest = errors.New("TEST")

// Inserts an element into a slice and shifts needed elements to right.
func insert(slice []interface{}, pos int, value interface{}) []interface{} {
	length := len(slice)
	if pos >= length {
		return append(slice, value)
	}
	slice = append(slice, value)
	copy(slice[pos+1:], slice[pos:])
	slice[pos] = value
	return slice
}

// replaces a key in object or array.
func putKey(doc JSON, tokens JSONPointer, value JSON) (JSON, error) {
	if tokens.IsRoot() {
		return value, nil
	}
	obj, key, err := tokens.Locate(doc)
	if err != nil {
		return nil, err
	}
	if obj == nil || key == nil {
		return doc, nil
	}
	switch container := obj.(type) {
	case map[string]JSON:
		container[*key] = value
	case []JSON:
		index, err := ParseTokenAsArrayIndex(*key, len(container))
		if err != nil {
			return nil, err
		}
		container[index] = value
	}
	return doc, nil
}

// Add `value` into `doc` at `tokens` location. `doc` and `value` params can
// be mutated, you need to clone them using `Copy` manually.
func Add(doc *JSON, tokens JSONPointer, value JSON) error {
	if tokens.IsRoot() {
		*doc = value
		return nil
	}
	parentTokens := tokens[:len(tokens)-1]
	containerPointer, err := parentTokens.Find(doc)
	if err != nil {
		return err
	}
	key := tokens[len(tokens)-1]
	containerInterface := *containerPointer
	switch container := containerInterface.(type) {
	case map[string]JSON:
		container[key] = value
	case []JSON:
		var index int = 0
		if key == "-" {
			index = len(container)
		} else {
			parsedIndex, err := ParseTokenAsArrayIndex(key, len(container))
			if err != nil {
				return err
			}
			index = parsedIndex
		}
		arr := insert(container, index, value)
		if len(tokens) == 1 {
			*doc = arr
			return nil
		}
		doc2, err := putKey(*doc, parentTokens, arr)
		if err != nil {
			return err
		}
		*doc = doc2
	}
	return nil
}

// Replace `value` into `doc` at `tokens` location. `doc` and `value` params can
// be mutated, you need to clone them using `Copy` manually.
func Replace(doc *JSON, tokens JSONPointer, value JSON) error {
	if tokens.IsRoot() {
		*doc = value
		return nil
	}
	parentTokens := tokens[:len(tokens)-1]
	obj, err := parentTokens.Find(doc)
	if err != nil {
		return err
	}
	key := tokens[len(tokens)-1]
	objInterface := *obj
	switch container := objInterface.(type) {
	case map[string]JSON:
		if _, ok := container[key]; !ok {
			return ErrNotFound
		}
		container[key] = value
	case []JSON:
		index, err := ParseTokenAsArrayIndex(key, len(container)-1)
		if err != nil {
			return err
		}
		if index >= len(container) {
			return ErrNotFound
		}
		container[index] = value
	}
	return nil
}

// Remove removes a value from JSON document.
func Remove(doc *JSON, tokens JSONPointer) (JSON, error) {
	if tokens.IsRoot() {
		*doc = nil
		return *doc, nil
	}
	parentTokens := tokens[:len(tokens)-1]
	obj, err := parentTokens.Find(doc)
	if err != nil {
		return nil, err
	}
	key := tokens[len(tokens)-1]
	objInterface := *obj
	switch container := objInterface.(type) {
	case map[string]JSON:
		value, ok := container[key]
		if !ok {
			return nil, ErrNotFound
		}
		delete(container, key)
		return value, nil
	case []JSON:
		index, err := ParseTokenAsArrayIndex(key, -1)
		if err != nil {
			return nil, err
		}
		if index >= len(container) {
			return nil, ErrNotFound
		}
		value := container[index]
		container = append(container[:index], container[index+1:]...)
		doc2, err := putKey(*doc, parentTokens, container)
		if err != nil {
			return nil, err
		}
		*doc = doc2
		return value, nil
	}
	return nil, nil
}

// Move executes JSON Patch "move" operation.
func Move(doc *JSON, from JSONPointer, to JSONPointer) error {
	value, err := Remove(doc, from)
	if err != nil {
		return err
	}
	return Add(doc, to, value)
}

// JSONPatchCopy executs JSON Patch "copy" operation.
func JSONPatchCopy(doc *JSON, from JSONPointer, to JSONPointer) error {
	value, err := from.Get(*doc)
	if err != nil {
		return err
	}
	return Add(doc, to, Copy(value))
}

// JSONPatchTest executes JSON Patch "test" operation.
func JSONPatchTest(doc *JSON, path JSONPointer, value JSON) error {
	target, err := path.Get(*doc)
	if err != nil {
		return err
	}
	if !DeepEqual(value, target) {
		return ErrTest
	}
	return nil
}

func insertString(src string, pos int, ins string) string {
	if pos > len(src) {
		pos = len(src)
	}
	return src[:pos] + ins + src[pos:]
}

// JSONPatchStrIns insert string into an existing string.
func JSONPatchStrIns(doc *JSON, tokens JSONPointer, pos int, ins string) error {
	if tokens.IsRoot() {
		str, ok := (*doc).(string)
		if !ok {
			return ErrNotAString
		}
		*doc = insertString(str, pos, ins)
		return nil
	}
	parentTokens := tokens[:len(tokens)-1]
	obj, err := parentTokens.Find(doc)
	if err != nil {
		return err
	}
	key := tokens[len(tokens)-1]
	switch container := (*obj).(type) {
	case map[string]JSON:
		value, ok := container[key]
		if !ok {
			if pos == 0 {
				container[key] = ""
				value = ""
			} else {
				return errors.New("POS")
			}
		}
		str, ok := value.(string)
		if !ok {
			return ErrNotAString
		}
		container[key] = insertString(str, pos, ins)
	case []JSON:
		index, err := ParseTokenAsArrayIndex(key, -1)
		if err != nil {
			return err
		}
		if index >= len(container) {
			return ErrNotFound
		}
		value := container[index]
		str, ok := value.(string)
		if !ok {
			return ErrNotAString
		}
		container[index] = insertString(str, pos, ins)
	}
	return nil
}

// ApplyOperation applies a single operation.
func ApplyOperation(doc *JSON, operation interface{}) error {
	switch op := operation.(type) {
	case *OpAdd:
		return op.apply(doc)
	case *OpReplace:
		return op.apply(doc)
	case *OpRemove:
		return op.apply(doc)
	case *OpMove:
		return op.apply(doc)
	case *OpCopy:
		return op.apply(doc)
	case *OpStrIns:
		return op.apply(doc)
	case *OpTest:
		err := op.apply(doc)
		if err != nil {
			return err
		}
	}
	return nil
}

// ApplyOps applies a JSON Patch to the document.
func ApplyOps(doc *JSON, ops []interface{}) error {
	for _, op := range ops {
		var err error
		err = ApplyOperation(doc, op)
		if err != nil {
			return err
		}
	}
	return nil
}

func (op *OpAdd) apply(doc *JSON) error {
	return Add(doc, op.path, Copy(op.value))
}

func (op *OpReplace) apply(doc *JSON) error {
	return Replace(doc, op.path, Copy(op.value))
}

func (op *OpRemove) apply(doc *JSON) error {
	_, err := Remove(doc, op.path)
	return err
}

func (op *OpMove) apply(doc *JSON) error {
	err := Move(doc, op.from, op.path)
	return err
}

func (op *OpCopy) apply(doc *JSON) error {
	err := JSONPatchCopy(doc, op.from, op.path)
	return err
}

func (op *OpTest) apply(doc *JSON) error {
	err := JSONPatchTest(doc, op.path, op.value)
	return err
}

func (op *OpStrIns) apply(doc *JSON) error {
	err := JSONPatchStrIns(doc, op.path, op.pos, op.str)
	return err
}

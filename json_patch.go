package jsonjoy

import (
	"errors"
	"strconv"
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

func jsonPatchStrOp(doc *JSON, tokens JSONPointer, fn func(str *string) (string, error)) error {
	if tokens.IsRoot() {
		str, ok := (*doc).(string)
		if !ok {
			return ErrNotAString
		}
		res, err := fn(&str)
		if err != nil {
			return err
		}
		*doc = res
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
			res, err := fn(nil)
			if err != nil {
				return err
			}
			container[key] = res
		} else {

			str, ok := value.(string)
			if !ok {
				return ErrNotAString
			}
			res, err := fn(&str)
			if err != nil {
				return err
			}
			container[key] = res
		}
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
		res, err := fn(&str)
		if err != nil {
			return err
		}
		container[index] = res
	}
	return nil
}

// JSONPatchStrIns insert string into an existing string.
func JSONPatchStrIns(doc *JSON, tokens JSONPointer, pos int, ins string) error {
	return jsonPatchStrOp(doc, tokens, func(str *string) (string, error) {
		if str == nil {
			if pos == 0 {
				return ins, nil
			}
			return "", errors.New("POS")
		}
		return insertString(*str, pos, ins), nil
	})
}

func deleteString(src string, pos int, length int) string {
	strLength := len(src)
	if pos >= strLength {
		return src
	}
	start := src[:pos]
	if pos+length >= strLength {
		return start
	}
	return start + src[pos+length:]
}

// JSONPatchStrDel deletes string from an existing string.
func JSONPatchStrDel(doc *JSON, tokens JSONPointer, pos int, deletionLength int) error {
	return jsonPatchStrOp(doc, tokens, func(str *string) (string, error) {
		if str == nil {
			return "", ErrNotFound
		}
		return deleteString(*str, pos, deletionLength), nil
	})
}

func flip(value interface{}) bool {
	if value == nil {
		return true
	}
	switch val := value.(type) {
	case string:
		if len(val) > 0 {
			return false
		}
		return false

	case float64:
		if val == 0.0 {
			return true
		}
		return false
	case bool:
		return !val

	}
	return false
}

// JSONPatchFlip flips a cell treating it as a boolean.
func jsonPatchFlip(doc *JSON, tokens JSONPointer) error {
	if tokens.IsRoot() {
		*doc = flip(*doc)
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
		val, ok := container[key]
		if !ok {
			return ErrNotFound
		}
		container[key] = flip(val)
	case []JSON:
		index, err := ParseTokenAsArrayIndex(key, len(container)-1)
		if err != nil {
			return err
		}
		if index >= len(container) {
			return ErrNotFound
		}
		container[index] = flip(container[index])
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
	case *OpStrDel:
		return op.apply(doc)
	case *OpFlip:
		return op.apply(doc)
	case *OpInc:
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

func (op *OpStrDel) apply(doc *JSON) error {
	err := JSONPatchStrDel(doc, op.path, op.pos, op.len)
	return err
}

func (op *OpFlip) apply(doc *JSON) error {
	err := jsonPatchFlip(doc, op.path)
	return err
}

func castToFloat64(val interface{}) float64 {
	if val == nil {
		return 0
	}
	switch f := val.(type) {
	case float64:
		return f
	case bool:
		if f {
			return 1
		}
		return 0
	case string:
		if res, err := strconv.ParseFloat(f, 64); err == nil {
			return res
		}
		return 0
	}
	return 1
}

func (op *OpInc) apply(doc *JSON) error {
	if op.path.IsRoot() {
		*doc = castToFloat64(*doc) + op.inc
		return nil
	}
	parentTokens := op.path[:len(op.path)-1]
	obj, err := parentTokens.Find(doc)
	if err != nil {
		return err
	}
	key := op.path[len(op.path)-1]
	objInterface := *obj
	switch container := objInterface.(type) {
	case map[string]JSON:
		val, ok := container[key]
		if !ok {
			return ErrNotFound
		}
		container[key] = castToFloat64(val) + op.inc
	case []JSON:
		index, err := ParseTokenAsArrayIndex(key, len(container)-1)
		if err != nil {
			return err
		}
		if index >= len(container) {
			return ErrNotFound
		}
		container[index] = castToFloat64(container[index]) + op.inc
	}
	return nil
}

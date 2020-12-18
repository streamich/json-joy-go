package jsonjoy

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

func updateKey(container JSON, key string, value JSON) error {
	switch container := container.(type) {
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

// Add `value` into `doc` at `tokens` location. `doc` and `value` params can
// be mutated, you need to clone them using `Copy` manually.
func Add(doc JSON, tokens JSONPointer, value JSON) (JSON, error) {
	if tokens.IsRoot() {
		return value, nil
	}
	parentTokens := tokens[:len(tokens)-1]
	container, err := parentTokens.Get(doc)
	if err != nil {
		return nil, err
	}
	key := tokens[len(tokens)-1]
	switch container := container.(type) {
	case map[string]JSON:
		container[key] = value
	case []JSON:
		var index int = 0
		if key == "-" {
			index = len(container)
		} else {
			parsedIndex, err := ParseTokenAsArrayIndex(key, len(container))
			if err != nil {
				return nil, err
			}
			index = parsedIndex
		}
		container = insert(container, index, value)
		if len(tokens) == 1 {
			return container, nil
		}
		return putKey(doc, parentTokens, container)
	}
	return doc, nil
}

// Replace `value` into `doc` at `tokens` location. `doc` and `value` params can
// be mutated, you need to clone them using `Copy` manually.
func Replace(doc JSON, tokens JSONPointer, value JSON) (JSON, error) {
	if tokens.IsRoot() {
		return value, nil
	}
	parentTokens := tokens[:len(tokens)-1]
	container, err := parentTokens.Get(doc)
	if err != nil {
		return nil, err
	}
	key := tokens[len(tokens)-1]
	return doc, updateKey(container, key, value)
}

// ApplyOperation applies a single operation.
func ApplyOperation(doc JSON, operation interface{}) (JSON, error) {
	switch op := operation.(type) {
	case *OpAdd:
		return op.apply(doc)
	case *OpReplace:
		return op.apply(doc)
	}
	return doc, nil
}

// ApplyOps applies a JSON Patch to the document.
func ApplyOps(doc JSON, operations []interface{}) (JSON, error) {
	doc2 := Copy(doc)
	for _, operation := range operations {
		var err error
		doc2, err = ApplyOperation(doc2, operation)
		if err != nil {
			return nil, err
		}
	}
	return doc2, nil
}

func (op *OpAdd) apply(doc JSON) (JSON, error) {
	return Add(doc, op.path, Copy(op.value))
}

func (op *OpReplace) apply(doc JSON) (JSON, error) {
	return Replace(doc, op.path, Copy(op.value))
}

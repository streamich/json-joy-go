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
func Remove(doc *JSON, tokens JSONPointer) error {
	if tokens.IsRoot() {
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
		delete(container, key)
	case []JSON:
		index, err := ParseTokenAsArrayIndex(key, len(container)-1)
		if err != nil {
			return err
		}
		if index >= len(container) {
			return ErrNotFound
		}
		container = append(container[:index], container[index+1:]...)
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
	}
	return nil
}

// ApplyOps applies a JSON Patch to the document.
func ApplyOps(doc *JSON, operations []interface{}) error {
	for _, operation := range operations {
		var err error
		err = ApplyOperation(doc, operation)
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
	return Remove(doc, op.path)
}

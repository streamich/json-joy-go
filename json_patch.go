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
func replaceKey(doc JSON, tokens JSONPointer, value JSON) (JSON, error) {
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
		return replaceKey(doc, parentTokens, container)
	}
	return doc, nil
}

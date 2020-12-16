package jsonjoy

// Add `value` into `doc` at `tokens` location.
func Add(doc JSON, tokens JSONPointer, value JSON) (JSON, error) {
	valueCopy := Copy(value)
	if tokens.IsRoot() {
		return valueCopy, nil
	}
	obj, key, err := tokens.Locate(doc)
	if err != nil {
		return nil, err
	}
	if obj == nil || key == nil {
		return doc, nil
	}
	switch container := obj.(type) {
	case map[string]interface{}:
		container[*key] = valueCopy
	case []interface{}:
		index, err := ParseTokenAsArrayIndex(*key, container)
		if err != nil {
			return nil, err
		}
		container[index] = valueCopy
	}
	return doc, nil
}

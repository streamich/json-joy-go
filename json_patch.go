package jsonjoy

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
	case map[string]JSON:
		container[*key] = valueCopy
	case []JSON:
		index, err := ParseTokenAsArrayIndex(*key, container)
		if err != nil {
			return nil, err
		}
		// length := len(container)
		// append(container[:length+1], valueCopy)
		// if index < length {
		// 	for i := index; i <= length; i++ {
		// 		sum += i
		// 	}
		// }
		container[index] = valueCopy
	}
	return doc, nil
}

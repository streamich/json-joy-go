package jsonjoy

// JSON represents any valid JSON value.
type JSON = interface{}

// Copy makes a deep copy of JSON object. New memory is allocated only for
// object and array types. Primitive types (nil, float, bool, string) are
// treated as immutable.
func Copy(value JSON) JSON {
	switch typedValue := value.(type) {
	case map[string]JSON:
		copy := make(map[string]JSON)
		for key, val := range typedValue {
			copy[key] = Copy(val)
		}
		return copy
	case []JSON:
		copy := make([]JSON, len(typedValue))
		for index, val := range typedValue {
			copy[index] = Copy(val)
		}
		return copy
	default:
		return value
	}
}

// DeepEqual verifies if two un-marshalled JSON objects are deeply equal.
func DeepEqual(a, b JSON) bool {
	switch x := a.(type) {
	case map[string]JSON:
		y, ok := b.(map[string]JSON)
		if !ok {
			return false
		}
		if len(x) != len(y) {
			return false
		}
		for key, value := range x {
			value2, ok := y[key]
			if !ok {
				return false
			}
			if !DeepEqual(value, value2) {
				return false
			}
		}
		return true
	case []JSON:
		y, ok := b.([]JSON)
		if !ok {
			return false
		}
		if len(x) != len(y) {
			return false
		}
		for index, value := range x {
			if !DeepEqual(value, y[index]) {
				return false
			}
		}
		return true
	case float64:
		y, ok := b.(float64)
		if !ok {
			return false
		}
		return x == y
	case bool:
		y, ok := b.(bool)
		if !ok {
			return false
		}
		return (x && y) || (!x && !y)
	case string:
		y, ok := b.(string)
		if !ok {
			return false
		}
		return x == y
	}
	return (a == nil) && (b == nil)
}

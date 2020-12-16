package jsonjoy

// JSON represents any valid JSON value.
type JSON = interface{}

// Copy makes a deep copy of JSON object. New memory is allocated only for
// object and array types. Primitive types (nil, float, bool, string) are
// treated as immutable.
func Copy(value JSON) JSON {
	switch typedValue := value.(type) {
	case map[string]interface{}:
		copy := make(map[string]JSON)
		for key, val := range typedValue {
			copy[key] = Copy(val)
		}
		return copy
	case []interface{}:
		copy := make([]JSON, len(typedValue))
		for index, val := range typedValue {
			copy[index] = Copy(val)
		}
		return copy
	default:
		return value
	}
}

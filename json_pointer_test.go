package jsonjoy

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UnescapeReferenceToken_ReturnsSameStringIfThereAreNoEscapedChars(t *testing.T) {
	decoded := UnescapeReferenceToken("foobar")
	assert.Equal(t, decoded, "foobar")
	decoded = UnescapeReferenceToken("foo/bar")
	assert.Equal(t, decoded, "foo/bar")
	decoded = UnescapeReferenceToken("foo~bar")
	assert.Equal(t, decoded, "foo~bar")
}

func Test_UnescapeReferenceToken_DecodesSpecialChars(t *testing.T) {
	decoded := UnescapeReferenceToken("foo~1bar")
	assert.Equal(t, decoded, "foo/bar")
	decoded = UnescapeReferenceToken("foo~0bar")
	assert.Equal(t, decoded, "foo~bar")
	decoded = UnescapeReferenceToken("~1~0foo~0bar~0~1")
	assert.Equal(t, decoded, "/~foo~bar~/")
	decoded = UnescapeReferenceToken("~1~0")
	assert.Equal(t, decoded, "/~")
	decoded = UnescapeReferenceToken("~0~1")
	assert.Equal(t, decoded, "~/")
	decoded = UnescapeReferenceToken("~0")
	assert.Equal(t, decoded, "~")
	decoded = UnescapeReferenceToken("~1")
	assert.Equal(t, decoded, "/")
}

func Test_EscapeReferenceToken_ReturnsSameStringIfThereAreNoScpecialChars(t *testing.T) {
	encoded := EscapeReferenceToken("foobar")
	assert.Equal(t, encoded, "foobar")
}

func Test_EscapeReferenceToken_EncodesSpecialChars(t *testing.T) {
	encoded := EscapeReferenceToken("foo/bar")
	assert.Equal(t, encoded, "foo~1bar")
	encoded = EscapeReferenceToken("foo~bar")
	assert.Equal(t, encoded, "foo~0bar")
	encoded = EscapeReferenceToken("~/")
	assert.Equal(t, encoded, "~0~1")
	encoded = EscapeReferenceToken("/~")
	assert.Equal(t, encoded, "~1~0")
	encoded = EscapeReferenceToken("~")
	assert.Equal(t, encoded, "~0")
	encoded = EscapeReferenceToken("/")
	assert.Equal(t, encoded, "~1")
}

func Test_NewJSONPointer_ReturnsEmptyArrayOnRootPointer(t *testing.T) {
	pointer, err := NewJSONPointer("")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 0)
}

func Test_NewJSONPointer_ReturnsErrorIfPointerDoesNotStartWithSlash(t *testing.T) {
	pointer, err := NewJSONPointer("foo/bar")
	assert.Nil(t, pointer)
	assert.NotNil(t, err)
}

func Test_NewJSONPointer_ParsesASingleStepPointer(t *testing.T) {
	pointer, err := NewJSONPointer("/foo")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 1)
	assert.Equal(t, (pointer)[0], "foo")
}

func Test_NewJSONPointer_ParsesAMultipleStepPointer(t *testing.T) {
	pointer, err := NewJSONPointer("/foo/bar/baz")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 3)
	assert.Equal(t, pointer[0], "foo")
	assert.Equal(t, pointer[1], "bar")
	assert.Equal(t, pointer[2], "baz")
}

func Test_NewJSONPointer_DecodesTokens(t *testing.T) {
	pointer, err := NewJSONPointer("/foo~1bar")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 1)
	assert.Equal(t, pointer[0], "foo/bar")
	pointer, err = NewJSONPointer("/foo~1/bar/ba~0~1~0z")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 3)
	assert.Equal(t, pointer[0], "foo/")
	assert.Equal(t, pointer[1], "bar")
	assert.Equal(t, pointer[2], "ba~/~z")
}

func Test_JSONPointer_Format_FormatsTokensIntoJsonPointer(t *testing.T) {
	tokens := JSONPointer{"foo", "bar", "baz"}
	str := tokens.Format()
	assert.Equal(t, str, "/foo/bar/baz")
}

func Test_JSONPointer_Format_FormatsASingleToken(t *testing.T) {
	tokens := JSONPointer{"aga"}
	str := tokens.Format()
	assert.Equal(t, str, "/aga")
}

func Test_JSONPointer_Format_FormatsARootPointer(t *testing.T) {
	tokens := JSONPointer{}
	str := tokens.Format()
	assert.Equal(t, str, "")
}

func Test_JSONPointer_Format_EncodesSpecialChars(t *testing.T) {
	tokens := JSONPointer{"foo/bar"}
	str := tokens.Format()
	assert.Equal(t, str, "/foo~1bar")
	tokens = JSONPointer{"foo/bar", "/", "~", "a~b/"}
	str = tokens.Format()
	assert.Equal(t, str, "/foo~1bar/~1/~0/a~0b~1")
}

func Test_JSONPointer_Get_ReturnsRootDocument(t *testing.T) {
	tokens := JSONPointer{}
	b := []byte(`{"foo": "bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.Nil(t, err)
	assert.Equal(t, value, doc)
}

func Test_JSONPointer_Get_ReturnsRootDocumentWhenRootIsArray(t *testing.T) {
	tokens := JSONPointer{}
	b := []byte(`["foo", "bar", "baz"]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.Nil(t, err)
	assert.Equal(t, value, doc)
}

func Test_JSONPointer_Get_ReturnsRootValueWhenRootIsBoolean(t *testing.T) {
	tokens := JSONPointer{}
	b := []byte(`true`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.Nil(t, err)
	assert.Equal(t, value, true)
}

func Test_JSONPointer_Get_ReturnsRootValueWhenRootIsBoolean2(t *testing.T) {
	tokens := JSONPointer{}
	b := []byte(`false`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.Nil(t, err)
	assert.Equal(t, value, false)
}

func Test_JSONPointer_Get_ReturnsRootValueWhenRootIsString(t *testing.T) {
	tokens := JSONPointer{}
	b := []byte(`"asdf"`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.Nil(t, err)
	assert.Equal(t, value, "asdf")
}

func Test_JSONPointer_Get_ReturnsRootValueWhenRootIsNumber(t *testing.T) {
	tokens := JSONPointer{}
	b := []byte(`-3`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.Nil(t, err)
	assert.Equal(t, value, -3.0)
}

func Test_JSONPointer_Get_ReturnsAFirstLevelKey(t *testing.T) {
	tokens := JSONPointer{"foo"}
	b := []byte(`{"foo": "bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.Nil(t, err)
	assert.Equal(t, "bar", value)
}

func Test_JSONPointer_Get_FindsADeepStringKeyInObjects(t *testing.T) {
	tokens := JSONPointer{"foo", "bar", "baz"}
	b := []byte(`{"foo": {"bar": {"baz": "qux"}}}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.Nil(t, err)
	assert.Equal(t, "qux", value)
}

func Test_JSONPointer_Get_ReturnsErrorNotFoundWhenLocatingMissingValue(t *testing.T) {
	tokens := JSONPointer{"foo2"}
	b := []byte(`{"foo": {"bar": {"baz": "qux"}}}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, value)
}

func Test_JSONPointer_Get_ReturnsOneLevelDeepArrayValue(t *testing.T) {
	tokens, _ := NewJSONPointer("/1")
	b := []byte(`["foo", "bar", "baz"]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.Nil(t, err)
	assert.Equal(t, "bar", value)
}

func Test_JSONPointer_Get_ReturnsThreeLevelsDeepArrayValue(t *testing.T) {
	tokens, _ := NewJSONPointer("/1/2/3")
	b := []byte(`["foo", [null, [], [true, true, false, "a"]], "baz"]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Get(doc)
	assert.Nil(t, err)
	assert.Equal(t, "a", value)
}

func Test_JSONPointer_Get_ReturnsJSONValuesOfVariousTypes(t *testing.T) {
	b := []byte(`{
		"a": null,
		"b": false,
		"c": true,
		"d": 0,
		"e": 1.1,
		"f": "",
		"g": "asdf",
		"h": [],
		"i": [null, false, true, -1, 3.4, "", "foo", [], {}],
		"j": {
			"a": null,
			"b": false,
			"c": true,
			"d": 0,
			"e": 1.1,
			"f": "",
			"g": "asdf",
			"h": [],
			"i": [null, false, true, -1, 3.4, "", "foo", [], {
				"a": null,
				"b": false,
				"c": true,
				"d": 0,
				"e": 1.1,
				"f": "",
				"g": "asdf",
				"h": [],
				"i": [null, false, true, -1, 3.4, "", "foo", [], {}]
			}]
		}
	}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	val, _ := (JSONPointer{"a"}).Get(doc)
	assert.Equal(t, nil, val)
	val, _ = (JSONPointer{"b"}).Get(doc)
	assert.Equal(t, false, val)
	val, _ = (JSONPointer{"c"}).Get(doc)
	assert.Equal(t, true, val)
	val, _ = (JSONPointer{"d"}).Get(doc)
	assert.Equal(t, 0.0, val)
	val, _ = (JSONPointer{"e"}).Get(doc)
	assert.Equal(t, 1.1, val)
	val, _ = (JSONPointer{"f"}).Get(doc)
	assert.Equal(t, "", val)
	val, _ = (JSONPointer{"g"}).Get(doc)
	assert.Equal(t, "asdf", val)
	val, _ = (JSONPointer{"h"}).Get(doc)
	assert.Equal(t, "[]", fmt.Sprint(val))
	val, _ = (JSONPointer{"i", "0"}).Get(doc)
	assert.Equal(t, nil, val)
	val, _ = (JSONPointer{"i", "1"}).Get(doc)
	assert.Equal(t, false, val)
	val, _ = (JSONPointer{"i", "2"}).Get(doc)
	assert.Equal(t, true, val)
	val, _ = (JSONPointer{"i", "3"}).Get(doc)
	assert.Equal(t, -1.0, val)
	val, _ = (JSONPointer{"i", "4"}).Get(doc)
	assert.Equal(t, 3.4, val)
	val, _ = (JSONPointer{"i", "5"}).Get(doc)
	assert.Equal(t, "", val)
	val, _ = (JSONPointer{"i", "6"}).Get(doc)
	assert.Equal(t, "foo", val)
	val, _ = (JSONPointer{"i", "7"}).Get(doc)
	assert.Equal(t, "[]", fmt.Sprint(val))
	val, _ = (JSONPointer{"j", "a"}).Get(doc)
	assert.Equal(t, nil, val)
	val, _ = (JSONPointer{"j", "b"}).Get(doc)
	assert.Equal(t, false, val)
	val, _ = (JSONPointer{"j", "c"}).Get(doc)
	assert.Equal(t, true, val)
	val, _ = (JSONPointer{"j", "d"}).Get(doc)
	assert.Equal(t, 0.0, val)
	val, _ = (JSONPointer{"j", "e"}).Get(doc)
	assert.Equal(t, 1.1, val)
	val, _ = (JSONPointer{"j", "f"}).Get(doc)
	assert.Equal(t, "", val)
	val, _ = (JSONPointer{"j", "g"}).Get(doc)
	assert.Equal(t, "asdf", val)
	val, _ = (JSONPointer{"j", "h"}).Get(doc)
	assert.Equal(t, "[]", fmt.Sprint(val))
	val, _ = (JSONPointer{"j", "i", "8", "i", "2"}).Get(doc)
	assert.Equal(t, true, val)
}

func Test_JSONPointer_Find_ReturnsRootDocument(t *testing.T) {
	tokens := JSONPointer{}
	b := []byte(`{"foo": "bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Find(&doc)
	assert.Nil(t, err)
	assert.Equal(t, *value, doc)
}

func Test_JSONPointer_Find_ReturnsJSONValuesOfVariousTypes(t *testing.T) {
	b := []byte(`{
		"a": null,
		"b": false,
		"c": true,
		"d": 0,
		"e": 1.1,
		"f": "",
		"g": "asdf",
		"h": [],
		"i": [null, false, true, -1, 3.4, "", "foo", [], {}],
		"j": {
			"a": null,
			"b": false,
			"c": true,
			"d": 0,
			"e": 1.1,
			"f": "",
			"g": "asdf",
			"h": [],
			"i": [null, false, true, -1, 3.4, "", "foo", [], {
				"a": null,
				"b": false,
				"c": true,
				"d": 0,
				"e": 1.1,
				"f": "",
				"g": "asdf",
				"h": [],
				"i": [null, false, true, -1, 3.4, "", "foo", [], {}]
			}]
		}
	}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	val, _ := (JSONPointer{"a"}).Find(&doc)
	assert.Equal(t, nil, *val)
	val, _ = (JSONPointer{"b"}).Find(&doc)
	assert.Equal(t, false, *val)
	val, _ = (JSONPointer{"c"}).Find(&doc)
	assert.Equal(t, true, *val)
	val, _ = (JSONPointer{"d"}).Find(&doc)
	assert.Equal(t, 0.0, *val)
	val, _ = (JSONPointer{"e"}).Find(&doc)
	assert.Equal(t, 1.1, *val)
	val, _ = (JSONPointer{"f"}).Find(&doc)
	assert.Equal(t, "", *val)
	val, _ = (JSONPointer{"g"}).Find(&doc)
	assert.Equal(t, "asdf", *val)
	val, _ = (JSONPointer{"h"}).Find(&doc)
	assert.Equal(t, "[]", fmt.Sprint(*val))
	val, _ = (JSONPointer{"i", "0"}).Find(&doc)
	assert.Equal(t, nil, *val)
	val, _ = (JSONPointer{"i", "1"}).Find(&doc)
	assert.Equal(t, false, *val)
	val, _ = (JSONPointer{"i", "2"}).Find(&doc)
	assert.Equal(t, true, *val)
	val, _ = (JSONPointer{"i", "3"}).Find(&doc)
	assert.Equal(t, -1.0, *val)
	val, _ = (JSONPointer{"i", "4"}).Find(&doc)
	assert.Equal(t, 3.4, *val)
	val, _ = (JSONPointer{"i", "5"}).Find(&doc)
	assert.Equal(t, "", *val)
	val, _ = (JSONPointer{"i", "6"}).Find(&doc)
	assert.Equal(t, "foo", *val)
	val, _ = (JSONPointer{"i", "7"}).Find(&doc)
	assert.Equal(t, "[]", fmt.Sprint(*val))
	val, _ = (JSONPointer{"j", "a"}).Find(&doc)
	assert.Equal(t, nil, *val)
	val, _ = (JSONPointer{"j", "b"}).Find(&doc)
	assert.Equal(t, false, *val)
	val, _ = (JSONPointer{"j", "c"}).Find(&doc)
	assert.Equal(t, true, *val)
	val, _ = (JSONPointer{"j", "d"}).Find(&doc)
	assert.Equal(t, 0.0, *val)
	val, _ = (JSONPointer{"j", "e"}).Find(&doc)
	assert.Equal(t, 1.1, *val)
	val, _ = (JSONPointer{"j", "f"}).Find(&doc)
	assert.Equal(t, "", *val)
	val, _ = (JSONPointer{"j", "g"}).Find(&doc)
	assert.Equal(t, "asdf", *val)
	val, _ = (JSONPointer{"j", "h"}).Find(&doc)
	assert.Equal(t, "[]", fmt.Sprint(*val))
	val, _ = (JSONPointer{"j", "i", "8", "i", "2"}).Find(&doc)
	assert.Equal(t, true, *val)
}

func Test_JSONPointer_Resolve_ReturnsNilForDocumentRoot(t *testing.T) {
	tokens := JSONPointer{}
	b := []byte(`{"foo": "bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	values, err := tokens.Resolve(doc)
	assert.Nil(t, err)
	assert.Nil(t, values)
}

func Test_JSONPointer_Resolve_ReturnsAFirstLevelKey(t *testing.T) {
	tokens := JSONPointer{"foo"}
	b := []byte(`{"foo": "bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	values, err := tokens.Resolve(doc)
	assert.Nil(t, err)
	assert.NotNil(t, values)
	assert.Equal(t, 1, len(values))
	assert.Equal(t, "bar", values[0])
}

func Test_JSONPointer_Resolve_ReturnsAllValuesOfDeepPointer(t *testing.T) {
	tokens, _ := NewJSONPointer("/foo/bar/baz")
	b := []byte(`{"foo": {"bar": {"baz": "qux"}}}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	values, err := tokens.Resolve(doc)
	value1, _ := json.Marshal(values[0])
	value2, _ := json.Marshal(values[1])
	assert.Nil(t, err)
	assert.Equal(t, 3, len(values))
	assert.Equal(t, `{"bar":{"baz":"qux"}}`, string(value1))
	assert.Equal(t, `{"baz":"qux"}`, string(value2))
	assert.Equal(t, "qux", values[2])
}

func Test_JSONPointer_Resolve_ResolvesPointerIntoArrays(t *testing.T) {
	tokens, _ := NewJSONPointer("/1/1/1")
	b := []byte(`[1, ["a", [null, "abc"]]]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	values, err := tokens.Resolve(doc)
	value1, _ := json.Marshal(values[0])
	value2, _ := json.Marshal(values[1])
	assert.Nil(t, err)
	assert.Equal(t, 3, len(values))
	assert.Equal(t, `["a",[null,"abc"]]`, string(value1))
	assert.Equal(t, `[null,"abc"]`, string(value2))
	assert.Equal(t, "abc", values[2])
}

func Test_JSONPointer_Resolve_ResolvesPointerIntoObjectsAndArrays(t *testing.T) {
	tokens, _ := NewJSONPointer("/abc/1/1")
	b := []byte(`{"abc": ["a", [null, "abc"]] }`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	values, err := tokens.Resolve(doc)
	value1, _ := json.Marshal(values[0])
	value2, _ := json.Marshal(values[1])
	assert.Nil(t, err)
	assert.Equal(t, 3, len(values))
	assert.Equal(t, `["a",[null,"abc"]]`, string(value1))
	assert.Equal(t, `[null,"abc"]`, string(value2))
	assert.Equal(t, "abc", values[2])
}

func Test_JSONPointer_Locate_ReturnsNilForDocumentRoot(t *testing.T) {
	tokens := JSONPointer{}
	b := []byte(`{"foo": "bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	obj, key, err := tokens.Locate(doc)
	assert.Nil(t, err)
	assert.Nil(t, key)
	assert.Nil(t, obj)
}

func Test_JSONPointer_Locate_CanLocateFirstLevelKeyInObject(t *testing.T) {
	tokens := JSONPointer{"foo"}
	b := []byte(`{"foo": "bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	obj, key, err := tokens.Locate(doc)
	assert.Nil(t, err)
	assert.Equal(t, "foo", *key)
	assert.Equal(t, "map[foo:bar]", fmt.Sprint(obj))
}

func Test_JSONPointer_Locate_CanLocateSecondLevelKeyInObject(t *testing.T) {
	tokens := JSONPointer{"foo", "c"}
	b := []byte(`{"foo": {"a": 1, "c": 2}}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	obj, key, err := tokens.Locate(doc)
	assert.Nil(t, err)
	assert.Equal(t, "c", *key)
	assert.Equal(t, "map[a:1 c:2]", fmt.Sprint(obj))
}

func Test_JSONPointer_Locate_CanLocateFirstLevelKeyInArray(t *testing.T) {
	tokens := JSONPointer{"1"}
	b := []byte(`["a", "b", "c"]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	obj, key, err := tokens.Locate(doc)
	assert.Nil(t, err)
	assert.Equal(t, "1", *key)
	assert.Equal(t, "[a b c]", fmt.Sprint(obj))
}

func Test_JSONPointer_Locate_CanLocateSecondLevelKeyInArray(t *testing.T) {
	tokens := JSONPointer{"1", "2"}
	b := []byte(`["a", [1,2,3], "c"]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	obj, key, err := tokens.Locate(doc)
	assert.Nil(t, err)
	assert.Equal(t, "2", *key)
	assert.Equal(t, "[1 2 3]", fmt.Sprint(obj))
}

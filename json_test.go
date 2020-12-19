package jsonjoy

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Json_Copy_ReturnsACopyOfJson(t *testing.T) {
	b := []byte(`{"foo":"bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	assert.Equal(t, "map[foo:bar]", fmt.Sprint(doc))
	assert.Equal(t, "map[foo:bar]", fmt.Sprint(doc2))
}

func Test_Json_Copy_ModifyingCopyDoesNotModifyOriginal(t *testing.T) {
	b := []byte(`{"foo":1}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	doc3, _ := doc2.(map[string]JSON)
	doc3["foo"] = 2.0
	assert.Equal(t, "map[foo:1]", fmt.Sprint(doc))
	assert.Equal(t, "map[foo:2]", fmt.Sprint(doc2))
	assert.Equal(t, "map[foo:2]", fmt.Sprint(doc3))
}

func Test_Json_Copy_ModifyingOriginalDoesNotModifyCopy(t *testing.T) {
	b := []byte(`{"foo":1}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	doc3, _ := doc.(map[string]interface{})
	doc3["foo"] = 3.0
	assert.Equal(t, "map[foo:3]", fmt.Sprint(doc))
	assert.Equal(t, "map[foo:1]", fmt.Sprint(doc2))
	assert.Equal(t, "map[foo:3]", fmt.Sprint(doc3))
}

func Test_Json_Copy_CopiesArrayOfVariousTypesOfElements(t *testing.T) {
	b := []byte(`[null, false, true, -1, 3.4, "", "foo", [], {}]`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	doc3, _ := doc2.([]JSON)
	doc3[0] = "asdf"
	assert.Equal(t, "[<nil> false true -1 3.4  foo [] map[]]", fmt.Sprint(doc))
	assert.Equal(t, "[asdf false true -1 3.4  foo [] map[]]", fmt.Sprint(doc2))
	assert.Equal(t, "[asdf false true -1 3.4  foo [] map[]]", fmt.Sprint(doc3))
}

func Test_Json_DeepEqual_TwoSimpleObjectsAreEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`{"foo": "bar"}`), &doc1)
	json.Unmarshal([]byte(`{"foo": "bar"}`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, true, isEqual)
}

func Test_Json_DeepEqual_TwoSimpleObjectsAreNotEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`{"foo": "bar"}`), &doc1)
	json.Unmarshal([]byte(`{"foo": "bar - 2"}`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, false, isEqual)
}

func Test_Json_DeepEqual_TwoNumbersAreEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`2.5`), &doc1)
	json.Unmarshal([]byte(`2.5`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, true, isEqual)
}

func Test_Json_DeepEqual_TwoNumbersAreNotEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`2.5`), &doc1)
	json.Unmarshal([]byte(`-2.5`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, false, isEqual)
}

func Test_Json_DeepEqual_TwoStringsAreEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`"asdf"`), &doc1)
	json.Unmarshal([]byte(`"asdf"`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, true, isEqual)
}

func Test_Json_DeepEqual_TwoStringsAreNotEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`"asdf"`), &doc1)
	json.Unmarshal([]byte(`"asdf - 2"`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, false, isEqual)
}

func Test_Json_DeepEqual_TwoTrueAreEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`true`), &doc1)
	json.Unmarshal([]byte(`true`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, true, isEqual)
}

func Test_Json_DeepEqual_TwoFalseAreEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`false`), &doc1)
	json.Unmarshal([]byte(`false`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, true, isEqual)
}

func Test_Json_DeepEqual_TwoBooleansAreNotEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`true`), &doc1)
	json.Unmarshal([]byte(`false`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, false, isEqual)
}

func Test_Json_DeepEqual_TwoNullsAreEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`null`), &doc1)
	json.Unmarshal([]byte(`null`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, true, isEqual)
}

func Test_Json_DeepEqual_TwoArrayAreEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`[1, 2, 3]`), &doc1)
	json.Unmarshal([]byte(`[1, 2, 3]`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, true, isEqual)
}

func Test_Json_DeepEqual_TwoArrayAreNotEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`[1, 2, 3, 4]`), &doc1)
	json.Unmarshal([]byte(`[1, 2, 3]`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, false, isEqual)
}

func Test_Json_DeepEqual_TwoArrayAreNotEqual2(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`[1, 2, 5]`), &doc1)
	json.Unmarshal([]byte(`[1, 2, 3]`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, false, isEqual)
}

func Test_Json_DeepEqual_TwoBigObjectsAreEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`{
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
	}`), &doc1)
	json.Unmarshal([]byte(`{
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
	}`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, true, isEqual)
}

func Test_Json_DeepEqual_TwoBigObjectsAreNotEqual(t *testing.T) {
	var doc1 interface{}
	var doc2 interface{}
	json.Unmarshal([]byte(`{
		"a": null,
		"b": false,
		"c": true,
		"d": 0,
		"e": 1.1,
		"f": "",
		"g": "THIS KEY IS NOT EQUAL",
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
	}`), &doc1)
	json.Unmarshal([]byte(`{
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
	}`), &doc2)
	isEqual := DeepEqual(doc1, doc2)
	assert.Equal(t, false, isEqual)
}

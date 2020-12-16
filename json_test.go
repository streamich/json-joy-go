package jsonjoy

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Copy_ReturnsACopyOfJson(t *testing.T) {
	b := []byte(`{"foo":"bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	doc2 := Copy(doc)
	assert.Equal(t, "map[foo:bar]", fmt.Sprint(doc))
	assert.Equal(t, "map[foo:bar]", fmt.Sprint(doc2))
}

func Test_Copy_ModifyingCopyDoesNotModifyOriginal(t *testing.T) {
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

func Test_Copy_ModifyingOriginalDoesNotModifyCopy(t *testing.T) {
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

func Test_Copy_CopiesArrayOfVariousTypesOfElements(t *testing.T) {
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

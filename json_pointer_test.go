package jsonjoy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DecodeReferenceToken_ReturnsSameStringIfThereAreNoEscapedChars(t *testing.T) {
	decoded := DecodeReferenceToken("foobar")
	assert.Equal(t, decoded, "foobar")
	decoded = DecodeReferenceToken("foo/bar")
	assert.Equal(t, decoded, "foo/bar")
	decoded = DecodeReferenceToken("foo~bar")
	assert.Equal(t, decoded, "foo~bar")
}

func Test_DecodeReferenceToken_DecodesSpecialChars(t *testing.T) {
	decoded := DecodeReferenceToken("foo~1bar")
	assert.Equal(t, decoded, "foo/bar")
	decoded = DecodeReferenceToken("foo~0bar")
	assert.Equal(t, decoded, "foo~bar")
	decoded = DecodeReferenceToken("~1~0foo~0bar~0~1")
	assert.Equal(t, decoded, "/~foo~bar~/")
	decoded = DecodeReferenceToken("~1~0")
	assert.Equal(t, decoded, "/~")
	decoded = DecodeReferenceToken("~0~1")
	assert.Equal(t, decoded, "~/")
	decoded = DecodeReferenceToken("~0")
	assert.Equal(t, decoded, "~")
	decoded = DecodeReferenceToken("~1")
	assert.Equal(t, decoded, "/")
}

func Test_EncodeReferenceToken_ReturnsSameStringIfThereAreNoScpecialChars(t *testing.T) {
	encoded := EncodeReferenceToken("foobar")
	assert.Equal(t, encoded, "foobar")
}

func Test_EncodeReferenceToken_EncodesSpecialChars(t *testing.T) {
	encoded := EncodeReferenceToken("foo/bar")
	assert.Equal(t, encoded, "foo~1bar")
	encoded = EncodeReferenceToken("foo~bar")
	assert.Equal(t, encoded, "foo~0bar")
	encoded = EncodeReferenceToken("~/")
	assert.Equal(t, encoded, "~0~1")
	encoded = EncodeReferenceToken("/~")
	assert.Equal(t, encoded, "~1~0")
	encoded = EncodeReferenceToken("~")
	assert.Equal(t, encoded, "~0")
	encoded = EncodeReferenceToken("/")
	assert.Equal(t, encoded, "~1")
}

func Test_ParseJSONPointer_ReturnsEmptyArrayOnRootPointer(t *testing.T) {
	pointer, err := ParseJSONPointer("")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 0)
}

func Test_ParseJSONPointer_ReturnsErrorIfPointerDoesNotStartWithSlash(t *testing.T) {
	pointer, err := ParseJSONPointer("foo/bar")
	assert.Nil(t, pointer)
	assert.NotNil(t, err)
}

func Test_ParseJSONPointer_ParsesASingleStepPointer(t *testing.T) {
	pointer, err := ParseJSONPointer("/foo")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 1)
	assert.Equal(t, (pointer)[0], "foo")
}

func Test_ParseJSONPointer_ParsesAMultipleStepPointer(t *testing.T) {
	pointer, err := ParseJSONPointer("/foo/bar/baz")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 3)
	assert.Equal(t, pointer[0], "foo")
	assert.Equal(t, pointer[1], "bar")
	assert.Equal(t, pointer[2], "baz")
}

func Test_ParseJSONPointer_DecodesTokens(t *testing.T) {
	pointer, err := ParseJSONPointer("/foo~1bar")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 1)
	assert.Equal(t, pointer[0], "foo/bar")
	pointer, err = ParseJSONPointer("/foo~1/bar/ba~0~1~0z")
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

func Test_JSONPointer_Find_ReturnsRootDocument(t *testing.T) {
	tokens := JSONPointer{}
	b := []byte(`{"foo": "bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Find(doc)
	assert.Nil(t, err)
	assert.Equal(t, value, doc)
}

func Test_JSONPointer_Find_ReturnsAFirstLevelKey(t *testing.T) {
	tokens := JSONPointer{"foo"}
	b := []byte(`{"foo": "bar"}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Find(doc)
	assert.Nil(t, err)
	assert.Equal(t, "bar", value)
}

func Test_JSONPointer_Find_FindsADeepStringKeyInObjects(t *testing.T) {
	tokens := JSONPointer{"foo", "bar", "baz"}
	b := []byte(`{"foo": {"bar": {"baz": "qux"}}}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Find(doc)
	assert.Nil(t, err)
	assert.Equal(t, "qux", value)
}

func Test_JSONPointer_Find_ReturnsErrorNotFoundWhenLocatingMissingValue(t *testing.T) {
	tokens := JSONPointer{"foo2"}
	b := []byte(`{"foo": {"bar": {"baz": "qux"}}}`)
	var doc interface{}
	json.Unmarshal(b, &doc)
	value, err := tokens.Find(doc)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, value)
}

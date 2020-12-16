package jsonjoy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeReferenceTokenReturnsSameStringIfThereAreNoEscapedChars(t *testing.T) {
	decoded := DecodeReferenceToken("foobar")
	assert.Equal(t, decoded, "foobar")
	decoded = DecodeReferenceToken("foo/bar")
	assert.Equal(t, decoded, "foo/bar")
	decoded = DecodeReferenceToken("foo~bar")
	assert.Equal(t, decoded, "foo~bar")
}

func TestDecodeReferenceTokenDecodesSpecialChars(t *testing.T) {
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

func TestEncodeReferenceTokenReturnsSameStringIfThereAreNoScpecialChars(t *testing.T) {
	encoded := EncodeReferenceToken("foobar")
	assert.Equal(t, encoded, "foobar")
}

func TestEncodeReferenceTokenEncodesSpecialChars(t *testing.T) {
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

func TestParseJSONPointerReturnsEmptyArrayOnRootPointer(t *testing.T) {
	pointer, err := ParseJSONPointer("")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 0)
}

func TestParseJSONPointerReturnsErrorIfPointerDoesNotStartWithSlash(t *testing.T) {
	pointer, err := ParseJSONPointer("foo/bar")
	assert.Nil(t, pointer)
	assert.NotNil(t, err)
}

func TestParseJSONPointerParsesASingleStepPointer(t *testing.T) {
	pointer, err := ParseJSONPointer("/foo")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 1)
	assert.Equal(t, (pointer)[0], "foo")
}

func TestParseJSONPointerParsesAMultipleStepPointer(t *testing.T) {
	pointer, err := ParseJSONPointer("/foo/bar/baz")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 3)
	assert.Equal(t, pointer[0], "foo")
	assert.Equal(t, pointer[1], "bar")
	assert.Equal(t, pointer[2], "baz")
}

func TestParseJSONPointerDecodesTokens(t *testing.T) {
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

func TestFormatJSONPointerFormatsTokensIntoJsonPointer(t *testing.T) {
	str := FormatJSONPointer([]string{"foo", "bar", "baz"})
	assert.Equal(t, str, "/foo/bar/baz")
}

func TestFormatJSONPointerFormatsASingleToken(t *testing.T) {
	str := FormatJSONPointer([]string{"aga"})
	assert.Equal(t, str, "/aga")
}

func TestFormatJSONPointerFormatsARootPointer(t *testing.T) {
	str := FormatJSONPointer([]string{})
	assert.Equal(t, str, "")
}

func TestFormatJSONPointerEncodesSpecialChars(t *testing.T) {
	str := FormatJSONPointer([]string{"foo/bar"})
	assert.Equal(t, str, "/foo~1bar")
	str = FormatJSONPointer([]string{"foo/bar", "/", "~", "a~b/"})
	assert.Equal(t, str, "/foo~1bar/~1/~0/a~0b~1")
}

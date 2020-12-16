package jsonjoy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReturnsEmptyArrayOnRootPointer(t *testing.T) {
	pointer, err := ParseJSONPointer("")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 0)
}

func TestReturnsErrorIfPointerDoesNotStartWithSlash(t *testing.T) {
	pointer, err := ParseJSONPointer("foo/bar")
	assert.Nil(t, pointer)
	assert.NotNil(t, err)
}

func TestParsesASingleStepPointer(t *testing.T) {
	pointer, err := ParseJSONPointer("/foo")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 1)
	assert.Equal(t, (pointer)[0], "foo")
}

func TestParsesAMultipleStepPointer(t *testing.T) {
	pointer, err := ParseJSONPointer("/foo/bar/baz")
	assert.Nil(t, err)
	assert.NotNil(t, pointer)
	assert.Equal(t, len(pointer), 3)
	assert.Equal(t, pointer[0], "foo")
	assert.Equal(t, pointer[1], "bar")
	assert.Equal(t, pointer[2], "baz")
}

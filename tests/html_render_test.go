package tests

import (
	"testing"

	"github.com/saasuke-labs/gsx"
	"github.com/stretchr/testify/assert"
)

func TestParseString_BasicComponent(t *testing.T) {
	src := `<Tag name="Go" />`
	output, warnings, err := gsx.ParseString("Tag", src)
	assert.NoError(t, err)
	assert.Equal(t, `{{ template "Tag" (props "name" "Go") }}`, output)
	assert.Empty(t, warnings)
}
func TestParseString_UnquotedAttribute(t *testing.T) {
	src := `<Tag name=Go />`
	output, warnings, err := gsx.ParseString("Tag", src)
	assert.NoError(t, err)
	assert.Equal(t, `{{ template "Tag" (props "name" "Go") }}`, output)
	assert.Len(t, warnings, 1)
	assert.Equal(t, "unquoted value for attribute \"name\"", warnings[0].Message)
}

func TestParseString_BasicComponentWithCurlyBraces(t *testing.T) {
	src := `<Tag name={ value } />`
	output, warnings, err := gsx.ParseString("Tag", src)
	assert.NoError(t, err)
	assert.Equal(t, `{{ template "Tag" (props "name" .value) }}`, output)
	assert.Empty(t, warnings)
}

func TestParseString_BasicComponentWithNumber(t *testing.T) {
	src := `<Tag name={ 5 } />`
	output, warnings, err := gsx.ParseString("Tag", src)
	assert.NoError(t, err)
	assert.Equal(t, `{{ template "Tag" (props "name" 5) }}`, output)
	assert.Empty(t, warnings)
}

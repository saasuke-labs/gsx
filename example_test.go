package gsx_test

import (
	"fmt"

	"github.com/saasuke-labs/gsx"
)

func ExampleParseString_basicComponent() {
	src := `<Tag name="Go" />`
	output, warnings, err := gsx.ParseString("Tag", src)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(output)
	fmt.Println(len(warnings) == 0)

	// Output:
	// {{ template "Tag" (dict "name" "Go") }}
	// true
}

func ExampleParseString_unquotedAttribute() {
	src := `<Tag name=Go />`
	output, warnings, err := gsx.ParseString("Tag", src)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(output)
	fmt.Println(warnings[0].Message)

	// Output:
	// {{ template "Tag" (dict "name" "Go") }}
	// unquoted value for attribute "name"
}

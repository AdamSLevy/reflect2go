// Copyright 2020 Adam S Levy
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package reflect2go_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/AdamSLevy/reflect2go"
)

func ExampleGoFile() {

	// We will generate a new type based on this predefined Original type.
	// The updated type will include an additional `env` struct field tag
	// based on the existing `json` tag.
	type Original struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	typ := reflect.TypeOf(Original{})
	var fields []reflect.StructField
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)

		val := strings.ToUpper(sf.Tag.Get("json"))
		sf.Tag += reflect.StructTag(fmt.Sprintf(" env:%q", val))

		fields = append(fields, sf)
	}

	newTyp := reflect.StructOf(fields)

	gofile, err := reflect2go.NewGoFile("pkg", reflect2go.DefineType("UpdatedTags", newTyp))
	if err != nil {
		panic(err)
	}

	if err := gofile.Render(os.Stdout); err != nil {
		panic(err)
	}

	// Output:
	// package pkg
	//
	// type UpdatedTags struct {
	// 	A int `json:"a" env:"A"`
	// 	B int `json:"b" env:"B"`
	// }
}

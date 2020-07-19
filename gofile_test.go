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

package reflect2go

import (
	"bytes"
	"go/format"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Import is an Option only for use in tests that adds an Import path.
func Import(paths ...string) Option {
	return func(gf *GoFile) error {
		gf.Imports = append(gf.Imports, paths...)
		return nil
	}
}

type GoFileTest struct {
	Name     string
	GoFile   GoFile
	Expected string
}

func mustFormat(src string) string {
	gofmt, err := format.Source([]byte(src))
	if err != nil {
		panic(err)
	}
	return string(gofmt)
}

var goFileTests = []GoFileTest{{
	Name:   "empty",
	GoFile: Must(NewGoFile("path/to/pkg")),
	Expected: mustFormat(`
package pkg

`),
}, {
	Name:   "one import",
	GoFile: Must(NewGoFile("path/to/pkg", Import("path/to/import"))),
	Expected: `
package pkg

import "path/to/import"
`,
}, {
	Name: "multiple imports",
	GoFile: Must(NewGoFile("path/to/pkg",
		Import("fmt"),
		Import("fmt"),
		Import("sync"),
		Import("github.com/example/example"),
	)),
	Expected: `
package pkg

import (
        "fmt"
        "github.com/example/example"
        "sync"
)
`,
}, {
	Name: "typedef",
	GoFile: Must(NewGoFile("github.com/AdamSLevy/reflect2go",
		CommentLineWrap(50),
		DefineType("S", reflect.TypeOf(struct {
			A int
			B *[]float64
			C []*float64 `json:"b" other:"tag"`
			D *struct {
				time.Time
			}
			E chan int
			F chan struct {
				X int
			}
			G map[string]interface {
				Bar()
				Foo()
			}
			H    chan<- int
			I    <-chan int
			Time time.Time
			sync.RWMutex
			GoFileTest
		}{}), Comment("S is a test type with a really really really really really really really really really really really really really really long comment.")))),
	Expected: `
package reflect2go

import (
        "time"
        "sync"
)

// S is a test type with a really really really
// really really really really really really
// really really really really really long
// comment.
type S struct {
        A int
        B *[]float64
        C []*float64 ` + "`" + `json:"b" other:"tag"` + "`" + `
        D *struct {
                time.Time
        }
	E chan int
	F chan struct {
		X int
	}
        G    map[string]interface{
                Bar()
                Foo()
        }
        H    chan<- int
        I    <-chan int
        Time time.Time
        sync.RWMutex
        GoFileTest
}`,
}, {
	Name: "typedef",
	GoFile: Must(NewGoFile("path/to/pkg", PkgName("foo"),
		PreambleComment("LICENSE GOES HERE"),
		PkgComment("Package foo defines things for completing tasks that involve problems to do with other things."),
		DefineType("S", reflect.TypeOf(struct {
			A int
			B *[]float64
			C []*float64 `json:"b" other:"tag"`
			D *struct {
				time.Time
			}
			Time time.Time
			sync.RWMutex
			GoFileTest
		}{})))),
	Expected: `
// LICENSE GOES HERE

// Package foo defines things for completing tasks that involve problems to do
// with other things.
package foo

import (
        "time"
        "github.com/AdamSLevy/reflect2go"
        "sync"
)

type S struct {
        A int
        B *[]float64
        C []*float64 ` + "`" + `json:"b" other:"tag"` + "`" + `
        D *struct {
                time.Time
        }
        Time time.Time
        sync.RWMutex
        reflect2go.GoFileTest
}`,
}}

func TestGoFile(t *testing.T) {

	for _, test := range goFileTests {
		t.Run(test.Name, func(t *testing.T) {
			var buf bytes.Buffer
			require.NoError(t, test.GoFile.Render(&buf))
			require.Equal(t, mustFormat(test.Expected), buf.String())
		})
	}

	_, err := NewGoFile("")
	require.Error(t, err, "NewGoFile: empty package name")
	gf, err := NewGoFile("pkg",
		DefineType("dup", reflect.TypeOf(struct{}{})),
		DefineType("dup", reflect.TypeOf(struct{}{})),
	)
	require.Error(t, err, "NewGoFile: duplicate type name")
	require.Panics(t, func() { Must(gf, err) }, "Must")

}

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
	"fmt"
	"reflect"
)

// typeDef is the data passed to the gofile template for a type declaration.
type typeDef struct {
	Comment    string
	Name       string
	Definition string
}

// TypeOption is an option for DefineType.
type TypeOption func(*typeDef)

// Comment adds a top level comment to a defined type.
func Comment(cmt string) TypeOption {
	return func(t *typeDef) {
		t.Comment = cmt
	}
}

// DefineType adds a new type declaration to the GoFile. An error is returned
// if the name has already been used on a type declaration within the GoFile.
func (f *GoFile) DefineType(name string, typ reflect.Type, opts ...TypeOption) error {
	if _, ok := f.Types[name]; ok {
		return fmt.Errorf("type already defined: %v", name)
	}

	typDef := typeDef{Name: name}
	for _, opt := range opts {
		opt(&typDef)
	}

	var buf bytes.Buffer
	f.render(&buf, typ)
	typDef.Definition = buf.String()

	f.Types[name] = typDef

	return nil
}

// render the definition of a single Go type into buf.
func (f *GoFile) render(buf *bytes.Buffer, typ reflect.Type) {

	if typ.Kind() == reflect.Ptr {
		buf.WriteString("*")
		f.render(buf, typ.Elem())
		return
	}

	// If a type is named, we can use that name as the definition.
	if name := typ.Name(); name != "" {

		// If the type does not belong to the same package as this
		// file, we must use the form of the type name that includes
		// the package identifier. e.g. pkg.Type vs Type
		if typ.PkgPath() != "" && typ.PkgPath() != f.pkgPath {
			name = fmt.Sprintf("%v", typ)
			f.Imports = append(f.Imports, typ.PkgPath())
		}

		buf.WriteString(name)
		return
	}

	// Since the type is not named, we will render the type directly.

	// Because fmt does not render struct tags properly, we need to handle
	// struct types manually. This also means manually processing any
	// composite types that may include unnamed struct types.

	switch typ.Kind() {
	case reflect.Slice:
		f.renderSlice(buf, typ)

	case reflect.Chan:
		f.renderChan(buf, typ)

	case reflect.Map:
		f.renderMap(buf, typ)

	case reflect.Struct:
		f.renderStruct(buf, typ)

	default:
		fmt.Fprintf(buf, "%v", typ)
	}

	return
}

func (f *GoFile) renderSlice(buf *bytes.Buffer, typ reflect.Type) {
	buf.WriteString("[]")
	f.render(buf, typ.Elem())
}

func (f *GoFile) renderChan(buf *bytes.Buffer, typ reflect.Type) {
	var chanDir string
	switch typ.ChanDir() {
	case reflect.RecvDir:
		chanDir = "<-chan "
	case reflect.SendDir:
		chanDir = "chan<- "
	case reflect.BothDir:
		chanDir = "chan "
	}
	buf.WriteString(chanDir)
	f.render(buf, typ.Elem())
}

func (f *GoFile) renderMap(buf *bytes.Buffer, typ reflect.Type) {
	buf.WriteString("map[")
	f.render(buf, typ.Key())
	buf.WriteString("]")
	f.render(buf, typ.Elem())
}

func (f *GoFile) renderStruct(buf *bytes.Buffer, typ reflect.Type) {
	buf.WriteString("struct {\n")
	for i := 0; i < typ.NumField(); i++ {
		// "\t<name> <type> `<tag>`\n"

		// <name>, omitted if field is embedded.
		sf := typ.Field(i)
		if sf.Anonymous {
			sf.Name = ""
		}
		fmt.Fprintf(buf, "\t%v ", sf.Name)

		// <type>
		f.render(buf, sf.Type)

		// `<tag>`, omitted if empty.
		if len(sf.Tag) > 0 {
			fmt.Fprintf(buf, " `%v`", sf.Tag)
		}

		buf.WriteString("\n")
	}
	buf.WriteString("}")
}

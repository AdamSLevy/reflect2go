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

// Package reflect2go generates the Go code for a given reflect.Type.
//
// This is useful when you want to generate the Go code for a dynamically
// defined type, like with reflect.StructOf.
//
// An example use case is a program that updates or manipulates the tags on the
// fields of an existing struct.
package reflect2go

import (
	"fmt"
	"path"
	"reflect"

	"github.com/bbrks/wrap"
)

// GoFile represents and renders a file with Golang code.
type GoFile struct {
	pkgPath         string
	wrapper         wrap.Wrapper
	commentLineWrap int
	goFileTmplInput
}

// goFileTmplInput is the input passed into the gofile template.
type goFileTmplInput struct {
	Preamble   string
	PkgComment string
	PkgName    string
	Imports    []string
	Types      map[string]typeDef
}

// Must returns the GoFile but panics if err is not nil.
//
//      err := Must(NewGoFile("pkg",
//              DefineType("X", reflect.TypeOf(struct{}{})),
//              DefineType("Y", reflect.TypeOf(struct{}{}), Comment("Y is a type.")),
//      )).Render(os.Stdout)
func Must(gf GoFile, err error) GoFile {
	if err != nil {
		panic(err)
	}
	return gf
}

// NewGoFile returns a GoFile that belongs to the package with import path
// pkgPath.
//
// The pkgPath is used to identify reflected types that belong to the same
// package, and thus don't require a package prefix in their declaration (e.g.
// pkg.Type vs Type).
//
// The file's package declaration will use path.Base(pkgPath) by default. Use
// the Option PkgName to override this.
//
// All comments are word wrapped at the limit set by CommentLineWrap, which
// defaults to DefaultCommentLineWrap.
//
// The order of Options has no significance or affect on the output. Comments
// are wrapped and defined types are sorted alphabetically prior to output.
func NewGoFile(pkgPath string, opts ...Option) (GoFile, error) {

	if pkgPath == "" {
		return GoFile{}, fmt.Errorf("package path is empty")
	}

	gf := GoFile{
		pkgPath,
		newCommentWrapper(),
		DefaultCommentLineWrap,
		goFileTmplInput{
			PkgName: path.Base(pkgPath),
			Types:   make(map[string]typeDef),
		},
	}

	for _, opt := range opts {
		if err := opt(&gf); err != nil {
			return GoFile{}, err
		}
	}

	return gf, nil
}

// newCommentWrapper returns a wrap.Wrapper initialized to wrap comments with a
// leading "// " and trim any user provided "// ".
func newCommentWrapper() wrap.Wrapper {
	wrapper := wrap.NewWrapper()
	wrapper.OutputLinePrefix = "// "
	wrapper.TrimInputPrefix = "// "
	wrapper.StripTrailingNewline = true
	return wrapper
}

// Option is an optional setting for NewGoFile.
type Option func(*GoFile) error

// PkgName overrides the package name, which by default is derived from the
// path.Base of the package path passed to NewGoFile.
func PkgName(name string) Option {
	return func(f *GoFile) error {
		f.PkgName = name
		return nil
	}
}

// PkgComment sets a top level package comment which immediately precedes the
// package declaration.
func PkgComment(cmnt string) Option {
	return func(f *GoFile) error {
		f.PkgComment = cmnt
		return nil
	}
}

// PreambleComment adds a top level comment that precedes all other content in
// the rendered GoFile. A blank line is added immediately after the comment to
// avoid it being confused with a package comment.
func PreambleComment(cmnt string) Option {
	return func(f *GoFile) error {
		f.Preamble = cmnt
		return nil
	}
}

// DefineType defines a type on the NewGoFile.
//
// It is equivalent to calling GoFile.DefineType on an existing GoFile.
func DefineType(name string, typ reflect.Type, opts ...TypeOption) Option {
	return func(f *GoFile) error {
		return f.DefineType(name, typ, opts...)
	}
}

// DefaultCommentLineWrap is the default CommentLineWrap limit used to wrap all
// comments.
const DefaultCommentLineWrap = 79

// CommentLineWrap sets the limit used to wrap all comments on word boundaries.
func CommentLineWrap(limit int) Option {
	return func(f *GoFile) error {
		f.commentLineWrap = limit
		return nil
	}
}

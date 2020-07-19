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
	"go/format"
	"io"
	"text/template"
)

const gofile = `
{{- if ne .PkgComment "" }}
{{ wrap .Preamble }}

{{end }}
{{- if ne .PkgComment "" }}
{{ wrap .PkgComment }}
{{- end }}
package {{ .PkgName }}
{{- if eq (len .Imports) 1 }}

import {{ printf "%q" (index .Imports 0) }}
{{- else if gt (len .Imports) 1 }}

import (
{{- range .Imports }}
        {{ printf "%q" . }}
{{- end}}
)
{{- end}}

{{- range .Types }}
{{- if ne .Comment "" }}
{{ wrap .Comment }}
{{- end }}
type {{ .Name }} {{ .Definition }}
{{- end}}
`

// Render writes the GoFile to w.
func (f GoFile) Render(w io.Writer) error {
	gofileTmpl := template.Must(template.New("gofile").Funcs(template.FuncMap{
		"wrap": func(s string) string {
			return f.wrapper.Wrap(s, f.commentLineWrap)
		},
	}).Parse(gofile))

	var buf bytes.Buffer
	if err := gofileTmpl.Execute(&buf, f); err != nil {
		// Prior validation should ensure the template always executes.
		panic(err)
	}

	gofmt, err := format.Source(buf.Bytes())
	if err != nil {
		// The template should always generate valid Go code.
		panic(fmt.Errorf("gofmt: %w \ninput:\n%v", err, buf.String()))
	}

	_, err = w.Write(gofmt)
	return err
}

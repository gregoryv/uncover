// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cover

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/tabwriter"
	"text/template"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

func WriteOutput(profiles []*Profile) error {
	var out *bufio.Writer
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()

	tabber := tabwriter.NewWriter(out, 1, 8, 4, ' ', 0)
	defer tabber.Flush()

	var total, covered int64
	for _, profile := range profiles {
		fn := profile.FileName
		file, err := findFile(fn)
		if err != nil {
			return err
		}
		funcs, err := findFuncs(file)
		if err != nil {
			return err
		}
		// Now match up functions and profile blocks.
		for _, f := range funcs {
			c, t := f.coverage(profile)
			// Only show uncovered funcs
			if percent(c, t) < 100 {
				// todo print the func signature
				sign := fmt.Sprintf("%sfunc %s(...) ...%s", green, f.Name, reset)
				fmt.Fprintf(tabber, "%s:%d\n%s ", fn, f.startLine, sign)
				Write(profile, tabber, f)
			}
			total += t
			covered += c
		}
	}
	fmt.Fprintf(tabber, "total:\t(statements)\t%.1f%%\n", percent(covered, total))

	return nil
}

// Write reads the profile data from profile and generates colored
// vt100 output to stdout.
func Write(profile *Profile, tabber *tabwriter.Writer, fe *FuncExtent) error {
	var d templateData

	// Read profile data
	fn := profile.FileName
	file, err := findFile(fn)
	if err != nil {
		return err
	}
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("can't read %q: %v", fn, err)
	}
	// Filter boundaries according to fe
	funcBoundaries := make([]Boundary, 0)
	for _, b := range profile.Boundaries(src) {
		if b.Offset > fe.End {
			break
		}
		if b.Offset >= fe.Offset {
			funcBoundaries = append(funcBoundaries, b)
		}
	}
	// Write colored source to buffer
	var buf bytes.Buffer
	err = vt100Gen(&buf, src, funcBoundaries)
	if err != nil {
		return err
	}
	d.Files = append(d.Files, &templateFile{
		Body: buf.String(),
	})

	err = colorTemplate.Execute(tabber, d)
	if err != nil {
		return err
	}

	return nil
}

// vt100Gen generates an coverage report with the provided filename,
// source code, and tokens, and writes it to the given Writer.
func vt100Gen(w io.Writer, src []byte, boundaries []Boundary) error {
	dst := bufio.NewWriter(w)
	var color string
	show := false
	for i := range src {
		for len(boundaries) > 0 && boundaries[0].Offset == i {
			show = true
			b := boundaries[0]
			if b.Start {
				color = red
				if b.Count >= 1 {
					color = green
				}
				fmt.Fprint(dst, color)
			} else {
				dst.WriteString(reset)
			}
			boundaries = boundaries[1:]
		}
		if show && len(boundaries) != 0 {
			dst.WriteByte(src[i])
		}
		if len(boundaries) == 0 {
			break
		}
	}
	return dst.Flush()
}

var colorTemplate = template.Must(template.New("").Parse(tpl))

type templateData struct {
	Files []*templateFile
	Set   bool
}

type templateFile struct {
	Name string
	Body string
}

const tpl = `{{range $i, $f := .Files}}{{$f.Body}}
{{end}}
`

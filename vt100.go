// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uncover

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/tabwriter"
	"text/template"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

func Report(profiles []*Profile, out io.Writer) (coverage float64, err error) {
	tabw := tabwriter.NewWriter(out, 1, 8, 4, ' ', 0)
	defer tabw.Flush()
	return WriteOutput(profiles, tabw)
}

var OnlyShow string

func WriteOutput(profiles []*Profile, out io.Writer) (coverage float64, err error) {
	var total, covered int64
	var file string
	var funcs []*FuncExtent
	for _, profile := range profiles {
		fn := profile.FileName
		file, err = findFile(fn)
		if err != nil {
			return
		}
		funcs, err = findFuncs(file)
		if err != nil {
			return
		}
		// filter funcs
		if OnlyShow != "" {
			tmp := make([]*FuncExtent, 0)
			for _, fe := range funcs {
				if strings.Index(fe.Name, OnlyShow) >= 0 {
					tmp = append(tmp, fe)
				}
			}
			funcs = tmp
		}

		// Match up functions and profile blocks.
		for _, f := range funcs {
			c, t := f.coverage(profile)
			// Only show uncovered funcs
			if percent(c, t) < 100 {
				// todo print the func signature
				sign := fmt.Sprintf("%sfunc %s(...) ...%s", green, f.Name, reset)
				fmt.Fprintf(out, "%s:%d\n%s ", fn, f.startLine, sign)
				Write(profile, out, f)
			}
			total += t
			covered += c
		}
	}
	coverage = percent(covered, total)
	fmt.Fprintf(out, "total:\t(statements)\t%.1f%%\n", coverage)
	return
}

// Write reads the profile data from profile and generates colored
// vt100 output to stdout.
func Write(profile *Profile, out io.Writer, fe *FuncExtent) error {
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

	err = colorTemplate.Execute(out, d)
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

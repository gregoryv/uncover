// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cover

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

func a() {
	// hepp
}

// Write reads the profile data from profile and generates colored
// vt100 output to stdout.
func Write(profile string) error {
	profiles, err := ParseProfiles(profile)
	if err != nil {
		return err
	}

	var d templateData

	for _, profile := range profiles {
		fn := profile.FileName
		file, err := findFile(fn)
		if err != nil {
			return err
		}
		src, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("can't read %q: %v", fn, err)
		}
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "FILE: %s\n-----------------------------------\n", file)
		err = vt100Gen(&buf, src, profile.Boundaries(src))
		if err != nil {
			return err
		}
		d.Files = append(d.Files, &templateFile{
			Name:     fn,
			Body:     template.HTML(buf.String()),
			Coverage: percentCovered(profile),
		})
	}

	err = htmlTemplate.Execute(os.Stdout, d)
	if err != nil {
		return err
	}

	return nil
}

// percentCovered returns, as a percentage, the fraction of the statements in
// the profile covered by the test run.
// In effect, it reports the coverage of a given source file.
func percentCovered(p *Profile) float64 {
	var total, covered int64
	for _, b := range p.Blocks {
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	if total == 0 {
		return 0
	}
	return float64(covered) / float64(total) * 100
}

// vt100Gen generates an coverage report with the provided filename,
// source code, and tokens, and writes it to the given Writer.
func vt100Gen(w io.Writer, src []byte, boundaries []Boundary) error {
	dst := bufio.NewWriter(w)
	var color string
	for i := range src {
		for len(boundaries) > 0 && boundaries[0].Offset == i {
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
		dst.WriteByte(src[i])
	}
	return dst.Flush()
}

var htmlTemplate = template.Must(template.New("html").Parse(tmplHTML))

type templateData struct {
	Files []*templateFile
	Set   bool
}

type templateFile struct {
	Name     string
	Body     template.HTML
	Coverage float64
}

const tmplHTML = `
{{range $i, $f := .Files}}
{{$f.Body}}
{{end}}
`

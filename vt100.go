// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uncover

import (
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/gregoryv/nexus"
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
		fn := profile.Filename
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
			p := percent(c, t)
			switch {
			case p == 0:
				fmt.Fprintf(out, "%s:%d\n", fn, f.startLine)
				fmt.Fprintf(out, "%s%s - UNCOVERED%s\n\n", red, f.Name, reset)
			case p < 100:
				fmt.Fprintf(out, "%s:%d\n", fn, f.startLine)
				Write(out, profile, f)
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
func Write(out io.Writer, profile *Profile, f *FuncExtent) error {
	// Read profile data
	fn := profile.Filename
	file, err := findFile(fn)
	if err != nil {
		return err
	}
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("can't read %q: %v", fn, err)
	}
	// Filter boundaries according to the extent
	funcBoundaries := make([]Boundary, 0)
	for _, b := range profile.Boundaries(src) {
		if b.Offset > f.End {
			break
		}
		if b.Offset >= f.Offset {
			funcBoundaries = append(funcBoundaries, b)
		}
	}
	// Write colored source to buffer
	err = vt100Gen(out, src, f.Name, funcBoundaries)
	return err
}

// vt100Gen generates an coverage report with the provided filename,
// source code, and tokens, and writes it to the given Writer.
// boundaries must contain pairs beginning and end
func vt100Gen(w io.Writer, src []byte, sign string, boundaries []Boundary) error {
	p, err := nexus.NewPrinter(w)
	p.Print(green, sign)
	for i := 0; i < len(boundaries); i += 2 {
		start := boundaries[i]
		end := boundaries[i+1]
		if start.Count == 0 {
			p.Print(red)
		} else {
			p.Print(green)
		}
		// handle empty blocks
		k := start.Offset
		if i > 0 {
			k = boundaries[i-1].Offset
		}
		p.Print(string(src[k:end.Offset]))
	}
	p.Printf("\n%s}%s\n\n", green, reset)
	p.Print(reset)
	return *err
}

// colors
var (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

// findFile finds the location of the named file in GOROOT, GOPATH
// etc.
func findFile(file string) (string, error) {
	dir, file := filepath.Split(file)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("can't find %q: %v", file, err)
	}
	return filepath.Join(pkg.Dir, file), nil
}

func percent(covered, total int64) float64 {
	if total == 0 {
		total = 1 // Avoid zero denominator.
	}
	return 100.0 * float64(covered) / float64(total)
}

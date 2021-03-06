// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements the visitor that computes the (line,
// column)-(line-column) range for each function.

package uncover

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

// coverage returns the fraction of the statements in the function
// that were covered, as a numerator and denominator.
func (f *FuncExtent) coverage(profile *Profile) (num, den int64) {
	// We could avoid making this n^2 overall by doing a single scan
	// and annotating the functions, but the sizes of the data
	// structures is never very large and the scan is almost
	// instantaneous.
	var covered, total int64
	// The blocks are sorted, so we can stop counting as soon as we
	// reach the end of the relevant block.
	for _, b := range profile.Blocks {
		if b.StartLine > f.endLine || (b.StartLine == f.endLine && b.StartCol >= f.endCol) {
			// Past the end of the function.
			break
		}
		if b.EndLine < f.startLine || (b.EndLine == f.startLine && b.EndCol <= f.startCol) {
			// Before the beginning of the function
			continue
		}
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	return covered, total
}

// findFuncs parses the file and returns a slice of FuncExtent descriptors.
func findFuncs(name string) ([]*FuncExtent, error) {
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return nil, err
	}
	visitor := &FuncVisitor{
		fset:    fset,
		name:    name,
		astFile: parsedFile,
	}
	ast.Walk(visitor, visitor.astFile)
	return visitor.funcs, nil
}

// FuncExtent describes a function's extent in the source by file and
// position.
type FuncExtent struct {
	Name      string
	Decl      *ast.FuncDecl
	startLine int
	startCol  int
	endLine   int
	endCol    int
	Offset    int
	End       int
}

// FuncVisitor implements the visitor that builds the function
// position list for a file.
type FuncVisitor struct {
	fset    *token.FileSet
	name    string // Name of file.
	astFile *ast.File
	funcs   []*FuncExtent
}

// Visit implements the ast.Visitor interface.
func (v *FuncVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body == nil {
			// Do not count declarations of assembly functions.
			break
		}
		start := v.fset.Position(n.Pos())
		end := v.fset.Position(n.End())

		fe := &FuncExtent{
			Name:      v.signature(n),
			Decl:      n,
			startLine: start.Line,
			startCol:  start.Column,
			endLine:   end.Line,
			endCol:    end.Column,
			Offset:    start.Offset,
			End:       end.Offset,
		}
		v.funcs = append(v.funcs, fe)
	}
	return v
}

func (v *FuncVisitor) signature(n *ast.FuncDecl) string {
	var w bytes.Buffer
	printer.Fprint(&w, v.fset, n)
	fn := w.String()
	// Using the printer may result in source being differently formatted
	// than the incoming declaration. This happens when we have aligned bodies
	//
	// func a()          {}
	// func longerName() {}
	//
	// printer.Fprint(...) of func a declaration would result in less whitespace
	// between ) {
	//
	i := strings.Index(fn, "\n")
	switch {
	case i == -1:
		// oneliner
		i = strings.Index(fn, "\t{")
		return fn[:i]
	default:
		i := strings.Index(fn, "{\n")
		// tidy multiline signatures
		sign := strings.TrimSpace(fn[:i])
		sign = strings.ReplaceAll(sign, "\n", "")
		return strings.ReplaceAll(sign, "\t", " ")
	}
}

package main

//go:generate stamp -clfile ../../CHANGELOG.md -go build_stamp.go

import (
	"flag"
	"github.com/gregoryv/cover"
	"github.com/gregoryv/stamp"
	"os"
)

func init() {
	stamp.InitFlags()
}

func main() {
	flag.Parse()
	stamp.AsFlagged()

	cover.Write(os.Args[1])
}

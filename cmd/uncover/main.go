package main

//go:generate stamp -clfile ../../CHANGELOG.md -go build_stamp.go

import (
	"flag"
	"github.com/gregoryv/stamp"
	"github.com/gregoryv/uncover"
	"os"
)

func init() {
	stamp.InitFlags()
}

func main() {
	flag.Parse()
	stamp.AsFlagged()
	profiles, err := uncover.ParseProfiles(flag.Arg(0))
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}
	uncover.Report(profiles, os.Stdout)
}

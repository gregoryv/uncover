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
	profiles, err := cover.ParseProfiles(flag.Arg(0))
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}
	cover.WriteOutput(profiles)
}

package main

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

import (
	"flag"
	"fmt"
	"os"

	"github.com/gregoryv/stamp"
	"github.com/gregoryv/uncover"
)

func init() {
	stamp.InitFlags()
}

func main() {
	min := flag.Float64("min", 0.0, "Fail if total coverage(%) is below min")
	flag.Parse()
	stamp.AsFlagged()
	profiles, err := uncover.ParseProfiles(flag.Arg(0))
	uncover.OnlyShow = flag.Arg(1)
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}
	coverage, _ := uncover.Report(profiles, os.Stdout)
	if coverage < *min {
		fmt.Printf("coverage to low: expected >= %v%%\n", *min)
		os.Exit(1)
	}
}

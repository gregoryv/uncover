package main

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

import (
	"fmt"
	"log"
	"os"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/uncover"
	"github.com/gregoryv/wolf"
)

var cmd wolf.Command = wolf.NewOSCmd()

func main() {
	var (
		cli = cmdline.NewBasicParser()
		min = cli.Option(
			"-min",
			"Fail if total coverage(%) is below min",
		).Float64(0.0)
		version  = cli.Flag("-v, --version")
		profile  = cli.NamedArg("PROFILE").String("")
		onlyShow = cli.NamedArg("FUNC").String("")
	)
	cli.Parse()

	uncover.OnlyShow = onlyShow
	log.SetFlags(0)
	log.SetOutput(cmd.Stderr())

	if version {
		fmt.Fprintln(cmd.Stdout(), uncover.Version())
		cmd.Exit(0)
	}

	profiles, err := uncover.ParseProfiles(profile)
	if err != nil {
		cmd.Fatal(err)
		return
	}

	coverage, err := uncover.Report(profiles, os.Stdout)
	if err != nil {
		log.Fatal("failed to find files; please set GOROOT")
	}

	if coverage < min {
		cmd.Fatal(fmt.Errorf("coverage to low: expected >= %v%%\n", min))
		return
	}
	cmd.Exit(0)
}

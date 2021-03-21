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
		cli  = cmdline.NewParser(cmd.Args()...)
		help = cli.Flag("-h, --help")
		min  = cli.Option(
			"-min",
			"Fail if total coverage(%) is below min",
		).Float64(0.0)

		profile = cli.Required("PROFILE").String("")
	)
	uncover.OnlyShow = cli.Optional("FUNC").String("")
	log.SetFlags(0)

	switch {
	case help:
		cli.WriteUsageTo(cmd.Stderr())
		cmd.Exit(0)
		return // return is here so we can test

	case !cli.Ok():
		cmd.Fatal(cli.Error())
		return
	}

	profiles, err := uncover.ParseProfiles(profile)
	if err != nil {
		cmd.Fatal(err)
		return
	}

	coverage, _ := uncover.Report(profiles, os.Stdout)
	if coverage < min {
		cmd.Fatal(fmt.Errorf("coverage to low: expected >= %v%%\n", min))
		return
	}
	cmd.Exit(0)
}

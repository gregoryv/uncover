package main

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

import (
	"fmt"
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

		profile  = cli.Required("PROFILE").String("")
		onlyShow = cli.Optional("FUNC").String("")
	)

	switch {
	case help:
		cli.WriteUsageTo(cmd.Stderr())
		cmd.Stop(0)
		return

	case !cli.Ok():
		fmt.Fprintln(cmd.Stderr(), cli.Error())
		cmd.Stop(1)
		return
	}

	c := Command{
		min:      min,
		profile:  profile,
		onlyShow: onlyShow,
	}
	err := c.Run()
	if err != nil {
		fmt.Fprintln(cmd.Stderr(), err)
		cmd.Stop(1)
		return
	}
}

type Command struct {
	min      float64
	profile  string
	onlyShow string
}

// Run
func (me *Command) Run() error {
	profiles, err := uncover.ParseProfiles(me.profile)
	if err != nil {
		return err
	}
	coverage, _ := uncover.Report(profiles, os.Stdout)
	if coverage < me.min {
		return fmt.Errorf("coverage to low: expected >= %v%%\n", me.min)
	}
	return nil
}

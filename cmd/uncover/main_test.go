package main

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/gregoryv/cmdline/clitest"
)

var homeDir, _ = os.Getwd()

func Test_main_missing_profile(t *testing.T) {
	tc := clitest.NewShellT("uncover")
	cmd = tc // inject it
	log.SetOutput(tc.Stderr())
	main()
	if tc.ExitCode != 1 {
		t.Error(tc)
	}
}

func Test_main_bad_profile(t *testing.T) {
	tc := clitest.NewShellT("uncover", "jibberish")
	cmd = tc // inject it
	log.SetOutput(tc.Stderr())
	main()
	if tc.ExitCode != 1 {
		t.Error(tc)
	}
}

func Test_main_help(t *testing.T) {
	tc := clitest.NewShellT("uncover", "-h")
	cmd = tc
	log.SetOutput(tc.Stderr())
	main()
	if tc.ExitCode != 0 {
		t.Error(tc.Dump())
	}
}
func Test_main_version(t *testing.T) {
	tc := clitest.NewShellT("uncover", "-v")
	copyProfile(t)
	cmd = tc
	log.SetOutput(tc.Stderr())
	main()
	if tc.ExitCode != 0 {
		t.Error(tc.ExitCode, tc.Dump())
	}
}
func Test_main_ok(t *testing.T) {
	tc := clitest.NewShellT("uncover", "profile.out")
	copyProfile(t)
	cmd = tc
	log.SetOutput(tc.Stderr())
	main()
	if tc.ExitCode != 0 {
		t.Error(tc.ExitCode, tc.Dump())
	}
}

func Test_main_fails_min(t *testing.T) {
	tc := clitest.NewShellT("uncover", "-min", "100", "profile.out")
	copyProfile(t)
	cmd = tc
	log.SetOutput(tc.Stderr())
	main()
	if tc.ExitCode != 1 {
		t.Error(tc.ExitCode, tc.Dump())
	}
}

func copyProfile(t *testing.T) {
	t.Helper()
	out, err := os.Create("profile.out")
	if err != nil {
		t.Fatal(err)
	}
	in, err := os.Open(homeDir + "/testdata/profile.out")
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(out, in)
	in.Close()
	out.Close()
}

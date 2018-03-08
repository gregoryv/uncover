package uncover_test

import (
	"github.com/gregoryv/uncover"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

var (
	profile  string
	profiles []*uncover.Profile
)

func init() {
	fh, err := ioutil.TempFile("", "uncover")
	if err != nil {
		panic(err)
	}
	profile = fh.Name()
	_, err = exec.Command("go", "test", "-coverprofile", profile,
		"github.com/gregoryv/uncover/test").Output()
	if err != nil {
		panic(err)
	}
	fh.Close()
	profiles, err = uncover.ParseProfiles(profile)
	if err != nil {
		panic(err)
	}

}

func TestReport(t *testing.T) {
	cov, _ := uncover.Report(profiles, os.Stdout)
	exp := 50.0
	if cov != exp {
		t.Error("Expected %v, got %v", exp, cov)
	}
}

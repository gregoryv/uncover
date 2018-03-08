package cover_test

import (
	"github.com/gregoryv/cover"
	"io/ioutil"
	"os/exec"
	"testing"
)

var (
	profile  string
	profiles []*cover.Profile
)

func init() {
	fh, err := ioutil.TempFile("", "uncover")
	if err != nil {
		panic(err)
	}
	profile = fh.Name()
	_, err = exec.Command("go", "test", "-coverprofile", profile,
		"github.com/gregoryv/cover/test").Output()
	if err != nil {
		panic(err)
	}
	fh.Close()
	profiles, err = cover.ParseProfiles(profile)
	if err != nil {
		panic(err)
	}

}

func TestWrite(t *testing.T) {
	cover.WriteOutput(profiles)
}

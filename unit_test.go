package uncover

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

var (
	profile  string
	profiles []*Profile
)

func init() {
	fh, err := ioutil.TempFile("", "uncover")
	if err != nil {
		panic(err)
	}
	profile = fh.Name()
	_, err = exec.Command("go", "test", "-coverprofile", profile,
		"github.com/gregoryv/uncover/testdata").Output()
	if err != nil {
		panic(err)
	}
	fh.Close()
	profiles, err = ParseProfiles(profile)
	if err != nil {
		panic(err)
	}

}

func TestReport(t *testing.T) {
	cov, _ := Report(profiles, os.Stdout)
	exp := 50.0
	if cov != exp {
		t.Errorf("Expected %v, got %v", exp, cov)
	}
}

func TestParseProfiles(t *testing.T) {
	exp := 2
	if len(profiles) != exp { // files a.go and b.go
		t.Errorf("Expected %v, got %v", exp, len(profiles))
	}
}

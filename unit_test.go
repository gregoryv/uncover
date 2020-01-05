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
	if len(profiles) != 2 {
		// files a.go and b.go
		t.Errorf("testdata has two files: len(profiles) = %v", len(profiles))
	}

	cov, _ := Report(profiles, os.Stdout)
	exp := 50.0
	if cov != exp {
		t.Errorf("Expected %v, got %v", exp, cov)
	}
}

func Test_percent(t *testing.T) {
	ok := func(got, exp float64) {
		t.Helper()
		if got != exp {
			t.Errorf("Got %v, expected %v", got, exp)
		}
	}
	ok(percent(1, 0), 100)
	ok(percent(21, 100), 21)
	ok(percent(2, 200), 1)
}

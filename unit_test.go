package uncover

import (
	"bytes"
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
		"github.com/gregoryv/uncover/test").Output()
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

func TestParseProfile2_badmodes(t *testing.T) {
	cases := []string{"", "mode", "mode:", "mode: setj"}
	for _, c := range cases {
		r := bytes.NewBufferString(c)
		err := ParseProfile2(r)
		if err == nil {
			t.Errorf("Expected to fail on %q with bad mode", c)
		}
	}
}

func TestParseProfile2_ok(t *testing.T) {
	cases := []string{
		"mode: set",
		`mode: set
github.com/gregoryv/uncover/test/a.go:3.16,5.2 1 0
github.com/gregoryv/uncover/test/b.go:3.10,5.2 1 1
`,
	}
	for _, c := range cases {
		r := bytes.NewBufferString(c)
		err := ParseProfile2(r)
		if err != nil {
			t.Errorf("Expected %q to be ok, got: %v", c, err)
		}
	}
}

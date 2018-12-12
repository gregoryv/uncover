package uncover

import (
	"bytes"
	"testing"
)

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

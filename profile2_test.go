package uncover

import (
	"bytes"
	"testing"
)

func TestParseProfile2_badmodes(t *testing.T) {
	cases := []string{"", "mode"}
	for _, c := range cases {
		r := bytes.NewBufferString(c)
		err := ParseProfile2(r)
		if err == nil {
			t.Errorf("Expected to fail on %q with bad mode", c)
		}
	}
}

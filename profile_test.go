package cover_test

import (
	"github.com/gregoryv/cover"
	"testing"
)

func TestParseProfiles(t *testing.T) {
	res, err := cover.ParseProfiles(profile)
	if err != nil {
		t.Fatal(err)
	}
	exp := 2
	if len(res) != exp { // files a.go and b.go
		t.Errorf("Expected %v, got %v", exp, len(res))
	}
}

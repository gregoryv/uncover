package cover_test

import (
	"testing"
)

func TestParseProfiles(t *testing.T) {
	exp := 2
	if len(profiles) != exp { // files a.go and b.go
		t.Errorf("Expected %v, got %v", exp, len(profiles))
	}
}

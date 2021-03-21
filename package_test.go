package uncover

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	got := Version()
	if got != "unreleased" && !strings.Contains(got, ".") {
		t.Error(got)
	}
}

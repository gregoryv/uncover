package cover_test

import (
	"github.com/gregoryv/cover"
	"testing"
)

func TestWrite(t *testing.T) {
	cover.WriteOutput("/tmp/c.out")
	t.Fail()
}

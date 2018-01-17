package cover_test

import (
	"github.com/gregoryv/cover"
	"testing"
)

func TestWrite(t *testing.T) {
	cover.Write("/tmp/c.out")
	t.Fail()
}

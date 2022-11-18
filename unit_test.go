package uncover

import (
	"io/ioutil"
	"testing"
)

func TestReport(t *testing.T) {
	profiles, err := ParseProfiles("testdata/profile.out")
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 2 {
		// files a.go and b.go
		t.Errorf("testdata has two files: len(profiles) = %v", len(profiles))
	}

	cov, err := Report(profiles, ioutil.Discard)
	if err != nil {
		t.Error(err)
	}
	if cov == 0.0 {
		t.Errorf("Unexpected coverage %v", cov)
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

func Test_findFile(t *testing.T) {
	_, err := findFile("x")
	if err == nil {
		t.Error("Should fail when looking for non existing file")
	}
	_, err = findFile("github.com/gregoryv/uncover/testdata/a.go")
	if err != nil {
		t.Error(err)
	}
}

// ----------------------------------------

func Benchmark_Write(b *testing.B) {
	profiles, err := ParseProfiles("testdata/profile.out")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		Report(profiles, ioutil.Discard)
	}
}

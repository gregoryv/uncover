package uncover

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestParseProfile2_badmodes(t *testing.T) {
	r := bytes.NewBufferString("")
	err := ParseProfile2(r)
	if err == nil {
		t.Error("Expected to fail on bad mode")
	}
}

func ParseProfile2(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(at(':'))
	scanner.Scan()
	if scanner.Text() != "mode" {
		return fmt.Errorf("First line is missing mode")
	}
	return nil
}

func at(b byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, b); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, dropCR(data[0:i]), nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), dropCR(data), nil
		}
		// Request more data.
		return 0, nil, nil
	}
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

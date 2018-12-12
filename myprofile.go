package uncover

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ParseProfile2(profile io.Reader) error {
	r := bufio.NewReader(profile)
	mode, err := parseMode(r)
	if err != nil {
		return err
	}
	_ = mode
	return nil
}

func parseMode(r *bufio.Reader) (mode string, err error) {
	mode, err = r.ReadString(':')
	if err != nil || err == io.EOF {
		return
	}
	if mode != "mode:" {
		err = fmt.Errorf("Missing mode: %v", mode)
		return
	}
	mode, err = r.ReadString('\n')
	if err != nil && err != io.EOF {
		return
	}
	mode = strings.TrimSpace(mode)
	if mode != "set" {
		err = fmt.Errorf("Unsupported mode: %v", mode)
		return
	}
	return mode, nil
}

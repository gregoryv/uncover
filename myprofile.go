package uncover

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

func ParseProfile2(profile io.Reader) error {
	r := bufio.NewReader(profile)
	mode, err := r.ReadString(':')
	log.Printf("%q", mode)
	if err != nil || err == io.EOF {
		return err
	}
	if mode != "mode:" {
		return fmt.Errorf("Missing mode: %v", mode)
	}
	return nil
}

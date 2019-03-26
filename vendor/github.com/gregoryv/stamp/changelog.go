package stamp

import (
	"fmt"
	"regexp"
)

type Versioner interface {
	Version() (string, error)
}

type Changelog struct {
	content []byte
}

func NewChangelog(content []byte) *Changelog {
	return &Changelog{
		content: content,
	}
}

func (cl *Changelog) Version() (version string, err error) {
	re := regexp.MustCompile(`## \[(.*)\]`)
	found := re.FindSubmatch(cl.content)
	if len(found) > 0 {
		return string(found[1]), nil
	}
	return "", fmt.Errorf("No version found, missing %q line", "## [VERSION]")
}

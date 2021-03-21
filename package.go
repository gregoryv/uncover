package uncover

import (
	_ "embed"
	"strings"
)

//go:embed "changelog.md"
var changelog string

func Version() string {
	from := strings.Index(changelog, "## [")
	to := strings.Index(changelog[from:], "]")
	return string(changelog[from+4 : from+to])
}

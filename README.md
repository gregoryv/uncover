[uncover](https://pkg.go.dev/github.com/gregoryv/uncover) - Generate coverage reports from coverprofiles

Generates colorized coverage report to stdout of uncovered funcs.
Source originates from the golang cover tool.

## Quick start

Install

    go install github.com/gregoryv/uncover/cmd/uncover@latest

In your project test with coverage and show result

    go test -coverprofile /tmp/c.out
    uncover /tmp/c.out [FuncName]

![screenshot](screenshot.png)

Expect a minimum coverage

    uncover -min 80.0 /tmp/c.out

## Difference from `go tool cover`

The purpose of uncover is to focus your work on what remains to be
verified. Thus it only shows uncovered lines. It also excludes
unreachable code, ie. `func _()` which is sometimes used for compile
checks.

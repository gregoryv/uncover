[![Build Status](https://travis-ci.org/gregoryv/cover.svg?branch=master)](https://travis-ci.org/gregoryv/uncover)

[uncover](https://godoc.org/github.com/gregoryv/uncover) - Generate coverage reports from coverprofiles

Generates colorized coverage report to stdout of uncovered funcs.
Source originates from the golang cover tool.

## Quick start

Install

    go get -u github.com/gregoryv/uncover/...

In your project test with coverage and show result

    go test -coverprofile /tmp/c.out
    uncover /tmp/c.out

![screenshot](screenshot.png)

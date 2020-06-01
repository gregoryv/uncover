[![Build Status](https://travis-ci.org/gregoryv/uncover.svg?branch=master)](https://travis-ci.org/gregoryv/uncover)
[![codecov](https://codecov.io/gh/gregoryv/uncover/branch/master/graph/badge.svg)](https://codecov.io/gh/gregoryv/uncover)
[![Maintainability](https://api.codeclimate.com/v1/badges/83083a5e52d4ffad3288/maintainability)](https://codeclimate.com/github/gregoryv/uncover/maintainability)

[uncover](https://godoc.org/github.com/gregoryv/uncover) - Generate coverage reports from coverprofiles

Generates colorized coverage report to stdout of uncovered funcs.
Source originates from the golang cover tool.

## Quick start

Install

    go get -u github.com/gregoryv/uncover/...

In your project test with coverage and show result

    go test -coverprofile /tmp/c.out
    uncover /tmp/c.out [FuncName]

![screenshot](screenshot.png)

Expect a minimum coverage

    uncover -min 80.0 /tmp/c.out

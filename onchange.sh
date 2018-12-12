#!/bin/bash -e

path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

case $extension in
    go)
        goimports -w $path
        gofmt -w $path
        ;;
esac

go install github.com/gregoryv/uncover/cmd/uncover
go test -v -coverprofile /tmp/c.out -run=TestParseProfile2* .
#go tool cover -o /tmp/coverage.html -html /tmp/c.out
uncover /tmp/c.out ParseProfile2

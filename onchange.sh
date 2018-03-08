#!/bin/bash -e

path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

case $extension in
    go)
        gofmt -w $path
        ;;
esac

go install github.com/gregoryv/cover/cmd/uncover
go test -coverprofile /tmp/c.out .
#go tool cover -o /tmp/coverage.html -html /tmp/c.out
#uncover /tmp/c.out

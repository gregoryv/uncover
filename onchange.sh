#!/bin/bash -e

path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

case $extension in
    go)
        goimports -w $path
        ;;
esac

go install ./cmd/uncover
go test -v -coverprofile /tmp/c.out .
uncover /tmp/c.out
go test -v -coverprofile /tmp/c.out2 ./test
uncover /tmp/c.out2

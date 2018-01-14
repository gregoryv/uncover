#!/bin/sh

path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

case $extension in
    go)
        gofmt -w $path
        go test -cover .
        go build cmd/cover
        ;;
esac

#!/bin/bash
unset dirs files
dirs=$(go list ./... | grep -v sagapi$)
set -x -e
go build main.go
for d in $dirs
do
	go test -v $d
done

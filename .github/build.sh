#!/usr/bin/env bash

set -xe

export GOARCH="amd64"
export GOOS="darwin"
/usr/local/go/bin/go build -o committer.amd64 committer.go

export GOARCH="arm64"
/usr/local/go/bin/go build -o committer.arm64 committer.go

/usr/bin/lipo committer.amd64 committer.arm64 -create -output committer

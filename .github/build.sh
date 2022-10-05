#!/usr/bin/env bash

set -xe

# Build for darwin
export GOOS="darwin"

export GOARCH="amd64"
/usr/local/go/bin/go build -o committer.amd64 committer.go

export GOARCH="arm64"
/usr/local/go/bin/go build -o committer.arm64 committer.go

/usr/bin/lipo committer.amd64 committer.arm64 -create -output committer

# Also build for linux
export GOOS="linux"

export GOARCH="arm64"
/usr/local/go/bin/go build -o committer.linux-arm64 committer.go
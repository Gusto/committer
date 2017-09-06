#!/bin/bash

VERSION=$(cat VERSION)
UPLOADED_PKG_FOUND=$(aws s3 --region 'us-west-2' ls 's3://vpc-access/' | grep committer-$VERSION)

if [ "$UPLOADED_PKG_FOUND" ]; then
  echo "Committer@${VERSION} already uploaded. Skipping build."
  exit 0
else
  set -xe

  export GOARCH="amd64"
  export GOOS="darwin"
  go build committer.go

  mv committer /go/bin/committer-$VERSION
  aws s3 cp /go/bin/committer-$VERSION s3://vpc-access/
  echo "Uploaded Committer@${VERSION} to vpc-access S3 bucket!"
fi

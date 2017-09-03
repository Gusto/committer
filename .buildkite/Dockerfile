FROM golang:1.8

RUN go get github.com/tools/godep

RUN groupadd -r buildkite-agent && useradd -u 9999 -r -g buildkite-agent buildkite-agent
RUN apt-get update && apt-get install -y --no-install-recommends \
		awscli \
	&& rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED 0
ENV WD /go/src/github.com/gusto/committer

RUN mkdir -p $WD
ADD . $WD/
WORKDIR $WD

USER buildkite-agent

RUN go-wrapper download

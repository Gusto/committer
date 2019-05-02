FROM golang:1.8

RUN printf 'deb http://deb.debian.org/debian jessie main\n\
deb http://security.debian.org/debian-security jessie/updates main' > /etc/apt/sources.list

RUN go get github.com/tools/godep

RUN apt-get update && apt-get install -y --no-install-recommends \
		awscli \
	&& rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED 0
ENV WD /go/src/github.com/gusto/committer

RUN mkdir -p $WD
ADD . $WD/
WORKDIR $WD

RUN go-wrapper download

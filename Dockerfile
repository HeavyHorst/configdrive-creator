FROM alpine

MAINTAINER Rene Kaufmann <kaufmann.r@gmail.com>

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV GO15VENDOREXPERIMENT 1

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
ADD . /go/src/github.com/HeavyHorst/configdrive-creator
WORKDIR /go/src/github.com/HeavyHorst/configdrive-creator

RUN set -ex && \
    apk --update add --no-cache cdrkit && \
	apk --update add --virtual build-deps go && \
	go build && \
    apk del build-deps

EXPOSE 3000

# Startup
ENTRYPOINT ./configdrive-creator

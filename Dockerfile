FROM golang:alpine

MAINTAINER Rene Kaufmann <kaufmann.r@gmail.com>

ENV GO15VENDOREXPERIMENT 1

RUN set -ex \
        && apk add --no-cache \
		cdrkit

ADD . /go/src/github.com/HeavyHorst/configdrive-creator
RUN go install github.com/HeavyHorst/configdrive-creator

# Expose 3000
EXPOSE 3000

# Startup
ENTRYPOINT /go/bin/configdrive-creator

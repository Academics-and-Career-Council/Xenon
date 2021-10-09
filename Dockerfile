FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git make build-essential
RUN mkdir /server
WORKDIR /server
COPY ./ ./
RUN make install
FROM scratch
LABEL MAINTAINER="Shivam Malhotra"
LABEL VERSION="0.0.1"

COPY --from=builder /server/build/xenon /go/bin/xenon
ENTRYPOINT ["/go/bin/xenon"]
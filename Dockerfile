FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git 
RUN mkdir /server
WORKDIR /server
COPY ./ ./
RUN go mod download
RUN	go build -o build/xenon main.go
FROM scratch
LABEL MAINTAINER="Shivam Malhotra"
LABEL VERSION="0.0.1"

COPY --from=builder /server/build/xenon /go/bin/xenon
ENTRYPOINT ["/go/bin/xenon"]
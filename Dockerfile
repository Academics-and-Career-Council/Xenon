FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git build-base
RUN apk add -U --no-cache ca-certificates && update-ca-certificates
RUN apk add --no-cache zeromq-dev musl-dev pkgconfig alpine-sdk libsodium-dev libzmq-static libsodium-static
RUN mkdir /server
WORKDIR /server
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/xenon
RUN CGO_LDFLAGS="$CGO_LDFLAGS -lstdc++ -lm -lsodium" \
  CGO_ENABLED=1 \
  GOOS=linux \
  go build -v -a --ldflags '-extldflags "-static" -v' \
  -o build/stargazer

FROM scratch
LABEL MAINTAINER="Shivam Malhotra"
LABEL VERSION="0.0.1"
COPY --from=builder /server/build/xenon /go/bin/xenon
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/xenon"]

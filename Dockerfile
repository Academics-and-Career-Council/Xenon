FROM golang:1.16-alpine
COPY ./ ./
RUN go mod download
RUN make build
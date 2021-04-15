FROM golang:latest
LABEL MAINTAINER="Shivam Malhotra"
LABEL VERSION="0.0.1"

# Build the server
RUN mkdir /server
WORKDIR /server
COPY ./ ./
RUN make install
run:
	go run main.go

build:
	go build -o build/xenon main.go 

install:
	go mod download
	go build -o build/xenon main.go 